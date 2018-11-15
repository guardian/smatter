# Smatter - Load Testing for Guardian Apps

Smatter is a tool that aims to get accurate saturation metrics for a given service, by safely load testing real production instances to the point where it begins to break its latency SLA.

Implementation-wise, you give smatter a stack/app/stage that identifies your (ec2-based) service, and it will detach a production instance from that services ELB, wait for it to drain, and then use the Vegeta library to load test it until it breaks a given latency.

## Usage

Smatter is a command line utility.

```smatter -config path/to/config.json```

## Example configuration

Smatter finds instances to test using the standard Guardian stack/app/stage
tags for services. It also assumes the stack name is the aws credentials profile
name to use when calling the aws api.

```
{
    "Target": {
        "Stack": "guardian stack",
        "App": "guardian app",
        "Stage": "guardian stage"
    },
    "MininumAllowedInstances": 2,
    "SecondsToDrain": 60,
    "Endpoint": ":9000/_healthcheck",
    "LatencyLimitSeconds": 5.0,
    "InitialConcurrencyLevel": 60
}
```
## Example Output

Smatter will iteratively test higher and higher concurrency levels until the target service breaches it's latency SLA.

![Image](https://i.imgur.com/utzJ4Ur.png)
