package sender

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeApplicationJSON = "application/json"
	ContentTypeTextHTML        = "text/html; charset=utf-8"
)

type (
	Sender struct {
		response http.ResponseWriter
		request  *http.Request
	}

	errorData struct {
		Error string `json:"error"`
	}

	panicData struct {
		Panic string `json:"panic"`
		Stack string `json:"stack"`
	}
)

func NewSender(w http.ResponseWriter, r *http.Request) Sender {
	return Sender{
		response: w,
		request:  r,
	}
}

func (s Sender) Response() http.ResponseWriter { return s.response }
func (s Sender) Request() *http.Request        { return s.request }

func (s Sender) SendStatus(statusCode int) {
	s.response.WriteHeader(statusCode)
}

func (s Sender) SetContentTypeJSON() {
	s.response.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
}

func (s Sender) SetContentTypeHTML() {
	s.response.Header().Set(HeaderContentType, ContentTypeTextHTML)
}

func (s Sender) SendReader(r io.Reader) (int64, error) {
	return io.Copy(s.response, r)
}

func (s Sender) SendBytes(data []byte) error {
	_, err := s.response.Write(data)
	if err != nil {
		return s.SendError(http.StatusInternalServerError, err)
	}

	return err
}

func (s Sender) SendJson(value any) error {
	s.SetContentTypeJSON()

	err := json.NewEncoder(s.response).Encode(value)
	if err != nil {
		return s.SendError(http.StatusInternalServerError, err)
	}

	return err
}

func (s Sender) SendError(statusCode int, err error) error {
	s.SetContentTypeJSON()

	s.SendStatus(statusCode)
	return json.NewEncoder(s.response).Encode(errorData{
		Error: err.Error(),
	})
}
