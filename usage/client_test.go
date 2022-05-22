package usage

import (
	"context"
	"encoding/json"
	"fmt"
	tgstat "github.com/helios-ag/tgstat-go"
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

var NewRestRequestStub = func(c *tgstat.Client, ctx context.Context, token, method, urlPath string, data map[string]string) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func TestClient_UsageStat(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		_, _, err := Stat(context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test Usage Stat response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.UsageStat, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.StatResult{
				Status:   "ok",
				Response: nil,
			})
		})

		response, _, err := Stat(context.Background())
		Expect(err).ToNot(HaveOccurred())

		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})

	t.Run("Test request trigger error", func(t *testing.T) {
		prepareClient("http://local123")
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub

		_, _, err := Stat(context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})
}
