package lib

import (
	"fmt"
	aws "github.com/aws/aws-sdk-go/aws"
)

func GetInstancesWithTags(profile string, app string) {

	svc := ec2.New(&aws.Config{
		Credentials: aws.credentials.NewSharedCredentials("", profile),
		Region:      "eu-west-1",
	})

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:App"),
				Values: []*string{
					aws.String(app),
				},
			},
		},
	}

	res, _ := svc.DescribeInstances(params)

	for _, i := range res.Reservations[0].Instances {
		var nt string
		for _, t := range i.Tags {
			if *t.Key == "App" {
				nt = *t.Value
				fmt.Println(nt, *i.InstanceID, *i.State.Name)
			}
		}

	}

}
