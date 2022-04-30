package validate

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type MessageOptions struct {
	Attribute string
	Value     any
	Arguments []string
}

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

func getMessage(ruleName string, options *MessageOptions) string {
	defaultMessage := func() string {
		if len(options.Arguments) == 0 {
			return ruleName
		}
		return ruleName + " " + strings.Join(options.Arguments, ", ")
	}
	message, ok := messages[ruleName]
	if !ok {
		return defaultMessage()
	}
	t, err := template.New(ruleName).Parse(message.String)
	if err != nil {
		log.Print(err)
		return defaultMessage()
	}
	buff := &bytes.Buffer{}
	err = t.Execute(buff, options)
	if err != nil {
		log.Print(err)
		return defaultMessage()
	}
	return buff.String()
}
