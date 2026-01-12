package ui

import (
	"strings"
	"testing"
)

func TestRenderListTable(t *testing.T) {
	items := []ListItem{
		{Key: "sam", Value: "005930"},
		{Key: "삼전", Value: "005930"},
	}

	output := RenderListTable(items)

	if !strings.Contains(output, "sam") {
		t.Errorf("Output should contain 'sam'")
	}
	if !strings.Contains(output, "삼전") {
		t.Errorf("Output should contain '삼전'")
	}
	if !strings.Contains(output, "005930") {
		t.Errorf("Output should contain '005930'")
	}
}
