package output

import (
	"encoding/json"
	"fmt"

	"mdnsmap/internal/asset"
)

func PrintJSON(assets []asset.Asset, pretty bool) error {
	var data []byte
	var err error
	if pretty {
		data, err = json.MarshalIndent(assets, "", "  ")
	} else {
		data, err = json.Marshal(assets)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
