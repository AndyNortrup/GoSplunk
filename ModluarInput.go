package splunk

import (
	"encoding/xml"
	"io"
)

// ModularInputHandler is an interface that has the methods required to handle
// the call from Splunk for a Modular input
type ModularInputHandler interface {
	ReturnScheme()
	ValidateScheme() (bool, string)
	StreamEvents()
}

/*
Example Modular Input Config
<input>
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
</input>
*/

// ModInputConfig holds information passed on Stdin to the modular input when it
// is invoked.
type ModInputConfig struct {
	ServerHost    string           `xml:"server_host"`
	ServerURI     string           `xml:"server_uri"`
	SessionKey    string           `xml:"session_key"`
	CheckpointDir string           `xml:"checkpoint_dir"`
	Stanzas       []ModInputStanza `xml:"configuration>stanza"`
}

// ModInputStanza holds parameters for a specific instance of the modular input
type ModInputStanza struct {
	StanzaName string          `xml:"name,attr"`
	Params     []ModInputParam `xml:"param"`
	ParamMap   map[string]string
}

//Adds a parameter to the stanza
func (stanza *ModInputStanza) AddParameter(name, value string) {
	stanza.Params = append(stanza.Params, ModInputParam{Name: name, Value: value})

	if stanza.ParamMap == nil {
		stanza.ParamMap = make(map[string]string)
	}
	stanza.ParamMap[name] = value
}

// ModInputParam are key value pairs for the input as defined in inputs.conf
type ModInputParam struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Scheme is the Scheme returned to Splunk to describe the required inputs for a
// Mod Input
type Scheme struct {
	XMLName               xml.Name   `xml:"scheme"`
	Title                 string     `xml:"title"`
	Description           string     `xml:"description"`
	UseExternalValidation bool       `xml:"use_external_validation"`
	StreamingMode         string     `xml:"streaming_mode"`
	Args                  []Argument `xml:"args"`
}

// Argument is an individual setting for a Modular input.  Returned as part of
// the scheme
type Argument struct {
	XMLName     xml.Name `xml:"arg"`
	Name        string   `xml:"name,attr"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	DataType    string   `xml:"data_type"`
}

// ReadModInputConfig takes a reader with XML data and Decodes it into a
// ModInputConfig.  The reader is probably wrapping os.Stdin in most cases.
func ReadModInputConfig(r io.Reader) (*ModInputConfig, error) {
	decoder := xml.NewDecoder(r)
	config := &ModInputConfig{}
	err := decoder.Decode(&config)

	for stanzaIndex, stanza := range config.Stanzas {
		stanza.ParamMap = make(map[string]string)
		for _, param := range stanza.Params {
			stanza.ParamMap[param.Name] = param.Value
		}
		config.Stanzas[stanzaIndex] = stanza
	}

	return config, err
}
