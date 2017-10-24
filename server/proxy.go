package server

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/istyle-inc/multimissile"
	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/jsonrpc"
	"github.com/istyle-inc/multimissile/wlog"
)

func sendHTTPRequest(wg *sync.WaitGroup, reqj jsonrpc.Request, forwardHeaders *http.Header, respj *jsonrpc.Response) {
	defer wg.Done()
	reqh, err := buildHTTPRequest(&reqj, forwardHeaders)
	if err != nil {
		*respj = buildJSONRPCErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, 0, 0)
		errorLog(wlog.Error, err.Error())
		return
	}
	start := time.Now()
	resp, err := HTTPClient.Do(reqh)
	end := time.Now()
	ptime := (end.Sub(start)).Seconds()
	if err != nil {
		*respj = buildJSONRPCErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime, resp.StatusCode)
		errorLog(wlog.Error, err.Error())
		return
	}

	defer resp.Body.Close()
	// Temporarily call here also in buildHTTPRequest, so skip error at this time
	ep, _ := config.FindEndpoint(msl.Config, reqj.Endpoint)
	if fail(ep.AcceptableHTTPStatuses, resp.StatusCode) && notSuccess(ep.ExceptableHTTPStatuses, resp.StatusCode) {
		*respj = buildHTTPError2JSONRPCErrorResponse(resp, reqj.ID, ptime)
		errorLog(wlog.Error, "%#v is failed: %s", reqj, resp.Status)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*respj = buildJSONRPCErrorResponse(jsonrpc.InternalError, err.Error(), reqj.ID, ptime, resp.StatusCode)
		errorLog(wlog.Error, err.Error())
		return
	}
	*respj = buildJSONRPCResponse(string(body), reqj.ID, ptime, resp.StatusCode)
}

func jsonRPC2HTTP(reqs *[]jsonrpc.Request, forwardHeaders *http.Header) ([]jsonrpc.Response, error) {
	wg := new(sync.WaitGroup)
	resps := make([]jsonrpc.Response, len(*reqs))
	// send requests to endpoint conccurrently
	for i, reqj := range *reqs {
		wg.Add(1)
		go sendHTTPRequest(wg, reqj, forwardHeaders, &resps[i])
	}

	wg.Wait()

	return resps, nil
}

func notSuccess(thresholds []int, value int) bool {
	for _, v := range thresholds {
		if v == value {
			return false
		}
	}
	return true
}

func fail(thresholds []int, value int) bool {
	for _, v := range thresholds {
		if v == value {
			return true
		}
	}
	return false
}
