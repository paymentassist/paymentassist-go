package pasdk

import (
	"testing"
	"time"
)

func Test_Plan(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	planID := 5

	request := PlanRequest{
		APISecret: getTestAPISecret(),
		APIKey:    getTestAPIKey(),
		Amount:    100000,
		PlanID:    &planID,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.PlanName != "4-Payment" {
		t.Error()
	}
	if response.Amount != 50000 {
		t.Error()
	}
	if response.Interest != 0 {
		t.Error()
	}
	if response.TotalRepayable != 50000 {
		t.Error()
	}

	date := time.Date(2019, 6, 12, 0, 0, 0, 0, time.UTC)

	if response.PaymentSchedule[3].Amount != 12500 {
		t.Error()
	}
	if !response.PaymentSchedule[3].Date.Equal(date) {
		t.Error()
	}
}

func Test_validatePlanRequest(t *testing.T) {
	request := PlanRequest{}

	if validatePlanRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.APIKey = "test"

	if validatePlanRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.APISecret = "test"

	if validatePlanRequest(request).Error() != "field Amount must be greater than 0" {
		t.Error()
	}

	request.Amount = 10000

	if validatePlanRequest(request) != nil {
		t.Error()
	}
}
