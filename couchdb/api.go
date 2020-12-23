package couchdb

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
)

// Create / UpdateById / GetById / DeleteById / GetByKey / Search

var (
	NoIdData = errors.New("no id data")
)

type CouchDB struct {
	HostPort string
	User     string
	Passwd   string
	db       string
	function string
}

func NewCouchDB(hostPort, user, passwd string) *CouchDB {
	c := &CouchDB{
		HostPort: hostPort,
		User:     user,
		Passwd:   passwd,
	}
	return c
}

func (c *CouchDB) basicAuth() map[string]string {
	auth := fmt.Sprintf("%s:%s", c.User, c.Passwd)
	enstr := base64.StdEncoding.EncodeToString([]byte(auth))
	headerVal := fmt.Sprintf("Basic %s", enstr)
	headers := make(map[string]string)
	headers["Authorization"] = headerVal
	return headers
}

func (c *CouchDB) SetDBName(name string) error {
	c.db = name

	// insure db exist, if not, create it
	url := c.reqURL()
	authHeader := c.basicAuth()
	resp, err := util.HttpGet(url, nil, authHeader)
	if err != nil {
		consolelog.Logger.Errorf("", "get db[%s] failed, %s", name, err.Error())
		return err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		// create it
		err = c.Create("", nil)
		if err != nil {
			consolelog.Logger.Errorf("", "put database[%s] failed, %s", name, err.Error())
			return err
		}
		// create success
		return nil
	} else {
		err = fmt.Errorf("create database[%s] failed, statuscode[%d]", name, resp.StatusCode)
		consolelog.Logger.Error("", err.Error())
		return err
	}
}

func (c *CouchDB) baseURI() string {
	return fmt.Sprintf("http://%s", c.HostPort)
}

func (c *CouchDB) reqURL() string {
	return fmt.Sprintf("%s/%s", c.baseURI(), c.db)
}

func (c *CouchDB) docURL(id string) string {
	return fmt.Sprintf("%s/%s", c.reqURL(), id)
}

func (c *CouchDB) Create(id string, data []byte) error {
	c.function = "Create"
	url := c.reqURL()
	if id != "" {
		url = c.docURL(id)
	}
	authHeader := c.basicAuth()
	resp, err := util.HttpPut(url, data, authHeader)
	if err != nil {
		consolelog.Logger.Errorf("", "db Create failed, %s", err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		consolelog.Logger.Errorf("", "db create failed,, read response body failed, %s", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		consolelog.Logger.Errorf("db create failed, reason: %s", string(body))
		return errors.New(string(body))
	}
	consolelog.Logger.Debugf("", "db create success, %s", string(body))

	return nil
}

func (c *CouchDB) UpdateById(id string, data []byte) ([]byte, error) {
	c.function = "UpdateById"
	url := c.docURL(id)
	authHeader := c.basicAuth()
	resp, err := util.HttpPut(url, data, authHeader)
	if err != nil {
		consolelog.Logger.Errorf("", "%s[%s] failed, %s", c.function, id, err.Error())
		return nil, err
	}
	body, err := c.processHttpResp(resp)
	return body, err
}

func (c *CouchDB) DeleteById(id, rev string) ([]byte, error) {
	c.function = "DeleteById"
	url := c.docURL(id)
	url = fmt.Sprintf("%s?rev=%s", url, rev)
	authHeader := c.basicAuth()
	resp, err := util.HttpDelete(url, nil, authHeader)
	if err != nil {
		consolelog.Logger.Errorf("", "%s[%s] failed, %s", c.function, id, err.Error())
		return nil, err
	}
	body, err := c.processHttpResp(resp)
	return body, err
}

func (c *CouchDB) GetById(id string) ([]byte, error) {
	c.function = "GetById"
	url := c.docURL(id)
	authHeader := c.basicAuth()
	resp, err := util.HttpGet(url, nil, authHeader)
	if err != nil {
		consolelog.Logger.Errorf("", "%s[%s] failed, %s", c.function, id, err.Error())
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, NoIdData
	}

	body, err := c.processHttpResp(resp)
	return body, err
}

func (c *CouchDB) GetByKey() {

}

func (c *CouchDB) Search() {

}

func (c *CouchDB) processHttpResp(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		consolelog.Logger.Errorf("", "%s read response body failed, %s", c.function, err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("%s response status code[%d]", c.function, resp.StatusCode)
		consolelog.Logger.Error("", err.Error())
		return nil, err
	}

	return body, nil
}
