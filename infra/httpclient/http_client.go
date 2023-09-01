package httpclient

import (
	"net/http"
)

type client struct{}

type HttpClientInterface interface {
    Get(string) (*http.Response, error)
}
var (
		Client HttpClientInterface = &client{}
)

func (h *client) Get(url string) (*http.Response, error) {
		request, err := http.NewRequest(http.MethodGet, url, nil)
		request.Header.Add("Accept", `application/json`)
		if err != nil {
				return nil, err
		}
		client := http.Client{}

		return client.Do(request)
}

