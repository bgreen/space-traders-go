package main

import (
	"context"
	"fmt"
	"os"

	stapi "github.com/bgreen/space-traders-go-sdk"
)

func main() {

	configuration := stapi.NewConfiguration()
	apiClient := stapi.NewAPIClient(configuration)
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiTUNNT09QMiIsInZlcnNpb24iOiJ2MiIsInJlc2V0X2RhdGUiOiIyMDIzLTA2LTI0IiwiaWF0IjoxNjg4MDU1NTU4LCJzdWIiOiJhZ2VudC10b2tlbiJ9.gaALYBZxBcP6FXFBxHLmAVwLeEnNc7CPELnsWK7COaGRlm8HZpzVfZB9AzBFUBjYPeqyYyXPukS62vO-ykeAkAo4ceL1DIihq72LQAbNtQvMMtS6F0m_UwX4HCQZqmSiCMkOc7cNRohoQHdt_jLdX8N4_V5P1-5wYCweF-Kqu0IoKbdoFpjHEpDUhcSmF-tUIGCleQ4wqu3pl3NX9HIjvJ2CSgUkcGXimEqpjcFPYpBC_bc73N0Coa7J9DhWEoCpkA1qnDNIvCOVYEpZIV2OYiKoGE_7bAciBE2BznEyPYw4ogZ1F8bK52oT4ka2kN9N0aHKnQnwzMmQkjMgir19Jw"
	auth := context.WithValue(context.Background(), stapi.ContextAccessToken, token)
	resp, r, err := apiClient.AgentsApi.GetMyAgent(auth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AgentsApi.GetMyAgent``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetMyAgent`: GetMyAgent200Response
	fmt.Fprintf(os.Stdout, "Response from `AgentsApi.GetMyAgent`: %v\n", resp)
}
