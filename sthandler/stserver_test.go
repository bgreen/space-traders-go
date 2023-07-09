package sthandler

import (
	"fmt"
	"testing"
)

func TestCallbackOnce(t *testing.T) {
	s := NewServer()

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
	s := NewServer()

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
