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

func NewSessionKey(username string, password string, baseURL string) (*SessionKey, error) {

	u, err := url.ParseRequestURI(baseURL)
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

	client := newSplunkHttpClient(false)
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

func GetEntities(baseURL string,
	path []string,
	namespace string,
	owner string,
	sessionKey string) (*RestResponse, error) {

	u, err := buildRequestPath(baseURL, path, namespace, owner)
	if err != nil {
		return &RestResponse{}, err
	}

	resp, err := makeGetRestRequest(u, sessionKey)
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
func KVStoreGetCollection(baseURL string,
	collection string,
	namespace string,
	owner string,
	sessionKey string) (io.ReadCloser, error) {

	u, err := buildRequestPath(baseURL, []string{"storage", "collections", "data", collection},
		namespace, owner)

	fmt.Printf("URL: %v", u)

	if err != nil {
		return nil, err
	}

	resp, err := makeGetRestRequest(u, sessionKey)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status: " + resp.Status + ":" + u.RequestURI())
	}

	return resp.Body, nil
}

func KVStoreUpdate(baseURL, collection, id string,
	payload interface{},
	namespace, owner, sessionKey string) error {
	u, err := buildRequestPath(baseURL,
		[]string{"storage", "collections", "data", collection, id},
		namespace, owner)

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
	req.Header.Add("Authorization", "Splunk "+sessionKey)

	if err != nil {
		return err
	}
	client := newSplunkHttpClient(false)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

func makeGetRestRequest(u *url.URL, sessionKey string) (*http.Response, error) {
	//Create the Request
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%v", u), nil)
	r.Header.Add("Authorization", "Splunk "+sessionKey)

	//Create a client
	client := newSplunkHttpClient(false)
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
func buildRequestPath(baseURL string,
	pieces []string,
	namespace string,
	owner string) (*url.URL, error) {

	if len(pieces) < 2 {
		return &url.URL{}, errors.New("Not enough path specifications.")
	}

	//Create Request urlStr
	u, _ := url.ParseRequestURI(baseURL)

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
