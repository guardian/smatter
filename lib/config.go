package lib

import (
	"encoding/json"
	"fmt"
)

// test config

var jsonBlob = []byte(`{
    "Target": { 
        "Stack": "frontend", 
        "App": "article", 
        "Stage": "prod"
    }
}`)

// types of json decoding

type ProdPunchTarget struct {
	Stack string
	App string
	Stage string
}
	
type ProdPunchConfig struct {
	Target ProdPunchTarget
}

// load config from json file

func LoadConfig() *ProdPunchConfig {
	
	var config *ProdPunchConfig

	err := json.Unmarshal(jsonBlob, &config)
	
	if err != nil {
		fmt.Println("error:", err)
	}

	return config
	
}
