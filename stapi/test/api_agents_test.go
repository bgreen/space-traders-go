/*
SpaceTraders API

Testing AgentsApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package stapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	openapiclient "github.com/bgreen/space-traders-go/stapi"
)

func Test_stapi_AgentsApiService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test AgentsApiService GetMyAgent", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		resp, httpRes, err := apiClient.AgentsApi.GetMyAgent(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
