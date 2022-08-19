package tgstat_go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/helios-ag/tgstat-go/endpoints"
	server "github.com/helios-ag/tgstat-go/testing"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test getting error response", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()

		newServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(ErrorResult{
				Status: "error",
				Error:  "empty_token",
			})
		})
		client, _ := newClient(newServer.URL)
		Token = "asd"
		ctx := context.Background()
		_, err := client.NewRestRequest(ctx, Token, http.MethodGet, endpoints.ChannelsGet, make(map[string]string))
		Expect(err).ShouldNot(HaveOccurred())
	})
}

func TestEmptyData(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test getting empty data response", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()

		newServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(ErrorResult{
				Status: "error",
				Error:  "empty_token",
			})
		})
		client, _ := newClient(newServer.URL)
		Token = "asd"
		ctx := context.Background()
		_, err := client.NewRestRequest(ctx, Token, http.MethodGet, endpoints.ChannelsGet, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("data is not initialised"))
	})
}

func TestWithEndpoint(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test With endpoint", func(t *testing.T) {
		val := "https://google.com"
		WithEndpoint(val)

		Expect(TGStatClient.Url).Should(Equal(val))
	})
}

func TestReader(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test Reader not nil", func(t *testing.T) {
		reader := reader
		Expect(reader).ShouldNot(Equal(nil))
	})

	t.Run("Test Reader return", func(t *testing.T) {
		reader := reader
		ioreader := io.Reader(bytes.NewReader([]byte{}))
		Expect(reader(ioreader)).ShouldNot(Equal(nil))
	})
}

func TestClientDo(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test client do with external api", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()

		newServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(ErrorResult{
				Status: "error",
				Error:  "empty_token",
			})
		})
		// Override internal reader func with
		reader = func(r io.Reader) (bytes []byte, e error) {
			return nil, errors.New("buf overflow")
		}
		ctx := context.Background()
		client, _ := newClient(newServer.URL)

		request, _ := client.NewRestRequest(ctx, Token, http.MethodGet, endpoints.ChannelsGet, make(map[string]string))
		_, err := client.Do(request, nil)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("buf overflow"))
		// restore reader
		reader = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}
	})

	t.Run("Test response body decode", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()

		newServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(ErrorResult{
				Status: "error",
				Error:  "empty_token",
			})
		})

		ctx := context.Background()
		client, _ := newClient(newServer.URL)
		request, _ := client.NewRestRequest(ctx, Token, http.MethodGet, endpoints.ChannelsGet, make(map[string]string))
		_, err := client.Do(request, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test Bad Response Code", func(t *testing.T) {
		newServer := server.NewServer()
		defer newServer.Teardown()

		newServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})
		client, _ := newClient(newServer.URL)
		ctx := context.Background()
		request, _ := client.NewRestRequest(ctx, "Token", http.MethodGet, endpoints.ChannelsGet, make(map[string]string))
		_, err := client.Do(request, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("tgstat server responded with status code 400"))
	})
}

func TestErrorFromResponse(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Expect application/json", func(t *testing.T) {
		resp := http.Response{
			Body:   io.NopCloser(bytes.NewBufferString("Hello World")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("abc")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json", func(t *testing.T) {
		resp := http.Response{
			Body:   io.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte("{\"test\": \"test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Expect wrong json header", func(t *testing.T) {
		resp := http.Response{
			Body:   io.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application_json")
		body := []byte("{\"test\": test\"}")
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Dont expect Error", func(t *testing.T) {
		resp := http.Response{
			Body:   io.NopCloser(bytes.NewBufferString("{\"test\": test\"}")),
			Header: make(http.Header),
		}
		resp.Header.Set("Content-Type", "application/json")
		body := []byte(`{"errorCode": 0, "errorMessage": ""}`)
		err := errorFromResponse(&resp, body)
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestNewRequest(t *testing.T) {
	RegisterTestingT(t)
	t.Run("URL is empty", func(t *testing.T) {
		_, err := newClient("")

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("URL is empty"))
	})

	t.Run("Client not configured", func(t *testing.T) {
		api := &Client{
			Url:        "url",
			httpClient: &http.Client{},
		}
		api = nil
		ctx := context.Background()

		_, err := api.NewRestRequest(ctx, Token, http.MethodGet, "htts://url", nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("client not configured"))
	})

	t.Run("Trigger NewRequest errors", func(t *testing.T) {
		api := GetAPI()
		ctx := context.Background()
		// Cyrillic M
		_, err := api.NewRestRequest(ctx, Token, "М", endpoints.ChannelsGet, make(map[string]string))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("net/http: invalid method \"М\""))
		//
		_, err = api.NewRestRequest(ctx, Token, http.MethodGet, "htt\\wrongUrl", make(map[string]string))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("Trigger NewRequest data values passed", func(t *testing.T) {
		api := GetAPI()
		ctx := context.Background()
		body := make(map[string]string)
		body["value"] = "some_value"

		_, err := api.NewRestRequest(ctx, "Token", http.MethodGet, "https://url", body)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Trigger NewRestRequest errors", func(t *testing.T) {
		api := GetAPI()
		ctx := context.Background()
		// Cyrillic M
		_, err := api.NewRestRequest(ctx, Token, "М", endpoints.ChannelsGet, make(map[string]string))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid method"))

		_, err = api.NewRestRequest(ctx, Token, http.MethodGet, "htt\\wrongUrl", make(map[string]string))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid character"))
	})

	t.Run("Test improper Config url", func(t *testing.T) {
		_, err := newClient("http//google.com")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse URL"))
	})
}

func TestString(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test string point", func(t *testing.T) {
		val := String("value")
		Expect(&val).To(HaveValue(Equal("value")))
	})
}

func TestBool(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test bool converter point", func(t *testing.T) {
		val := Bool(false)
		Expect(&val).To(HaveValue(Equal(false)))
	})
}

func TestInt(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test int converter point", func(t *testing.T) {
		val := Int(100)
		Expect(&val).To(HaveValue(Equal(100)))
	})
}
func TestValidateDate(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test date is in text format triggers error", func(t *testing.T) {
		err := ValidateDate(String("blalbla"))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must be numeric"))
	})

	t.Run("Test date in numeric format (timestamp)", func(t *testing.T) {
		err := ValidateDate(String("123123123"))
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test date with nil", func(t *testing.T) {
		err := ValidateDate(nil)
		Expect(err).ToNot(HaveOccurred())
	})
}
