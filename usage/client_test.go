package usage

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
		Url:   "http://local",
	}
	tgstat.SetConfig(cfg)
	tgstat.WithEndpoint(URL)
}

func TestClient_UsageStat(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		_, _, err := Stat(context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test Usage Stat response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.UsageStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.StatResponse{
				Status: "ok",
				Response: []struct {
					ServiceKey    string `json:"serviceKey"`
					Title         string `json:"title"`
					SpentChannels string `json:"spentChannels,omitempty"`
					SpentRequests string `json:"spentRequests"`
					ExpiredAt     int64  `json:"expiredAt"`
					SpentWords    string `json:"spentWords,omitempty"`
				}{{
					ServiceKey:    "api_stat_l",
					Title:         "Доступ к Stat API (тариф L)",
					SpentChannels: "1989/2500",
					SpentRequests: "89152/400000",
					ExpiredAt:     1542732689,
					SpentWords:    "111/11",
				}},
			})
		})

		response, _, err := Stat(context.Background())
		Expect(err).ToNot(HaveOccurred())

		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}
