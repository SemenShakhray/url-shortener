package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SemenShakhray/url-shortener/internal/config"
	"github.com/SemenShakhray/url-shortener/internal/handlers"
	"github.com/SemenShakhray/url-shortener/internal/handlers/mock_service"
	"github.com/SemenShakhray/url-shortener/internal/router"
	"github.com/SemenShakhray/url-shortener/internal/storage"
	"github.com/SemenShakhray/url-shortener/pkg/api"
	"github.com/SemenShakhray/url-shortener/pkg/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSaveURL(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		url          string
		expectedCode int
		respError    string
		mockError    string
		user         string
		password     string
	}{
		{
			name:         "Success",
			alias:        "alias",
			url:          "https://example.com",
			expectedCode: http.StatusOK,
			user:         "123",
			password:     "123",
		},
		{
			name:         "Empty alias",
			alias:        "",
			url:          "https://example.com",
			expectedCode: http.StatusOK,
			user:         "123",
			password:     "123",
		},
		{
			name:         "Empty URL",
			alias:        "alias",
			url:          "",
			expectedCode: http.StatusBadRequest,
			respError:    "field URL is required field",
			user:         "123",
			password:     "123",
		},
		{
			name:         "Service error",
			alias:        "alias",
			url:          "https://example.com",
			expectedCode: http.StatusInternalServerError,
			respError:    "internal error",
			mockError:    "internal error",
			user:         "123",
			password:     "123",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		Server: config.Server{
			User:     "123",
			Password: "123",
		},
	}
	log := slogdiscard.NewDiscardLogger()
	mockService := mock_service.NewMockServicer(ctrl)
	handler := handlers.NewHandler(log, mockService)
	r := router.NewRouter(handler, cfg)

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {

			if tc.respError == "" {
				mockService.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			}
			if tc.mockError != "" {
				mockService.EXPECT().
					SaveURL(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New(tc.mockError))
			}
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			req.SetBasicAuth(tc.user, tc.password)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			var resp handlers.Response

			body := w.Body.String()
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestRedirect(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		mockError  error
		url        string
		statusCode int
		respError  string
	}{
		{
			name:       "Success",
			alias:      "example",
			url:        "https://example.com",
			statusCode: http.StatusFound,
			respError:  "",
		},
		{
			name:       "Empty alias",
			alias:      "",
			statusCode: http.StatusBadRequest,
			respError:  "alias cannot be empty",
		},
		{
			name:       "URL not found",
			alias:      "nonexistent",
			mockError:  storage.ErrURLNotFound,
			statusCode: http.StatusBadRequest,
			respError:  "url not found",
		},
		{
			name:       "Internal server error",
			alias:      "error",
			mockError:  errors.New("internal error"),
			statusCode: http.StatusInternalServerError,
			respError:  "internal error",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		Server: config.Server{
			User:     "123",
			Password: "123",
		},
	}
	log := slogdiscard.NewDiscardLogger()
	mockService := mock_service.NewMockServicer(ctrl)
	handler := handlers.NewHandler(log, mockService)
	r := router.NewRouter(handler, cfg)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.alias == "" {
				mockService.EXPECT().GetURL(gomock.Any(), tc.alias).Times(0)
			} else if tc.mockError == nil {
				mockService.EXPECT().GetURL(gomock.Any(), tc.alias).Return(tc.url, nil)
			} else {
				mockService.EXPECT().GetURL(gomock.Any(), tc.alias).Return("", tc.mockError)
			}

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := api.GetRedirect(ts.URL + "/url/" + tc.alias)

			if tc.respError == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.url, resp.Header.Get("Location"))
			} else if tc.respError != "" {
				require.Error(t, err)
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				require.NoError(t, err)
				assert.Contains(t, string(body), tc.respError)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		mockError  error
		statusCode int
		respError  string
		user       string
		password   string
	}{
		{
			name:       "Success",
			alias:      "example",
			statusCode: http.StatusOK,
			respError:  "",
			user:       "123",
			password:   "123",
		},
		{
			name:       "Empty alias",
			alias:      "",
			statusCode: http.StatusBadRequest,
			respError:  "alias cannot be empty",
			user:       "123",
			password:   "123",
		},
		{
			name:       "Internal server error",
			alias:      "error",
			mockError:  errors.New("failed delete url"),
			statusCode: http.StatusInternalServerError,
			respError:  "failed delete url",
			user:       "123",
			password:   "123",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		Server: config.Server{
			User:     "123",
			Password: "123",
		},
	}
	log := slogdiscard.NewDiscardLogger()
	service := mock_service.NewMockServicer(ctrl)
	handler := handlers.NewHandler(log, service)
	r := router.NewRouter(handler, cfg)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.alias == "" {
				service.EXPECT().
					DeleteURL(gomock.Any(), tc.alias).Times(0)
			} else if tc.mockError == nil {
				service.EXPECT().
					DeleteURL(gomock.Any(), tc.alias).Return(nil)
			} else {
				service.EXPECT().
					DeleteURL(gomock.Any(), tc.alias).Return(tc.mockError)
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			req.SetBasicAuth(tc.user, tc.password)
			require.NoError(t, err)
			r.ServeHTTP(w, req)

			var resp handlers.Response

			body := w.Body.String()
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
			require.Equal(t, tc.statusCode, w.Code)

		})
	}
}
