package rituals

import (
	"os"
	"testing"
)

func TestFilterEnvironment(t *testing.T) {
	// Prepare environment slice with sensitive and non-sensitive variables
	env := []string{
		"PATH=/usr/bin",
		"HOME=/home/user",
		"RUNE_DEBUG=true",
		"RUNE_ENV=development",
		"RUNE_SENTRY_DSN=https://public@sentry.io/1",
		"GITHUB_TOKEN=abc",
		"AWS_SECRET_ACCESS_KEY=shhh",
		"MY_APP_PASSWORD=secret",
		"CUSTOM_TOKEN=tok",
		"SAFE_VAR=value",
	}

	filtered := filterEnvironment(env)

	// Helper to check presence
	contains := func(key string) bool {
		for _, kv := range filtered {
			if len(kv) >= len(key)+1 && kv[:len(key)+1] == key+"=" {
				return true
			}
		}
		return false
	}

	// Should keep safe and allowed vars
	if !contains("PATH") || !contains("HOME") || !contains("SAFE_VAR") {
		t.Fatalf("expected safe vars to be present: %v", filtered)
	}
	if !contains("RUNE_DEBUG") || !contains("RUNE_ENV") {
		t.Fatalf("expected allowed Rune vars to be present: %v", filtered)
	}

	// Should remove sensitive ones
	sensitiveKeys := []string{
		"RUNE_SENTRY_DSN", "GITHUB_TOKEN", "AWS_SECRET_ACCESS_KEY", "MY_APP_PASSWORD", "CUSTOM_TOKEN",
	}
	for _, k := range sensitiveKeys {
		if contains(k) {
			t.Fatalf("expected %s to be filtered out", k)
		}
	}
}

func TestFilterEnvironmentFromOS(t *testing.T) {
	// Set some env vars in process and ensure they are filtered
	_ = os.Setenv("RUNE_SENTRY_DSN", "dsn")
	_ = os.Setenv("NPM_TOKEN", "tok")
	_ = os.Setenv("SAFE", "1")
	defer os.Unsetenv("RUNE_SENTRY_DSN")
	defer os.Unsetenv("NPM_TOKEN")
	defer os.Unsetenv("SAFE")

	filtered := filterEnvironment(os.Environ())
	for _, disallowed := range []string{"RUNE_SENTRY_DSN", "NPM_TOKEN"} {
		for _, kv := range filtered {
			if len(kv) >= len(disallowed)+1 && kv[:len(disallowed)+1] == disallowed+"=" {
				t.Fatalf("%s should have been filtered", disallowed)
			}
		}
	}
}
