package handler

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParsePathParams(t *testing.T) {
	subtests := []struct {
		name        string
		urlPath     string
		base        string
		result      string
		errExpected bool
	}{
		{
			name:        "path valid",
			urlPath:     "/users/123",
			base:        "users",
			result:      "123",
			errExpected: false,
		},
		{
			name:        "path invalid",
			urlPath:     "us/ers/123",
			base:        "users",
			result:      "",
			errExpected: true,
		},
		{
			name:        "base not matching",
			urlPath:     "/users/123",
			base:        "items",
			result:      "",
			errExpected: true,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{Path: subtest.urlPath},
			}

			result, err := parsePathParam(req, subtest.base)

			if subtest.result != result {
				t.Errorf("expected: %v, got: %v", subtest.result, result)
			}

			if (err != nil) != subtest.errExpected {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
