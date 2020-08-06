package nextengine

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	authHost = "https://base.next-engine.org"
	apiHost  = "https://api.next-engine.org"
	// Success is API Result constant
	Success = "success"
	// Error is API Result constant
	Error = "error"
	// Redirect is API Result constant
	Redirect = "redirect"
)

type commonResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// APIResponse is Structure that represents API response
type APIResponse struct {
	commonResult
	Token
	Count string
	Data  []map[string]interface{}
}

// Token is Structure that represents a api token
type Token struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenEndDate  string `json:"access_token_end_date"`
	RefreshTokenEndDate string `json:"refresh_token_end_date"`
}

// TokenRepository is API token write/read interface
// If you want to change the storage location of API token to DB or session, you need to implement this interface
type TokenRepository interface {
	Token(context.Context) (Token, error)
	Save(context.Context, Token) error
}

// Config is Structure holding the settings of NextEngine API client
type Config struct {
	clientID        string
	clientSecret    string
	redirectURI     string
	HTTPClient      *http.Client
	tokenRepository TokenRepository
}

// NewDefaultClient Returns an instance of the API client with default settings
func NewDefaultClient(clientID string, clientSecret string, redirectURI string, accessToken string, refreshToken string) *Config {
	cli := &http.Client{}
	tr := &defaultTokenRepository{
		t: Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	return NewClient(clientID, clientSecret, redirectURI, cli, tr)
}

// NewClient Returns an instance of the API client
func NewClient(clientID string, clientSecret string, redirectURI string, httpClient *http.Client, tr TokenRepository) *Config {
	return &Config{
		clientID:        clientID,
		clientSecret:    clientSecret,
		redirectURI:     redirectURI,
		HTTPClient:      httpClient,
		tokenRepository: tr,
	}
}

// SignInURI Returns the URI of the authentication screen of Nexe Engine
func (c *Config) SignInURI(extraParam url.Values) string {
	u, _ := url.Parse(authHost + "/users/sign_in/")

	v := url.Values{}
	v.Add("client_id", c.clientID)
	v.Add("redirect_uri", c.redirectURI)
	for key, vals := range extraParam {
		for _, val := range vals {
			v.Add(key, val)
		}
	}
	u.RawQuery = v.Encode()
	return u.String()
}

func newRequest(ctx context.Context, method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// Authorize Fetch API token using uid and state
func (c *Config) Authorize(ctx context.Context, uid string, state string) (*APIResponse, error) {
	v := url.Values{}
	v.Add("client_id", c.clientID)
	v.Add("client_secret", c.clientSecret)
	v.Add("uid", uid)
	v.Add("state", state)

	return c.request(ctx, "/api_neauth", nil, v)
}

// APIExecute is Execute the API and return the result
// Please specify a path starting with / for endpoint
func (c *Config) APIExecute(ctx context.Context, endpoint string, params map[string]string) (*APIResponse, error) {
	v := url.Values{}

	tok, err := c.tokenRepository.Token(ctx)
	if err != nil {
		return nil, err
	}
	v.Add("access_token", tok.AccessToken)
	v.Add("refresh_token", tok.RefreshToken)

	return c.request(ctx, endpoint, params, v)
}

// APIExecuteNoRequiredLogin is Execute API that does not require login and return the result
// Please specify a path starting with / for endpoint
func (c *Config) APIExecuteNoRequiredLogin(ctx context.Context, endpoint string, params map[string]string) (*APIResponse, error) {
	v := url.Values{}
	v.Add("client_id", c.clientID)
	v.Add("client_secret", c.clientSecret)
	return c.request(ctx, endpoint, params, v)
}

func (c *Config) request(ctx context.Context, endpoint string, params map[string]string, extraParams url.Values) (*APIResponse, error) {
	u, _ := url.Parse(apiHost + endpoint)

	v := url.Values{}
	for key, val := range params {
		v.Add(key, val)
	}

	for key, vals := range extraParams {
		for _, val := range vals {
			v.Add(key, val)
		}
	}

	httpRequest, err := newRequest(ctx, "POST", u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	httpResponse, err := c.HTTPClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return c.responseHandler(ctx, httpResponse.Body)
}

func (c *Config) responseHandler(ctx context.Context, body io.ReadCloser) (*APIResponse, error) {
	defer body.Close()

	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	res := &APIResponse{}
	if err = json.Unmarshal(payload, &res); err != nil {
		return nil, err
	}

	if res.AccessToken != "" && res.RefreshToken != "" {
		if err := c.tokenRepository.Save(ctx, res.Token); err != nil {
			return nil, err
		}
	}

	if res.Result != Success {
		return nil, &APIError{
			commonResult: commonResult{
				Code:    res.Code,
				Message: res.Message,
				Result:  res.Result,
			},
		}
	}

	return res, nil
}
