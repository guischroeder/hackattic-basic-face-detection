package hackattic

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type HackatticClient struct {
    BaseUrl string
    HttpClient http.Client
}

type Problem struct {
    ImageUrl string `json:"image_url"`
}

func (hm HackatticClient) GetProblem(accessToken string) (Problem, error) {
    url := fmt.Sprintf("%s/problem?access_token=%s", hm.BaseUrl, accessToken)
    request, err := http.NewRequest(http.MethodGet, url, nil)

    problem := Problem{}
    if err != nil {
        return problem, err
    }

    response, err := hm.HttpClient.Do(request)
    if err != nil {
        return problem, err
    }

    defer response.Body.Close()

    err = json.NewDecoder(response.Body).Decode(&problem)

    if err != nil {
        return problem, err
    }

    return problem, nil
}
