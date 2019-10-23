package main

import (
	"testing"

	"github.com/cloudflare/cloudflare-go"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Toggle(t *testing.T) {
	Convey("Test toggle", t, func() {
		var rule cloudflare.PageRule
		provider := PageRuleProvider{}
		mock := &PageRuleMock{}
		provider.request = mock

		Convey("Disable()", func() {
			provider.request.Disable(rule)
			So(mock.toggleWasCalled, ShouldBeTrue)
			So(mock.ruleStatus, ShouldEqual, "disabled")
		})

		Convey("Enable()", func() {
			provider.request.Enable(rule)
			So(mock.toggleWasCalled, ShouldBeTrue)
			So(mock.ruleStatus, ShouldEqual, "active")
		})
	})
}

func Test_NewPageRuleRequest(t *testing.T) {
	Convey("newPageRuleRequest()", t, func() {
		zoneId := "42"
		mockApi := &cloudflare.API{}
		req := newPageRuleRequest(zoneId, mockApi)
		So(zoneId, ShouldEqual, req.zoneID)
		So(req, ShouldNotBeNil)
	})
}

type PageRuleMock struct {
	toggleWasCalled bool
	ruleStatus      string
}

func (p *PageRuleMock) Enable(rule cloudflare.PageRule) {
	p.toggleWasCalled = true
	p.ruleStatus = "active"
}

func (p *PageRuleMock) Disable(rule cloudflare.PageRule) {
	p.toggleWasCalled = true
	p.ruleStatus = "disabled"
}
