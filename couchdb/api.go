package couchdb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/leyle/go-api-starter/httpclient"
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
	DBName   string

	// method name
	function string
}

func NewCouchDB(hostPort, user, passwd, dbName string) *CouchDB {
	c := &CouchDB{
		HostPort: hostPort,
		User:     user,
		Passwd:   passwd,
		DBName:   dbName,
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

func (c *CouchDB) SetDBName(ctx context.Context, name string) error {
	c.DBName = name

	// insure DBName exist, if not, create it
	url := c.reqURL()
	authHeader := c.basicAuth()
	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeader,
		Debug:   true,
	}

	resp := httpclient.Get(cReq)
	if resp.Err != nil {
		resp.Logger.Error().Msg("SetDBName failed")
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return nil
	}

	if resp.Code == http.StatusNotFound {
		err := c.Create(ctx, "", nil)
		if err != nil {
			resp.Logger.Error().Str("dbname", name).Err(err).Msg("create database name failed")
			return err
		}
		return nil
	}

	// something wrong happend, logmiddleware and return err
	resp.Logger.Error().Str("dbname", name).Msg("create failed")
	return errors.New(string(resp.Body))
}

func isHttpStatusCodeOK(code int) bool {
	if code == http.StatusOK || code == http.StatusCreated {
		return true
	}
	return false
}

func (c *CouchDB) baseURI() string {
	return fmt.Sprintf("http://%s", c.HostPort)
}

func (c *CouchDB) reqURL() string {
	return fmt.Sprintf("%s/%s", c.baseURI(), c.DBName)
}

func (c *CouchDB) docURL(id string) string {
	return fmt.Sprintf("%s/%s", c.reqURL(), id)
}

// create DBName or item
func (c *CouchDB) Create(ctx context.Context, id string, data []byte) error {
	c.function = "Create"
	url := c.reqURL()
	if id != "" {
		url = c.docURL(id)
	}
	authHeader := c.basicAuth()

	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeader,
		Body:    data,
		Debug:   true,
	}

	resp := httpclient.Put(cReq)
	if resp.Err != nil {
		resp.Logger.Err(resp.Err).Str("id", id).Msg("Create data failed")
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return nil
	}

	resp.Logger.Error().Str("id", id).Msg("create failed")
	return errors.New(string(resp.Body))
}

func (c *CouchDB) UpdateById(ctx context.Context, id string, data []byte) ([]byte, error) {
	c.function = "UpdateById"
	url := c.docURL(id)
	authHeader := c.basicAuth()

	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeader,
		Body:    data,
		Debug:   true,
	}

	resp := httpclient.Put(cReq)
	if resp.Err != nil {
		resp.Logger.Error().Str("method", c.function).Err(resp.Err).Msg("update failed")
		return resp.Body, resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return resp.Body, nil
	}

	return resp.Body, errors.New(string(resp.Body))
}

func (c *CouchDB) DeleteById(ctx context.Context, id, rev string) ([]byte, error) {
	c.function = "DeleteById"
	url := c.docURL(id)
	url = fmt.Sprintf("%s?rev=%s", url, rev)
	authHeader := c.basicAuth()

	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeader,
		Debug:   true,
	}

	resp := httpclient.Delete(cReq)
	if resp.Err != nil {
		resp.Logger.Error().Str("method", c.function).Err(resp.Err).Msg("delete failed")
		return resp.Body, resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return resp.Body, nil
	}

	// error occurred
	return resp.Body, errors.New(string(resp.Body))
}

func (c *CouchDB) GetById(ctx context.Context, id string, v interface{}) ([]byte, error) {
	c.function = "GetById"
	url := c.docURL(id)
	authHeader := c.basicAuth()

	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeader,
		V:       v,
		Debug:   true,
	}

	resp := httpclient.Get(cReq)
	if resp.Err != nil {
		resp.Logger.Error().Err(resp.Err).Str("method", c.function).Str("id", id).Msg("get failed")
		return resp.Body, resp.Err
	}

	if resp.Code == http.StatusNotFound {
		return resp.Body, NoIdData
	}

	if isHttpStatusCodeOK(resp.Code) {
		return resp.Body, nil
	}

	return resp.Body, errors.New(string(resp.Body))
}

func (c *CouchDB) GetByKey() {

}

func (c *CouchDB) Search() {

}
