package test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sloppyio/cli/pkg/api"
)

type Helper struct {
	t *testing.T
}

func NewHelper(t *testing.T) *Helper {
	return &Helper{
		t: t,
	}
}

func (h *Helper) LoadFile(f string) []byte {
	cwd, err := os.Getwd()
	if err != nil {
		h.t.Fatal(err)
	}
	b, err := ioutil.ReadFile(filepath.Join(cwd, "testdata", f))
	if err != nil {
		h.t.Fatal(err)
	}
	return b
}

func (h *Helper) NewAPIServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func (h *Helper) NewClient(host net.Addr) *api.Client {
	client := api.NewClient()
	err := client.SetBaseURL(fmt.Sprintf("http://%s/", host.String()))
	if err != nil {
		h.t.Fatal(err)
	}
	return client
}

// Creates a http handler which returns a response with the given content
// on the given path path. Otherwise return 404 - not found.
func (h *Helper) NewHTTPTestHandler(content []byte, path string) http.HandlerFunc {
	c := content
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, path) {
			_, err := w.Write(c)
			if err != nil {
				h.t.Error(err)
				return
			}
		} else {
			http.NotFound(w, r)
		}
	}
}

// Creates a http handler which explicit returns a 404 json error response for the given path with given content.
func (h *Helper) NewHTTPNotFoundHandler(content []byte, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, path) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, content)
			if err != nil {
				h.t.Error(err)
			}
		} else {
			http.NotFound(w, r)
		}
	}
}
