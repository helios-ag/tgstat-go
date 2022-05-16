package posts

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
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

func TestClient_ChannelPosts(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsPosts, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]schema.ChannelPostsResponseItem, 0)
			items = append(items, schema.ChannelPostsResponseItem{
				ID:            7377,
				Date:          1540123429,
				Views:         148382,
				Link:          "t.me/breakingmash",
				ChannelID:     0,
				ForwardedFrom: "",
				IsDeleted:     0,
				Text:          "",
				Media: schema.Media{
					MediaType: "",
					MimeType:  "",
					Size:      0,
				},
			})
			channel := schema.Channel{
				ID:                7377,
				Link:              "t.me/breakingmash",
				Username:          "@breakingmash",
				Title:             "Mash",
				About:             "Помахаться и обсудить новости - @mash_chat ...",
				Image100:          "//static2.tgstat.com/public/images/channels/_100/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				Image640:          "//static2.tgstat.com/public/images/channels/_0/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				ParticipantsCount: 0,
			}

			json.NewEncoder(w).Encode(schema.ChannelPosts{
				Status: "ok",
				Response: schema.ChannelPostsResponse{
					Channel:    channel,
					Count:      50,
					TotalCount: 150,
					Items:      items,
				},
			})
		})

		request := channels.PostsRequest{
			ChannelId: "",
		}
		_, _, err := channels.Posts(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel posts response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsPosts, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]schema.ChannelPostsResponseItem, 0)
			items = append(items, schema.ChannelPostsResponseItem{
				ID:            7377,
				Date:          1540123429,
				Views:         148382,
				Link:          "t.me/breakingmash",
				ChannelID:     0,
				ForwardedFrom: "",
				IsDeleted:     0,
				Text:          "",
				Media: schema.Media{
					MediaType: "",
					MimeType:  "",
					Size:      0,
				},
			})
			channel := schema.Channel{
				ID:                7377,
				Link:              "t.me/breakingmash",
				Username:          "@breakingmash",
				Title:             "Mash",
				About:             "Помахаться и обсудить новости - @mash_chat ...",
				Image100:          "//static2.tgstat.com/public/images/channels/_100/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				Image640:          "//static2.tgstat.com/public/images/channels/_0/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				ParticipantsCount: 0,
			}

			json.NewEncoder(w).Encode(schema.ChannelPosts{
				Status: "ok",
				Response: schema.ChannelPostsResponse{
					Channel:    channel,
					Count:      50,
					TotalCount: 150,
					Items:      items,
				},
			})
		})
		request := channels.PostsRequest{
			ChannelId:    "t.me/varlam",
			Limit:        nil,
			Offset:       nil,
			StartTime:    nil,
			EndTime:      nil,
			HideForwards: nil,
			HideDeleted:  nil,
		}
		response, _, err := channels.Posts(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
