package tests

import (
	"net/http"
	"net/url"
	"testing"
	"time"
	"url-shortener/internal/http/handlers/links/shorten"
	"url-shortener/internal/lib/api"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const host = "localhost:8080"

type testCase struct {
	name      string
	sourceURL string
	expiresAt *time.Time
	error     string
	status    int
}

func newHTTPExpect(t *testing.T) *httpexpect.Expect {
	u := url.URL{Scheme: "http", Host: host}
	return httpexpect.Default(t, u.String())
}

func TestURLShortener_ShortenRedirect(t *testing.T) {
	expiresAtFuture := gofakeit.FutureDate()
	expiresAtPast := gofakeit.PastDate()

	testCases := []testCase{
		{
			"Valid request",
			gofakeit.URL(),
			&expiresAtFuture,
			"",
			http.StatusOK},
		{
			"Invalid sourceURL",
			"invalid_url",
			&expiresAtFuture,
			"invalid request: Key: 'Request.SourceURL' Error:Field validation for 'SourceURL' failed on the 'url' tag",
			http.StatusUnprocessableEntity},
		{
			"Expired expiresAt",
			gofakeit.URL(),
			&expiresAtPast,
			"date cannot be expired",
			http.StatusUnprocessableEntity},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := newHTTPExpect(t)
			testOrchestrator(t, e, tc)
		})
	}
}

func testOrchestrator(t *testing.T, e *httpexpect.Expect, tc testCase) {
	shortURL := testShorten(t, e, tc)
	if shortURL != "" {
		testRedirect(t, shortURL, tc.sourceURL)
	}
}

func testShorten(t *testing.T, e *httpexpect.Expect, tc testCase) string {
	resp := e.POST("/links").
		WithJSON(shorten.Request{SourceURL: tc.sourceURL, ExpiresAt: tc.expiresAt}).
		Expect().
		Status(tc.status)

	if tc.error != "" {
		resp.Body().Contains(tc.error) // Проверяем наличие ошибки
		return ""
	}

	return resp.JSON().Object().Value("shortUrl").String().Raw()
}

func testRedirect(t *testing.T, shortURL, expectedURL string) {
	redirectedToURL, err := api.GetRedirect(shortURL)
	require.NoError(t, err)
	require.Equal(t, expectedURL, redirectedToURL)
}
