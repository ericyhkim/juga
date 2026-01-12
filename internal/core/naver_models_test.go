package core

import (
	"encoding/json"
	"testing"
)

func TestNaverResponseUnmarshalling(t *testing.T) {
	sampleJSON := `{
		"datas": [
			{
				"itemCode": "005930",
				"stockName": "삼성전자",
				"closePrice": "140,800",
				"compareToPreviousClosePrice": "-200",
				"fluctuationsRatio": "-0.14",
				"marketStatus": "OPEN",
				"compareToPreviousPrice": {
					"code": "5",
					"text": "하락",
					"name": "FALLING"
				}
			}
		],
		"time": "20260108123842"
	}`

	var resp NaverResponse
	err := json.Unmarshal([]byte(sampleJSON), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal sample JSON: %v", err)
	}

	if len(resp.Datas) != 1 {
		t.Errorf("Expected 1 data item, got %d", len(resp.Datas))
	}

	data := resp.Datas[0]
	if data.ItemCode != "005930" {
		t.Errorf("Expected ItemCode 005930, got %s", data.ItemCode)
	}
	if data.ClosePrice != "140,800" {
		t.Errorf("Expected ClosePrice 140,800, got %s", data.ClosePrice)
	}
	if data.CompareToPreviousPrice.Name != "FALLING" {
		t.Errorf("Expected CompareToPreviousPrice Name FALLING, got %s", data.CompareToPreviousPrice.Name)
	}
}
