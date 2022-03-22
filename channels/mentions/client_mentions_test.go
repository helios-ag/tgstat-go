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
	t.Run("Test channel mentions request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelMentionsRequest{
			ChannelId: "",
			Limit:     nil,
			Offset:    nil,
		}
		_, _, err := channels.Mentions(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel mentions response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsMentions, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]schema.MentionItem, 0)
			items = append(items, schema.MentionItem{
				MentionID:   48258272,
				MentionType: "channel",
				PostID:      4375814870,
				PostLink:    "https://t.me/Heath_Ledger_media/51932",
				PostDate:    1543487975,
				ChannelID:   197080,
			})

			response := schema.ChannelMentionsResponse{
				Items: items,
			}
			json.NewEncoder(w).Encode(schema.ChannelMentions{
				Status:   "ok",
				Response: response,
			})
		})
		request := channels.ChannelMentionsRequest{
			ChannelId: "/tme/123",
			Limit:     nil,
			Offset:    nil,
			StartDate: nil,
			EndDate:   nil,
		}
		response, _, err := channels.Mentions(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})

	t.Run("Test channel mentions extended response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsMentions, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]schema.MentionItem, 0)
			items = append(items, schema.MentionItem{
				MentionID:   48258272,
				MentionType: "channel",
				PostID:      4375814870,
				PostLink:    "https://t.me/Heath_Ledger_media/51932",
				PostDate:    1543487975,
				ChannelID:   197080,
			})

			chans := make([]schema.Channel, 0)
			chans = append(chans, schema.Channel{
				ID:                7377,
				Link:              "t.me/breakingmash",
				Username:          "@breakingmash",
				Title:             "Mash",
				About:             "Помахаться и обсудить новости - @mash_chat ...",
				Image100:          "//static2.tgstat.com/public/images/channels/_100/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				Image640:          "//static2.tgstat.com/public/images/channels/_0/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				ParticipantsCount: 0,
			})

			response := schema.ChannelMentionsResponseExtended{
				Items:    items,
				Channels: chans,
			}
			json.NewEncoder(w).Encode(schema.ChannelMentionsExtended{
				Status:   "ok",
				Response: response,
			})
		})

		request := channels.ChannelMentionsRequest{
			ChannelId: "t.me/varlam",
			Limit:     nil,
			Offset:    nil,
		}
		response, _, err := channels.Mentions(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})

}
