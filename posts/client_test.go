package posts

import (
	"context"
	"encoding/json"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	server "github.com/helios-ag/tgstat-go/testing"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"testing"
)

func prepareClient(URL string) {
	tgstat.Token = "token"
	tgstat.WithEndpoint(URL)
}

var NewRestRequestStub = func(c *tgstat.Client, ctx context.Context, token, method, urlPath string, data map[string]string) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func TestClient_PostsGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://local123")
		_, _, err := Get(context.Background(), "t.me/123")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		_, _, err := Get(context.Background(), "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("postId can not be empty"))
	})

	t.Run("Test newrest request triggers error", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub
		_, _, err := Get(context.Background(), "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("postId can not be empty"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.PostResult{
				Status: "ok",
				Response: tgstat.PostResponse{
					ID:            0,
					Date:          0,
					Views:         0,
					Link:          "",
					ChannelID:     0,
					ForwardedFrom: nil,
					IsDeleted:     0,
					Text:          "",
					Media:         tgstat.Media{},
				},
			})
		})

		response, _, err := Get(context.Background(), "t.me/123")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_PostsStat(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test PostStat request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostStatRequest{
			PostId: "",
			Group:  nil,
		}
		_, _, err := PostStat(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("PostId: cannot be blank"))
	})

	t.Run("Test PostStat group isn't nil and not valid", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostStatRequest{
			PostId: "t.me/123/123",
			Group:  tgstat.String("test"),
		}
		_, _, err := PostStat(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("must be a valid value"))
	})

	t.Run("Test PostStat group isn't nil", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.PostStatResult{
				Status:   "ok",
				Response: tgstat.PostStatResponse{},
			})
		})

		req := PostStatRequest{
			PostId: "t.me/123/123",
			Group:  tgstat.String("day"),
		}
		_, _, err := PostStat(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test rest request triggers error", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub
		req := PostStatRequest{
			PostId: "321",
			Group:  nil,
		}

		_, _, err := PostStat(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test PostsStat response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.PostStatResult{
				Status:   "ok",
				Response: tgstat.PostStatResponse{},
			})
		})

		req := PostStatRequest{
			PostId: "321",
			Group:  nil,
		}
		response, _, err := PostStat(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_PostsSearch(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		req := PostSearchRequest{
			Q: "Query",
		}
		_, _, err := PostSearch(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test search validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostSearchRequest{
			Q: "",
		}
		_, _, err := PostSearch(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("cannot be blank"))
	})

	t.Run("Test PostsSearch response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsSearch, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.PostSearchResult{
				Status:   "ok",
				Response: tgstat.PostSearchResultResponse{},
			})
		})

		req := PostSearchRequest{
			Q: "test",
		}
		response, _, err := PostSearch(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_PostsSearchExtended(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		req := PostSearchRequest{
			Q: "Test",
		}
		_, _, err := PostSearchExtended(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsSearch, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.PostStatResult{
				Status:   "ok",
				Response: tgstat.PostStatResponse{},
			})
		})

		req := PostSearchRequest{
			Q: "Search",
		}
		response, _, err := PostSearchExtended(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}
