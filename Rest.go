package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const rest_base_url = "https://localhost:8089"
const login_endpoint = "services/auth/login"

func GetSessionKey(username string, password string) (*SplunkSessionKey, error) {

	u, err := url.ParseRequestURI(rest_base_url)
	if err != nil {
		return &SplunkSessionKey{}, err
	}

	u.Path = login_endpoint

	urlStr := fmt.Sprintf("%v", u)

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)

	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return &SplunkSessionKey{}, err
	}
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := newSplunkHttpClient(false)
	resp, err := client.Do(r)

	if err != nil {
		return &SplunkSessionKey{}, err
	}

	defer resp.Body.Close()

	restResp := &SplunkSessionKey{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(restResp)
	if err != nil {
		return restResp, err
	}

	return restResp, nil
}

func GetEntities(path []string,
	namespace string,
	owner string,
	sessionKey string) (*SplunkRestResponse, error) {

	u, err := buildRequestPath(path, namespace, owner)
	if err != nil {
		return &SplunkRestResponse{}, err
	}

	//Create the Request
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%v", u), nil)
	r.Header.Add("Authorization", "Splunk "+sessionKey)

	//Create a client
	client := newSplunkHttpClient(false)
	resp, err := client.Do(r)

	if resp.StatusCode != http.StatusOK {
		return &SplunkRestResponse{}, errors.New(resp.Status)
	}

	if err != nil {
		return &SplunkRestResponse{}, err
	}

	defer resp.Body.Close()

	//Decode the response from XML
	decoder := xml.NewDecoder(resp.Body)
	result := &SplunkRestResponse{}

	err = decoder.Decode(result)
	if err != nil {
		return &SplunkRestResponse{}, err
	}

	return result, nil
}

//buildRequestPath builds a path for the REST request
func buildRequestPath(pieces []string, namespace string, owner string) (*url.URL, error) {
	if len(pieces) < 2 {
		return &url.URL{}, errors.New("Not enough path specifications.")
	}

	//Create Request urlStr
	u, _ := url.ParseRequestURI(rest_base_url)

	//Build the address
	if len(namespace) > 0 {
		u.Path += "/servicesNS"

		//Add the user
		if len(owner) > 0 {
			u.Path += "/" + owner
		}

		u.Path += "/" + namespace
	}

	for _, item := range pieces {
		u.Path += "/" + item
	}
	return u, nil
}

func newSplunkHttpClient(validateTLS bool) *http.Client {

	//Splunk ships with self signed certificates and these run on a lot of instances
	// this makes it really hard to do certificate validation
	tr := &http.Transport{}

	if !validateTLS {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &http.Client{Transport: tr}
}
