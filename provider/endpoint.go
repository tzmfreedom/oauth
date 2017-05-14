package provider

type Endpoint struct {
	AuthURL  string
	TokenURL string
}

var endpoints = map[string]*Endpoint{}

func registerEndpoint(p string, e *Endpoint) {
	endpoints[p] = e
}

func GetEndpoint(p string) *Endpoint {
	return endpoints[p]
}
