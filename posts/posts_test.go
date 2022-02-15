package posts

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	server "github.com/helios-ag/tgstat-go/testing"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"testing"
)

func prepareClient(URL string) {
	cfg := tgstat.ClientConfig{
		"token",
		false,
		"http://local",
	}
	tgstat.SetConfig(cfg)
	tgstat.WithEndpoint(URL)
}

func TestClient_PostsGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {

		_, _, err := PostGet(context.Background(), "t.me/123")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.PostResponse{
				Status: "ok",
				Response: struct {
					ID            int         `json:"id"`
					Date          int         `json:"date"`
					Views         int         `json:"views"`
					Link          string      `json:"link"`
					ChannelID     int         `json:"channel_id"`
					ForwardedFrom interface{} `json:"forwarded_from"`
					IsDeleted     int         `json:"is_deleted"`
					Text          string      `json:"text"`
					Media         struct {
						MediaType string `json:"media_type"`
						Caption   string `json:"caption"`
					} `json:"media"`
				}{
					ID:            0,
					Date:          0,
					Views:         0,
					Link:          "",
					ChannelID:     0,
					ForwardedFrom: nil,
					IsDeleted:     0,
					Text:          "",
					Media: struct {
						MediaType string `json:"media_type"`
						Caption   string `json:"caption"`
					}{
						MediaType: "mediaPhoto",
						Caption:   "",
					},
				},
			})
		})

		response, _, err := PostGet(context.Background(), "t.me/123")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_PostsStat(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostStatRequest{
			PostId: "",
			Group:  nil,
		}
		_, _, err := PostStat(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.PostStatResponse{
				Status: "",
				Response: struct {
					ViewsCount    int `json:"viewsCount"`
					ForwardsCount int `json:"forwardsCount"`
					MentionsCount int `json:"mentionsCount"`
					Forwards      []struct {
						PostID    string `json:"postId"`
						PostLink  string `json:"postLink"`
						PostDate  string `json:"postDate"`
						ChannelID int    `json:"channelId"`
					} `json:"forwards"`
					Mentions []struct {
						PostID    string `json:"postId,omitempty"`
						PostLink  string `json:"postLink,omitempty"`
						PostDate  string `json:"postDate,omitempty"`
						ChannelID int    `json:"channelId,omitempty"`
					} `json:"mentions"`
					Views []struct {
						Date        string `json:"date"`
						ViewsGrowth int    `json:"viewsGrowth"`
					} `json:"views"`
				}{
					ViewsCount:    0,
					ForwardsCount: 0,
					MentionsCount: 0,
					Forwards:      nil,
					Mentions:      nil,
					Views:         nil,
				},
			})
		})

		req := PostStatRequest{
			PostId: "",
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
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostSearchRequest{
			Q:              "Query",
			Limit:          nil,
			Offset:         nil,
			PeerType:       nil,
			StartDate:      nil,
			EndDate:        nil,
			HideForwards:   nil,
			HideDeleted:    nil,
			StrongSearch:   nil,
			MinusWords:     nil,
			ExtendedSyntax: nil,
		}
		_, _, err := PostSearch(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsSearch, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.PostSearchResponse{
				Status: "",
				Response: struct {
					Count      int `json:"count"`
					TotalCount int `json:"total_count"`
					Items      []struct {
						ID            int64       `json:"id"`
						Date          int         `json:"date"`
						Views         int         `json:"views"`
						Link          string      `json:"link"`
						ChannelID     int         `json:"channel_id"`
						ForwardedFrom interface{} `json:"forwarded_from"`
						IsDeleted     int         `json:"is_deleted"`
						Text          string      `json:"text"`
						Snippet       string      `json:"snippet"`
						Media         struct {
							MediaType string `json:"media_type"`
							MimeType  string `json:"mime_type"`
							Size      int    `json:"size"`
						} `json:"media"`
					} `json:"items"`
				}{
					Count:      112,
					TotalCount: 23,
					Items:      nil,
				},
			})
		})

		req := PostStatRequest{
			PostId: "t.me/123",
			Group:  nil,
		}
		response, _, err := PostStat(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_PostsSearchExtended(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := PostSearchRequest{
			Q:              "",
			Limit:          nil,
			Offset:         nil,
			PeerType:       nil,
			StartDate:      nil,
			EndDate:        nil,
			HideForwards:   nil,
			HideDeleted:    nil,
			StrongSearch:   nil,
			MinusWords:     nil,
			ExtendedSyntax: nil,
		}
		_, _, err := PostSearchExtended(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test PostsGet response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.PostsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.PostStatResponse{
				Status: "",
				Response: struct {
					ViewsCount    int `json:"viewsCount"`
					ForwardsCount int `json:"forwardsCount"`
					MentionsCount int `json:"mentionsCount"`
					Forwards      []struct {
						PostID    string `json:"postId"`
						PostLink  string `json:"postLink"`
						PostDate  string `json:"postDate"`
						ChannelID int    `json:"channelId"`
					} `json:"forwards"`
					Mentions []struct {
						PostID    string `json:"postId,omitempty"`
						PostLink  string `json:"postLink,omitempty"`
						PostDate  string `json:"postDate,omitempty"`
						ChannelID int    `json:"channelId,omitempty"`
					} `json:"mentions"`
					Views []struct {
						Date        string `json:"date"`
						ViewsGrowth int    `json:"viewsGrowth"`
					} `json:"views"`
				}{
					ViewsCount:    0,
					ForwardsCount: 0,
					MentionsCount: 0,
					Forwards:      nil,
					Mentions:      nil,
					Views:         nil,
				},
			})
		})

		req := PostStatRequest{
			PostId: "",
			Group:  nil,
		}
		response, _, err := PostStat(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}
