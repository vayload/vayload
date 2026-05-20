package client

import (
	"net/http"

	httpi "github.com/vayload/vayload/pkg/http"
)

func NewClient() *httpi.HttpClient {
	return httpi.NewHttpClient(httpi.HttpClientConfig{
		BaseURL: "http://unix",
		Transport: &http.Transport{
			DialContext: GetDialer("vayload"),
		},
	})
}
