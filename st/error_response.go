package st

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    string `json:"data"`
}

func DecodeApiError(r *http.Response) (APIError, error) {
	apiErr := ErrorResponse{}
	localVarBody, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(localVarBody, &apiErr)

	return apiErr.Error, err
}

func (e APIError) Error() string {
	return e.Message
}

//lint:file-ignore U1000 Ignore all unused code

// General Error Codes
const cooldownConflictError = 4000
const waypointNoAccessError = 4001

// Account Error Codes
const tokenEmptyError = 4100
const tokenMissingSubjectError = 4101
const tokenInvalidSubjectError = 4102
const missingTokenRequestError = 4103
const invalidTokenRequestError = 4104
const invalidTokenSubjectError = 4105
const accountNotExistsError = 4106
const agentNotExistsError = 4107
const accountHasNoAgentError = 4108
const registerAgentExistsError = 4109

// Ship Error Codes
const navigateInTransitError = 4200
const navigateInvalidDestinationError = 4201
const navigateOutsideSystemError = 4202
const navigateInsufficientFuelError = 4203
const navigateSameDestinationError = 4204
const shipExtractInvalidWaypointError = 4205
const shipExtractPermissionError = 4206
const shipJumpNoSystemError = 4207
const shipJumpSameSystemError = 4208
const shipJumpMissingModuleError = 4210
const shipJumpNoValidWaypointError = 4211
const shipJumpMissingAntimatterError = 4212
const shipInTransitError = 4214
const shipMissingSensorArraysError = 4215
const purchaseShipCreditsError = 4216
const shipCargoExceedsLimitError = 4217
const shipCargoMissingError = 4218
const shipCargoUnitCountError = 4219
const shipSurveyVerificationError = 4220
const shipSurveyExpirationError = 4221
const shipSurveyWaypointTypeError = 4222
const shipSurveyOrbitError = 4223
const shipSurveyExhaustedError = 4224
const shipRefuelDockedError = 4225
const shipRefuelInvalidWaypointError = 4226
const shipMissingMountsError = 4227
const shipCargoFullError = 4228
const shipJumpFromGateToGateError = 4229
const waypointChartedError = 4230
const shipTransferShipNotFound = 4231
const shipTransferAgentConflict = 4232
const shipTransferSameShipConflict = 4233
const shipTransferLocationConflict = 4234
const warpInsideSystemError = 4235
const shipNotInOrbitError = 4236
const shipInvalidRefineryGoodError = 4237
const shipInvalidRefineryTypeError = 4238
const shipMissingRefineryError = 4239
const shipMissingSurveyorError = 4240

// Contract Error Codes
const acceptContractNotAuthorizedError = 4500
const acceptContractConflictError = 4501
const fulfillContractDeliveryError = 4502
const contractDeadlineError = 4503
const contractFulfilledError = 4504
const contractNotAcceptedError = 4505
const contractNotAuthorizedError = 4506
const shipDeliverTermsError = 4508
const shipDeliverFulfilledError = 4509
const shipDeliverInvalidLocationError = 4510

// Market Error Codes
const marketTradeInsufficientCreditsError = 4600
const marketTradeNoPurchaseError = 4601
const marketTradeNotSoldError = 4602
const marketNotFoundError = 4603
const marketTradeUnitLimitError = 4604
