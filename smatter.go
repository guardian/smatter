package main

// Smatter is a tool that aims to get accurate saturation metrics
// for a given service, by safely load testing real production instances
// to the point where it begins to break its latency SLA.
//
// Implementation-wise, you give smatter a stack/app/stage that
// identifies your (ec2-based) service, and it will detach a production
// instance from that services ELB, wait for it to drain, and then
// use the Vegeta library to load test it until it breaks a given
// latency.

import (
    "os"
    "log"
    "time"
    "bufio"

    smatter "github.com/MatthewJWalls/smatter/lib"
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

    config, err := smatter.LoadConfig("config.json")

    handleErr(err)

	instances := smatter.GetInstancesWithTags(
        config.Target.Stack,
        config.Target.App,
        config.Target.Stage,
    )

    if len(instances) > config.MininumAllowedInstances {

        instance := instances[0]

        log.Printf("Using instance: %s\n", instance.InstanceId)

        elb, err := smatter.GetLoadBalancerForInstance(
            config.Target.Stack,
            instance,
        )

        handleErr(err)

        log.Printf("Using elb: %s\n", elb.Name)

        asg, err := smatter.GetAutoScalingGroupForInstance(
            config.Target.Stack,
            instance,
        )

        handleErr(err)

        log.Printf("Using asg: %s\n", asg.Name)

        confirmationMessage(
            "Going to detach instance from its ELB and ASG, OK?",
        )

        err = smatter.DetachInstanceFromELB(
            config.Target.Stack,
            elb,
            instance,
        )

        handleErr(err)

        err = smatter.DetachInstanceFromASG(
            config.Target.Stack,
            asg,
            instance,
        )

        handleErr(err)

        log.Println("Waiting for connections to drain")
        time.Sleep(time.Duration(config.SecondsToDrain) * time.Second)

        url := "http://" + instance.PublicDnsName + config.Endpoint

        metrics := smatter.LoadTest(
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
