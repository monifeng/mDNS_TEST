package config

import (
	"time"
)

type Config struct {
	CIDR        string
	Ports       string
	Services    string
	Timeout     time.Duration
	Pretty      bool
	DatasetMode bool
}

var DefaultServices = []string{
	"_workstation._tcp",
	"_http._tcp",
	"_https._tcp",
	"_smb._tcp",
	"_qdiscover._tcp",
	"_device-info._tcp",
	"_afpovertcp._tcp",
	"_ssh._tcp",
	"_ipp._tcp",
}
