package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

const (
	spotifyAuthUrl  = "https://accounts.spotify.com/authorize"
	spotifyTokenUrl = "https://accounts.spotify.com/api/token"
)

type SpotifyAuthRequestPKCE struct {
	ClientId        string
	ResponseType    string
	Redirect        string
	State           string
	Scope           string
	ChallengeMethod string
	CodeChallenge   string
}

type SpotifyTokenRequestPKCE struct {
	GrantType    string
	Code         string
	Redirect     string
	ClientId     string
	CodeVerifier string
}

func main() {
	v := NewVerifier()
	c := GenerateCodeChallenge(v)
	req, err := NewSpotifyAuthRequestPKCE("test", "null", "test test2", c, "test.com")
	if err != nil {
		fmt.Println(err)
	}
	url, err := req.GenerateSpotifyAuthUrl()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(url)
}

func (s *SpotifyAuthRequestPKCE) GenerateSpotifyAuthUrl() (string, error) {
	baseUrl, err := url.Parse(spotifyAuthUrl)
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"client_id":             s.ClientId,
		"response_type":         s.ResponseType,
		"redirect_url":          s.Redirect,
		"state":                 s.State,
		"scope":                 s.Scope,
		"code_challenge_method": s.ChallengeMethod,
		"code_challenge":        s.CodeChallenge,
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}
	baseUrl.RawQuery = values.Encode()

	return baseUrl.String(), nil
}

func NewSpotifyAuthRequestPKCE(clientId string, state string, scope string, codeChallenge string, redirect string) (*SpotifyAuthRequestPKCE, error) {
	if clientId == "" || redirect == "" || codeChallenge == "" {
		err := errors.New("one or more required params are empty strings [clientId, codeChallenge, redirect]")
		return nil, err
	} else {
		return &SpotifyAuthRequestPKCE{
			ClientId:        clientId,
			ResponseType:    "code",
			Redirect:        redirect,
			State:           state,
			Scope:           scope,
			ChallengeMethod: "SA256",
			CodeChallenge:   codeChallenge,
		}, nil
	}
}

func NewVerifier() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))
}

func GenerateCodeChallenge(verifier string) string {
	hash := sha256.New()
	hash.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
