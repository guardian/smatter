package main

import (
	"fmt"
	prodpunch "github.com/MatthewJWalls/prodpunch/lib"	
)

func main() {
	
	config := prodpunch.LoadConfig()

	prodpunch.GetInstancesWithTags("frontend", "article")
	
	//metrics := prodpunch.Run_test("http://example.com/")

	fmt.Printf("Running against: %s\n", config.Target.App)	
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	
}
