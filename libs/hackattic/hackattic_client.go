package hackattic

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type HackatticClient struct {
    AccessToken string
    HttpClient http.Client
}

type Problem struct {
    ImageUrl string `json:"image_url"`
}

func (h HackatticClient) GetProblem() (Problem, error) {
    baseUrl := "https://hackattic.com/challenges/basic_face_detection"
    url := fmt.Sprintf("%s/problem?access_token=%s", baseUrl, h.AccessToken)
    request, err := http.NewRequest(http.MethodGet, url, nil)

    problem := Problem{}
    if err != nil {
        return problem, err
    }

    response, err := h.HttpClient.Do(request)
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
