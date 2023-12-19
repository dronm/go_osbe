package validRusTel

import (
	"testing"
	"os"
)

const (
	TEST_VAR_TEL_TRUE = "TEL_TRUE"
	TEST_VAR_TEL_FALSE = "TEL_FALSE"
)

func getTestVar(t *testing.T, n string) *string {
	v := os.Getenv(n)
	if v == "" {
		t.Fatalf("getTestVar() failed: %s environment variable is not set", n)
	}
	return &v
}

func TestCheck(t *testing.T) {	
	tel := *getTestVar(t, TEST_VAR_TEL_TRUE)
	if !Check(tel) {
		t.Fatalf("Tel: %s expeced to be correct, but it is not", tel)	
	}
	
	tel = *getTestVar(t, TEST_VAR_TEL_FALSE)
	if Check(tel) {
		t.Fatalf("Tel: %s expeced to be incorrect, but it is correct", tel)	
	}
	
}

