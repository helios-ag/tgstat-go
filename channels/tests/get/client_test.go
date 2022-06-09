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

func TestClient_ChannelGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		channelId := ""

		_, _, err := channels.Get(context.Background(), channelId)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel Get response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.ChannelResponseResult{
				Status: "ok",
				Response: tgstat.ChannelResponse{
					Id:                321,
					Link:              "t.me/varlamov",
					Username:          "@varlamov",
					Title:             "Varlamov.ru",
					About:             "Илья Варламов. Make Russia warm again! ...",
					Image100:          "//static.tgstat.ru/public/images/channels/_100/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					Image640:          "//static.tgstat.ru/public/images/channels/_0/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					ParticipantsCount: 100,
					TGStatRestriction: tgstat.TGStatRestrictions{
						RedLabel:   true,
						BlackLabel: true,
					},
				},
			})
		})

		channelId := "test"

		response, _, err := channels.Get(context.Background(), channelId)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
			"Response": MatchFields(IgnoreExtras, Fields{
				"Title": ContainSubstring("Varlam"),
			}),
		})))
	})
}
