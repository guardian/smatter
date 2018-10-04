package lib

import (
	"time"
	
	vegeta "github.com/tsenart/vegeta/lib"
	
)

// uses the vegeta library to run a load test against a target

func Run_test(target string) vegeta.Metrics {

	rate := vegeta.Rate{Freq: 10, Per: time.Second}
	
	duration := 1 * time.Second
	
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    target,
	})
	
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}

	metrics.Close()

	return metrics

}
