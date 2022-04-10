package validate

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func Run(requestHttp *http.Request, requestStruct any) error {
	var decoder = schema.NewDecoder()
	decoder.SetAliasTag("query")
	err := decoder.Decode(requestStruct, requestHttp.URL.Query())
	if err, ok := err.(schema.MultiError); ok {
		return fromSchemaMultiError(err)
	} else if err != nil {
		return errors.Wrap(err, "Could decode query string")
	}

	if requestHttp.Body != http.NoBody {
		defer requestHttp.Body.Close()

		err := json.NewDecoder(requestHttp.Body).Decode(requestStruct)
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return fromJsonUnmarshalTypeError(err, requestStruct)
		} else if err != nil {
			return errors.Wrap(err, "Could decode body")
		}
	}
	return nil
}
