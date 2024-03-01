package pasdk

import (
	"encoding/json"
	"errors"
	"time"
)

// PASDKError is a custom error type that provides detailed information
// about a Payment Assist SDK error.
type PASDKError struct {
	// IsRequestRefusedError is true if the API received the request but refused
	// to process it. This can happen if the request was ill-formed, for example. Retrying
	// the same request again is unlikely to result in a different outcome.
	IsRequestRefusedError bool

	// IsValidationFailedError is true if the pre-request checks determined that your
	// request was invalid. This can happen if you have missing parameters, for example.
	// Retrying the same request again is guaranteed to have the same outcome.
	IsValidationFailedError bool

	// IsUnexpectedError is a catch-all for any unclassified or unexpected error.
	// This indicates that something unexpected has happened, such as a connection
	// failure or any kind of panic. You may wish to retry the request if you
	// recieve this kind of error.
	IsUnexpectedError bool

	ErrorMessage string
}

// Wrap wraps the error message in this error with another error message.
func (err *PASDKError) Wrap(errorString string) *PASDKError {
	err.ErrorMessage = errorString + err.ErrorMessage
	return err
}

// Error returns the error message.
func (err PASDKError) Error() string {
	return err.ErrorMessage
}

// GetErrorType returns the type of error as a string. You may find this helpful
// for debugging/logging purposes.
func (err PASDKError) GetErrorType() string {
	if err.IsRequestRefusedError {
		return "RequestRefusedError"
	}
	if err.IsValidationFailedError {
		return "ValidationFailedError"
	}
	if err.IsUnexpectedError {
		return "UnexpectedError"
	}

	return ""
}

type Plan struct {
	ID                 int    `json:"plan_id"`              // The ID of this plan.
	Name               string `json:"name"`                 // The name of this plan.
	Instalments        int    `json:"instalments"`          // The number of instalments in this plan.
	Deposit            bool   `json:"deposit"`              // Whether a deposit is required by this plan (first payment taken immediately).
	APR                string `json:"apr"`                  // The annual percentage interest rate of this plan.
	Frequency          string `json:"frequency"`            // How often payments are made on this plan.
	MinAmount          *int   `json:"min_amount"`           // The minimum amount allowed under this plan in pence, if any.
	MaxAmount          *int   `json:"max_amount"`           // The maximum amount allowed under this plan in pence, if any.
	CommissionRate     string `json:"commission_rate"`      // The Payment Assist commission rate charged under this plan as a percentage.
	CommissionFixedFee *int   `json:"commission_fixed_fee"` // The Payment Assist fixed commission fee charged under this plan in pence.
}

func (plan *Plan) UnmarshalJSON(data []byte) error {
	type Alias Plan

	tmp := struct {
		APR json.Number `json:"apr"`
		*Alias
	}{
		Alias: (*Alias)(plan), // Cast plan to Alias type, to unmarshal other fields normally.
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.New("couldn't unmarshal Plan: " + err.Error())
	}

	// API returns it as a number but we need it as a string since letting
	// it be a floating point could cause issues.
	plan.APR = tmp.APR.String()

	return nil
}

type Repayment struct {
	Date   time.Time `json:"date"`   // The due date of this repayment.
	Amount int       `json:"amount"` // The amount of this repayment, in pence.
}

func (repayment *Repayment) UnmarshalJSON(data []byte) error {
	type Alias Repayment

	tmp := struct {
		Date interface{} `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(repayment), // Cast repayment to Alias type, to unmarshal other fields normally.
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.New("couldn't unmarshal Repayment: " + err.Error())
	}

	// Go fails to decode the date format returned by the API by default, so we need
	// to explicitly tell it how to.
	date, err := time.Parse("2006-01-02", tmp.Date.(string))

	if err != nil {
		return errors.New("failed unmarshaling Repayment object: " + err.Error())
	}

	repayment.Date = date

	return nil
}
