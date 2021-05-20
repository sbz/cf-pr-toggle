package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ApiKey string `env:"TOKEN,required"`
	Email  string `env:"EMAIL,required"`
	Domain string `env:"DOMAIN,default=6dev.net"`
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

func newPageRuleRequest(zoneID string, api *cloudflare.API) *PageRuleRequest {
	return &PageRuleRequest{
		zoneID: zoneID,
		api:    *api,
	}
}

func main() {

	var config Config
	ctx := context.Background()

	prefixLookup := envconfig.PrefixLookuper("CLOUDFLARE_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, &config, prefixLookup); err != nil {
		log.Fatalf("Error reading config: %s\n", err)
	}

	api, err := cloudflare.New(config.ApiKey, config.Email)
	if err != nil {
		log.Fatal(err)
	}

	zoneId, err := api.ZoneIDByName(config.Domain)
	if err != nil {
		log.Fatal(err)
	}

	pageRules, err := api.ListPageRules(zoneId)
	if err != nil {
		log.Fatal(err)
	}

	request := newPageRuleRequest(zoneId, api)
	provider := PageRuleProvider{request: request}

	for _, rule := range pageRules {
		fmt.Printf("Current Rule: %+v\n", rule)
		if rule.Status == "disabled" {
			log.Printf("Found page rule disabled, will active.\n")
			provider.request.Enable(rule)
		} else {
			log.Printf("Found page rule active, will disable.\n")
			provider.request.Disable(rule)
		}
	}

    os.Exit(0)
}
