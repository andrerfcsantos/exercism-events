package exgo_test

import (
	"testing"

	"github.com/andrerfcsantos/exercism-events/exgo"
)

func TestGetAllMentoringRequests(t *testing.T) {

	client, err := exgo.New()
	if err != nil {
		t.Errorf("error creating client: %v", err)
	}

	res, err := client.GetAllMentoringRequests("rust")
	if err != nil {
		t.Errorf("error getting mentoring requests: %v", err)
	}

	t.Logf("mentoring requests: %d | meta: %+v", len(res.Results), res.Meta)
}
