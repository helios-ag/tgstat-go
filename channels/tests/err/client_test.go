package channels

import (
	"context"
	"encoding/json"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"github.com/helios-ag/tgstat-go/endpoints"
	"github.com/helios-ag/tgstat-go/schema"
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

func TestClient_ChannelErr(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		request := channels.ChannelViewsRequest{
			ChannelId: "",
		}

		_, _, err := channels.Err(context.Background(), request)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("ChannelId: cannot be blank"))
	})

	t.Run("Test channel Err response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.ChannelErr, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := schema.ChannelErrResponse{
				Period: "2021-11-26",
				Err:    41,
			}

			responses := make([]schema.ChannelErrResponse, 0)
			responses = append(responses, response)

			json.NewEncoder(w).Encode(schema.ChannelErr{
				Status:   "ok",
				Response: responses,
			})
		})

		request := channels.ChannelViewsRequest{
			ChannelId: "t.me/varlam",
		}

		response, _, err := channels.Err(context.Background(), request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": Equal("ok"),
		})))
	})
}
