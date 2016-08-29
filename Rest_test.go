package splunk

import (
	"fmt"
	"testing"
)

func TestBuildPath(t *testing.T) {
	c := &Client{
		BaseURL:   LocalSplunkMgmntURL,
		Owner:     "nobody",
		Namespace: "TA-GoogleFitness",
	}
	result, err := c.buildRequestPath([]string{"storage", "passwords"})
	if err != nil {
		t.Fatalf("Error building namespaced path: %v", err)
	}

	if fmt.Sprintf("%v", result) != "https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords" {
		t.Fatalf("Failed to build namespaced REST request.  "+
			"Expected: https://localhost:8089/servicesNS/nobody/TA-GoogleFitness/storage/passwords"+
			"Received: %v:\n", result)
	}

	c.Namespace = ""
	c.Owner = ""
	result, err = c.buildRequestPath([]string{"services", "properties"})
	if err != nil {
		t.Fatalf("Error building non-namespace path: %v", err)
	}

	if fmt.Sprintf("%v", result) != "https://localhost:8089/services/properties" {
		t.Fatalf("Failed to build namespaced REST request.  "+
			"Expected: https://localhost:8089/services/properties"+
			"Received: %v:\n", result)
	}

	c.Namespace = "fitness_for_splunk"
	c.Owner = "nobody"
	result, err = c.buildRequestPath([]string{"storage", "collections", "data", "google_tokens"})
	if err != nil {
		t.Fatalf("Error building namespace path: %v", err)
	}
	if fmt.Sprintf("%v", result) != "https://localhost:8089/servicesNS/nobody/fitness_for_splunk/storage/collections/data/google_tokens" {
		t.Fatalf("Failed to build namespaced REST request for KV Store. \n"+
			"Expected: https://localhost:8089/services/properties\n"+
			"Received: %v:\n", result)
	}
}
