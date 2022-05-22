package callback

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

func TestClient_SetCallback(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://local123")

		_, _, err := SetCallback(context.Background(), "t.me/123")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse URL"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		prepareClient("http://local123")
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub
		_, _, err := SetCallback(context.Background(), "https://t.me/123")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test validate callback url", func(t *testing.T) {
		err := validateCallbackUrl("")
		Expect(err.Error()).To(Equal("CallbackUrl must be set"))
	})

	t.Run("Test validate callback", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)
		testServer.Mux.HandleFunc(endpoints.SetCallbackURL, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.SetCallbackVerificationResult{
				Status:     "error",
				Error:      "wrong verify code",
				VerifyCode: "TGSTAT_VERIFY_CODE_123456",
			})
		})
		_, _, err := SetCallback(context.Background(), "hi/there?")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to parse URL"))
	})

	t.Run("Test SetCallback response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.SetCallbackURL, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.SetCallbackVerificationResult{
				Status:     "error",
				Error:      "wrong verify code",
				VerifyCode: "TGSTAT_VERIFY_CODE_123456",
			})
		})

		response, _, err := SetCallback(context.Background(), "https://myserver.me/callback")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("error"),
		})))
	})
}

func TestClient_GetCallbackInfo(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient("http://local123.tld")

		_, _, err := GetCallbackInfo(context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub
		_, _, err := GetCallbackInfo(context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test GetCallbackInfo response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.GetCallbackURL, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.GetCallbackResponse{
				Status: "ok",
				Response: tgstat.CallbackResponse{
					Url:                "https://test.ru/callback.php",
					PendingUpdateCount: 2,
					LastErrorDate:      1571562358,
					LastErrorMessage:   "Timeout expired",
				},
			})
		})

		response, _, err := GetCallbackInfo(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_SubscribeChannel(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		req := SubscribeChannelRequest{
			SubscriptionId: String("blabla"),
			ChannelId:      "t.me/username",
			EventTypes:     "new_post",
		}
		_, _, err := SubscribeChannel(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub

		req := SubscribeChannelRequest{
			SubscriptionId: String("blabla"),
			ChannelId:      "t.me/username",
			EventTypes:     "new_post",
		}

		_, _, err := SubscribeChannel(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test validation", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		req := SubscribeChannelRequest{
			SubscriptionId: String("blabla"),
			ChannelId:      "t.me/username",
			EventTypes:     "new_post1",
		}

		_, _, err := SubscribeChannel(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("EventTypes: must be a valid value"))
	})

	t.Run("Test SubscribeChannel response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.SubscribeChannel, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.SubscribeResponse{
				SubscriptionId: 123,
			})
		})

		req := SubscribeChannelRequest{
			SubscriptionId: String("blabla"),
			ChannelId:      "t.me/username",
			EventTypes:     "new_post",
		}

		response, _, err := SubscribeChannel(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"SubscriptionId": Equal(123),
		})))
	})
}

func TestClient_SubscribeWord(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient("http://localhost123.tld")

		req := SubscribeWordRequest{
			Q:          "Test",
			EventTypes: "new_post",
		}
		_, _, err := SubscribeWord(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub

		req := SubscribeWordRequest{
			Q:          "Test",
			EventTypes: "new_post",
		}
		_, _, err := SubscribeWord(context.Background(), req)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test SubscribeWord response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.SubscribeWord, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.Subscribe{
				Status: "ok",
				Response: tgstat.SubscribeResponse{
					SubscriptionId: 123,
				},
			})
		})

		req := SubscribeWordRequest{
			Q:          "Test",
			EventTypes: "new_post",
		}
		response, _, err := SubscribeWord(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_SubscriptionList(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient("http://localhost123.tld")

		req := SubscriptionsListRequest{}
		_, _, err := SubscriptionsList(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub

		req := SubscriptionsListRequest{}
		_, _, err := SubscriptionsList(context.Background(), req)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test SubscriptionsList response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.SubscriptionsList, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.SubscriptionList{
				Status: "ok",
				Response: tgstat.SubscriptionListResponse{
					TotalCount:    0,
					Subscriptions: nil,
				},
			})
		})

		req := SubscriptionsListRequest{}
		response, _, err := SubscriptionsList(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func TestClient_Unsubscribe(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient("http://localhost123.tld")

		_, _, err := Unsubscribe(context.Background(), "123")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test rest request failed", func(t *testing.T) {
		oldNewRequest := tgstat.NewRestRequest
		tgstat.NewRestRequest = NewRestRequestStub

		_, _, err := Unsubscribe(context.Background(), "123")

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("error happened"))
		tgstat.NewRestRequest = oldNewRequest
	})

	t.Run("Test Unsubscribe response Mapping", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.Unsubscribe, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tgstat.SuccessResult{
				Status: "ok",
			})
		})

		response, _, err := Unsubscribe(context.Background(), "123")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
		})))
	})
}

func String(v string) *string {
	return &v
}
