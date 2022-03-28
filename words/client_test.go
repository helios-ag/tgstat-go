package words

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
	server "github.com/helios-ag/tgstat-go/testing"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func prepareClient(URL string) {
	tgstat.Token = "token"
	tgstat.WithEndpoint(URL)
}

func makeStrP(s string) *string {
	return &s
}

func makeBoolP(b bool) *bool {
	return &b
}

func TestClient_MentionsByPeriod(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test words request validation", func(t *testing.T) {
		prepareClient("localhost")
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
			items := make([]schema.WordsMentionsItem, 0)
			items = append(items, schema.WordsMentionsItem{
				Period:        "2018-11-04",
				MentionsCount: 1000,
				ViewsCount:    3985,
			})

			response := schema.WordsMentionsResponse{
				Items: items,
			}
			json.NewEncoder(w).Encode(schema.WordsMentions{
				Status:   "ok",
				Response: response,
			})
		})

		req := MentionPeriodRequest{
			Q:        "",
			PeerType: makeStrP("5"),
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
			items := make([]schema.WordsMentionsItem, 0)
			items = append(items, schema.WordsMentionsItem{
				Period:        "2018-11-04",
				MentionsCount: 1000,
				ViewsCount:    3985,
			})

			response := schema.WordsMentionsResponse{
				Items: items,
			}
			json.NewEncoder(w).Encode(schema.WordsMentions{
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
		prepareClient("http://nonexistinghost")

		req := MentionPeriodRequest{
			Q:        "q",
			PeerType: makeStrP("p"),
		}
		_, _, err := MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:         "q",
			PeerType:  makeStrP("p"),
			StartDate: makeStrP("2020-01-19 03:14:07"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     makeStrP("p"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     makeStrP("p"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(false),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     makeStrP("p"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(false),
			MinusWords:   makeStrP("something"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:            "q",
			PeerType:     makeStrP("p"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(true),
			MinusWords:   makeStrP("something"),
			Group:        makeStrP("hope"),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		req = MentionPeriodRequest{
			Q:              "q",
			PeerType:       makeStrP("p"),
			EndDate:        makeStrP("2020-01-19 03:14:07"),
			HideForwards:   makeBoolP(true),
			StrongSearch:   makeBoolP(true),
			MinusWords:     makeStrP("something"),
			Group:          makeStrP("hope"),
			ExtendedSyntax: makeBoolP(true),
		}
		_, _, err = MentionsByPeriod(context.Background(), req)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
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
			json.NewEncoder(w).Encode(schema.WordsMentions{
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
		prepareClient("http://nonexisting")

		req := MentionsByChannelRequest{
			Q:        "q",
			PeerType: makeStrP("p"),
		}
		_, _, err := MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:         "q",
			PeerType:  makeStrP("chat"),
			StartDate: makeStrP("2020-01-19 03:14:07"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     makeStrP("all"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
		}
		_, _, err = MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     makeStrP("all"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(false),
		}
		_, _, err = MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     makeStrP("all"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(false),
			MinusWords:   makeStrP("something"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:            "q",
			PeerType:     makeStrP("chat"),
			EndDate:      makeStrP("2020-01-19 03:14:07"),
			HideForwards: makeBoolP(true),
			StrongSearch: makeBoolP(true),
			MinusWords:   makeStrP("something"),
		}
		_, _, err = MentionsByChannels(context.Background(), req)

		req = MentionsByChannelRequest{
			Q:              "q",
			PeerType:       makeStrP("chat"),
			EndDate:        makeStrP("2020-01-19 03:14:07"),
			HideForwards:   makeBoolP(true),
			StrongSearch:   makeBoolP(true),
			MinusWords:     makeStrP("something"),
			ExtendedSyntax: makeBoolP(true),
		}
		_, _, err = MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})
}
