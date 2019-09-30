package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/kelseyhightower/envconfig"
	"log"
)

const (
	domain = "6dev.net"
)

type Config struct {
	ApiKey string `envconfig:"TOKEN" required:"true"`
	Email  string `envconfig:"EMAIL" required:"true"`
}

func main() {

	var config Config
	err := envconfig.Process("cloudflare", &config)
	if err != nil {
		log.Fatalf("Error reading config: %s\n", err)
	}

	api, err := cloudflare.New(config.ApiKey, config.Email)
	if err != nil {
		log.Fatal(err)
	}

	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		log.Fatal(err)
	}

	pageRules, err := api.ListPageRules(zoneID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(pageRules))

	for _, r := range pageRules {
		var newRule cloudflare.PageRule
		if r.Status == "disabled" {
			log.Printf("Found page rule disabled, will active.\n")
			newRule = cloudflare.PageRule{Status: "active"}
		} else {
			log.Printf("Found page rule active, will disable.\n")
			newRule = cloudflare.PageRule{Status: "disabled"}
		}
		err := api.ChangePageRule(zoneID, r.ID, newRule)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("NewRule: %+v\n", newRule)
		pageRule, _ := api.ListPageRules(zoneID)
		fmt.Printf("Current Rule: %v\n", pageRule[0])
	}
}
