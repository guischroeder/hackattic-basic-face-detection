package hackattic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"hackattic-basic-face-detection/infra/httpclient"
)

type HackatticClient struct {
    AccessToken string
    HttpClient http.Client
}

type Problem struct {
    ImageUrl string `json:"image_url"`
}

const BASE_URL = "https://hackattic.com/challenges/basic_face_detection"

func (h HackatticClient) GetProblem() (Problem, error) {
    url := fmt.Sprintf("%s/problem?access_token=%s", BASE_URL, h.AccessToken)
    response, err := httpclient.Client.Get(url)
    problem := Problem{}
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

type Solution struct {
    FaceTiles [][2]int `json:"face_tiles"`
}

type Result struct {
    Result string `json:"result"`
    Rejected string `json:"rejected"`
}

func (h HackatticClient) SubmitSolution(faceTiles [][2]int) (Result, error) {
    url := fmt.Sprintf("%s/solve?access_token=%s", BASE_URL, h.AccessToken)
    result := Result{}

    buffer := new(bytes.Buffer)
    solution := Solution{FaceTiles: faceTiles}
    json.NewEncoder(buffer).Encode(solution)

    request, _ := http.NewRequest(http.MethodPost, url, buffer)
    request.Header.Add("Content-Type", "application/json")

    response, err := h.HttpClient.Do(request)
    if err != nil {
        return result, err
    }

    defer response.Body.Close()
    
    err = json.NewDecoder(response.Body).Decode(&result)

    if err != nil {
        return result, err
    }

    return result, nil
}
