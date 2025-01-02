package myservice

import (
	"context"
	"io"
	"log"
	"net/http"

	"goeasy.dev/errors"
)

type Service interface {
	Get(ctx context.Context, key string) (dest string, err error)
	Put(ctx context.Context, key, value string) error
}

type Controller struct {
	service Service
}

func NewServiceController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// is this a comment group?
func (c *Controller) POST(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	key := r.URL.Path
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = c.service.Put(r.Context(), key, string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) GET(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	dest, err := c.service.Get(r.Context(), key)
	if err != nil {
		if err == errors.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dest))
}
