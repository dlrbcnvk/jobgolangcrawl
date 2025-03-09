package testing

import (
	"jobgolangcrawl/database"
	"log"
	"testing"
)

func TestGetRDSSecret(t *testing.T) {
	secret, err := database.GetRDSSecret()
	if err != nil {
		t.Error(err)
	}
	log.Printf("secret is: %s", secret)
}
