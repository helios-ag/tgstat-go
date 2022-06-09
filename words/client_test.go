package words

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	server "github.com/helios-ag/tgstat-go/testing"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func prepareClient(URL string) {
	tgstat.Token = "token"
	tgstat.WithEndpoint(URL)
}

func TestClient_MentionsByPeriod(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test words request validation", func(t *testing.T) {
		prepareClient("http://nonexisting")
		req := MentionPeriodRequest{
			Q: "",
		}
		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Q: cannot be blank"))
	})

	t.Run("Test words request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.WordsMentionsByPeriod, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]tgstat.WordsMentionsResponseItem, 0)
			items = append(items, tgstat.WordsMentionsResponseItem{
				Period:        "2018-11-04",
				MentionsCount: 1000,
				ViewsCount:    3985,
			})

			response := tgstat.WordsMentionsResponse{
				Items: items,
			}
			json.NewEncoder(w).Encode(tgstat.WordsMentions{
				Status:   "ok",
				Response: response,
			})
		})

		req := MentionPeriodRequest{
			Q:        "",
			PeerType: tgstat.String("5"),
		}
		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Q: cannot be blank"))
	})

	t.Run("Test mention order response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.WordsMentionsByPeriod, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]tgstat.WordsMentionsResponseItem, 0)
			items = append(items, tgstat.WordsMentionsResponseItem{
				Period:        "2018-11-04",
				MentionsCount: 1000,
				ViewsCount:    3985,
			})

			response := tgstat.WordsMentionsResponse{
				Items: items,
			}
			json.NewEncoder(w).Encode(tgstat.WordsMentions{
				Status:   "ok",
				Response: response,
			})
		})

		req := MentionPeriodRequest{
			Q: "TgStat",
		}

		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test mention order validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.WordsMentionsByPeriod, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.WordsMentions{
				Status: "ok",
			})
		})

		req := MentionPeriodRequest{
			Q:        "q",
			PeerType: tgstat.String("all"),
		}
		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:         "q",
			PeerType:  tgstat.String("all"),
			StartDate: tgstat.String("2020-01-19 03:14:07"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(false),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(false),
			MinusWords:   tgstat.String("something"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(true),
			MinusWords:   tgstat.String("something"),
			Group:        tgstat.String("day"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionPeriodRequest{
			Q:              "q",
			PeerType:       tgstat.String("all"),
			EndDate:        tgstat.String("2020-01-19 03:14:07"),
			HideForwards:   tgstat.Bool(true),
			StrongSearch:   tgstat.Bool(true),
			MinusWords:     tgstat.String("something"),
			Group:          tgstat.String("week"),
			ExtendedSyntax: tgstat.Bool(true),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		Expect(err).ToNot(HaveOccurred())
	})

}

func TestClient_MentionsByChannels(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Test mention by channels order response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.WordsMentionsByChannels, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.WordsMentions{
				Status: "ok",
			})
		})

		req := MentionsByChannelRequest{
			Q: "Term",
		}

		_, _, err := MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Test mention order response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.WordsMentionsByChannels, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.WordsMentions{
				Status: "ok",
			})
		})

		req := MentionsByChannelRequest{
			Q:        "q",
			PeerType: tgstat.String("chat"),
		}
		_, _, err := MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:         "q",
			PeerType:  tgstat.String("chat"),
			StartDate: tgstat.String("2020-01-19 03:14:07"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(false),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     tgstat.String("all"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(false),
			MinusWords:   tgstat.String("something"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     tgstat.String("chat"),
			EndDate:      tgstat.String("2020-01-19 03:14:07"),
			HideForwards: tgstat.Bool(true),
			StrongSearch: tgstat.Bool(true),
			MinusWords:   tgstat.String("something"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		req = MentionsByChannelRequest{
			Q:              "q",
			PeerType:       tgstat.String("chat"),
			EndDate:        tgstat.String("2020-01-19 03:14:07"),
			HideForwards:   tgstat.Bool(true),
			StrongSearch:   tgstat.Bool(true),
			MinusWords:     tgstat.String("something"),
			ExtendedSyntax: tgstat.Bool(true),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})
}
