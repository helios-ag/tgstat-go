package channels

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

			item := schema.ForwardItem{
				ForwardID: 123,
				PostID:    0,
				PostLink:  "t.me/123123",
				PostDate:  0,
				ChannelID: 0,
			}

			items := make([]schema.ForwardItem, 0)
			items = append(items, item)

			json.NewEncoder(w).Encode(schema.ChannelForwards{
				Status: "ok",
				Response: schema.ChannelForwardsResponse{
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
