package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const HttpClientErrCode = 0

const reqTimeout time.Duration = 10 // seconds

// HTTP METHOD WRAPPER
type ClientRequest struct {
	Ctx     context.Context
	Url     string
	Query   map[string][]string
	Headers map[string]string
	Body    []byte
	Timeout int

	V interface{} // response body unmarshal struct

	Debug  bool // if true, logmiddleware response body
	method string
}

type ClientResponse struct {
	Code   int   // http status code and default err code(0)
	Err    error // when program err occurred
	Body   []byte
	Raw    *http.Response // be careful, response.Body can be read exactly once
	Logger *zerolog.Logger
}

func Get(req *ClientRequest) *ClientResponse {
	req.method = http.MethodGet
	resp := httpRequest(req)
	return resp
}

func Post(req *ClientRequest) *ClientResponse {
	req.method = http.MethodPost
	resp := httpRequest(req)
	return resp
}

func Put(req *ClientRequest) *ClientResponse {
	req.method = http.MethodPut
	resp := httpRequest(req)
	return resp
}

func Delete(req *ClientRequest) *ClientResponse {
	req.method = http.MethodDelete
	resp := httpRequest(req)
	return resp
}

func API(req *ClientRequest, method string) *ClientResponse {
	req.method = method
	resp := httpRequest(req)
	return resp
}

func httpRequest(req *ClientRequest) *ClientResponse {
	var err error
	startT := time.Now()
	resp := &ClientResponse{
		Code: HttpClientErrCode,
		Err:  err,
	}

	lctx := zerolog.Ctx(req.Ctx)
	logger := lctx.With().Str("method", req.method).Str("url", req.Url).Logger()
	resp.Logger = &logger
	logger.Debug().Msg("start http request...")

	// generate req
	var newReq *http.Request
	if req.Body != nil {
		newReq, err = http.NewRequest(req.method, req.Url, bytes.NewBuffer(req.Body))
	} else {
		newReq, err = http.NewRequest(req.method, req.Url, nil)
	}
	if err != nil {
		logger.Error().Err(err).Send()
		resp.Err = err
		return resp
	}

	// process url query string
	if req.Query != nil {
		urlV := url.Values{}
		for k, vs := range req.Query {
			if len(vs) == 1 {
				urlV.Set(k, vs[0])
			} else if len(vs) > 1 {
				for _, v := range vs {
					urlV.Add(k, v)
				}
			}
		}
		newReq.URL.RawQuery = urlV.Encode()
		logger.Debug().Str("fullUrl", newReq.URL.String()).Send()
	}

	// process headers
	for k, v := range req.Headers {
		newReq.Header.Add(k, v)
	}

	// timeout
	timeout := reqTimeout
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout)
	}
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	doResp, err := client.Do(newReq)
	if err != nil {
		logger.Error().Err(err).Send()
		resp.Err = err
		return resp
	}
	resp.Raw = doResp
	resp.Code = doResp.StatusCode

	// print debug info
	logger.Debug().Int("statusCode", doResp.StatusCode).Str("elapsed", time.Since(startT).String()).Msg("http response")

	respBody, err := ioutil.ReadAll(doResp.Body)
	if err != nil {
		logger.Error().Err(err).Send()
		resp.Err = err
		return resp
	}
	defer doResp.Body.Close()
	resp.Body = respBody

	if req.Debug {
		// print response body
		logger.Debug().RawJSON("responseBody", respBody).Send()
	}

	// check if need unmarshal response body
	if req.V != nil {
		err = json.Unmarshal(respBody, &req.V)
		if err != nil {
			logger.Error().Err(err).Send()
			resp.Err = err
			return resp
		}
	}

	return resp
}
