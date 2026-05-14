package asset

import (
	"testing"
)

func TestBuildBannerQNAP(t *testing.T) {
	txt := []string{
		"accessType=https",
		"accessPort=86",
		"model=TS-X64",
		"displayModel=TS-464C",
		"fwVer=5.2.9",
		"fwBuildNum=20260214",
	}
	banner := BuildBanner(txt, "_qdiscover._tcp.local.")
	if banner.Fingerprint.Vendor != "QNAP" {
		t.Errorf("expected vendor QNAP, got %s", banner.Fingerprint.Vendor)
	}
	if banner.Fingerprint.Model != "TS-X64" {
		t.Errorf("expected model TS-X64, got %s", banner.Fingerprint.Model)
	}
	if banner.Fingerprint.DisplayModel != "TS-464C" {
		t.Errorf("expected display model TS-464C, got %s", banner.Fingerprint.DisplayModel)
	}
	if banner.Depth.Level != "fingerprinted" {
		t.Errorf("expected depth fingerprinted, got %s", banner.Depth.Level)
	}
	if banner.Depth.Score < 60 {
		t.Errorf("expected score >= 60, got %d", banner.Depth.Score)
	}
}

func TestBuildBannerHTTP(t *testing.T) {
	txt := []string{
		"path=/",
	}
	banner := BuildBanner(txt, "_http._tcp.local.")
	if banner.Fingerprint.Path != "/" {
		t.Errorf("expected path /, got %s", banner.Fingerprint.Path)
	}
	if banner.Fingerprint.ProtocolHint != "http" {
		t.Errorf("expected protocol hint http, got %s", banner.Fingerprint.ProtocolHint)
	}
}

func TestBuildBannerDeviceInfo(t *testing.T) {
	txt := []string{
		"model=Xserve",
	}
	banner := BuildBanner(txt, "_device-info._tcp.local.")
	if banner.Fingerprint.Model != "Xserve" {
		t.Errorf("expected model Xserve, got %s", banner.Fingerprint.Model)
	}
	if banner.Fingerprint.ProtocolHint != "device-info" {
		t.Errorf("expected protocol hint device-info, got %s", banner.Fingerprint.ProtocolHint)
	}
}

func TestBuildBannerEmpty(t *testing.T) {
	txt := []string{}
	banner := BuildBanner(txt, "_workstation._tcp.local.")
	if banner.Depth.Level != "none" {
		t.Errorf("expected depth none, got %s", banner.Depth.Level)
	}
	if banner.Depth.Score != 0 {
		t.Errorf("expected score 0, got %d", banner.Depth.Score)
	}
}
