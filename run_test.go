package validate

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_unmarshals_data_from_query_string(t *testing.T) {
	type Request struct {
		Foo string `query:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?foo=bar", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, "bar", structRequest.Foo)
}

func TestRun_unmarshals_json_data_from_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{ "foo": "bar" }`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, "bar", structRequest.Foo)
}
