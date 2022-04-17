package rules

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

func initStringRules() {
	AddStringRule("alpha", func(value string, args []string) error {
		for _, c := range value {
			if !alpha(c) {
				return fmt.Errorf("")
			}
		}
		return nil
	})
	AddStringRule("alpha_dash", func(value string, args []string) error {
		for _, c := range value {
			if !alpha(c) && c != '-' {
				return fmt.Errorf("")
			}
		}
		return nil
	})
	AddStringRule("alpha_num", func(value string, args []string) error {
		for _, c := range value {
			if !alpha(c) && !numeric(c) {
				return fmt.Errorf("")
			}
		}
		return nil
	})
	AddStringRule("numeric", func(value string, args []string) error {
		for _, c := range value {
			if !numeric(c) {
				return fmt.Errorf("")
			}
		}
		return nil
	})
	AddStringRule("email", func(value string, args []string) error {
		_, err := mail.ParseAddress(value)
		return err
	})
	AddStringRule("ends_with", func(value string, args []string) error {
		if len(args) < 1 {
			log.Print("end_with must have 1 argument")
			return nil
		}
		if !strings.HasSuffix(value, args[0]) {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("starts_with", func(value string, args []string) error {
		if len(args) < 1 {
			log.Print("starts_with must have 1 argument")
			return nil
		}
		if !strings.HasPrefix(value, args[0]) {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("ip_address", func(value string, args []string) error {
		if net.ParseIP(value) == nil {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("json", func(value string, args []string) error {
		var v any
		err := json.Unmarshal([]byte(value), &v)
		if err != nil {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("mac_address", func(value string, args []string) error {
		_, err := net.ParseMAC(value)
		if err != nil {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("not_regex", func(value string, args []string) error {
		if len(args) < 1 {
			log.Print("not_regex must have 1 argument")
			return nil
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			log.Printf("not_regex arg is not valid regex: %v", err)
			return nil
		}
		if re.MatchString(value) {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("regex", func(value string, args []string) error {
		if len(args) < 1 {
			log.Print("regex must have 1 argument")
			return nil
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			log.Printf("regex arg is not valid regex: %v", err)
			return nil
		}
		if !re.MatchString(value) {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("timezone", func(value string, args []string) error {
		_, err := time.LoadLocation(value)
		if err != nil {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("url", func(value string, args []string) error {
		_, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("")
		}
		return nil
	})
	AddStringRule("uuid", func(value string, args []string) error {
		_, err := uuid.Parse(value)
		if err != nil {
			return fmt.Errorf("")
		}
		return nil
	})
}

func AddStringRule(key string, cb func(value string, args []string) error) {
	AddRule(key, func(options *ValidationOptions) error {
		value, ok := options.Value.(string)
		if !ok {
			return nil
		}
		return cb(value, options.Arguments)
	})
}

func alpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func numeric(c rune) bool {
	return c >= '0' && c <= '9'
}
