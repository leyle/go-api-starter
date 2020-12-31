package couchdb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leyle/go-api-starter/httpclient"
	"github.com/rs/zerolog"
	"net/http"
)

// Create / UpdateById / GetById / DeleteById / Search
var (
	NoIdData = errors.New("no id data")
)

type CouchDBOption struct {
	HostPort string
	User     string
	Passwd   string

	// http or https
	Protocol string
}

type CouchDBClient struct {
	Opt *CouchDBOption
	db  string
	// method name, used by logger
	method string
}

func New(opt *CouchDBOption, db string) *CouchDBClient {
	if opt.Protocol == "" {
		opt.Protocol = "http"
	}
	return &CouchDBClient{
		Opt: opt,
		db:  db,
	}
}

func (c *CouchDBOption) basicAuth() map[string]string {
	auth := fmt.Sprintf("%s:%s", c.User, c.Passwd)
	enstr := base64.StdEncoding.EncodeToString([]byte(auth))
	headerVal := fmt.Sprintf("Basic %s", enstr)

	headers := map[string]string{
		"Authorization": headerVal,
		"Content-Type":  "application/json",
	}
	return headers
}

func isHttpStatusCodeOK(code int) bool {
	if code == http.StatusOK || code == http.StatusCreated {
		return true
	}
	return false
}

func (c *CouchDBClient) basicAuth() map[string]string {
	return c.Opt.basicAuth()
}

func (c *CouchDBClient) dbURL() string {
	return fmt.Sprintf("%s://%s/%s", c.Opt.Protocol, c.Opt.HostPort, c.db)
}

func (c *CouchDBClient) docIdURL(docId string) string {
	return fmt.Sprintf("%s/%s", c.dbURL(), docId)
}

func (c *CouchDBClient) createIndexURL() string {
	return fmt.Sprintf("%s/%s", c.dbURL(), "_index")
}

func (c *CouchDBClient) searchURL() string {
	return fmt.Sprintf("%s/%s", c.dbURL(), "_find")
}

func (c *CouchDBClient) CreateDatabase(ctx context.Context) error {
	c.method = "CreateDatabase"
	headers := c.basicAuth()
	url := c.dbURL()

	req := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: headers,
		Debug:   true,
	}

	resp := httpclient.Get(req)
	if resp.Err != nil {
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Send()
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		resp.Logger.Debug().Str("action", c.method).Str("database", c.db).Msg("database already exist")
		return nil
	}

	if resp.Code == http.StatusNotFound {
		// db not exist, create it
		err := c.CreateDoc(ctx, "", nil)
		if err != nil {
			resp.Logger.Error().Err(err).Str("action", c.method).Str("database", c.db).Msg("create database failed")
			return err
		}

		resp.Logger.Debug().Str("action", c.method).Str("database", c.db).Msg("create database success")
		return nil
	} else {
		// other errors
		err := fmt.Errorf("statusCode[%d], body:%s", resp.Code, string(resp.Body))
		resp.Logger.Error().Err(err).Str("action", c.method).Str("database", c.db).Msg("create database failed")
		return err
	}
}

func (c *CouchDBClient) CreateIndex(ctx context.Context, fields []string) error {
	for _, field := range fields {
		err := c.createIndex(ctx, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CouchDBClient) createIndex(ctx context.Context, field string) error {
	c.method = "CreateIndex"
	headers := c.basicAuth()
	url := c.createIndexURL()

	name := "index-" + field
	body := map[string]interface{}{
		"name": name,
		"type": "json",
		"index": map[string][]string{
			"fields": {field},
		},
	}
	data, _ := json.Marshal(body)

	req := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: headers,
		Body:    data,
		Debug:   true,
	}

	resp := httpclient.Post(req)
	if resp.Err != nil {
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Str("database", c.db).Send()
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		resp.Logger.Info().Str("action", c.method).Str("database", c.db).Send()
		return nil
	}

	err := fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Str("database", c.db).Send()
	return err
}

// create DBName or item
func (c *CouchDBClient) CreateDoc(ctx context.Context, id string, data []byte) error {
	c.method = "Create"
	url := c.dbURL()
	if id != "" {
		url = c.docIdURL(id)
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
		resp.Logger.Err(resp.Err).Str("action", c.method).Str("id", id).Msg("Create data failed")
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return nil
	}

	err := fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Str("id", id).Send()
	return err
}

func (c *CouchDBClient) UpdateById(ctx context.Context, id string, data []byte) ([]byte, error) {
	c.method = "UpdateById"
	url := c.docIdURL(id)
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
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Str("id", id).Msg("update failed")
		return resp.Body, resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return resp.Body, nil
	}

	err := fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Str("id", id).Send()
	return resp.Body, err
}

func (c *CouchDBClient) DeleteById(ctx context.Context, id, rev string) error {
	c.method = "DeleteById"
	url := c.docIdURL(id)
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
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Str("id", id).Msg("delete failed")
		return resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		return nil
	}

	// error occurred
	err := fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Str("id", id).Send()
	return err
}

func (c *CouchDBClient) GetById(ctx context.Context, id string, v interface{}) ([]byte, error) {
	c.method = "GetById"
	url := c.docIdURL(id)
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
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Str("id", id).Msg("get failed")
		return resp.Body, resp.Err
	}

	if resp.Code == http.StatusNotFound {
		return resp.Body, NoIdData
	}

	if isHttpStatusCodeOK(resp.Code) {
		return resp.Body, nil
	}

	err := fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Str("id", id).Send()
	return resp.Body, err
}

type SearchRequest struct {
	Selector       interface{} `json:"selector"`
	Sort           interface{} `json:"sort,omitempty"`
	Limit          int         `json:"limit"`
	Skip           int         `json:"skip"`
	ExecutionStats bool        `json:"execution_stats,omitempty"`
}

func (sr *SearchRequest) marshal() []byte {
	data, err := json.Marshal(sr)
	if err != nil {
		return []byte(err.Error())
	}
	return data
}

type SearchResponse struct {
	Docs     json.RawMessage `json:"docs"`
	Bookmark string          `json:"bookmark"`
}

func (c *CouchDBClient) Search(ctx context.Context, searchReq *SearchRequest, v interface{}) (*SearchResponse, error) {
	c.method = "SearchByKey"
	url := c.searchURL()
	authHeaders := c.basicAuth()

	const (
		minLimit = 1
		minSkip  = 0
	)

	if searchReq.Limit < minLimit {
		searchReq.Limit = minLimit
	}
	if searchReq.Skip < minSkip {
		searchReq.Skip = minSkip
	}
	logger := zerolog.Ctx(ctx)

	logger.Debug().Str("action", c.method).RawJSON("searchRequest", searchReq.marshal()).Send()

	data, err := json.Marshal(searchReq)
	if err != nil {
		logger.Error().Err(err).Str("action", c.method).Msg("marshal input search request failed")
		return nil, err
	}

	req := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     url,
		Headers: authHeaders,
		Body:    data,
		Timeout: 30,
		V:       v,
		Debug:   true,
	}

	resp := httpclient.Post(req)
	if resp.Err != nil {
		resp.Logger.Error().Err(resp.Err).Str("action", c.method).Send()
		return nil, resp.Err
	}

	if isHttpStatusCodeOK(resp.Code) {
		// parse data
		var sr *SearchResponse
		err = json.Unmarshal(resp.Body, &sr)
		if err != nil {
			resp.Logger.Error().Err(err).Str("action", c.method).Msg("unmarshal search result failed")
			return nil, err
		}
		return sr, nil
	}

	err = fmt.Errorf("statusCode[%d], body[%s]", resp.Code, string(resp.Body))
	resp.Logger.Error().Err(err).Str("action", c.method).Send()
	return nil, err
}
