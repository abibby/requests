package validate

import (
	_ "embed"
	"encoding/json"

	"github.com/pkg/errors"
)

type Message struct {
	Array   string `json:"array"`
	String  string `json:"string"`
	Numeric string `json:"numeric"`
}

func (m *Message) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		s := ""
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		m.Array = s
		m.String = s
		m.Numeric = s
		return nil
	} else {
		type localMessage Message
		return json.Unmarshal(b, (*localMessage)(m))
	}
}

var messages map[string]*Message

//go:embed lang.json
var lang []byte

func init() {
	err := json.Unmarshal(lang, &messages)
	if err != nil {
		panic(errors.Wrap(err, "could not parse lang.json"))
	}
}
