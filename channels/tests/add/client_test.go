package channels

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

func TestClient_ChannelAdd(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelAddRequest{
			ChannelName: "",
			Country:     nil,
			Language:    nil,
			Category:    nil,
		}

		_, _, err := channels.Add(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelName: cannot be blank"))
	})

	t.Run("Test channel add response pending", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsForwards, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.ChannelAddPending{
				Status: "pending",
			})
		})

		request := channels.ChannelForwardRequest{
			ChannelId: "t.me/varlamov",
			Limit:     nil,
			Offset:    nil,
			StartDate: nil,
			EndDate:   nil,
		}
		response, _, err := channels.Forwards(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("pending"),
		})))
	})

	t.Run("Test channel add response success", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelsAdd, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.ChannelAddSuccess{
				Status: "ok",
				Response: struct {
					ChannelId int `json:"channelId"`
				}{
					ChannelId: 1234,
				},
			})
		})

		request := channels.ChannelAddRequest{
			ChannelName: "t.me/@varlamov",
			Country:     nil,
			Language:    nil,
			Category:    nil,
		}
		response, _, err := channels.Add(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
