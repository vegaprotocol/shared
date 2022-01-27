package errors

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Errors map[string][]error

func NewErrors() Errors {
	return Errors{}
}

func (e Errors) Error() string {
	if len(e) <= 0 {
		return ""
	}

	propMessages := make([]string, 0, len(e))
	for prop, errs := range e {
		errMessages := make([]string, 0, len(errs))
		for _, err := range errs {
			errMessages = append(errMessages, err.Error())
		}
		propMessageFmt := fmt.Sprintf("%v (%v)", prop, strings.Join(errMessages, ", "))
		propMessages = append(propMessages, propMessageFmt)
	}

	sort.Strings(propMessages)
	return strings.Join(propMessages, ", ")
}

func (e Errors) Empty() bool {
	return len(e) == 0
}

// AddForProperty adds an error for a given property.
func (e Errors) AddForProperty(prop string, err error) {
	errs, ok := e[prop]
	if !ok {
		errs = []error{}
	}

	e[prop] = append(errs, err)
}

// FinalAddForProperty behaves like AddForProperty, but is meant to be called in
// a "return" statement. This helper is usually used for terminal errors.
func (e Errors) FinalAddForProperty(prop string, err error) Errors {
	e.AddForProperty(prop, err)
	return e
}

// Add adds a general error that is not related to a specific property.
func (e Errors) Add(err error) {
	e.AddForProperty("*", err)
}

// FinalAdd behaves like Add, but is meant to be called in a "return" statement.
// This helper is usually used for terminal errors.
func (e Errors) FinalAdd(err error) Errors {
	e.Add(err)
	return e
}

func (e Errors) Merge(oth Errors) {
	if oth == nil {
		return
	}

	for prop, errs := range oth {
		for _, err := range errs {
			e.AddForProperty(prop, err)
		}
	}
}

func (e Errors) Get(prop string) []error {
	messages, ok := e[prop]
	if !ok {
		return nil
	}
	return messages
}

func (e Errors) ErrorOrNil() error {
	if len(e) <= 0 {
		return nil
	}
	return e
}

func (e Errors) MarshalJSON() ([]byte, error) {
	out := map[string][]string{}
	for prop, errs := range e {
		messages := []string{}
		for _, err := range errs {
			messages = append(messages, err.Error())
		}
		out[prop] = messages
	}
	return json.Marshal(out)
}
