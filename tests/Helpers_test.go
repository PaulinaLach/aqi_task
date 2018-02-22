package tests

import (
	"aqi/helpers"
	"testing"
)

//TestFetchEnvExisting tests if an existing key is returned correctly.
func TestFetchEnvExisting(t *testing.T) {
	if helpers.FetchEnv("ALWAYS_EXISTING_KEY") != "true" {
		t.Errorf("Env ALWAYS_EXISTING_KEY is not correctly returned")
	}
}

//TestFetchEnvEmptyPanic tests if a non existing key caused panic.
func TestFetchEnvEmptyPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FetchEnv didn't panic")
		}
	}()

	helpers.FetchEnv("NON_EXISTING")
}

//TestFetchEnvEmptyNoPanic tests if non existing key with "false" flag did not panic.
func TestFetchEnvEmptyNoPanic(t *testing.T) {
	if helpers.FetchEnv("NON_EXISTING", false) != "" {
		t.Errorf("Env NON_EXISTING didn't returned empty string")
	}
}
