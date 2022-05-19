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

func TestClient_ChannelSearch(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.SearchRequest{
			Q:                   "",
			SearchByDescription: 0,
			Country:             "russia",
			Language:            nil,
			Category:            "",
			Limit:               nil,
		}
		_, _, err := channels.Search(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Category: Either query or category  is required.; Q: Either query or category is required"))
	})

	t.Run("Test channel Search response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsSearch, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]tgstat.ChannelSearchItem, 0)
			items = append(items, tgstat.ChannelSearchItem{
				Id:                0,
				Link:              "t.me/varlamov",
				Username:          "varlamov",
				Title:             "Varlam",
				About:             "abc",
				Image100:          "",
				Image640:          "",
				ParticipantsCount: 5,
			})

			json.NewEncoder(w).Encode(tgstat.ChannelSearchResult{
				Status: "ok",
				Response: tgstat.ChannelSearch{
					Count: 0,
					Items: items,
				},
			})
		})
		request := channels.SearchRequest{
			Q:                   "test",
			SearchByDescription: 0,
			Country:             "russia",
			Language:            nil,
			Category:            "",
			Limit:               nil,
		}
		response, _, err := channels.Search(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
