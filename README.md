# pingo

The application pingo is a network scanner that uses ping and arp cache.

## Usage

```console
$ go build cmd\main.go

$ .\main.exe --help
Usage of .\main.exe:
  -cidr string
        targets in CIDR notation, e.g. 192.168.0.0/24
  -knownDevicesFile string
        file with known devices (line format: <mac address>;<name>) (default "devices.csv")
  -logLevel int
        log level (debug=0, info=1, error=2) (default 2)

$ .\main.exe -cidr 192.168.1.0/24
IP           MAC address
192.168.1.1  ab:12:cd:e3:fa:b4
192.168.1.23 c5:67:de:f8:a9:bc
192.168.1.45 -
192.168.1.67 de:12:3f:45:67:ab

$ .\main.exe -cidr 192.168.1.0/24 -knownDevicesFile my-devices.csv
IP           MAC address       Device name
192.168.1.1  ab:12:cd:e3:fa:b4 Router
192.168.1.23 c5:67:de:f8:a9:bc Raspberry Pi
192.168.1.45 -                 Current device
192.168.1.67 de:12:3f:45:67:ab Smartphone
```
