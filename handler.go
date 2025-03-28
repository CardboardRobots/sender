package sender

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

func Handler(handler func(context.Context, Sender)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := NewSender(w, r)
		defer func() {
			if result := recover(); result != nil {
				stack := debug.Stack()
				s.SendStatus(http.StatusInternalServerError)
				_ = s.SendJson(panicData{
					Panic: fmt.Sprint(result),
					Stack: string(stack),
				})
			}
		}()
		handler(r.Context(), s)
	}
}

func BasicHandler(errorMap func(error) int, handler func(context.Context, Sender) error) func(w http.ResponseWriter, r *http.Request) {
	return Handler(func(ctx context.Context, s Sender) {
		err := handler(ctx, s)
		if err != nil {
			s.SendError(errorMap(err), err)
		}
	})
}

func RedirectHandler(errorMap func(error) int, handler func(context.Context, Sender) (string, int, error)) func(w http.ResponseWriter, r *http.Request) {
	return Handler(func(ctx context.Context, s Sender) {
		url, code, err := handler(ctx, s)
		if err != nil {
			s.SendError(errorMap(err), err)
		}

		s.Redirect(url, code)
	})
}

func JsonHandler[T any](errorMap func(error) int, handler func(context.Context, Sender) (T, error)) func(w http.ResponseWriter, r *http.Request) {
	return Handler(func(ctx context.Context, s Sender) {
		result, err := handler(ctx, s)
		if err == nil {
			err1 := s.SendJson(result)
			if err1 != nil {
				err = errors.Join(err, err1)
			} else {
				return
			}
		}

		s.SendError(errorMap(err), err)
	})
}

func StreamHandler(errorMap func(error) int, handler func(context.Context, Sender) (io.Reader, error)) func(w http.ResponseWriter, r *http.Request) {
	return Handler(func(ctx context.Context, s Sender) {
		result, err := handler(ctx, s)
		if err == nil {
			_, err1 := s.SendReader(result)
			if err1 != nil {
				err = errors.Join(err, err1)
			} else {
				return
			}
		}

		s.SendError(errorMap(err), err)
	})
}

func TemplateHandler[T any](errorMap func(error) int, template *Template, handler func(context.Context, Sender) (T, error)) func(w http.ResponseWriter, r *http.Request) {
	return Handler(func(ctx context.Context, s Sender) {
		result, err := handler(ctx, s)
		if err == nil {
			err1 := template.ExecuteTemplate(s.response, result)
			if err1 != nil {
				err = errors.Join(err, err1)
			} else {
				return
			}
		}

		s.SendError(errorMap(err), err)
	})
}
