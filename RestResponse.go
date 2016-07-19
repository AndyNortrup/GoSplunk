package splunk

type SessionKey struct {
	SessionKey string `xml:"sessionKey"`
	Message    string `xml:"messages>msg"`
}

type RestResponse struct {
	Title   string      `xml:"title"`
	Id      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Entries []RestEntry `xml:"entry"`
}

type RestEntry struct {
	Messages string         `xml:"messages"`
	Title    string         `xml:"title"`
	ID       string         `xml:"id"`
	Updated  string         `xml:"updated"`
	Link     []string       `xml:"link"`
	Author   string         `xml:"author"`
	Contents RestDictionary `xml:"content>dict"`
}

type RestDictionary struct {
	Keys []RestKey `xml:"key"`
}

type RestKey struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"key"`
}
