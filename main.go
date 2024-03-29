package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	err := p.api.ChangePageRule(context.Background(), p.zoneID, targetRule.ID, changedRule)
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

func renderRule(rule cloudflare.PageRule) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%-30s", rule.ID))
	sb.WriteString("\t")
	sb.WriteString(fmt.Sprintf("%-20s", strings.TrimSpace(rule.Targets[0].Constraint.Value)))
	sb.WriteString("\t")
	sb.WriteString(fmt.Sprintf("%-20s", rule.Status))
	sb.WriteString("\t")
	sb.WriteString(fmt.Sprintf("%-15s", rule.ModifiedOn.String()))

	return sb.String()
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

	pageRules, err := api.ListPageRules(context.Background(), zoneId)
	if err != nil {
		log.Fatal(err)
	}

	request := newPageRuleRequest(zoneId, api)
	provider := PageRuleProvider{request: request}

	if len(os.Args) > 2 {
		fmt.Printf("usage: %s [<ruleId>]\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Printf("%-30s\t\t%-20s\t%-20s\t%-15s\n", "Rule Id", "URL", "Status", "Last Updated")

		for _, rule := range pageRules {
			fmt.Println(renderRule(rule))
		}

		fmt.Printf("%d existing rules.\n", len(pageRules))

		os.Exit(0)
	}

	ruleId := os.Args[1]
	rule, err := api.PageRule(context.Background(), zoneId, ruleId)

	if err != nil {
		log.Fatalf("Rule Id %q not valid\n", ruleId)
	}

	if rule.Status == "disabled" {
		log.Printf("Found page rule disabled, will active.\n")
		provider.request.Enable(rule)
	} else {
		log.Printf("Found page rule active, will disable.\n")
		provider.request.Disable(rule)
	}

	os.Exit(0)
}
