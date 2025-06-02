package main

import (
	"crypto/rand"
	"encoding/base64"
)

type SpotifyAuthToken struct {
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
	Expiration int    `json:"expires_in"`
}

func (s *SpotifyAuthToken) RetrieveTokenPKCE(clientId string, clientSecret string) error {
	// body := bytes.NewBufferString(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientId, clientSecret))

	// res, err := http.Post("https://accounts.spotify.com/api/token", "application/x-www-form-urlencoded", body)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()

	// var responseBody []byte
	// i, err := res.Body.Read(responseBody)
	// if err != nil {
	// 	return err
	// }

	// if i == 0 {
	// 	return errors.New("the authorization token request returned no data")
	// }

	// err = json.Unmarshal(responseBody, s)
	// if err != nil {
	// 	return err
	// }

	randStr := make([]byte, 64)
	_, err := rand.Read(randString)
	if err != nil {
		return err
	}

	randStrEncoded64 := base64.URLEncoding.EncodeToString(randStr)

	codeVerifier := randStrEncoded64[:len(randStrEncoded64)]

	return nil
}
