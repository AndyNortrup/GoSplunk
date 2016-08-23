package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const LocalSplunkMgmntURL = "https://localhost:8089"
const login_endpoint = "services/auth/login"

type Client struct {
	SessionKey  string
	Namespace   string
	Owner       string
	BaseURL     string
	ValidateTLS bool
}

// NewClientFromSessionKey creates a Client object when the user already has
// a session key.  For instance when you are writing a modular input which comes
// with a session key provided in the configuration.
func NewClientFromSessionKey(sessionKey, namespace, owner, baseURL string,
	validateTLS bool) *Client {
	return &Client{
		SessionKey:  sessionKey,
		Namespace:   namespace,
		Owner:       owner,
		BaseURL:     baseURL,
		ValidateTLS: validateTLS,
	}
}

// NewClientFromLogin creates a Client object is used when the user must provide
// their credentials in order to log in.
func NewClientFromLogin(username, password, namespace, owner, baseURL string,
	validateTLS bool) (*Client, error) {

	c := Client{
		Namespace:   namespace,
		Owner:       owner,
		BaseURL:     baseURL,
		ValidateTLS: validateTLS,
	}

	key, err := c.getSessionKey(username, password)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		SessionKey:  key.SessionKey,
		Namespace:   namespace,
		Owner:       owner,
		BaseURL:     baseURL,
		ValidateTLS: validateTLS,
	}, nil
}

func (c *Client) getSessionKey(username, password string) (*SessionKey, error) {

	u, err := url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return &SessionKey{}, err
	}

	u.Path = login_endpoint

	urlStr := fmt.Sprintf("%v", u)

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)

	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return &SessionKey{}, err
	}
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := c.newSplunkHttpClient()
	resp, err := client.Do(r)

	if err != nil {
		return &SessionKey{}, err
	}

	defer resp.Body.Close()

	restResp := &SessionKey{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(restResp)
	if err != nil {
		return restResp, err
	}

	return restResp, nil
}

func (c *Client) GetEntities(path []string) (*RestResponse, error) {

	u, err := c.buildRequestPath(path)
	if err != nil {
		return &RestResponse{}, err
	}

	resp, err := c.makeGetRestRequest(u)
	if err != nil {
		return &RestResponse{}, err
	}

	defer resp.Body.Close()

	//Decode the response from XML
	decoder := xml.NewDecoder(resp.Body)
	result := &RestResponse{}

	err = decoder.Decode(result)
	if err != nil {
		return &RestResponse{}, err
	}

	return result, nil
}

// KVStoreGetCollection returns values from a KV Store collection.  Result is a
// io.ReadCloser so that it can be JSON decoder.
func (c *Client) KVStoreGetCollection(collection string) (io.ReadCloser, error) {

	u, err := c.buildRequestPath([]string{"storage", "collections", "data", collection})

	if err != nil {
		return nil, err
	}

	resp, err := c.makeGetRestRequest(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status: " + resp.Status + ":" + u.RequestURI())
	}

	return resp.Body, nil
}

func (c *Client) KVStoreUpdate(collection, id string, payload interface{}) error {
	u, err := c.buildRequestPath([]string{"storage", "collections", "data", collection, id})

	if err != nil {
		return err
	}

	//Encode the payload into JSON
	b, err := json.Marshal(payload)

	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v", u), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Splunk "+c.SessionKey)

	if err != nil {
		return err
	}
	client := c.newSplunkHttpClient()
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) makeGetRestRequest(u *url.URL) (*http.Response, error) {
	//Create the Request
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%v", u), nil)
	r.Header.Add("Authorization", "Splunk "+c.SessionKey)

	//Create a client
	client := c.newSplunkHttpClient()
	resp, err := client.Do(r)

	if resp.StatusCode != http.StatusOK {
		return &http.Response{}, errors.New(resp.Status)
	}

	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}

//buildRequestPath builds a path for the REST request
func (c *Client) buildRequestPath(pieces []string) (*url.URL, error) {

	if len(pieces) < 2 {
		return &url.URL{}, errors.New("Not enough path specifications.")
	}

	//Create Request urlStr
	u, _ := url.ParseRequestURI(c.BaseURL)

	//Build the address
	if len(c.Namespace) > 0 {
		u.Path += "/servicesNS"

		//Add the user
		if len(c.Owner) > 0 {
			u.Path += "/" + c.Owner
		}

		u.Path += "/" + c.Namespace
	}

	for _, item := range pieces {
		u.Path += "/" + item
	}
	return u, nil
}

func (c *Client) newSplunkHttpClient() *http.Client {

	//Splunk ships with self signed certificates and these run on a lot of instances
	// this makes it really hard to do certificate validation
	tr := &http.Transport{}

	if !c.ValidateTLS {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &http.Client{Transport: tr}
}
