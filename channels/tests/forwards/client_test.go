package channels

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
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

func TestClient_ChannelForwards(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelForwardRequest{
			ChannelId: "",
		}

		_, _, err := channels.Forwards(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel Forward response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsForwards, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			item := tgstat.ForwardItem{
				ForwardID: 123,
				PostID:    0,
				PostLink:  "t.me/123123",
				PostDate:  0,
				ChannelID: 0,
			}

			items := make([]tgstat.ForwardItem, 0)
			items = append(items, item)

			json.NewEncoder(w).Encode(tgstat.ChannelForwards{
				Status: "ok",
				Response: tgstat.ChannelForwardsResponse{
					Items: items,
				},
			})
		})

		request := channels.ChannelForwardRequest{
			ChannelId: "t.me/varlamov",
		}
		response, _, err := channels.Forwards(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}

func TestClient_ChannelForwardsExtended(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelForwardRequest{
			ChannelId: "",
		}

		_, _, err := channels.ForwardsExtended(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel Forward response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsForwards, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			item := tgstat.ForwardItem{
				ForwardID: 123,
				PostID:    0,
				PostLink:  "t.me/123123",
				PostDate:  0,
				ChannelID: 0,
			}

			items := make([]tgstat.ForwardItem, 0)
			items = append(items, item)

			channel := tgstat.Channel{
				ID:                7377,
				Link:              "t.me/breakingmash",
				Username:          "@breakingmash",
				Title:             "Mash",
				About:             "Помахаться и обсудить новости - @mash_chat ...",
				Image100:          "//static2.tgstat.com/public/images/channels/_100/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				Image640:          "//static2.tgstat.com/public/images/channels/_0/a7/a76c0abe2b7b1b79e70f0073f43c3b44.jpg",
				ParticipantsCount: 0,
			}
			channelItems := make([]tgstat.Channel, 0)
			channelItems = append(channelItems, channel)

			json.NewEncoder(w).Encode(tgstat.ChannelForwardsExtended{
				Status: "ok",
				Response: tgstat.ChannelForwardsResponseExtended{
					Items:    items,
					Channels: channelItems,
				},
			})
		})

		request := channels.ChannelForwardRequest{
			ChannelId: "t.me/varlamov",
		}
		response, _, err := channels.ForwardsExtended(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
