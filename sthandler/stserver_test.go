package sthandler

import (
	"fmt"
	"testing"

	"github.com/bgreen/space-traders-go/stapi"
)

func TestCallbackOnce(t *testing.T) {
	s := NewServer()

	s.Start()
	defer s.Stop()

	// Print specific fields
	c := NewCallbackOnceChannel[stapi.Agent](s)
	s.GetMyAgent()
	agent := <-c

	fmt.Println("Agent:")
	fmt.Printf("Name: %v\nCredits: %v\n",
		agent.Symbol,
		agent.Credits)
}
