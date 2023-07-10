package main

import (
	"fmt"

	"github.com/bgreen/space-traders-go/st"
)

func main() {

	s := st.NewServer()

	s.Start()
	defer s.Stop()

	// Print specific fields
	agent, err := s.GetMyAgent()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Agent:")
	fmt.Printf("Name: %v\nCredits: %v\n",
		agent.Symbol,
		agent.Credits)

	ships, err := s.GetMyShips()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Ships:")
	for i, v := range ships {
		fmt.Printf("%v:\t%v\t%v\n", i, v.Symbol, v.Frame.Name)
	}

	contracts, err := s.GetContracts()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Contracts:")
	for i, v := range contracts {
		fmt.Printf("%v:\t%v\t%v\n", i, v.Id, v.Type)
	}

	systems, err := s.GetSystems()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Systems:")
	for i, v := range systems {
		fmt.Printf("%v:\t%v\t%v\n", i, v.Symbol, v.Type)
	}
}
