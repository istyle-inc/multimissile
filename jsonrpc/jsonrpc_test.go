package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequest(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	assert.Nil(validateRequest(&Request{
		Version: "2.0",
		Path:    "/resource/get",
		Params:  params,
		ID:      "1",
	}))
}

func TestValidateRequestVersion(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	assert.NotNil(validateRequest(&Request{
		Version: "1.0",
		Path:    "/resource/get",
		Params:  params,
		ID:      "1",
	}))
}

func TestValidateRequestMethod(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	assert.NotNil(validateRequest(&Request{
		Version: "2.0",
		Path:    "",
		Params:  params,
		ID:      "1",
	}))

	assert.Nil(validateRequest(&Request{
		Version: "2.0",
		Path:    "/reousrce/get",
		Params:  params,
		ID:      "1",
	}))
}

func TestValidateRequestHTTPMethod(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	assert.Nil(validateRequest(&Request{
		Version:    "2.0",
		Path:       "/resource/get",
		HTTPMethod: "DELETE",
		Params:     params,
		ID:         "1",
	}))

	assert.Nil(validateRequest(&Request{
		Version:    "2.0",
		Path:       "/resource/get",
		HTTPMethod: "",
		Params:     params,
		ID:         "1",
	}))

	assert.Nil(validateRequest(&Request{
		Version:    "2.0",
		Path:       "/resource/get",
		HTTPMethod: "GET",
		Params:     params,
		ID:         "1",
	}))

	assert.Nil(validateRequest(&Request{
		Version:    "2.0",
		Path:       "/resource/update",
		HTTPMethod: "POST",
		Params:     params,
		ID:         "1",
	}))
}

func TestValidateRequestID(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	assert.NotNil(validateRequest(&Request{
		Version: "2.0",
		Path:    "/resource/get",
		Params:  params,
		ID:      "",
	}))

	assert.Nil(validateRequest(&Request{
		Version: "2.0",
		Path:    "/reousrce/get",
		Params:  params,
		ID:      "1",
	}))
}

func TestValidateRequests(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	reqs := make([]Request, 0)
	reqs = append(reqs, Request{
		Version: "2.0",
		Path:    "/resource/get1",
		Params:  params,
		ID:      "1",
	})
	reqs = append(reqs, Request{
		Version: "2.0",
		Path:    "/resource/get2",
		Params:  params,
		ID:      "2",
	})

	assert.Nil(ValidateRequests(&reqs))
}

func TestValidateRequestsIDDup(t *testing.T) {
	assert := assert.New(t)

	params := make(map[string]interface{})
	params["id"] = 1

	reqs := make([]Request, 0)
	reqs = append(reqs, Request{
		Version: "2.0",
		Path:    "/resource/get1",
		Params:  params,
		ID:      "1",
	})
	reqs = append(reqs, Request{
		Version: "2.0",
		Path:    "/resource/get2",
		Params:  params,
		ID:      "1",
	})

	assert.NotNil(ValidateRequests(&reqs))
}
