package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	smatter "github.com/MatthewJWalls/smatter/lib"
)

func confirmationMessage(msg string) {

	log.Printf(msg + " [y/N]")
	input, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')

	if string([]byte(input)[0]) != "y" {
		log.Fatal("Exiting")
	}

}

func main() {

	configLocation := flag.String(
		"config",
		"config.json",
		"Path to the smatter config file",
	)

	flag.Parse()

	config, err := smatter.LoadConfig(*configLocation)

	if err != nil {
		log.Fatal(err)
	}

	instances := smatter.GetInstancesWithTags(
		config.Target.Stack,
		config.Target.App,
		config.Target.Stage,
	)

	if len(instances) > config.MininumAllowedInstances {

		instance := instances[0]

		confirmationMessage(
			fmt.Sprintf(
				"Going to detach instance %s from its ELB and ASG, OK?",
				instance.InstanceId,
			),
		)

		err = smatter.DetachAndDrain(
			config.Target.Stack,
			instance,
			time.Duration(config.SecondsToDrain)*time.Second,
		)

		url := "http://" + instance.PublicDnsName + config.Endpoint

		metrics := smatter.LoadTest(
			url,
			1*time.Second,
			10,
		)

		log.Printf("99th percentile: %s\n", metrics.Latencies.P99)

	} else {

		log.Printf(
			"ELB must have > %d instances to proceed with test\n",
			config.MininumAllowedInstances,
		)

	}

}
