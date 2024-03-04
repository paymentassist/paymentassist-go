package pasdk

import (
	"encoding/json"
	"testing"
)

func Test_Repayment_UnmarshalJSON_UsesCorrectDateFormat(t *testing.T) {
	jsonData := `{  
		"date": "2019-07-24",
		"amount": 12500
	}`

	var repayment Repayment

	err := json.Unmarshal([]byte(jsonData), &repayment)

	if err != nil {
		t.Error()
	}

	if repayment.Date.Day() != 24 {
		t.Error()
	}
	if repayment.Date.Month() != 7 {
		t.Error()
	}
	if repayment.Date.Year() != 2019 {
		t.Error()
	}
}
