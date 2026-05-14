package asset

type Asset struct {
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	Host     string    `json:"host"`
	Service  string    `json:"service"`
	Name     string    `json:"name"`
	IPv4     []string  `json:"ipv4,omitempty"`
	IPv6     []string  `json:"ipv6,omitempty"`
	Hostname string    `json:"hostname"`
	TTL      uint32    `json:"ttl,omitempty"`
	TXT      []string  `json:"txt,omitempty"`
	Banner   Banner    `json:"banner"`
}

type Banner struct {
	Raw         []string          `json:"raw,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
	Summary     string            `json:"summary"`
	Fingerprint Fingerprint       `json:"fingerprint"`
	Depth       BannerDepth       `json:"depth"`
}

type Fingerprint struct {
	Vendor           string `json:"vendor,omitempty"`
	Product          string `json:"product,omitempty"`
	Model            string `json:"model,omitempty"`
	DisplayModel     string `json:"display_model,omitempty"`
	FirmwareVersion  string `json:"firmware_version,omitempty"`
	FirmwareBuild    string `json:"firmware_build,omitempty"`
	AccessType       string `json:"access_type,omitempty"`
	AccessPort       string `json:"access_port,omitempty"`
	Path             string `json:"path,omitempty"`
	ProtocolHint     string `json:"protocol_hint,omitempty"`
}

type BannerDepth struct {
	Level          string   `json:"level"`
	Score          int      `json:"score"`
	MatchedFields  []string `json:"matched_fields"`
	Evidence       []string `json:"evidence"`
}
