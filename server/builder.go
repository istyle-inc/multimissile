package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/istyle-inc/multimissile"
	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/jsonrpc"
)

func buildRequestURI(ep, method, qs string) string {
	if strings.HasPrefix(ep, "http://") {
		return fmt.Sprintf("%s%s%s", ep, method, qs)
	}
	if strings.HasPrefix(ep, "https://") {
		return fmt.Sprintf("%s%s%s", ep, method, qs)
	}
	return fmt.Sprintf("http://%s%s%s", ep, method, qs)
}

func buildURLEncodedString(params jsonrpc.RequestParams, method string) (string, error) {
	values := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			values.Set(k, v.(string))
		case json.Number:
			values.Set(k, fmt.Sprintf("%s", v))
		default:
			return "", fmt.Errorf("msl supports only string and number")
		}
	}

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		return values.Encode(), nil
	}

	return fmt.Sprintf("?%s", values.Encode()), nil
}

func buildJSONRPCResponse(body, id string, time float64, status int) jsonrpc.Response {
	return jsonrpc.Response{
		Version: jsonrpc.Version,
		Result:  body,
		Status:  status,
		ID:      id,
		Time:    time,
	}
}

func buildJSONRPCErrorResponse(code int, msg, id string, time float64, status int) jsonrpc.Response {
	jsonRPCError := &jsonrpc.Error{
		Code:    code,
		Status:  status,
		Message: msg,
	}

	return jsonrpc.Response{
		Version: jsonrpc.Version,
		Error:   jsonRPCError,
		Status:  status,
		ID:      id,
		Time:    time,
	}
}

func buildHTTPError2JSONRPCErrorResponse(resp *http.Response, id string, time float64) jsonrpc.Response {
	switch resp.StatusCode {
	case http.StatusNotFound:
		return buildJSONRPCErrorResponse(jsonrpc.MethodNotFoundError, resp.Status, id, time, resp.StatusCode)
	}
	return buildJSONRPCErrorResponse(jsonrpc.InternalError, resp.Status, id, time, resp.StatusCode)
}

func buildHTTPRequest(reqj *jsonrpc.Request, forwardHeaders *http.Header) (*http.Request, error) {
	var reqh *http.Request

	ep, err := config.FindEndpoint(msl.Config, reqj.Endpoint)
	if err != nil {
		return reqh, err
	}

	es, err := buildURLEncodedString(reqj.Params, reqj.HTTPMethod)
	if err != nil {
		return reqh, err
	}

	switch reqj.HTTPMethod {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		uri := buildRequestURI(ep.URL, reqj.Path, "")
		reqh, err = http.NewRequest(reqj.HTTPMethod, uri, strings.NewReader(es))
		if err != nil {
			return reqh, err
		}
		reqh.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		uri := buildRequestURI(ep.URL, reqj.Path, es)
		reqh, err = http.NewRequest(reqj.HTTPMethod, uri, nil)
		if err != nil {
			return reqh, err
		}
	}

	ua := forwardHeaders.Get("User-Agent")
	if ua == "" {
		reqh.Header.Set("User-Agent", msl.ServerHeader())
	} else {
		reqh.Header.Set("User-Agent", ua)
	}

	xForwardedFor := forwardHeaders.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		reqh.Header.Set("X-Forwarded-For", xForwardedFor)
	}

	for _, headers := range ep.ProxySetHeaders {
		if len(headers) < 2 {
			continue
		}
		key := headers[0]
		value := strings.Join(headers[1:], ",")
		if key == "Host" {
			reqh.Host = value
		} else {
			reqh.Header.Set(key, value)
		}
	}

	for _, passHeaders := range ep.ProxyPassHeaders {
		length := len(passHeaders)
		if length < 2 {
			continue
		}
		key := passHeaders[0]
		realIndex := 0
		passedValues := make([]string, length)
		for _, headerKey := range passHeaders[1:] {
			headerValue := forwardHeaders.Get(headerKey)
			if len(headerValue) > 0 {
				passedValues[realIndex] = headerValue
				realIndex++
			}
		}
		value := strings.Join(passedValues[:realIndex], ",")
		reqh.Header.Set(key, value)
	}

	return reqh, nil
}
