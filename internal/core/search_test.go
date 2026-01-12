package core

import (
	"reflect"
	"testing"
)

func TestFindTickers(t *testing.T) {
	// Sample data subset
	tickers := []Ticker{
		{Code: "005930", Name: "삼성전자", Market: "KOSPI"},
		{Code: "005935", Name: "삼성전자우", Market: "KOSPI"},
		{Code: "006400", Name: "삼성SDI", Market: "KOSPI"},
		{Code: "010140", Name: "삼성중공업", Market: "KOSPI"},
		{Code: "009150", Name: "삼성전기", Market: "KOSPI"},
		{Code: "035720", Name: "카카오", Market: "KOSPI"},
	}

	tests := []struct {
		name     string
		query    string
		expected []string // Expecting Names in this order
	}{
		{
			name:     "Exact match",
			query:    "카카오",
			expected: []string{"카카오"},
		},
		{
			name:     "Partial match preference (Shortest string wins)",
			query:    "삼전",
			expected: []string{"삼성전자", "삼성전기", "삼성전자우"}, 
			// "삼성전자" (4 chars) vs "삼성전기" (4 chars) -> Tie broken by original index or alpha?
			// "삼성전자" (match indices 0, 2) score vs "삼성전기" (match indices 0, 2).
			// Depending on scoring, "삼성전자" should be top.
			// Actually "삼전" in "삼성전자" matches '삼'(0) '전'(2).
			// "삼전" in "삼성전기" matches '삼'(0) '전'(2).
			// Same length. Stable sort preserves order?
		},
		{
			name:     "No match",
			query:    "XYZ",
			expected: nil,
		},
		{
			name:     "Empty query",
			query:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FindTickers(tickers, tt.query)
			
			var resultNames []string
			for _, r := range results {
				resultNames = append(resultNames, r.Name)
			}
			
			// We only check the top N matches where N is len(expected)
			if len(resultNames) > len(tt.expected) {
				resultNames = resultNames[:len(tt.expected)]
			}

			if !reflect.DeepEqual(resultNames, tt.expected) {
				t.Errorf("FindTickers() = %v, want %v", resultNames, tt.expected)
			}
		})
	}
}
