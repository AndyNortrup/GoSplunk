package splunk

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

const modInputConfigExample string = `<input>
		<server_host>myHost</server_host>
		<server_uri>https://127.0.0.1:8089</server_uri>
		<session_key>123102983109283019283</session_key>
		<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
		<configuration>
			<stanza name="myScheme://aaa">
					<param name="param1">value1</param>
					<param name="param2">value2</param>
					<param name="disabled">0</param>
					<param name="index">default</param>
			</stanza>
		</configuration>
	</input>`

func TestDecode(t *testing.T) {

	result := &ModInputConfig{}
	err := xml.Unmarshal([]byte(modInputConfigExample), result)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	modInputConfigCases(result, t)
}

func TestDecodeMethod(t *testing.T) {
	reader := bytes.NewReader([]byte(modInputConfigExample))
	result, err := ReadModInputConfig(reader)
	if err != nil {
		t.Logf("Error reading Mod Input Config: %v\n", err)
	}
	modInputConfigCases(result, t)
}

func modInputConfigCases(result *ModInputConfig, t *testing.T) {
	if result.ServerHost != "myHost" {
		t.Logf("Incorrect Server Host returned.  Expecting: myHost Got: %v",
			result.ServerHost)
		t.Fail()
	}

	if result.Stanzas[0].StanzaName != "myScheme://aaa" {
		t.Logf("Incorrect Stanza Name Returned.  Expecting: myScheme://aaa Got: %v",
			result.Stanzas[0].StanzaName)
		t.Fail()
	}

	if result.Stanzas[0].Params[0].Value != "value1" {
		t.Logf("Incorrect Stanza Value Returned.  Expecting: value1 Got: %v",
			result.Stanzas[0].Params[0].Value)
		t.Fail()
	}

	if result.Stanzas[0].Params[0].Name != "param1" {
		t.Logf("Incorrect Stanza Name Returned.  Expecting: param1 Got: %v",
			result.Stanzas[0].Params[0].Name)
		t.Fail()
	}
}

func TestReadConfig(t *testing.T) {
	result, err := ReadModInputConfig(strings.NewReader(modInputConfigExample))
	if err != nil {
		t.Logf("Unable to read ModInputConfig: %v\n", err)
		t.Fail()
	}

	if result.SessionKey != "123102983109283019283" {
		t.Logf("Incorrect Session Key returned. Expected: 123102983109283019283\t Received: %v\n", result.SessionKey)
		t.Fail()
	}

	if result.Stanzas[0].StanzaName != "myScheme://aaa" {
		t.Logf("Incorrect Stanza Name returned. Expected: myScheme://aaa\t Received: %v\n", result.Stanzas[0].StanzaName)
		t.Fail()
	}
}
