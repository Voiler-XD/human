package quality

import (
	"strings"
	"testing"

	"github.com/barun-bash/human/internal/ir"
)

func TestGenerateSecurityTests_Basic(t *testing.T) {
	app := &ir.Application{
		Name: "TestApp",
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
		},
	}

	script, count := generateSecurityTests(app)
	if count == 0 {
		t.Fatal("expected count > 0")
	}
	if !strings.HasPrefix(script, "#!/usr/bin/env bash") {
		t.Error("missing shebang")
	}
	if !strings.Contains(script, "assert_status()") {
		t.Error("missing assert_status helper")
	}
	if !strings.Contains(script, "assert_not_status()") {
		t.Error("missing assert_not_status helper")
	}
	if !strings.Contains(script, "BASE_URL") {
		t.Error("missing BASE_URL variable")
	}
}

func TestGenerateSecurityTests_NoAPIs(t *testing.T) {
	app := &ir.Application{Name: "Empty"}

	script, count := generateSecurityTests(app)
	if count != 0 {
		t.Errorf("expected 0 count for no APIs, got %d", count)
	}
	if script != "" {
		t.Error("expected empty script for no APIs")
	}
}

func TestGenerateSecurityTests_AuthBypass(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{Name: "UpdateTask", Auth: true},
		},
	}

	script, count := generateSecurityTests(app)
	if !strings.Contains(script, "without auth token") {
		t.Error("missing 'without auth token' test")
	}
	if !strings.Contains(script, "malformed token") {
		t.Error("missing 'malformed token' test")
	}
	if !strings.Contains(script, "assert_status") {
		t.Error("missing assert_status call for auth bypass")
	}
	// 2 auth bypass tests + 1 CORS = 3 minimum
	if count < 2 {
		t.Errorf("expected at least 2 auth bypass tests, got total %d", count)
	}
}

func TestGenerateSecurityTests_NoAuth(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
			{Name: "CreateTask"},
		},
	}

	script, _ := generateSecurityTests(app)
	if strings.Contains(script, "Auth Bypass") {
		t.Error("should not have auth bypass section for public-only endpoints")
	}
}

func TestGenerateSecurityTests_SQLInjection(t *testing.T) {
	app := &ir.Application{
		Data: []*ir.DataModel{
			{Name: "Task", Fields: []*ir.DataField{
				{Name: "title", Type: "text"},
			}},
		},
		APIs: []*ir.Endpoint{
			{Name: "CreateTask", Params: []*ir.Param{{Name: "title"}}},
		},
	}

	script, count := generateSecurityTests(app)
	if !strings.Contains(script, "SQL Injection") {
		t.Error("missing SQL Injection section")
	}
	if !strings.Contains(script, "OR") {
		t.Error("missing SQL OR payload")
	}
	if !strings.Contains(script, "UNION SELECT") {
		t.Error("missing UNION SELECT payload")
	}
	if !strings.Contains(script, "DROP TABLE") {
		t.Error("missing DROP TABLE payload")
	}
	// 3 SQLi + 3 XSS + 1 CORS minimum
	if count < 3 {
		t.Errorf("expected at least 3 SQLi tests, got total %d", count)
	}
}

func TestGenerateSecurityTests_XSS(t *testing.T) {
	app := &ir.Application{
		Data: []*ir.DataModel{
			{Name: "Post", Fields: []*ir.DataField{
				{Name: "body", Type: "text"},
			}},
		},
		APIs: []*ir.Endpoint{
			{Name: "CreatePost", Params: []*ir.Param{{Name: "body"}}},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "XSS") {
		t.Error("missing XSS section")
	}
	if !strings.Contains(script, "<script>") {
		t.Error("missing script tag payload")
	}
	if !strings.Contains(script, "onerror") {
		t.Error("missing img onerror payload")
	}
	if !strings.Contains(script, "javascript:") {
		t.Error("missing javascript: payload")
	}
}

func TestGenerateSecurityTests_ValidationBypass(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{
				Name:   "CreateTask",
				Params: []*ir.Param{{Name: "title"}},
				Validation: []*ir.ValidationRule{
					{Field: "title", Rule: "not_empty"},
				},
			},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "Validation Bypass") {
		t.Error("missing Validation Bypass section")
	}
	if !strings.Contains(script, "assert_status") {
		t.Error("missing assert_status for validation bypass")
	}
	// Should expect 400
	if !strings.Contains(script, "400") {
		t.Error("validation bypass should expect HTTP 400")
	}
}

func TestGenerateSecurityTests_IDOR(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{
				Name: "GetTask",
				Steps: []*ir.Action{
					{Type: "query", Text: "find task by id"},
				},
			},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "IDOR") {
		t.Error("missing IDOR section")
	}
	if !strings.Contains(script, "99999") {
		t.Error("missing foreign resource ID")
	}
	if !strings.Contains(script, "assert_not_status") {
		t.Error("missing assert_not_status for IDOR")
	}
}

func TestGenerateSecurityTests_RateLimiting(t *testing.T) {
	// With rate limit rule
	app := &ir.Application{
		Auth: &ir.Auth{
			Rules: []*ir.Action{
				{Type: "configure", Text: "rate limit 100 requests per minute"},
			},
		},
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "Rate Limiting") {
		t.Error("missing Rate Limiting section")
	}
	if !strings.Contains(script, "429") {
		t.Error("missing 429 status check")
	}
	if !strings.Contains(script, "seq 1 120") {
		t.Error("missing 120 request loop")
	}
}

func TestGenerateSecurityTests_RateLimiting_Absent(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
		},
	}

	script, _ := generateSecurityTests(app)
	if strings.Contains(script, "Rate Limiting") {
		t.Error("rate limiting section should be absent when no rate limit rule exists")
	}
}

func TestGenerateSecurityTests_CORS(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "CORS") {
		t.Error("missing CORS section")
	}
	if !strings.Contains(script, "evil.example.com") {
		t.Error("missing evil origin")
	}
	if !strings.Contains(script, "OPTIONS") {
		t.Error("missing OPTIONS preflight request")
	}
}

func TestGenerateSecurityTests_Footer(t *testing.T) {
	app := &ir.Application{
		APIs: []*ir.Endpoint{
			{Name: "GetTasks"},
		},
	}

	script, _ := generateSecurityTests(app)
	if !strings.Contains(script, "PASS") {
		t.Error("missing PASS counter in footer")
	}
	if !strings.Contains(script, "FAIL") {
		t.Error("missing FAIL counter in footer")
	}
	if !strings.Contains(script, "exit 1") {
		t.Error("missing exit 1 for failures")
	}
}

func TestShellEscapeSingleQuote(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"hello", "hello"},
		{"it's", `it'\''s`},
		{"'", `'\''`},
		{"no quotes", "no quotes"},
		{"a'b'c", `a'\''b'\''c`},
	}

	for _, tt := range tests {
		got := shellEscapeSingleQuote(tt.input)
		if got != tt.expect {
			t.Errorf("shellEscapeSingleQuote(%q) = %q, want %q", tt.input, got, tt.expect)
		}
	}
}

func TestCurlJSONBody(t *testing.T) {
	// Single field
	body := curlJSONBody(map[string]string{"name": "test"})
	if !strings.Contains(body, `"name":"test"`) {
		t.Errorf("expected JSON with name field, got %s", body)
	}
	if !strings.HasPrefix(body, "{") || !strings.HasSuffix(body, "}") {
		t.Errorf("expected JSON object, got %s", body)
	}

	// Empty
	body = curlJSONBody(map[string]string{})
	if body != "{}" {
		t.Errorf("expected empty JSON object, got %s", body)
	}

	// Special characters
	body = curlJSONBody(map[string]string{"val": `he said "hi"`})
	if !strings.Contains(body, `\"hi\"`) {
		t.Errorf("expected escaped quotes, got %s", body)
	}
}
