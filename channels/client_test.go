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

		_, _, err := Get(context.Background(), channelId)
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
				Response: schema.ChannelWithRestriction{
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

		response, _, err := Get(context.Background(), channelId)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
			"Response": MatchFields(IgnoreExtras, Fields{
				"Title": ContainSubstring("Varlam"),
			}),
		})))
	})
}

func TestClient_ChannelSearch(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := SearchRequest{
			Q:                   "",
			SearchByDescription: 0,
			Country:             "russia",
			Language:            nil,
			Category:            "",
			Limit:               nil,
		}
		_, _, err := Search(context.Background(), request)
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
			items := make([]schema.ChannelSearchItem, 0)
			items = append(items, schema.ChannelSearchItem{
				Id:                0,
				Link:              "t.me/varlamov",
				Username:          "varlamov",
				Title:             "Varlam",
				About:             "abc",
				Image100:          "",
				Image640:          "",
				ParticipantsCount: 5,
			})

			json.NewEncoder(w).Encode(schema.ChannelSearchResponse{
				Status: "ok",
				Response: schema.ChannelSearch{
					Count: 0,
					Items: items,
				},
			})
		})
		request := SearchRequest{
			Q:                   "test",
			SearchByDescription: 0,
			Country:             "russia",
			Language:            nil,
			Category:            "",
			Limit:               nil,
		}
		response, _, err := Search(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}

func TestClient_ChannelStat(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel stat validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		channelId := ""

		_, _, err := Stat(context.Background(), channelId)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelID must be set"))
	})

	t.Run("Test channel Stat response mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.ChannelStatResponse{
				Status: "ok",
				Response: schema.ChannelStat{
					Id:                0,
					Title:             "Varlam",
					Username:          "Varlam",
					ParticipantsCount: 0,
					AvgPostReach:      0,
					ErrPercent:        0,
					DailyReach:        0,
					CiIndex:           0,
				},
			})
		})

		channelId := "test"

		response, _, err := Stat(context.Background(), channelId)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
