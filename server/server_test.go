package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/istyle-inc/multimissile"
	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/jsonrpc"
	"github.com/stretchr/testify/assert"
)

func responseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprintf("%s", r.URL))
}

func TestmultimissileHandler(t *testing.T) {
	assert := assert.New(t)

	muxGET := http.NewServeMux()
	muxGET.HandleFunc("/user/get", responseHandler)
	muxGET.HandleFunc("/item/get", responseHandler)
	go http.ListenAndServe(":30001", muxGET)
	muxPOST := http.NewServeMux()
	muxPOST.HandleFunc("/item/update", responseHandler)
	go http.ListenAndServe(":30002", muxPOST)

	var err error
	msl.Config, err = config.Load("../config/example.toml")
	assert.Nil(err)

	ts := httptest.NewServer(http.HandlerFunc(multimissileHandler))
	defer ts.Close()

	payload := `[
    {"jsonrpc": "2.0", "endpoint": "ep-1",                        "path": "/user/get",    "params": { "user_id": 1 },                   "id": "1"},
    {"jsonrpc": "2.0", "endpoint": "ep-1", "http_method": "GET",  "path": "/item/get",    "params": { "item_id": 2 },                   "id": "2"},
    {"jsonrpc": "2.0", "endpoint": "ep-2", "http_method": "POST", "path": "/item/update", "params": { "item_id": 2, "desc": "update" }, "id": "3"}
]
`
	res, err := http.Post(ts.URL, "application/json", strings.NewReader(payload))
	assert.Nil(err)
	defer res.Body.Close()

	assert.Equal(200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)

	var resj []jsonrpc.Response
	err = json.Unmarshal(body, &resj)
	assert.Nil(err)

	for _, res := range resj {
		assert.Equal("2.0", res.Version)
		assert.Nil(res.Error)
		switch res.ID {
		case "1":
			assert.Equal("/user/get?user_id=1", res.Result)
		case "2":
			assert.Equal("/item/get?item_id=2", res.Result)
		case "3":
			assert.Equal("/item/update", res.Result)
		}
	}
}
