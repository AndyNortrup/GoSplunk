package splunk

import (
	"fmt"
	"testing"
)

func TestBuildPath(t *testing.T) {
	result, err := buildRequestPath([]string{"storage", "passwords"}, "TA-GoogleFitness", "nobody")
	if err != nil {
		t.Fatalf("Error building namespaced path: %v", err)
	}

	if fmt.Sprintf("%v", result) != "https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords" {
		t.Fatalf("Failed to build namespaced REST request.  "+
			"Expected: https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords"+
			"Received: %v:\n", result)
	}

	result, err = buildRequestPath([]string{"services", "properties"}, "", "")
	if err != nil {
		t.Fatalf("Error building non-namespace path: %v", err)
	}

	if fmt.Sprintf("%v", result) != "https://localhost:8089/services/properties" {
		t.Fatalf("Failed to build namespaced REST request.  "+
			"Expected: https://localhost:8089/services/properties"+
			"Received: %v:\n", result)
	}
}
