package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

func TestCheckDomain(t *testing.T) {
	type testCase struct {
		name           string
		setupServer    func() *httptest.Server
		inputDomains   func(serverURL string) []string
		expectedResult func(serverURL string) []models.Domain
	}

	tests := []testCase{
		{
			name: "Success - Status OK",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			inputDomains: func(serverURL string) []string {
				return []string{serverURL}
			},
			expectedResult: func(serverURL string) []models.Domain {
				return []models.Domain{
					{
						FullUrl:     serverURL,
						StatusCode:  http.StatusOK,
						RedirectUrl: "",
					},
				}
			},
		},
		{
			name: "Success - Status Not Found",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			inputDomains: func(serverURL string) []string {
				return []string{serverURL}
			},
			expectedResult: func(serverURL string) []models.Domain {
				return []models.Domain{
					{
						FullUrl:     serverURL,
						StatusCode:  http.StatusNotFound,
						RedirectUrl: "",
					},
				}
			},
		},
		{
			name: "Success - With Redirect",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/" {
						http.Redirect(w, r, "/target", http.StatusMovedPermanently)
						return
					}
					w.WriteHeader(http.StatusOK)
				}))
			},
			inputDomains: func(serverURL string) []string {
				return []string{serverURL}
			},
			expectedResult: func(serverURL string) []models.Domain {
				return []models.Domain{
					{
						FullUrl:     serverURL,
						StatusCode:  http.StatusOK, // 200 after following redirect
						RedirectUrl: serverURL + "/target",
					},
				}
			},
		},
		{
			name: "Failure - Invalid URL",
			setupServer: func() *httptest.Server {
				return nil
			},
			inputDomains: func(_ string) []string {
				return []string{"http://invalid-url-that-does-not-exist"}
			},
			expectedResult: func(_ string) []models.Domain {
				return []models.Domain{
					{
						FullUrl:     "http://invalid-url-that-does-not-exist",
						StatusCode:  0, // 0 because request failed
						RedirectUrl: "",
					},
				}
			},
		},
		{
			name: "Success - Multiple Domains",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/ok":
						w.WriteHeader(http.StatusOK)
					case "/404":
						w.WriteHeader(http.StatusNotFound)
					default:
						w.WriteHeader(http.StatusOK)
					}
				}))
			},
			inputDomains: func(serverURL string) []string {
				return []string{serverURL + "/ok", serverURL + "/404"}
			},
			expectedResult: func(serverURL string) []models.Domain {
				return []models.Domain{
					{
						FullUrl:     serverURL + "/ok",
						StatusCode:  http.StatusOK,
						RedirectUrl: "",
					},
					{
						FullUrl:     serverURL + "/404",
						StatusCode:  http.StatusNotFound,
						RedirectUrl: "",
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var serverURL string
			var server *httptest.Server

			if tc.setupServer != nil {
				server = tc.setupServer()
			}

			if server != nil {
				defer server.Close()
				serverURL = server.URL
			}

			svc := &domainCheckerService{
				apiClient: http.DefaultClient,
			}

			input := tc.inputDomains(serverURL)
			inputChan := make(chan string, len(input))
			for _, url := range input {
				inputChan <- url
			}
			close(inputChan)

			resultChan := svc.checkDomain(context.Background(), inputChan)

			var got []models.Domain
			for result := range resultChan {
				got = append(got, result)
			}

			want := tc.expectedResult(serverURL)

			// Sort both slices to ensure consistent comparison regardless of processing order
			sort.Slice(got, func(i, j int) bool {
				return got[i].FullUrl < got[j].FullUrl
			})
			sort.Slice(want, func(i, j int) bool {
				return want[i].FullUrl < want[j].FullUrl
			})

			if !reflect.DeepEqual(got, want) {
				t.Errorf("checkDomain() = %v, want %v", got, want)
			}
		})
	}
}
