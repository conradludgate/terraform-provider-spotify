package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// ClientConfigurer for spotify API access
func ClientConfigurer(d *schema.ResourceData) (interface{}, error) {
	server, err := url.Parse(d.Get("auth_server").(string))
	if err != nil {
		return nil, fmt.Errorf("auth_server was not a valid url: %w", err)
	}
	server.Path = path.Join(server.Path, "api/v1/token")
	server.Path = path.Join(server.Path, d.Get("token_id").(string))

	transport := &transport{
		Endpoint: server.String(),
		Username: d.Get("username").(string),
		APIKey:   d.Get("api_key").(string),
	}

	client := spotify.NewClient(&http.Client{
		Transport: transport,
	})
	return &client, nil
}

type transport struct {
	Endpoint string
	Username string
	APIKey   string
	Base     http.RoundTripper
	token    *oauth2.Token
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.token.Valid() {
		if err := t.getToken(); err != nil {
			return nil, err
		}
	}

	t.token.SetAuthHeader(req)

	return t.base().RoundTrip(req)
}

func (t *transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func (t *transport) getToken() error {
	req, err := http.NewRequest("POST", t.Endpoint, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.Username, t.APIKey)
	resp, err := t.base().RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	t.token = &oauth2.Token{}
	if err := json.Unmarshal(body, t.token); err != nil {
		return err
	}

	if !t.token.Valid() {
		return errors.New("could not get a valid token")
	}

	return nil
}
