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

func TestClient_Mentions(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel subscribers request validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelSubscribersRequest{
			ChannelId: "",
			StartDate: nil,
			EndDate:   nil,
			Group:     nil,
		}
		_, _, err := channels.Subscribers(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel mentions response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsSubscribers, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			items := make([]tgstat.ChannelSubscribersResponse, 0)
			items = append(items, tgstat.ChannelSubscribersResponse{
				Period:            "2018-11-04",
				ParticipantsCount: 1156,
			})

			json.NewEncoder(w).Encode(tgstat.ChannelSubscribers{
				Status:   "ok",
				Response: items,
			})
		})
		request := channels.ChannelSubscribersRequest{
			ChannelId: "/tme/123",
			StartDate: nil,
			EndDate:   nil,
		}
		response, _, err := channels.Subscribers(context.Background(), request)

		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
