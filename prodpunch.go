package main

// Prodpunch is a tool that aims to get accurate saturation metrics
// for a given service, by safely load testing real production instances
// to the point where it begins to break its latency SLA.
//
// Implementation-wise, you give prodpunch a stack/app/stage that
// identifies your (ec2-based) service, and it will detach a production
// instance from that services ELB, wait for it to drain, and then
// use the Vegeta library to load test it until it breaks a given
// latency.

import (
    "os"
    "log"
    "time"
    "bufio"

	prodpunch "github.com/MatthewJWalls/prodpunch/lib"
)

func confirmationMessage(msg string) {

    log.Printf(msg + " [y/N]")
    input, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')

    if string([]byte(input)[0]) != "y" {
        log.Fatal("Exiting")
    }

}

func handleErr(err error) {

    if err != nil {
        log.Fatal(err)
    }

}

func main() {

    config, err := prodpunch.LoadConfig("config.json")

    handleErr(err)

	instances := prodpunch.GetInstancesWithTags(
        config.Target.Stack,
        config.Target.App,
        config.Target.Stage,
    )

    if len(instances) > config.MininumAllowedInstances {

        instance := instances[0]

        log.Printf("Using instance: %s\n", instance.InstanceId)

        elb, err := prodpunch.GetLoadBalancerForInstance(
            config.Target.Stack,
            instance,
        )

        handleErr(err)

        confirmationMessage(
            "Going to detach instance from its ELB, OK?",
        )

        err = prodpunch.DetachInstanceFromELB(
            config.Target.Stack,
            elb,
            instance,
        )

        handleErr(err)

        log.Println("Waiting for connections to drain")
        time.Sleep(60 * time.Second)

        url := "http://" + instance.PublicDnsName + ":9000/_healthcheck"

        metrics := prodpunch.LoadTest(
            url,
            1 * time.Second,
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
