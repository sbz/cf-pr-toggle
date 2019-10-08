package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	ApiKey string `envconfig:"TOKEN" required:"true"`
	Email  string `envconfig:"EMAIL" required:"true"`
	Domain string `envconfig:"DOMAIN" default:"6dev.net"`
}

type PageRule interface {
	Enable(cloudflare.PageRule)
	Disable(cloudflare.PageRule)
}

type PageRuleRequest struct {
	zoneID string
	api    cloudflare.API
}

type PageRuleProvider struct {
	request PageRule
}

func (p *PageRuleRequest) toggle(targetRule cloudflare.PageRule, ruleStatus string) {
	changedRule := cloudflare.PageRule{Status: ruleStatus}

	err := p.api.ChangePageRule(p.zoneID, targetRule.ID, changedRule)
	if err != nil {
		log.Fatal(err)
	}

}

func (p *PageRuleRequest) Enable(targetRule cloudflare.PageRule) {
	p.toggle(targetRule, "active")
}

func (p *PageRuleRequest) Disable(targetRule cloudflare.PageRule) {
	p.toggle(targetRule, "disabled")
}

func newPageRuleRequest(zoneID string, api *cloudflare.API) PageRule {
	return &PageRuleRequest{
		zoneID: zoneID,
		api:    *api,
	}
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

	zoneID, err := api.ZoneIDByName(config.Domain)
	if err != nil {
		log.Fatal(err)
	}

	pageRules, err := api.ListPageRules(zoneID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(pageRules))

	request := newPageRuleRequest(zoneID, api)
	provider := PageRuleProvider{request: request}

	for _, rule := range pageRules {
		if rule.Status == "disabled" {
			log.Printf("Found page rule disabled, will active.\n")
			provider.request.Enable(rule)
		} else {
			log.Printf("Found page rule active, will disable.\n")
			provider.request.Disable(rule)
		}

		fmt.Printf("Current Rule: %+v\n", rule)
	}
}
