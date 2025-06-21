package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/SteinerLabs/lms/libs/log"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

var ErrShutdown = errors.New("shutdown initiated")

type App struct {
	mux     *http.ServeMux
	mw      []Middleware
	log     *log.Logger
	origins []string
}

func NewApp(log *log.Logger, mw ...Middleware) *App {
	mux := http.NewServeMux()

	return &App{
		log: log,
		mux: mux,
		mw:  mw,
	}
}

func (a *App) EnableCORS(origins []string) {
	a.origins = origins

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	}
	handler = wrapMiddleware([]Middleware{a.corsHandler}, handler)

	a.Handle("OPTIONS", "", "/", handler)
}

func (a *App) corsHandler(handler Handler) Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		return handler(ctx, w, r)
	}

	return h
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {
	if a.origins != nil {
		handler = wrapMiddleware([]Middleware{a.corsHandler}, handler)
	}
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}

		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.log.Error("web-respond", "error", err, "path", path, "method", method, "group", group)
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
