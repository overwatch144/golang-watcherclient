package watcherclient

import (
	"fmt"
	"net/url"
)

// buildQueryString builds query string from ListOptions
func buildQueryString(opts *ListOptions) string {
	if opts == nil {
		return ""
	}

	params := url.Values{}

	if opts.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", opts.Limit))
	}

	if opts.Marker != "" {
		params.Add("marker", opts.Marker)
	}

	if opts.SortKey != "" {
		params.Add("sort_key", opts.SortKey)
	}

	if opts.SortDir != "" {
		params.Add("sort_dir", opts.SortDir)
	}

	if len(params) == 0 {
		return ""
	}

	return "?" + params.Encode()
}

// StringPtr returns a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a bool value
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns a pointer to an int value
func IntPtr(i int) *int {
	return &i
}
