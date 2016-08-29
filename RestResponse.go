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

func (resp *RestResponse) string() string {
	result := "Title: " + resp.Title +
		"Id: " + resp.Id +
		"Updated: " + resp.Updated
	return result
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

//TODO: add a map[sting]string to get easier access to values.
type RestDictionary struct {
	Keys []RestKey `xml:"key"`
}

type RestKey struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func (key *RestKey) GoString() string {
	return "Name: " + key.Name + "Value: " + key.Value
}
