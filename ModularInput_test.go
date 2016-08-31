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

func TestEncodeScheme(t *testing.T) {

	//Sample taken from
	//http://docs.splunk.com/Documentation/Splunk/6.4.3/AdvancedDev/ModInputsScripts
	testScheme := "<scheme><title>Amazon S3</title>" +
		"<description>Get data from Amazon S3.</description>" +
		"<use_external_validation>true</use_external_validation>" +
		"<streaming_mode>xml</streaming_mode><endpoint><args><arg name=\"name\">" +
		"<title>Resource name</title><description>An S3 resource name without the leading s3://. " +
		"For example, for s3://bucket/file.txt specify bucket/file.txt. " +
		"You can also monitor a whole bucket (for example by specifying 'bucket'), " +
		"or files within a sub-directory of a bucket " +
		"(for example 'bucket/some/directory/'; note the trailing slash)." +
		"</description><data_type>string</data_type><required_on_create>true" +
		"</required_on_create><required_on_edit>false</required_on_edit>" +
		"</arg><arg name=\"key_id\"><title>Key ID</title><description>Your Amazon key ID." +
		"</description><data_type>string</data_type><required_on_create>true</required_on_create>" +
		"<required_on_edit>false</required_on_edit></arg><arg name=\"secret_key\">" +
		"<title>Secret key</title><description>Your Amazon secret key.</description>" +
		"<data_type>string</data_type><required_on_create>true</required_on_create>" +
		"<required_on_edit>false</required_on_edit></arg></args></endpoint></scheme>"

	scheme := NewModInputScheme("Amazon S3",
		"Get data from Amazon S3.", true, StreamingModeXML)
	scheme.AddArgument("name", "Resource name",
		"An S3 resource name without the leading s3://. "+
			"For example, for s3://bucket/file.txt specify bucket/file.txt. "+
			"You can also monitor a whole bucket (for example by specifying 'bucket'), "+
			"or files within a sub-directory of a bucket "+
			"(for example 'bucket/some/directory/'; note the trailing slash).",
		ModInputArgString, true, false)
	scheme.AddArgument("key_id", "Key ID", "Your Amazon key ID.",
		ModInputArgString, true, false)
	scheme.AddArgument("secret_key", "Secret key", "Your Amazon secret key.",
		ModInputArgString, true, false)

	mTestScheme := &Scheme{}
	err := xml.Unmarshal([]byte(testScheme), mTestScheme)

	if err != nil {
		t.Logf("Unable to Marshal testScheme to Scheme object: %v", err)
		t.FailNow()
	}

	bTestScheme, err := xml.Marshal(mTestScheme)
	if err != nil {
		t.Log("Unable to Unmarshal testScheme to bytes.")
		t.FailNow()
	}

	marshaled, err := xml.Marshal(scheme)
	if err != nil {
		t.Logf("Unable to marshal scheme to XML: %v", err)
		t.FailNow()
	}

	if bytes.Compare(marshaled, bTestScheme) != 0 {
		t.Logf("Marshalled scheme does not match expected scheme."+
			"\nExpected: %s\nReceived: %s",
			bTestScheme, marshaled)
		t.Fail()
	}
}
