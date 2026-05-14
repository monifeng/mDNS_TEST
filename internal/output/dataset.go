package output

import (
	"encoding/json"
	"fmt"

	"mdnsmap/internal/asset"
)

type DatasetEntry struct {
	IP              string            `json:"ip"`
	Port            int               `json:"port"`
	Host            string            `json:"host"`
	Service         string            `json:"service"`
	BannerRaw       []string          `json:"banner_raw"`
	ExpectedDepth   string            `json:"expected_depth"`
	IdentifiedFields map[string]string `json:"identified_fields"`
	Score           int               `json:"score"`
}

type Dataset struct {
	DatasetType string        `json:"dataset_type"`
	Assets      []DatasetEntry `json:"assets"`
}

func PrintDataset(assets []asset.Asset, pretty bool) error {
	entries := make([]DatasetEntry, 0, len(assets))
	for _, a := range assets {
		entries = append(entries, DatasetEntry{
			IP:              a.IP,
			Port:            a.Port,
			Host:            a.Host,
			Service:         a.Service,
			BannerRaw:       a.Banner.Raw,
			ExpectedDepth:   a.Banner.Depth.Level,
			IdentifiedFields: a.Banner.Fields,
			Score:           a.Banner.Depth.Score,
		})
	}
	ds := Dataset{
		DatasetType: "mdns_banner_depth_validation",
		Assets:      entries,
	}
	var data []byte
	var err error
	if pretty {
		data, err = json.MarshalIndent(ds, "", "  ")
	} else {
		data, err = json.Marshal(ds)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
