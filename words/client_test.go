package words

import (
	"context"
	"encoding/json"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
	"tgstat"
	"tgstat/endpoints"
	"tgstat/schema"
)

func TestClient_ChannelMentionsByPeriod(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test mention order validation", func(t *testing.T) {
		client, _ := tgstat.prepareClient()
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
		_, _, err := client.ChannelMentionsByPeriod(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test mention order response Mapping", func(t *testing.T) {
		server := tgstat_go.newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.WordsMentionsByPeriod, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.WordsMentionsResponse{
				Status: "ok",
			})
		})

		_, _, err := server.Client.UsageStat(context.Background())
		Expect(err).ToNot(HaveOccurred())
	})
}

func TestClient_ChannelMentionsByChannels(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test mention order validation", func(t *testing.T) {
		client, _ := tgstat.prepareClient()
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
		_, _, err := client.ChannelMentionsByPeriod(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test mention order response Mapping", func(t *testing.T) {
		server := tgstat_go.newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.WordsMentionsByPeriod, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.WordsMentionsResponse{
				Status: "ok",
			})
		})

		_, _, err := server.Client.UsageStat(context.Background())
		Expect(err).ToNot(HaveOccurred())
	})
}
