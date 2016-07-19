package splunk

import (
	"encoding/xml"
	"log"
	"strings"
	"testing"
)

const accountName = "testing_user"
const password = "TestAccount"

func TestGetSessionKey(t *testing.T) {
	sessionKey, err := NewSessionKey(accountName, password)
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
	sessionKey, err := NewSessionKey(accountName, password)
	if err != nil {
		t.Fatal("Unable to get access token.  Check that Splunk is running.")
	}

	response, err := GetEntities([]string{"services", "properties"}, "", "", sessionKey.SessionKey)
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
