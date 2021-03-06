package posts

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

func TestClient_Views(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel views request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelViewsRequest{
			ChannelId: "",
		}
		_, _, err := channels.Views(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel views response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsViews, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]tgstat.ChannelViewsResponse, 0)
			items = append(items, tgstat.ChannelViewsResponse{
				Period:     "2018-11-04",
				ViewsCount: 3985,
			})

			json.NewEncoder(w).Encode(tgstat.ChannelViews{
				Status:   "ok",
				Response: items,
			})
		})
		request := channels.ChannelViewsRequest{
			ChannelId: "/tme/123",
			StartDate: nil,
			EndDate:   nil,
		}
		response, _, err := channels.Views(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
