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

		// load test at ever-increasing concurrency levels until we get to
		// the point where we break our latency SLA.

		latencySLABreached := false
		concurrency := 10

		for latencySLABreached == false {

			log.Printf(
				"Running load test at concurrency level %d\n",
				concurrency,
			)

			metrics := smatter.LoadTest(
				url,
				30*time.Second,
				concurrency,
			)

			log.Printf(
				"99th percentile latency %s\n",
				metrics.Latencies.P99,
			)

			if metrics.Latencies.P99 > time.Duration(config.LatencyLimitSeconds)*time.Second {
				log.Println("Latency limit breached. Finishing.")
				break
			}

			concurrency = int(float32(concurrency) * 1.2)
			time.Sleep(30 * time.Second)

		}

	} else {

		log.Printf(
			"ELB must have > %d instances to proceed with test\n",
			config.MininumAllowedInstances,
		)

	}

}
