package tests

import (
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/SemenShakhray/url-shortener/internal/handlers"
	"github.com/SemenShakhray/url-shortener/pkg/api"
	"github.com/SemenShakhray/url-shortener/pkg/random"
	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPatg(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())
	t.Log("URL:", u.String())

	alias, err := random.NewRandomString(10)
	require.NoError(t, err)
	e.POST("/url").
		WithJSON(handlers.Request{
			URL:   gofakeit.URL(),
			Alias: alias,
		}).WithBasicAuth("Sam", "qwerty").
		Expect().
		JSON().Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		alias  string
		error  string
		status int
	}{
		{
			name:   "Valid URL",
			url:    gofakeit.URL(),
			alias:  gofakeit.Word() + gofakeit.Word(),
			status: http.StatusOK,
		},
		{
			name:   "Invalid URL",
			url:    "invalid url",
			alias:  gofakeit.Word(),
			error:  "field URL is not valid URL",
			status: http.StatusBadRequest,
		},
		{
			name:   "Empty URL",
			url:    "",
			alias:  gofakeit.Word(),
			error:  "field URL is required field",
			status: http.StatusBadRequest,
		},
		{
			name:   "Empty alias",
			url:    gofakeit.URL(),
			alias:  "",
			status: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			//test save
			resp := e.POST("/url").
				WithJSON(handlers.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("Sam", "qwerty").
				Expect().Status(tc.status).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			//test redirect
			testRedirect(t, alias, tc.url)

			//test remove
			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("Sam", "qwerty").
				Expect().Status(tc.status).
				JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			testRedirectNotFound(t, alias)

		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path.Join("url/", alias),
	}

	resp, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, resp.Header.Get("Location"))
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path.Join("url/", alias),
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)

}
