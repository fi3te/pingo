package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/fi3te/pingo/pkg/logging"
	probing "github.com/prometheus-community/pro-bing"
)

const MAX_NUMBER_OF_TARGETS = 4096

type PingResult struct {
	Target  string
	Count   int
	Timeout time.Duration
	Stats   *probing.Statistics
	Err     error
}

func PingConcurrent(targets []string, packetCount int, singleActionTimeout time.Duration, ls *logging.LogSetup) ([]PingResult, error) {
	numberOfTargets := len(targets)
	if numberOfTargets > MAX_NUMBER_OF_TARGETS {
		return nil, fmt.Errorf("number of targets to ping is greater than maximum of %d", MAX_NUMBER_OF_TARGETS)
	}

	var wg sync.WaitGroup
	wg.Add(numberOfTargets)

	results := make([]PingResult, numberOfTargets)
	for i, t := range targets {
		go func(index int, target string) {
			results[index] = Ping(target, packetCount, singleActionTimeout, ls)
			wg.Done()
		}(i, t)
	}

	wg.Wait()

	return results, nil
}

func Ping(target string, packetCount int, timeout time.Duration, ls *logging.LogSetup) (result PingResult) {
	result.Target = target
	result.Count = packetCount
	result.Timeout = timeout

	var pinger *probing.Pinger
	pinger, result.Err = probing.NewPinger(target)
	if result.Err != nil {
		return
	}
	pinger.SetPrivileged(true)
	pinger.Count = packetCount
	pinger.Timeout = timeout
	pinger.Size = 548
	result.Err = pinger.Run()
	if result.Err != nil {
		ls.Debug.Printf("Ping for target %s failed: %v", target, result.Err)
		return
	}
	result.Stats = pinger.Statistics()

	ls.Debug.Printf("Ping for target %s: Sent=%d, Received=%d", target, result.Stats.PacketsSent, result.Stats.PacketsRecv)

	return
}
