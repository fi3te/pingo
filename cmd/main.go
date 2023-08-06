package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/fi3te/pingo/pkg/file"
	"github.com/fi3te/pingo/pkg/logging"
	"github.com/fi3te/pingo/pkg/network"
)

func main() {
	logLevelPtr := flag.Int("logLevel", int(logging.LevelError), fmt.Sprintf("log level (debug=%d, info=%d, error=%d)", logging.LevelDebug, logging.LevelInfo, logging.LevelError))
	cidrPtr := flag.String("cidr", "", "targets in CIDR notation, e.g. 192.168.0.0/24")
	knownDevicesFilePtr := flag.String("knownDevicesFile", "", "file with known devices (line format: <mac address>;<name>)")
	flag.Parse()

	ls := logging.New(logging.LogLevel(*logLevelPtr))

	mapping := readKnownDevicesFile(*knownDevicesFilePtr, ls)
	scanResult := scanNetwork(*cidrPtr, ls)
	filteredResult := filterResults(scanResult, ls)
	printResult(filteredResult, mapping)
}

func readKnownDevicesFile(filePath string, ls *logging.LogSetup) map[string]string {
	macAddressToName := make(map[string]string)
	if filePath == "" {
		return macAddressToName
	}
	content, err := file.ReadCsvFile(filePath)
	if err != nil {
		ls.Error.Fatalf("Cannot read file '%s' with known devices: %v", filePath, err)
	}
	for index, line := range content {
		if len(line) != 2 {
			ls.Error.Fatalf("File with known devices is invalid (line: %d)", index)
		}
		macAddressToName[line[0]] = line[1]
	}
	return macAddressToName
}

func scanNetwork(cidr string, ls *logging.LogSetup) []network.PingResult {
	ips, err := network.GetIPsForCidr(cidr)
	if err != nil {
		ls.Error.Fatal(err)
	}
	numberOfIps := len(ips)
	ls.Debug.Printf("CIDR block contains %d IP addresses.", numberOfIps)

	ls.Info.Println("Starting ICMP requests...")
	results, err := network.PingConcurrent(ips, 1, time.Second, ls)
	if err != nil {
		ls.Error.Fatal(err)
	}
	return results
}

func filterResults(results []network.PingResult, ls *logging.LogSetup) []network.ArpTableEntry {
	var availableTargets []string
	numberOfErrors := 0
	for _, result := range results {
		if result.Err != nil {
			ls.Error.Println(result.Err)
			numberOfErrors++
		}
		if result.Stats != nil && result.Stats.PacketLoss < 100 {
			availableTargets = append(availableTargets, result.Target)
		}
	}
	if numberOfErrors > 0 {
		ls.Debug.Printf("%d errors have occurred.", numberOfErrors)
	}

	ls.Info.Println("Reading ARP cache...")
	return network.FilterArpTable(availableTargets)
}

func printResult(entries []network.ArpTableEntry, macAddressMapping map[string]string) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "IP\tMAC address\tDevice name\n")
	for _, entry := range entries {
		if entry.MacAddress == "" {
			entry.MacAddress = "-"
		}
		mappingValue := macAddressMapping[entry.MacAddress]
		if mappingValue == "" {
			mappingValue = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", entry.IpAddress, entry.MacAddress, mappingValue)
	}
}
