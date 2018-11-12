package lib

import (
	"encoding/json"
    "io/ioutil"
)

// types for json decoding

type ProdPunchTarget struct {
	Stack string
	App string
	Stage string
}

type ProdPunchConfig struct {
	Target ProdPunchTarget
    MininumAllowedInstances int
}

// load config from json file

func LoadConfig(path string) (*ProdPunchConfig, error) {

    jsonBlob, readErr := ioutil.ReadFile(path)

    if readErr != nil {
        return &ProdPunchConfig{}, readErr
    }

	var config *ProdPunchConfig

	jsonErr := json.Unmarshal(jsonBlob, &config)

	if jsonErr != nil {
        return &ProdPunchConfig{}, jsonErr
	}

	return config, nil

}
