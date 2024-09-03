package goweb

import (
	"fmt"
	"net/http"
	"time"

	"github.com/trwk76/goweb/log"
)

func New(errHandler ErrorHandler) *Web {
	if errHandler == nil {
		errHandler = defaultErrorHandler
	}

	return &Web{
		mux:  http.NewServeMux(),
		errh: errHandler,
	}
}

func (w *Web) Mux() *http.ServeMux {
	return w.mux
}

type (
	Web struct {
		mux  *http.ServeMux
		errh ErrorHandler
	}

	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

	respBuffer struct {
		w  http.ResponseWriter
		s  int
		cl uint64
	}
)

func (i *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tim := time.Now()
	buf := respBuffer{w: w}
	corrID := newCorrID()

	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if !ok {
				err = fmt.Errorf("request panic: %v", p)
			}

			log.Fatal(time.Now(), corrID, err.Error())

			if buf.s == 0 {
				i.errh(&buf, r, err)
			}
		}

		if buf.s == 0 {
			DefaultResponse(&buf, r, http.StatusNotImplemented)
		}

		log.Access(tim, corrID, r, buf.s, buf.cl, time.Since(tim))
	}()

	setCorrID(r, corrID)

	i.mux.ServeHTTP(&buf, r)
}

func defaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	DefaultResponse(w, r, http.StatusInternalServerError)
}

func (b *respBuffer) Header() http.Header {
	return b.w.Header()
}

func (b *respBuffer) WriteHeader(status int) {
	b.s = status
	b.w.WriteHeader(status)
}

func (b *respBuffer) Write(p []byte) (int, error) {
	done, err := b.w.Write(p)
	b.cl += uint64(done)
	return done, err
}

var (
	_ http.Handler        = (*Web)(nil)
	_ http.ResponseWriter = (*respBuffer)(nil)
)
