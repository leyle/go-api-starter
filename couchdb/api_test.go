package couchdb

import (
	"context"
	"encoding/json"
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
	couchdbHost   = "192.168.2.40:5984"
	couchdbUser   = "admin"
	couchdbPasswd = "passwd"
	couchdbName   = "dev"
)

var id = "f829bd13-e816-4b32-90fd-d956a65913c1"

func getContext() context.Context {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	bctx := context.Background()
	ctx := logger.WithContext(bctx)
	return ctx
}

func getClient() *CouchDB {
	return NewCouchDB(couchdbHost, couchdbUser, couchdbPasswd, couchdbName)
}

func TestCouchDB_Create(t *testing.T) {
	ctx := getContext()
	c := getClient()
	err := c.SetDBName(ctx, couchdbName)
	if err != nil {
		t.Fatal(err)
	}

	user := &CaUser{
		EnrollId: util.GenerateDataId(),
		Secret:   util.GetCurTime().HumanTime,
		Created:  util.GetCurTime(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Create(ctx, user.EnrollId, data)
	if err != nil {
		t.Fatal(err)
	}

	user.EnrollId = util.GenerateDataId()
	err = c.Create(ctx, user.EnrollId, data)
	if err != nil {
		t.Fatal(err)
	}

	id = user.EnrollId
}

func TestCouchDB_GetById(t *testing.T) {
	ctx := getContext()
	c := getClient()
	err := c.SetDBName(ctx, couchdbName)
	if err != nil {
		t.Fatal(err)
	}

	var user *CaUser
	_, err = c.GetById(ctx, id, &user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("id", user.EnrollId, "rev", user.Rev, "t", user.Created.Second)
}

func TestCouchDB_UpdateById(t *testing.T) {
	ctx := getContext()
	c := getClient()
	err := c.SetDBName(ctx, couchdbName)
	if err != nil {
		t.Fatal(err)
	}
	var user *CaUser
	_, err = c.GetById(ctx, id, &user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("id", user.EnrollId, "rev", user.Rev, "t", user.Created.Second)
	t.Log("Prepare for update")

	user.EnrollId = user.EnrollId + "UPDATE"
	user.Created = util.GetCurTime()
	updData, err := json.Marshal(&user)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.UpdateById(ctx, id, updData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp))

	id = user.EnrollId
}

func TestCouchDB_DeleteById(t *testing.T) {
	ctx := getContext()
	c := getClient()
	err := c.SetDBName(ctx, couchdbName)
	if err != nil {
		t.Fatal(err)
	}

	var user *CaUser
	_, err = c.GetById(ctx, id, &user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("id", user.EnrollId, "rev", user.Rev, "t", user.Created.Second)
	t.Log("Prepare for delete")

	resp, err := c.DeleteById(ctx, id, user.Rev)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(resp))

	// recheck if deleted
	checkResp, err := c.GetById(ctx, id, &user)
	if err == NoIdData {
		t.Log("Delete data success")
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(checkResp))
}
