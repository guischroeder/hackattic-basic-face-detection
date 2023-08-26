package hackattic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProblem(t *testing.T) {
    response := `{
        "image_url": "http://hackattic.com/image"
    }`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        _, _ = w.Write([]byte(response))
    }))

    defer server.Close()

    hm := HackatticClient{
        BaseUrl: server.URL,
    }
    result, err := hm.GetProblem("abc123")
    if err != nil {
        t.Error("TestGetProblem failed")
        return
    }

    wait := Problem{
        ImageUrl: "http://hackattic.com/image",
    }

    assert.Equal(t, result, wait)
}
