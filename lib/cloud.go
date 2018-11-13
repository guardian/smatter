package lib

import (
    "log"
    "errors"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/elb"
    "github.com/aws/aws-sdk-go/service/autoscaling"
)

type EC2Instance struct {
    PublicDnsName string
    InstanceId string
}

type ELB struct {
    Name string
}

type ASG struct {
    Name string
}

func isInstanceInELB(elb *elb.LoadBalancerDescription, instance EC2Instance) bool {
    for idx, _ := range elb.Instances {
        if *elb.Instances[idx].InstanceId == instance.InstanceId {
           return true
        }
    }
    return false
}

func isInstanceInASG(asg *autoscaling.Group, instance EC2Instance) bool {
    for idx, _ := range asg.Instances {
        if *asg.Instances[idx].InstanceId == instance.InstanceId {
           return true
        }
    }
    return false
}

func GetLoadBalancerForInstance(profile string, instance EC2Instance) (ELB, error) {

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("eu-west-1"),
        Credentials: credentials.NewSharedCredentials("", profile),
    })

    if err != nil {
        return ELB{}, err
    }

    svc := elb.New(sess)

    // aws provides no mechanism to get ELB by tag or instance id, so we have to
    // iterate through every single ELB and find it manually.

	params := &elb.DescribeLoadBalancersInput{
        LoadBalancerNames: []*string{ },
	}

	resp, err := svc.DescribeLoadBalancers(params)

    if err != nil {
        return ELB{}, err
    }

    // TODO: Paging
    for idx, _ := range resp.LoadBalancerDescriptions {
        elb := resp.LoadBalancerDescriptions[idx]
        if isInstanceInELB(elb, instance) {
            return ELB{*elb.LoadBalancerName}, nil
        }
    }

    return ELB{}, errors.New("Could not find ELB for instance")

}

func GetInstancesWithTags(profile string, app string, stage string) []EC2Instance {

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("eu-west-1"),
        Credentials: credentials.NewSharedCredentials("", profile),
    })

    if err != nil {
        log.Fatal(err)
    }

    svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:App"),
				Values: []*string{
					aws.String("article"),
				},
			},
            &ec2.Filter{
				Name: aws.String("tag:Stage"),
				Values: []*string{
					aws.String(stage),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)

    if err != nil {
       log.Fatal(err)
    }

    var instanceIds = []EC2Instance{}

    for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
            if *inst.State.Name == "running" {
                instanceIds = append(
                    instanceIds,
                    EC2Instance{
                        InstanceId: *inst.InstanceId,
                        PublicDnsName: *inst.PublicDnsName,
                    },
                )
            }
        }
    }

    return instanceIds

}

func DetachInstanceFromELB(profile string, loadBalancer ELB, instance EC2Instance) error {

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("eu-west-1"),
        Credentials: credentials.NewSharedCredentials("", profile),
    })

    if err != nil {
        return err
    }

    svc := elb.New(sess)

    input := &elb.DeregisterInstancesFromLoadBalancerInput{
        Instances: []*elb.Instance{
            { InstanceId: aws.String(instance.InstanceId) },
        },
        LoadBalancerName: aws.String(loadBalancer.Name),
    }

    _, detachErr := svc.DeregisterInstancesFromLoadBalancer(input)

    return detachErr

}

func DetachInstanceFromASG(profile string, scalingGroup ASG, instance EC2Instance) error {

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("eu-west-1"),
        Credentials: credentials.NewSharedCredentials("", profile),
    })

    if err != nil {
        return err
    }

    svc := autoscaling.New(sess)

    input := &autoscaling.DetachInstancesInput{
        AutoScalingGroupName: aws.String(scalingGroup.Name),
        InstanceIds: []*string{
            aws.String(instance.InstanceId),
        },
        ShouldDecrementDesiredCapacity: aws.Bool(false),
    }

    _, detachErr := svc.DetachInstances(input)

    return detachErr

}

func GetAutoScalingGroupForInstance(profile string, instance EC2Instance) (ASG, error) {

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("eu-west-1"),
        Credentials: credentials.NewSharedCredentials("", profile),
    })

    if err != nil {
        return ASG{}, err
    }

    svc := autoscaling.New(sess)

    // aws provides no mechanism to get ASG by tag or instance id, so we have to
    // iterate through every single ASG and find it manually.

	params := &autoscaling.DescribeAutoScalingGroupsInput{
        AutoScalingGroupNames: []*string{ },
	}

	resp, err := svc.DescribeAutoScalingGroups(params)

    if err != nil {
        return ASG{}, err
    }

    // TODO: Paging
    for idx, _ := range resp.AutoScalingGroups {
        asg := resp.AutoScalingGroups[idx]
        if isInstanceInASG(asg, instance) {
            return ASG{*asg.AutoScalingGroupName}, nil
        }
    }

    return ASG{}, errors.New("Could not find ASG for instance")

}
