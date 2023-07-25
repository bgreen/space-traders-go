package st

import (
	"fmt"
	"testing"
)

func TestCallbackOnce(t *testing.T) {
	s := NewClient()

	s.Start()
	defer s.Stop()

	// Print specific fields
	agent, _ := s.GetMyAgent()

	fmt.Println("Agent:")
	fmt.Printf("Name: %v\nCredits: %v\n",
		agent.Symbol,
		agent.Credits)
}

func TestCallbackTwice(t *testing.T) {
	s := NewClient()

	s.Start()
	defer s.Stop()

	// Print specific fields
	agent, _ := s.GetMyAgent()

	fmt.Println("Agent:")
	fmt.Printf("Name: %v\nCredits: %v\n",
		agent.Symbol,
		agent.Credits)

	// Starting System
	ships, _ := s.GetMyShips()

	fmt.Println("Ships:")
	for _, v := range ships {
		fmt.Printf("Name: %v\n", v.Symbol)
	}
}

func TestSystemsPage(t *testing.T) {
	s := NewClient()

	s.Start()
	defer s.Stop()

	// Print specific fields
	systems, err := s.GetMoreSystems(0, 10)

	if err != nil {
		t.Errorf("Couldn't retrieve page %v", 1)
	} else {
		t.Logf("Retrieved %v systems", len(systems))
	}

	// Starting System
	systems, err = s.GetMoreSystems(10, 10)

	if err != nil {
		t.Errorf("Couldn't retrieve page %v", 2)
	} else {
		t.Logf("Retrieved %v systems", len(systems))
	}

}
