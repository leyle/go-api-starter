package httpclient

import (
	"context"
	"encoding/json"
	"github.com/leyle/go-api-starter/logmiddleware"
	"testing"
)

// https://httpbin.org/get
// https://httpbin.org/post
// docker run -p 8000:80 kennethreitz/httpbin

func TestGet(t *testing.T) {
	url := "http://192.168.2.40:8000/get"
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	ctx := context.Background()
	lctx := logger.WithContext(ctx)

	query := make(map[string][]string)
	query["name"] = []string{"jack", "telsa"}
	query["age"] = []string{"12", "23"}

	req := &ClientRequest{
		Ctx:   lctx,
		Url:   url,
		Query: query,
		Debug: true,
	}

	resp := Get(req)
	if resp.Err != nil {
		t.Fatal(resp.Err)
	}

	t.Log(resp.Code)
	t.Log(string(resp.Body))
}

func TestPost(t *testing.T) {
	url := "http://192.168.2.40:8000/post"
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	ctx := context.Background()
	lctx := logger.WithContext(ctx)

	query := make(map[string][]string)
	query["name"] = []string{"jack", "telsa"}
	query["age"] = []string{"12", "23"}

	type Form struct {
		Name   string `json:"name"`
		Passwd string `json:"passwd"`
	}

	form := &Form{
		Name:   "jack",
		Passwd: "passwd",
	}

	type RespForm struct {
		Json *Form `json:"json"`
	}

	reqData, err := json.Marshal(form)
	if err != nil {
		t.Fatal(err)
	}

	var respForm *RespForm

	req := &ClientRequest{
		Ctx:   lctx,
		Url:   url,
		Query: query,
		Debug: true,
		Body:  reqData,
		V:     &respForm,
	}

	resp := Post(req)
	if resp.Err != nil {
		t.Fatal(resp.Err)
	}

	t.Log(resp.Code)
	t.Log(string(resp.Body))

	t.Log("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	t.Log(respForm.Json.Name, respForm.Json.Passwd)
}
