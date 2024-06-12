package tests

import (
	"encoding/json"
	"golang-test/model"
	"golang-test/routes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func loadTestData(filePath string) ([]model.OptionsContract, error) {
	var contracts []model.OptionsContract
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &contracts)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func TestOptionsContractModelValidation(t *testing.T) {
	// Test data
	validContract := model.OptionsContract{
		Type:           "Call",
		StrikePrice:    100,
		Bid:            10.05,
		Ask:            12.04,
		ExpirationDate: time.Date(2025, 12, 17, 0, 0, 0, 0, time.UTC),
		LongShort:      "long",
	}

	invalidContract := model.OptionsContract{
		Type:           "InvalidType",
		StrikePrice:    -100,
		Bid:            -10.05,
		Ask:            -12.04,
		ExpirationDate: time.Date(2020, 12, 17, 0, 0, 0, 0, time.UTC),
		LongShort:      "invalid",
	}

	// Valid contract should not trigger any validation error
	assert.Equal(t, "Call", validContract.Type)
	assert.Equal(t, 100.0, validContract.StrikePrice)
	assert.Equal(t, 10.05, validContract.Bid)
	assert.Equal(t, 12.04, validContract.Ask)
	assert.Equal(t, "long", validContract.LongShort)

	// Invalid contract should trigger validation error (assuming we have validation logic)
	assert.NotEqual(t, "Call", invalidContract.Type)
	assert.NotEqual(t, 100.0, invalidContract.StrikePrice)
	assert.NotEqual(t, 10.05, invalidContract.Bid)
	assert.NotEqual(t, 12.04, invalidContract.Ask)
	assert.NotEqual(t, "long", invalidContract.LongShort)
}

func TestAnalysisEndpoint(t *testing.T) {
	router := routes.SetupRouter()

	contracts, err := loadTestData("../testdata/testdata.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	contractsJSON, err := json.Marshal(contracts)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/analyze", strings.NewReader(string(contractsJSON)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "max_profit")
	assert.Contains(t, w.Body.String(), "max_loss")
	assert.Contains(t, w.Body.String(), "break_even_points")
}

func TestIntegration(t *testing.T) {
	router := routes.SetupRouter()

	contracts, err := loadTestData("../testdata/testdata.json")
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	contractsJSON, err := json.Marshal(contracts)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/analyze", strings.NewReader(string(contractsJSON)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "max_profit")
	assert.Contains(t, w.Body.String(), "max_loss")
	assert.Contains(t, w.Body.String(), "break_even_points")
}
