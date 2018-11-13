# Smatter

Smatter is a tool that aims to get accurate saturation metrics for a given service, by safely load testing real production instances to the point where it begins to break its latency SLA.

Implementation-wise, you give smatter a stack/app/stage that identifies your (ec2-based) service, and it will detach a production instance from that services ELB, wait for it to drain, and then use the Vegeta library to load test it until it breaks a given latency.

## Usage

Smatter is a command line utility.

```smatter -config path/to/config.json```
