package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//NOTE: Get Access Token from gh using code

func GetAccessToken(code string) (string, error) {

	request := map[string]string{
		"client_id":     ClientID,
		"client_secret": ClientSecret,
		"code":          code,
	}

	requestBytes, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestBytes))

	if err != nil {

		log.Printf("Error sending github:%s\n", err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.Client{}

	response, err := client.Do(req)

	if err != nil {

		log.Printf("Error sending request:%s\n", err.Error())

		return "", err
	}

	body, _ := io.ReadAll(response.Body)

	var ghResponse ghAccessTokenResponse

	json.Unmarshal(body, &ghResponse)

	return ghResponse.AccessToken, nil
}

//NOTE: Get user data from gh using token

func GetUserData(token string) (string, error) {

	reqq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)

	headerValue := fmt.Sprintf("token %s", token)

	reqq.Header.Set("Authorization", headerValue)

	resp, _ := http.DefaultClient.Do(reqq)

	responseBody, _ := io.ReadAll(resp.Body)

	return string(responseBody), nil
}
