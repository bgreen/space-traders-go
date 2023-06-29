/*
SpaceTraders API

SpaceTraders is an open-universe game and learning platform that offers a set of HTTP endpoints to control a fleet of ships and explore a multiplayer universe.  The API is documented using [OpenAPI](https://github.com/SpaceTradersAPI/api-docs). You can send your first request right here in your browser to check the status of the game server.  ```json http {   \"method\": \"GET\",   \"url\": \"https://api.spacetraders.io/v2\", } ```  Unlike a traditional game, SpaceTraders does not have a first-party client or app to play the game. Instead, you can use the API to build your own client, write a script to automate your ships, or try an app built by the community.  We have a [Discord channel](https://discord.com/invite/jh6zurdWk5) where you can share your projects, ask questions, and get help from other players.   

API version: 2.0.0
Contact: joel@spacetraders.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package stapi

import (
	"encoding/json"
)

// checks if the PurchaseCargoRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PurchaseCargoRequest{}

// PurchaseCargoRequest struct for PurchaseCargoRequest
type PurchaseCargoRequest struct {
	Symbol TradeSymbol `json:"symbol"`
	// Amounts of units to purchase.
	Units int32 `json:"units"`
}

// NewPurchaseCargoRequest instantiates a new PurchaseCargoRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPurchaseCargoRequest(symbol TradeSymbol, units int32) *PurchaseCargoRequest {
	this := PurchaseCargoRequest{}
	this.Symbol = symbol
	this.Units = units
	return &this
}

// NewPurchaseCargoRequestWithDefaults instantiates a new PurchaseCargoRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPurchaseCargoRequestWithDefaults() *PurchaseCargoRequest {
	this := PurchaseCargoRequest{}
	return &this
}

// GetSymbol returns the Symbol field value
func (o *PurchaseCargoRequest) GetSymbol() TradeSymbol {
	if o == nil {
		var ret TradeSymbol
		return ret
	}

	return o.Symbol
}

// GetSymbolOk returns a tuple with the Symbol field value
// and a boolean to check if the value has been set.
func (o *PurchaseCargoRequest) GetSymbolOk() (*TradeSymbol, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Symbol, true
}

// SetSymbol sets field value
func (o *PurchaseCargoRequest) SetSymbol(v TradeSymbol) {
	o.Symbol = v
}

// GetUnits returns the Units field value
func (o *PurchaseCargoRequest) GetUnits() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Units
}

// GetUnitsOk returns a tuple with the Units field value
// and a boolean to check if the value has been set.
func (o *PurchaseCargoRequest) GetUnitsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Units, true
}

// SetUnits sets field value
func (o *PurchaseCargoRequest) SetUnits(v int32) {
	o.Units = v
}

func (o PurchaseCargoRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PurchaseCargoRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["symbol"] = o.Symbol
	toSerialize["units"] = o.Units
	return toSerialize, nil
}

type NullablePurchaseCargoRequest struct {
	value *PurchaseCargoRequest
	isSet bool
}

func (v NullablePurchaseCargoRequest) Get() *PurchaseCargoRequest {
	return v.value
}

func (v *NullablePurchaseCargoRequest) Set(val *PurchaseCargoRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePurchaseCargoRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePurchaseCargoRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePurchaseCargoRequest(val *PurchaseCargoRequest) *NullablePurchaseCargoRequest {
	return &NullablePurchaseCargoRequest{value: val, isSet: true}
}

func (v NullablePurchaseCargoRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePurchaseCargoRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


