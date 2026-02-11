package model

// WebsiteConfig is supposed to be same with webserver/model/website.go
type WebsiteConfig struct {
	Id           int      `json:"id"`
	ServerNames  []string `json:"server_names"`
	Ports        []string `json:"ports"`
	Upstreams    []string `json:"upstreams"`
	CertFilename string   `json:"cert_filename"`
	KeyFilename  string   `json:"key_filename"`
}

// LoadBalancerConfig is used for upstream load balancing
type LoadBalancerConfig struct {
	Enabled   bool             `json:"enabled"`
	Algorithm string           `json:"algorithm"` // round_robin, least_conn, ip_hash
	Upstreams []UpstreamConfig `json:"upstreams"`
}

// UpstreamConfig is used for upstream server configuration
type UpstreamConfig struct {
	URL         string `json:"url"`
	Weight      int    `json:"weight,omitempty"`
	MaxFails    int    `json:"max_fails,omitempty"`
	FailTimeout int    `json:"fail_timeout,omitempty"`
}
