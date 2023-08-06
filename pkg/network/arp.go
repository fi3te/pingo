package network

import (
	"github.com/mostlygeek/arp"
)

type ArpTableEntry struct {
	IpAddress  string
	MacAddress string
}

func FilterArpTable(ipAddresses []string) []ArpTableEntry {
	var result []ArpTableEntry
	count := len(ipAddresses)
	if count == 0 {
		return result
	}

	arpTable := arp.Table()

	result = make([]ArpTableEntry, count)
	for index, ipAddress := range ipAddresses {
		result[index] = ArpTableEntry{IpAddress: ipAddress, MacAddress: arpTable[ipAddress]}

	}

	return result
}
