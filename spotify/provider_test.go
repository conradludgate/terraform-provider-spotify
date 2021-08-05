package spotify_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/conradludgate/terraform-provider-spotify/spotify"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jarcoal/httpmock"
	"golang.org/x/oauth2"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = spotify.Provider()
	testAccProviders = map[string]*schema.Provider{
		"spotify": testAccProvider,
	}
}

func RegisterAuthResponse(apiKey, accessToken string) {
	httpmock.RegisterResponder("POST", "http://localhost:27228/api/v1/token/terraform",
		RespondWith(
			JSON(oauth2.Token{AccessToken: accessToken, Expiry: time.Now().Add(time.Hour)}),
			VerifyBasicAuth("SpotifyAuthProxy", apiKey),
		),
	)
}

type Verifier func(req *http.Request) error

func VerifyBearer(accessToken string) Verifier {
	return func(req *http.Request) error {
		if req.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", accessToken) {
			return errors.New("invalid access token")
		}
		return nil
	}
}

func VerifyBasicAuth(username, password string) Verifier {
	return func(req *http.Request) error {
		user, pass, ok := req.BasicAuth()
		if !ok {
			return errors.New("missing auth")
		}
		if user != username || pass != password {
			return errors.New("invalid auth")
		}
		return nil
	}
}

type object map[string]interface{}
type array []interface{}

func VerifyJSONBody(expected interface{}) Verifier {
	return func(req *http.Request) error {
		if req.Body == nil {
			return errors.New("no body")
		}

		if req.Header.Get("content-type") != "application/json" {
			return errors.New("no json body")
		}

		reqBody := make(map[string]interface{})
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			return errors.New("could not read body")
		}

		expJson, err := json.Marshal(expected)
		if err != nil {
			return errors.New("could not encode expected json")
		}

		expBody := make(map[string]interface{})
		if err := json.Unmarshal(expJson, &expBody); err != nil {
			return errors.New("could not decode expected json")
		}

		if diff := deep.Equal(reqBody, expBody); diff != nil {
			return fmt.Errorf("unexpected request:\n\t%s", strings.Join(diff, "\n\t"))
		}

		return nil
	}
}

func RespondWith(responder httpmock.Responder, verifiers ...Verifier) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		for _, verifier := range verifiers {
			if err := verifier(req); err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, err.Error()), nil
			}
		}

		return responder(req)
	}
}

func JSON(response interface{}) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(http.StatusOK, response)
	}
}
