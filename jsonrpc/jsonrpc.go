package jsonrpc

import (
	"fmt"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

const Version = "2.0"

const (
	ParseError          = -32700
	InvalidRequestError = -32600
	MethodNotFoundError = -32601
	InvalidParamsError  = -32602
	InternalError       = -32603
)

// RequestParams parameter map
type RequestParams map[string]interface{}

// Request define JSON-RPC Request
type Request struct {
	Version    string        `json:"jsonrpc" validate:"eq=2.0"`
	Path       string        `json:"path" validate:"required"`
	HTTPMethod string        `json:"http_method" validate:""`
	Params     RequestParams `json:"params,omitempty"`
	ID         string        `json:"id" validate:"required"`
	// extention
	Endpoint string `json:"endpoint"`
}

// Response define JSON-RPC Response
type Response struct {
	Version string  `json:"jsonrpc"`
	Status  int     `json:"status"`
	Result  string  `json:"result,omitempty"`
	Error   *Error  `json:"error,omitempty"`
	ID      string  `json:"id"`
	Time    float64 `json:"time,omitempty"`
}

// Error define Request-Error
type Error struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var validate = validator.New()

// ValidateRequests validate request info
func ValidateRequests(reqs *[]Request) error {
	idMap := make(map[string]bool, len(*reqs))
	for _, r := range *reqs {
		if _, ok := idMap[r.ID]; ok {
			return fmt.Errorf("ID:%s is duplicated", r.ID)
		}
		if err := validateRequest(&r); err != nil {
			return err
		}
		idMap[r.ID] = true
	}
	return nil
}

func validateRequest(r *Request) (err error) {

	if err = validate.Struct(r); err != nil {
		return err
	}
	// empty method is treated as GET.
	if !validMethod(r.HTTPMethod) {
		return fmt.Errorf("malformed HTTP method: %s", r.HTTPMethod)
	}
	return nil
}

func validMethod(m string) bool {
	switch m {
	case "", http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}
