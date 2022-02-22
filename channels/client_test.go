package channels

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
		Token: "token",
		URL:   "http://local",
	}
	tgstat.SetConfig(cfg)
	tgstat.WithEndpoint(URL)
}

func TestClient_ChannelGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		prepareClient("localhost")

		channelId := ""

		_, _, err := ChannelGet(context.Background(), channelId)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelID must be set"))
	})

	t.Run("Test channel Get response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.ChannelResponse{
				Status: "ok",
				Response: schema.Channel{
					Id:                321,
					Link:              "t.me/varlamov",
					Username:          "@varlamov",
					Title:             "Varlamov.ru",
					About:             "Илья Варламов. Make Russia warm again! ...",
					Image100:          "//static.tgstat.ru/public/images/channels/_100/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					Image640:          "//static.tgstat.ru/public/images/channels/_0/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					ParticipantsCount: 100,
					TGStatRestriction: schema.TGStatRestriction{
						RedLabel:   true,
						BlackLabel: true,
					},
				},
			})
		})

		channelId := "test"

		response, _, err := ChannelGet(context.Background(), channelId)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
			"Response": MatchFields(IgnoreExtras, Fields{
				"Title": ContainSubstring("Varlam"),
			}),
		})))
	})
}
