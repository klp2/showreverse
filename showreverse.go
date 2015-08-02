package main

import (
	//net
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to provide an IP or CIDR block.\n")
	} else {
		for _, cidr := range os.Args[1:] {
			fmt.Println("You entered:", cidr)
			ip, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Fatal(err)
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				records, err := net.LookupAddr(ip.String())
				if err != nil {
					fmt.Println(ip, "=> NXDOMAIN, (got", err, ")")
				}
				for _, ptr := range records[0:] {
					fmt.Println(ip, "=>", ptr)
				}
			}
		}
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
