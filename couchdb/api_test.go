package couchdb

import (
	"context"
	"encoding/json"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/go-api-starter/util"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

// fabric ca user
type CaUser struct {
	// enrollId equals User.Id
	EnrollId string        `json:"enrollId"`
	Rev      string        `json:"_rev,omitempty"`
	Secret   string        `json:"secret"`
	Created  *util.CurTime `json:"created"`
}

var (
	couchdbHost   = "localhost:5984"
	couchdbUser   = "admin"
	couchdbPasswd = "passwd"
	couchdbName   = "dev"
)

var opt = &CouchDBOption{
	HostPort: couchdbHost,
	User:     couchdbUser,
	Passwd:   couchdbPasswd,
	Protocol: "http",
}

func getContext() context.Context {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	bctx := context.Background()
	ctx := logger.WithContext(bctx)
	return ctx
}

func TestClient_CreateDatabase(t *testing.T) {
	ctx := getContext()
	client := New(opt, couchdbName)
	err := client.CreateDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateIndex(t *testing.T) {
	ctx := getContext()
	client := New(opt, couchdbName)

	fields := []string{"enrollId", "secret", "created.second"}
	err := client.CreateIndex(ctx, fields)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateDoc(t *testing.T) {
	ctx := getContext()
	client := New(opt, couchdbName)

	ua := &CaUser{
		EnrollId: logmiddleware.GenerateReqId(),
		Secret:   "secret",
		Created:  util.GetCurTime(),
	}

	data, _ := json.Marshal(&ua)

	err := client.CreateDoc(ctx, ua.EnrollId, data)
	if err != nil {
		t.Fatal(err)
	}

}

func TestClient_SearchByKey(t *testing.T) {
	ctx := getContext()
	client := New(opt, couchdbName)

	/*
		selector := map[string]interface{}{
			"enrollId": "d7ee535a-445d-4476-bf54-0b3169ba798f",
		}
	*/

	/*
		selector := map[string]interface{}{
			"enrollId": map[string]string{
				"$regex": "5fed8",
			},
		}
	*/

	selector := map[string]interface{}{
		"created.second": map[string]interface{}{
			"$lte": 1609404283,
		},
	}

	/*
		searchReq := &SearchRequest{
			Selector: selector,
		}
	*/

	sort := []map[string]string{
		{
			"created.second": "asc",
		},
	}
	searchReq := &SearchRequest{
		Selector: selector,
		Sort:     sort,
		Limit:    10,
		Skip:     0,
	}

	type UAResp struct {
		Docs []*CaUser `json:"docs"`
	}

	var uaResp *UAResp

	resp, err := client.Search(ctx, searchReq, &uaResp)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.Bookmark)
}

func TestClient_GetById(t *testing.T) {
	ctx := getContext()
	client := New(opt, couchdbName)
	id := "5fed8eb3310e433cf278f5ba"
	var ua *CaUser
	_, err := client.GetById(ctx, id, &ua)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ua.EnrollId, ua.Created.HumanTime)
}

func TestClient_UpdateById(t *testing.T) {
	// first we need to get rev id
	ctx := getContext()
	client := New(opt, couchdbName)
	id := "5fed8eb3310e433cf278f5ba"
	var ua *CaUser
	_, err := client.GetById(ctx, id, &ua)
	if err != nil {
		t.Fatal(err)
	}

	ua.EnrollId = util.GenerateDataId()
	// ua.EnrollId = id
	ua.Created = util.GetCurTime()
	ua.Secret = logmiddleware.GenerateReqId()

	data, _ := json.Marshal(&ua)

	_, err = client.UpdateById(ctx, id, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteById(t *testing.T) {
	// first we need to get rev id
	ctx := getContext()
	client := New(opt, couchdbName)
	id := "5fed8eb3310e433cf278f5ba"
	var ua *CaUser
	_, err := client.GetById(ctx, id, &ua)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteById(ctx, id, ua.Rev)
	if err != nil {
		t.Fatal(err)
	}

}
