package tgstat_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test Environment
type Server struct {
	Server *httptest.Server
	Mux    *http.ServeMux
}

func (server *Server) Teardown() {
	server.Server.Close()
	server.Server = nil
	server.Mux = nil
}

func getCfg() *ClientConfig {
	cfg := ClientConfig{
		Token: "token",
	}
	return &cfg
}

func newServer() Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return Server{
		Server: server,
		Mux:    mux,
	}
}

func TestNewClient(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test trailing slashes remove", func(t *testing.T) {
		client, _ := NewClient(getCfg())
		if strings.HasSuffix(client.Config.Endpoint, "/") {
			t.Fatalf("endpoint has trailing slashes: %q", client.Config.Endpoint)
		}
	})
	t.Run("Test getting error response", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Status: "error",
				Error:  "empty_token",
			})
		})
		client, _ := NewClient(getCfg())

		ctx := context.Background()
		_, err := client.NewRestRequest(ctx, http.MethodGet, endpoints.ChannelsGet, nil)
		Expect(err).ShouldNot(HaveOccurred())
	})
}

func TestClientDo(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test client do with external api", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Status: "error",
				Error:  "empty_token",
			})
		})
		// Override internal reader func with
		reader = func(r io.Reader) (bytes []byte, e error) {
			return nil, errors.New("buf overflow")
		}
		ctx := context.Background()
		client, _ := NewClient(getCfg())

		request, _ := client.NewRestRequest(ctx, http.MethodGet, endpoints.ChannelsGet, nil)
		_, err := client.Do(request, nil)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("buf overflow"))
		// restore reader
		reader = func(r io.Reader) ([]byte, error) {
			return ioutil.ReadAll(r)
		}
	})

	t.Run("Test response body decode", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Status: "error",
				Error:  "empty_token",
			})
		})

		ctx := context.Background()
		client, _ := NewClient(getCfg())
		request, _ := client.NewRestRequest(ctx, http.MethodGet, endpoints.ChannelsGet, nil)
		_, err := client.Do(request, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test Bad Response Code", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		})
		client, _ := NewClient(getCfg())
		ctx := context.Background()
		request, _ := client.NewRestRequest(ctx, http.MethodGet, endpoints.ChannelsGet, nil)
		_, err := client.Do(request, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("tgstat server responded with status code 400"))
	})
}

func TestErrorFromResponse(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Expect application/json", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("Hello World")),
			Header: make(http.Header, 0),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("abc")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header, 0),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("{\"test\": \"test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json header", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header, 0),
		}
		resp.Header.Set("Content-Type", "application_json")
		body := []byte("{\"test\": test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect Error", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header, 0),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte(`{"errorCode": 5, "errorMessage": "Ошибка"}`)
		err := errorFromResponse(&resp, body)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Ошибка"))
	})

	t.Run("Dont expect Error", func(t *testing.T) {
		resp := http.Response{
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header, 0),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte(`{"errorCode": 0, "errorMessage": ""}`)
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestNewRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Add rest to path", func(t *testing.T) {
		client, _ := NewClient(
			&ClientConfig{
				Token: "token",
			},
		)
		ctx := context.Background()
		_, err := client.NewRequest(ctx, http.MethodGet, endpoints.ChannelsGet, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("path contains rest request, use NewRestRequest instead"))
	})

	t.Run("Invalid config test", func(t *testing.T) {
		_, err := NewClient(
			&ClientConfig{},
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to validate given config"))
	})

	t.Run("Empty config test", func(t *testing.T) {
		_, err := NewClient(nil)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("passed in config cannot be nil"))
	})

	t.Run("Trigger NewRequest errors", func(t *testing.T) {
		client, _ := NewClient(
			&ClientConfig{
				Token: "token",
			},
		)
		ctx := context.Background()
		// Cyrillic M
		_, err := client.NewRequest(ctx, "М", endpoints.ChannelsGet, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid method"))

		_, err = client.NewRequest(ctx, "GET", "htt\\wrongUrl", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("Trigger NewRestRequest errors", func(t *testing.T) {
		client, _ := NewClient(
			&ClientConfig{
				Token: "token",
			},
		)
		ctx := context.Background()
		// Cyrillic M
		_, err := client.NewRestRequest(ctx, "М", endpoints.ChannelsGet, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid method"))

		_, err = client.NewRestRequest(ctx, "GET", "htt\\wrongUrl", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("Test improper Config url", func(t *testing.T) {
		_, err := NewClient(
			&ClientConfig{
				Endpoint: "http\\:wrongUrl",
			},
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse URL"))
	})
}
