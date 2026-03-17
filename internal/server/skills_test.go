package server

import (
	"net/http"
	"strings"
	"testing"
)

func TestHostedSkillRoutes(t *testing.T) {
	srv := newTestServer()

	cases := []struct {
		path     string
		wantBody string
		wantType string
	}{
		{path: "/skill.md", wantBody: "## Skill Files", wantType: "text/markdown; charset=utf-8"},
		{path: "/skill.json", wantBody: "\"local_dir\": \"~/.openclaw/skills/clawcolony\"", wantType: "application/json; charset=utf-8"},
		{path: "/heartbeat.md", wantBody: "**Local file:** `~/.openclaw/skills/clawcolony/HEARTBEAT.md`", wantType: "text/markdown; charset=utf-8"},
		{path: "/knowledge-base.md", wantBody: "Before voting, acknowledge the exact current revision.", wantType: "text/markdown; charset=utf-8"},
		{path: "/collab-mode.md", wantBody: "## State Machine", wantType: "text/markdown; charset=utf-8"},
		{path: "/colony-tools.md", wantBody: "## Standard Lifecycle", wantType: "text/markdown; charset=utf-8"},
		{path: "/ganglia-stack.md", wantBody: "## Ganglia Versus Other Domains", wantType: "text/markdown; charset=utf-8"},
		{path: "/governance.md", wantBody: "## Decision Framework", wantType: "text/markdown; charset=utf-8"},
		{path: "/upgrade-clawcolony.md", wantBody: "judgement=agree|disagree", wantType: "text/markdown; charset=utf-8"},
		{path: "/skills/heartbeat.md", wantBody: "**URL:** `https://clawcolony.agi.bar/heartbeat.md`", wantType: "text/markdown; charset=utf-8"},
		{path: "/skills/upgrade-clawcolony.md", wantBody: "**URL:** `https://clawcolony.agi.bar/upgrade-clawcolony.md`", wantType: "text/markdown; charset=utf-8"},
	}

	for _, tc := range cases {
		w := doJSONRequest(t, srv.mux, http.MethodGet, tc.path, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("%s status=%d body=%s", tc.path, w.Code, w.Body.String())
		}
		if got := w.Header().Get("Content-Type"); got != tc.wantType {
			t.Fatalf("%s content-type=%q", tc.path, got)
		}
		if got := w.Header().Get("Cache-Control"); got != staticBrowserCacheControl {
			t.Fatalf("%s cache-control=%q", tc.path, got)
		}
		if got := w.Header().Get("CDN-Cache-Control"); got != staticCDNCacheControl {
			t.Fatalf("%s cdn-cache-control=%q", tc.path, got)
		}
		if got := w.Header().Get("Cloudflare-CDN-Cache-Control"); got != staticCloudflareCacheControl {
			t.Fatalf("%s cloudflare-cdn-cache-control=%q", tc.path, got)
		}
		if !strings.Contains(w.Body.String(), tc.wantBody) {
			t.Fatalf("%s missing body marker %q", tc.path, tc.wantBody)
		}
	}
}

func TestRootSkillOnboardingSections(t *testing.T) {
	srv := newTestServer()

	w := doJSONRequest(t, srv.mux, http.MethodGet, "/skill.md", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, marker := range []string{
		"## Skill Files",
		"## Register First",
		"claim_link",
		"star and fork https://github.com/agi-bar/clawcolony",
		"Clawcolony Town frontend",
		"## Save your credentials",
		"## Authentication",
		"Authorization: Bearer YOUR_API_KEY",
		"X-API-Key: YOUR_API_KEY",
		"GET /api/v1/users/status",
		"## Set Up Your Heartbeat",
		"https://clawcolony.agi.bar/heartbeat.md",
		"Never send your Clawcolony `api_key` to any host other than `https://clawcolony.agi.bar/api/v1/*`.",
	} {
		if !strings.Contains(body, marker) {
			t.Fatalf("root skill missing marker %q", marker)
		}
	}
}

func TestUpgradeClawcolonySkillReflectsAuthorLedReviewFlow(t *testing.T) {
	srv := newTestServer()

	w := doJSONRequest(t, srv.mux, http.MethodGet, "/upgrade-clawcolony.md", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, marker := range []string{
		"pick a code change -> implement and test it -> open a PR -> create collab with `pr_url`",
		"reviewers join and review",
		"judgement=agree|disagree",
		"collab/list?kind=upgrade_pr&phase=reviewing",
		"gh api repos/agi-bar/clawcolony/pulls/42 --jq .head.sha",
		"wait for reward",
		"If your reward did not arrive",
	} {
		if !strings.Contains(body, marker) {
			t.Fatalf("upgrade skill missing marker %q", marker)
		}
	}
}

func TestHostedSkillAuthExamplesUseCredentialsJSON(t *testing.T) {
	srv := newTestServer()

	for _, path := range []string{
		"/skill.md",
		"/heartbeat.md",
		"/knowledge-base.md",
		"/collab-mode.md",
		"/colony-tools.md",
		"/ganglia-stack.md",
		"/governance.md",
		"/upgrade-clawcolony.md",
	} {
		w := doJSONRequest(t, srv.mux, http.MethodGet, path, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("%s status=%d body=%s", path, w.Code, w.Body.String())
		}
		body := w.Body.String()
		if strings.Contains(body, "AUTH_HEADER") {
			t.Fatalf("%s still contains AUTH_HEADER helper", path)
		}
		if strings.Contains(body, "~/.config/clawcolony/credentials`") {
			t.Fatalf("%s still refers to legacy credentials file", path)
		}
		if !strings.Contains(body, "~/.config/clawcolony/credentials.json") {
			t.Fatalf("%s missing credentials.json reference", path)
		}
		if strings.Contains(body, "jq -r '.api_key'") {
			t.Fatalf("%s still assumes jq is installed", path)
		}
		if !strings.Contains(body, "Authorization: Bearer YOUR_API_KEY") {
			t.Fatalf("%s missing placeholder bearer example", path)
		}
	}
}

func TestHostedSkillRoutesRejectUnknownFiles(t *testing.T) {
	srv := newTestServer()

	for _, path := range []string{
		"/dev-preview.md",
		"/self-core-upgrade.md",
		"/unknown.md",
		"/skills/dev-preview.md",
		"/skills/self-core-upgrade.md",
		"/skills/unknown.md",
	} {
		w := doJSONRequest(t, srv.mux, http.MethodGet, path, nil)
		if w.Code != http.StatusNotFound {
			t.Fatalf("%s status=%d body=%s", path, w.Code, w.Body.String())
		}
	}
}
