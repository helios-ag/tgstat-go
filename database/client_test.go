package database

import (
	"context"
	"encoding/json"
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

func TestClient_CountriesGet(t *testing.T) {
	RegisterTestingT(t)
	t.Run("Test host not reachable", func(t *testing.T) {
		prepareClient("http://localhost123")

		lang := "ru"
		_, _, err := CountriesGet(context.Background(), lang)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseCountries, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			russia := tgstat.Country{
				Code: "Ru",
				Name: "Russia",
			}

			us := tgstat.Country{
				Code: "US",
				Name: "United States",
			}
			json.NewEncoder(w).Encode(tgstat.CountryResult{
				Status: "ok",
				Response: []tgstat.Country{
					russia,
					us,
				},
			})
		})

		response, _, err := CountriesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(tgstat.Country{
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
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseCategories, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			tech := tgstat.Category{
				Code: "tech",
				Name: "Технологии",
			}

			news := tgstat.Category{
				Code: "news",
				Name: "Новости",
			}
			json.NewEncoder(w).Encode(tgstat.CategoryResult{
				Status: "ok",
				Response: []tgstat.Category{
					tech,
					news,
				},
			})
		})

		response, _, err := CategoriesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(tgstat.Category{
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
		Expect(err.Error()).To(ContainSubstring("dial tcp"))
	})

	t.Run("Test Countries get data", func(t *testing.T) {
		testServer := server.NewServer()
		defer testServer.Teardown()
		prepareClient(testServer.URL)

		testServer.Mux.HandleFunc(endpoints.DatabaseLanguages, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			russia := tgstat.Language{
				Code: "Ru",
				Name: "Russia",
			}

			english := tgstat.Language{
				Code: "US",
				Name: "United States",
			}

			json.NewEncoder(w).Encode(tgstat.LanguageResult{
				Status: "ok",
				Response: []tgstat.Language{
					russia,
					english,
				},
			})
		})

		response, _, err := LanguagesGet(context.Background(), "ru")
		Expect(err).ToNot(HaveOccurred())
		Expect(response).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Status": ContainSubstring("ok"),
			"Response": ContainElement(tgstat.Language{
				Code: "Ru",
				Name: "Russia",
			}),
		})))
	})
}
