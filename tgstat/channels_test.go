package tgstat

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"testing"
	"tgstat/endpoints"
	"tgstat/schema"
)

func prepareClient() (*Client, error) {
	cfg := ClientConfig{
		token: "test",
	}

	client, err := NewClient(&cfg, WithEndpoint("http://api-tgstat///"))

	return client, err
}

var NewRestRequestStub = func(
	c *Client,
	ctx context.Context,
	method,
	urlPath string,
	data map[string]string,
	jsonParams map[string]string) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

var NewRequestStub = func(
	c *Client,
	ctx context.Context,
	method,
	urlPath string,
	data interface{},
) (*http.Request, error) {
	return nil, fmt.Errorf("error happened")
}

func TestClient_ChannelGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test channel validation", func(t *testing.T) {
		client, _ := prepareClient()

		channelId := ""

		_, _, err := client.ChannelGet(context.Background(), channelId)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("return url must be set"))
	})

	t.Run("Test channel Get response Mapping", func(t *testing.T) {
		server := newServer()
		defer server.Teardown()

		server.Mux.HandleFunc(endpoints.ChannelsGet, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(schema.ChannelResponse{
				Status: "ok",
				Response: schema.Channel {
					Id: 321,
					Link: "t.me/varlamov",
					Username: "@varlamov",
					Title: "Varlamov.ru",
					About: "Илья Варламов. Make Russia warm again! ...",
					Image100: "//static.tgstat.ru/public/images/channels/_100/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					Image640: "//static.tgstat.ru/public/images/channels/_0/ca/caf1a3dfb505ffed0d024130f58c5cfa.jpg",
					ParticipantsCount: 100,
					TGStatRestriction: schema.TGStatRestriction {
						RedLabel: true,
						BlackLabel: true,
					},
				},
			})
		})

		channelId := "test"

		response, _, err := server.Client.ChannelGet(context.Background(), channelId)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainSubstring("varlam"),
		})))

	})
}