package lib

import (
	"encoding/json"
    "io/ioutil"
)

// types for json decoding

type LoadTestTarget struct {
	Stack string
	App string
	Stage string
}

type LoadTestConfig struct {
	Target LoadTestTarget
    MininumAllowedInstances int
    SecondsToDrain int
    Endpoint string
}

// load config from json file

func LoadConfig(path string) (*LoadTestConfig, error) {

    jsonBlob, readErr := ioutil.ReadFile(path)

    if readErr != nil {
        return &LoadTestConfig{}, readErr
    }

	var config *LoadTestConfig

	jsonErr := json.Unmarshal(jsonBlob, &config)

	if jsonErr != nil {
        return &LoadTestConfig{}, jsonErr
	}

	return config, nil

}
