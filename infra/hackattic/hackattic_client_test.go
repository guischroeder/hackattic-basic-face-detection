package hackattic

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"hackattic-basic-face-detection/infra/httpclient"
)

var (
	getRequestFunc func(url string) (*http.Response, error)
)

type httpClientGetMock struct {}
func (h *httpClientGetMock) Get(request string) (*http.Response, error) {
    return getRequestFunc(request)
}

func TestGetProblem(t *testing.T) {
    response := `{
        "image_url": "http://hackattic.com/image"
    }`
    getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(response)),
		}, nil
	}
    httpclient.Client = &httpClientGetMock{}

    hm := HackatticClient{
        AccessToken: "token",
    }
    result, err := hm.GetProblem()
    if err != nil {
        t.Error("TestGetProblem failed")
        return
    }

    wait := Problem{
        ImageUrl: "http://hackattic.com/image",
    }

    assert.Equal(t, result, wait)
}
