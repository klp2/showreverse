package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	//	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to provide an IP or CIDR block.\n")
	} else {
		ipPtr := make(map[string]string)
		//ipPtr := make(map[ip]string)

		for _, cidr := range os.Args[1:] {
			fmt.Println("You entered:", cidr)
			ip, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Fatal(err)
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				ipPtr[ip.String()] = ""
				//ipPtr[ip] = ""
			}
		}
		var wg sync.WaitGroup
		for ip, _ := range ipPtr {
			wg.Add(1)
			go func(ip string) {
				//go func(ip ip) {
				defer wg.Done()
				records, err := net.LookupAddr(ip)
				if err != nil {
					//words := make([]string)
					var words []string
					words = append(words, " => NXDOMAIN, got(", err.Error(), ")")
					ipPtr[ip] = strings.Join(words, "")
					//ipPtr[ip] = fmt.Printf("%s => NXDOMAIN, got (%s)"), ip, err)
				} else {
					ipPtr[ip] = strings.Join(records, ",")
					//fmt.Println(ipPtr[ip])
				}
			}(ip)
		}
		//wg.Done()
		fmt.Println("Performing all lookups, please be patient")
		wg.Wait()

		var keys []string
		for k := range ipPtr {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Println(k, "=>", ipPtr[k])
		}
		/*
			for ip, ptr := range ipPtr {
				fmt.Println(ip, "=>", ptr)
			}
		*/
		fmt.Println("The end!")
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

/*
func ptrLookup(ip string, c chan string) {
	records, err := net.LookupAddr(ip)
	if err != nil {
		//words := make([]string)
		var words []string
		words = append(words, " => NXDOMAIN, got(", err.Error(), ")")
		ipPtr[ip] = strings.Join(words, "")
		fmt.Println(ipPtr[ip])
		//ipPtr[ip] = fmt.Printf("%s => NXDOMAIN, got (%s)"), ip, err)
	}
}
*/
