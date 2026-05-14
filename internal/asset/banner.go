package asset

import (
	"strings"
)

func BuildBanner(txt []string, serviceType string) Banner {
	fields := make(map[string]string)
	for _, t := range txt {
		eqIdx := strings.Index(t, "=")
		if eqIdx > 0 {
			key := strings.TrimSpace(t[:eqIdx])
			value := strings.TrimSpace(t[eqIdx+1:])
			fields[key] = value
		}
	}
	fp := Fingerprint{}
	matchedFields := []string{}
	evidence := []string{}

	if v, ok := fields["accessType"]; ok && v != "" {
		fp.AccessType = v
		matchedFields = append(matchedFields, "accessType")
		evidence = append(evidence, "TXT:accessType="+v)
	}
	if v, ok := fields["accessPort"]; ok && v != "" {
		fp.AccessPort = v
		matchedFields = append(matchedFields, "accessPort")
		evidence = append(evidence, "TXT:accessPort="+v)
	}
	if v, ok := fields["model"]; ok && v != "" {
		fp.Model = v
		matchedFields = append(matchedFields, "model")
		evidence = append(evidence, "TXT:model="+v)
	}
	if v, ok := fields["displayModel"]; ok && v != "" {
		fp.DisplayModel = v
		matchedFields = append(matchedFields, "displayModel")
		evidence = append(evidence, "TXT:displayModel="+v)
	}
	if v, ok := fields["fwVer"]; ok && v != "" {
		fp.FirmwareVersion = v
		matchedFields = append(matchedFields, "fwVer")
		evidence = append(evidence, "TXT:fwVer="+v)
	}
	if v, ok := fields["fwBuildNum"]; ok && v != "" {
		fp.FirmwareBuild = v
		matchedFields = append(matchedFields, "fwBuildNum")
		evidence = append(evidence, "TXT:fwBuildNum="+v)
	}
	if v, ok := fields["path"]; ok && v != "" {
		fp.Path = v
		matchedFields = append(matchedFields, "path")
		evidence = append(evidence, "TXT:path="+v)
	}
	if strings.Contains(serviceType, "_qdiscover") {
		fp.ProtocolHint = "qnap-qdiscover"
		fp.Vendor = "QNAP"
	} else if strings.Contains(serviceType, "_smb") {
		fp.ProtocolHint = "smb"
	} else if strings.Contains(serviceType, "_afpovertcp") {
		fp.ProtocolHint = "afp"
	} else if strings.Contains(serviceType, "_http") {
		fp.ProtocolHint = "http"
	} else if strings.Contains(serviceType, "_device-info") {
		fp.ProtocolHint = "device-info"
	} else if strings.Contains(serviceType, "_workstation") {
		fp.ProtocolHint = "workstation"
	}

	depthLevel := "none"
	score := 0
	if len(txt) > 0 {
		depthLevel = "basic"
		score = 20
	}
	if len(fields) > 0 {
		depthLevel = "txt_fields"
		score = 40
	}
	if len(matchedFields) >= 2 {
		depthLevel = "fingerprinted"
		score = 60 + len(matchedFields)*8
		if score > 100 {
			score = 100
		}
	}

	summaryParts := []string{}
	if fp.ProtocolHint != "" {
		summaryParts = append(summaryParts, fp.ProtocolHint)
	}
	if fp.Model != "" {
		summaryParts = append(summaryParts, fp.Model)
	}
	if fp.Vendor != "" {
		summaryParts = append(summaryParts, fp.Vendor)
	}
	summary := strings.Join(summaryParts, " / ")
	if summary == "" {
		summary = "mDNS service"
	}

	return Banner{
		Raw:         txt,
		Fields:      fields,
		Summary:     summary,
		Fingerprint: fp,
		Depth: BannerDepth{
			Level:         depthLevel,
			Score:         score,
			MatchedFields: matchedFields,
			Evidence:      evidence,
		},
	}
}
