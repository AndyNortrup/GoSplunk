package splunk

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

const accountName = "testing_user"
const password = "TestAccount"

func TestGetSessionKey(t *testing.T) {
	c := &Client{BaseURL: LocalSplunkMgmntURL}
	sessionKey, err := c.getSessionKey(accountName, password)
	if err != nil {
		t.Fatalf("Failed to get session key: %v\n", err)
	}
	if len(sessionKey.Message) > 0 {
		t.Log("Recieved Message instead of seassion key.")
		t.Fail()
	}

	if len(sessionKey.SessionKey) == 0 {
		t.Log("Recieved zero length session key.")
		t.Fail()
	}
}

/**
* TestRestResponse is an integration test to establish if a rest response can
* be pulled from a local test server
* Requires a running Splunk server and a configured accountName and password.
 */

func TestRestResponse(t *testing.T) {
	c, err := NewClientFromLogin(accountName,
		password, "", "", LocalSplunkMgmntURL, false)
	if err != nil {
		t.Fatal("Unable to get access token.  Check that Splunk is running.")
	}

	response, err := c.GetEntities([]string{"services", "properties"})
	if err != nil {
		t.Fatal("Error querying endpoint, check that Splunk is running")
	}

	const expecting = "properties"
	if response.Title != expecting {
		t.Logf("Expecting title: %v\t Received: %v", expecting, response.Title)
		t.Fail()
	}
}

const passwordRestResponse string = `<?xml-stylesheet type="text/xml" href="/static/atom.xsl"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:s="http://dev.splunk.com/ns/rest" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">
  <title>passwords</title>
  <id>https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords</id>
  <updated>2016-07-18T22:00:32-07:00</updated>
  <generator build="debde650d26e" version="6.4.1"/>
  <author>
    <name>Splunk</name>
  </author>
  <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/_new" rel="create"/>
  <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/_reload" rel="_reload"/>
  <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/_acl" rel="_acl"/>
  <opensearch:totalResults>1</opensearch:totalResults>
  <opensearch:itemsPerPage>30</opensearch:itemsPerPage>
  <opensearch:startIndex>0</opensearch:startIndex>
  <s:messages/>
  <entry>
    <title>:616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com:</title>
    <id>https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A</id>
    <updated>2016-07-18T22:00:32-07:00</updated>
    <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A" rel="alternate"/>
    <author>
      <name>admin</name>
    </author>
    <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A" rel="list"/>
    <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A/_reload" rel="_reload"/>
    <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A" rel="edit"/>
    <link href="/servicesNS/nobody/TA-GoogleFitness/storage/passwords/%3A616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com%3A" rel="remove"/>
    <content type="text/xml">
      <s:dict>
        <s:key name="clear_password">clear_password_value</s:key>
        <s:key name="eai:acl">
          <s:dict>
            <s:key name="app">TA-GoogleFitness</s:key>
            <s:key name="can_change_perms">1</s:key>
            <s:key name="can_list">1</s:key>
            <s:key name="can_share_app">1</s:key>
            <s:key name="can_share_global">1</s:key>
            <s:key name="can_share_user">1</s:key>
            <s:key name="can_write">1</s:key>
            <s:key name="modifiable">1</s:key>
            <s:key name="owner">admin</s:key>
            <s:key name="perms">
              <s:dict>
                <s:key name="read">
                  <s:list>
                    <s:item>admin</s:item>
                  </s:list>
                </s:key>
                <s:key name="write">
                  <s:list>
                    <s:item>admin</s:item>
                  </s:list>
                </s:key>
              </s:dict>
            </s:key>
            <s:key name="removable">1</s:key>
            <s:key name="sharing">app</s:key>
          </s:dict>
        </s:key>
        <s:key name="encr_password">$1$Yh95L+vjruRO5+RCGkSBik8MGBFZ3yiTfw==</s:key>
        <s:key name="password">********</s:key>
        <s:key name="realm"></s:key>
        <s:key name="username">616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com</s:key>
      </s:dict>
    </content>
  </entry>
</feed>`

func TestRestUnmarshal(t *testing.T) {
	decoder := xml.NewDecoder(strings.NewReader(passwordRestResponse))
	resp := &RestResponse{}
	err := decoder.Decode(resp)
	if err != nil {
		log.Fatalf("Unable to decode REST response. \n%v\n", err)
	}

	expectedKeyZeroName := "clear_password"
	result := resp.Entries[0].Contents.Keys[0].Name
	if result != expectedKeyZeroName {
		log.Printf("Incorrect Key Name. Expected: %v, Recived: %v",
			expectedKeyZeroName,
			result)
		t.Fail()
	}

	expectedKeyZeroValue := "clear_password_value"
	result = resp.Entries[0].Contents.Keys[0].Value
	if result != expectedKeyZeroValue {
		log.Printf("Incorrect Key Name. Expected: %v, Recived: %v",
			expectedKeyZeroValue,
			result)
		t.Fail()
	}
}

func TestUpdateKVStore(t *testing.T) {
	c, err := NewClientFromLogin(accountName,
		password,
		"fitness_for_splunk",
		"nobody",
		LocalSplunkMgmntURL,
		false)

	if err != nil {
		t.Fatal("Unable to get access token.  Check that Splunk is running.")
	}

	type User struct {
		Name         string `json:"name"`
		UserID       string `json:"id"`
		Scope        []string
		oauth2.Token `json:"token"`
		TokenExpiry  string `json:"token_expiry"`
	}

	token := oauth2.Token{
		AccessToken:  "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiI0VDVaVzYiLCJhdWQiOiIyMjdNVkoiLCJpc3MiOiJGaXRiaXQiLCJ0eXAiOiJhY2Nlc3NfdG9rZW4iLCJzY29wZXMiOiJyc29jIHJzZXQgcmFjdCBybG9jIHJ3ZWkgcmhyIHJudXQgcnBybyByc2xlIiwiZXhwIjoxNDcwODIwMzkwLCJpYXQiOjE0NzA3OTE1OTB9.gXKGs_FsTXNC6ayyAfLbPGIMytgV-Jz4XdiilUrbQwU",
		RefreshToken: "4a724b82c52fd0f556402e270f27674be5ae5dfb2646b213ced87512e06aa9a8",
		TokenType:    "Bearer",
	}
	payload := &User{
		UserID: "4T5ZW6",
		Name:   "Andrew CHANGED Nortrup",
		Scope:  []string{"settings nutrition activity heartrate weight sleep location social profile"},
		Token:  token,
	}

	err = c.KVStoreUpdate("fitbit_tokens", "4T5ZW6", payload)

	if err != nil {
		log.Printf("Error updating KV Store: %v", err)
	}

	reader, err := c.KVStoreGetCollection("fitbit_tokens")
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var result []User
	err = decoder.Decode(&result)
	if err != nil {
		b, _ := ioutil.ReadAll(reader)
		log.Fatalf("Failed to decode updated KVStore Key\nError:%v\nResponse:%s", err, b)
	}

	if result[0].Name != "Andrew CHANGED Nortrup" {
		t.Fail()
		t.Logf("Failed to update KV Store record: %v", result[0])
	}
}

func TestDecodePasswordStore(t *testing.T) {
	input := `<?xml-stylesheet type="text/xml" href="/static/atom.xsl"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:s="http://dev.splunk.com/ns/rest" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">
  <title>passwords</title>
  <id>https://localhost:8089/servicesNS/nobody/TA-FitnessTrackers/storage/passwords</id>
  <updated>2016-08-10T06:43:30-07:00</updated>
  <generator build="f2c836328108" version="6.4.0"/>
  <author>
    <name>Splunk</name>
  </author>
  <link href="/servicesNS/nobody/TA-FitnessTrackers/storage/passwords/_new" rel="create"/>
  <link href="/servicesNS/nobody/TA-FitnessTrackers/storage/passwords/_reload" rel="_reload"/>
  <link href="/servicesNS/nobody/TA-FitnessTrackers/storage/passwords/_acl" rel="_acl"/>
  <opensearch:totalResults>3</opensearch:totalResults>
  <opensearch:itemsPerPage>30</opensearch:itemsPerPage>
  <opensearch:startIndex>0</opensearch:startIndex>
  <s:messages/>
  
</feed>`
	conf := &ModInputConfig{}
	err := xml.Unmarshal([]byte(input), conf)
	if err != nil {
		t.Fatalf("Failed to unmarshal input. %v", err)
	}
}
