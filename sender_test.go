package sender

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	e := errors.New("new error")
	data, err := marshalError(e)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "" {
		t.Error(string(data))
	}

	e2 := jsonError{}
	data, err = marshalError(e2)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "" {
		t.Error(string(data))
	}
}

type jsonError struct {
}

func (e jsonError) Error() string {
	return "new error"
}

func (e jsonError) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func marshalError(err error) ([]byte, error) {
	switch v := err.(type) {
	case json.Marshaler:
		return json.Marshal(v)
	default:
		return json.Marshal(errorData{
			Error: err.Error(),
		})
	}
}
