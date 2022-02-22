package database

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
		URL:   "http://local",
	}
	tgstat.SetConfig(cfg)
	tgstat.WithEndpoint(URL)
}

func TestClient_CountriesGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		lang := "ru"
		_, _, err := CountriesGet(context.Background(), lang)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseCountries, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			russia := schema.Country{
				Code: "Ru",
				Name: "Russia",
			}

			us := schema.Country{
				Code: "US",
				Name: "United States",
			}
			json.NewEncoder(w).Encode(schema.CountryResponse{
				Status: "ok",
				Response: []schema.Country{
					russia,
					us,
				},
			})
		})

		response, _, err := CountriesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(schema.Country{
				Code: "Ru",
				Name: "Russia",
			}),
		})))
	})
}

func TestClient_CategoriesGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		lang := "ru"
		_, _, err := CategoriesGet(context.Background(), lang)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseCategories, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			tech := schema.Category{
				Code: "tech",
				Name: "Технологии",
			}

			news := schema.Category{
				Code: "news",
				Name: "Новости",
			}
			json.NewEncoder(w).Encode(schema.CategoryResponse{
				Status: "ok",
				Response: []schema.Category{
					tech,
					news,
				},
			})
		})

		response, _, err := CategoriesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(schema.Category{
				Code: "tech",
				Name: "Технологии",
			}),
		})))
	})
}

func TestClient_LanguagesGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost23")

		lang := "ru"
		_, _, err := LanguagesGet(context.Background(), lang)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no such host"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseLanguages, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			russia := schema.Language{
				Code: "Ru",
				Name: "Russia",
			}

			english := schema.Language{
				Code: "US",
				Name: "United States",
			}

			json.NewEncoder(w).Encode(schema.LanguageResponse{
				Status: "ok",
				Response: []schema.Language{
					russia,
					english,
				},
			})
		})

		response, _, err := LanguagesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(schema.Language{
				Code: "Ru",
				Name: "Russia",
			}),
		})))
	})
}
