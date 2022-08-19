package lib

import (
	"context"
	"net/http"
	"net/url"

	cc "golang.org/x/oauth2/clientcredentials"
)

func ClientCredentials(id string, secret string, scope string, tokenUrl string) *http.Client {

	config := cc.Config{
		ClientID:     id,
		ClientSecret: secret,
		Scopes:       []string{scope},
		TokenURL:     tokenUrl,
		EndpointParams: url.Values{
			"grant_type": {"client_credentials"},
		}}

	client := config.Client(context.Background())

	return client
}
