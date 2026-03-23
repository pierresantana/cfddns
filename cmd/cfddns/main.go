package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pierresantana/cfddns/internal/cloudflare"
	"github.com/pierresantana/cfddns/internal/env"
	"github.com/pierresantana/cfddns/internal/ip"
)

func flagOrEnv(name, envKey, defaultVal, usage string) *string {
	return flag.String(name, "", fmt.Sprintf("%s (env: %s, default: %s)", usage, envKey, defaultVal))
}

func resolve(flagVal *string, envKey, defaultVal string) string {
	if *flagVal != "" {
		return *flagVal
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return defaultVal
}

func main() {
	if err := env.LoadFile(".env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	hostFlag := flagOrEnv("host", "CF_HOST", "", "DNS record name (e.g. home)")
	domainFlag := flagOrEnv("domain", "CF_DOMAIN", "", "Domain / zone name (e.g. sstec.com.br)")
	ttlFlag := flagOrEnv("ttl", "CF_TTL", "300", "TTL for the DNS record in seconds")
	tokenFlag := flagOrEnv("token", "CF_API_TOKEN", "", "Cloudflare API token")
	flag.Parse()

	host := resolve(hostFlag, "CF_HOST", "")
	domain := resolve(domainFlag, "CF_DOMAIN", "")
	ttlStr := resolve(ttlFlag, "CF_TTL", "300")
	apiToken := resolve(tokenFlag, "CF_API_TOKEN", "")

	if host == "" {
		fmt.Fprintln(os.Stderr, "error: host is required (-host or CF_HOST)")
		flag.Usage()
		os.Exit(1)
	}
	if domain == "" {
		fmt.Fprintln(os.Stderr, "error: domain is required (-domain or CF_DOMAIN)")
		flag.Usage()
		os.Exit(1)
	}
	if apiToken == "" {
		fmt.Fprintln(os.Stderr, "error: Cloudflare API token is required (-token or CF_API_TOKEN)")
		os.Exit(1)
	}

	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		log.Fatalf("invalid TTL %q: %v", ttlStr, err)
	}

	hostname := host + "." + domain

	// 1. Detect public IP
	currentIP, err := ip.GetPublicIP()
	if err != nil {
		log.Fatalf("failed to detect public IP: %v", err)
	}
	log.Printf("detected public IP: %s", currentIP)

	// 2. Resolve zone
	client := cloudflare.NewClient(apiToken)

	zoneID, err := client.GetZoneID(domain)
	if err != nil {
		log.Fatalf("failed to resolve zone %q: %v", domain, err)
	}
	log.Printf("resolved zone %q → %s", domain, zoneID)

	// 3. Check existing records
	records, err := client.ListDNSRecords(zoneID, hostname)
	if err != nil {
		log.Fatalf("failed to list DNS records: %v", err)
	}

	// 4. Create, update, or noop
	if len(records) == 0 {
		log.Printf("no existing A record for %s, creating...", hostname)
		if err := client.CreateDNSRecord(zoneID, hostname, currentIP, ttl); err != nil {
			log.Fatalf("failed to create DNS record: %v", err)
		}
		log.Printf("created A record: %s → %s (TTL %d)", hostname, currentIP, ttl)
		return
	}

	record := records[0]
	if record.Content == currentIP {
		log.Printf("already up to date: %s → %s", hostname, currentIP)
		return
	}

	log.Printf("updating %s: %s → %s", hostname, record.Content, currentIP)
	if err := client.UpdateDNSRecord(zoneID, record.ID, hostname, currentIP, ttl); err != nil {
		log.Fatalf("failed to update DNS record: %v", err)
	}
	log.Printf("updated A record: %s → %s (TTL %d)", hostname, currentIP, ttl)
}
