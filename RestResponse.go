package splunk

type SplunkSessionKey struct {
	SessionKey string `xml:"sessionKey"`
	Message    string `xml:"messages>msg"`
}

type SplunkRestResponse struct {
	Title   string            `xml:"title"`
	Id      string            `xml:"id"`
	Updated string            `xml:"updated"`
	Entries []SplunkRestEntry `xml:"entry"`
}

type SplunkRestEntry struct {
	Messages string            `xml:"messages"`
	Title    string            `xml:"title"`
	ID       string            `xml:"id"`
	Updated  string            `xml:"updated"`
	Link     []string          `xml:"link"`
	Author   string            `xml:"author"`
	Contents SplunkRestContent `xml:"content"`
}
type SplunkRestContent struct {
	Dictionary SplunkRestDictionary `xml:"dict"`
}

type SplunkRestDictionary struct {
	Keys []SplunkRestKey `xml:"key"`
}

type SplunkRestKey struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
