package tgstat

import (
	"context"
	"encoding/json"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"testing"
	"tgstat/endpoints"
	"tgstat/schema"
)

func TestClient_CountriesGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		client, _ := prepareClient()
		lang := "ru"
		_, _, err := client.CountriesGet(context.Background(), lang)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.UsageStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.StatResponse{
				Status: "ok",
				Response: []struct {
					ServiceKey    string `json:"serviceKey"`
					Title         string `json:"title"`
					SpentChannels string `json:"spentChannels,omitempty"`
					SpentRequests string `json:"spentRequests"`
					ExpiredAt     int    `json:"expiredAt"`
					SpentWords    string `json:"spentWords,omitempty"`
				}{{ServiceKey: "api_stat_l", Title: "Доступ к Stat API (тариф L)", SpentChannels: "1989/2500", SpentRequests: "89152/400000", ExpiredAt: 1542732689, SpentWords: "111/11"}},
			})
		})

		response, _, err := server.Client.UsageStat(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"ServiceKey": ContainSubstring("stat"),
			"Title":      ContainSubstring("Stat"),
		})))
	})
}
