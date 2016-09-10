package ds

import "time"

type OpenshiftClient struct {
	host string
	//authUrl string
	oapiUrl string
	kapiUrl string

	namespace   string
	username    string
	password    string
	bearerToken string
}

type Event struct {
	Type   string
	Object interface{}
}

type EventObject struct {
	Kind           string
	ApiVersion     string
	InvolvedObject involvedObject
	Reason         string
	Message        string
	Source         source
	FirstTimestamp time.Time
	LastTimestamp  time.Time
	Count          int
	Type           string
}

type involvedObject struct {
	Kind      string
	Namespace string
	Name      string
	Uid       string
}

type source struct {
	Component string
	Host      string
}

type SlackMsg struct {
	Channel    string `json:"channel"`
	Username   string `json:"username"`
	Text       string `json:"text"`
	Icon_emoji string `json:"icon_emoji"`
}

type ResourcequotasObject struct {
	Kind       string
	ApiVersion string
	Metadata   metadata
	Spec       spec
	Status     status
}

type metadata struct {
	Name              string
	Namespace         string
	SelfLink          string
	Uid               string
	ResourceVersion   string
	CreationTimestamp time.Time
}

type spec struct {
	Hard hard
}

type hard struct {
	LimitsCpu      string `json:"limits.cpu"`
	LimitsMemory   string `json:"limits.memory"`
	RequestsCpu    string `json:"requests.cpu"`
	RequestsMemory string `json:"requests.memory"`
}

type status struct {
	Hard hard
	Used used
}

type used struct {
	LimitsCpu      string `json:"limits.cpu"`
	LimitsMemory   string `json:"limits.memory"`
	RequestsCpu    string `json:"requests.cpu"`
	RequestsMemory string `json:"requests.memory"`
}
