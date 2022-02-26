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
	cfg := tgstat.ClientConfig{
		Token: "token",
		Url:   "http://local",
	}
	tgstat.SetConfig(cfg)
	tgstat.WithEndpoint(URL)
}

func TestClient_ChannelMentionsByPeriod(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test mention order validation", func(t *testing.T) {
		prepareClient("localhost")
		req := MentionPeriodRequest{
			Q:              "",
			PeerType:       nil,
			StartDate:      nil,
			EndDate:        nil,
			HideForwards:   nil,
			StrongSearch:   nil,
			MinusWords:     nil,
			Group:          nil,
			ExtendedSyntax: nil,
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
			json.NewEncoder(w).Encode(schema.WordsMentionsResponse{
				Status: "ok",
				Response: struct {
					Items []struct {
						Period        string `json:"period"`
						MentionsCount int    `json:"mentions_count"`
						ViewsCount    int    `json:"views_count"`
					} `json:"items"`
				}{
					Items: nil,
				},
			})
		})

		req := MentionPeriodRequest{
			Q: "TgStat",
		}

		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_WordsMentionsByPeriod(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test mention order validation", func(t *testing.T) {
		prepareClient("http://nonexistinghost")

		req := MentionPeriodRequest{
			Q: "q",
		}
		_, _, err := MentionsByPeriod(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test mention order response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.WordsMentionsByChannels, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.WordsMentionsResponse{
				Status: "ok",
			})
		})

		req := MentionsByChannelRequest{
			Q: "Term",
		}

		_, _, err := MentionsByChannels(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
	})
}
