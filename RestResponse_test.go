package splunk

import "testing"

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
* Requires a running Splunk server and that TA-GoogleFitness has a local/passwords.conf
* in the TA-GoogleFitness app.
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
