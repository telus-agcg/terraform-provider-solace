/*
Terraform Private Registry

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package tfclient

import (
	"encoding/json"
)

// CreateRegistryProviderResponseDataAttributes struct for CreateRegistryProviderResponseDataAttributes
type CreateRegistryProviderResponseDataAttributes struct {
	Name         *string                                                  `json:"name,omitempty"`
	Namespace    *string                                                  `json:"namespace,omitempty"`
	CreatedAt    *string                                                  `json:"created-at,omitempty"`
	UpdatedAt    *string                                                  `json:"updated-at,omitempty"`
	RegistryName *string                                                  `json:"registry-name,omitempty"`
	Permissions  *CreateRegistryProviderResponseDataAttributesPermissions `json:"permissions,omitempty"`
}

// NewCreateRegistryProviderResponseDataAttributes instantiates a new CreateRegistryProviderResponseDataAttributes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateRegistryProviderResponseDataAttributes() *CreateRegistryProviderResponseDataAttributes {
	this := CreateRegistryProviderResponseDataAttributes{}
	return &this
}

// NewCreateRegistryProviderResponseDataAttributesWithDefaults instantiates a new CreateRegistryProviderResponseDataAttributes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateRegistryProviderResponseDataAttributesWithDefaults() *CreateRegistryProviderResponseDataAttributes {
	this := CreateRegistryProviderResponseDataAttributes{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *CreateRegistryProviderResponseDataAttributes) SetName(v string) {
	o.Name = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetNamespace() string {
	if o == nil || o.Namespace == nil {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetNamespaceOk() (*string, bool) {
	if o == nil || o.Namespace == nil {
		return nil, false
	}
	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasNamespace() bool {
	if o != nil && o.Namespace != nil {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *CreateRegistryProviderResponseDataAttributes) SetNamespace(v string) {
	o.Namespace = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetCreatedAt() string {
	if o == nil || o.CreatedAt == nil {
		var ret string
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetCreatedAtOk() (*string, bool) {
	if o == nil || o.CreatedAt == nil {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasCreatedAt() bool {
	if o != nil && o.CreatedAt != nil {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given string and assigns it to the CreatedAt field.
func (o *CreateRegistryProviderResponseDataAttributes) SetCreatedAt(v string) {
	o.CreatedAt = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetUpdatedAt() string {
	if o == nil || o.UpdatedAt == nil {
		var ret string
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetUpdatedAtOk() (*string, bool) {
	if o == nil || o.UpdatedAt == nil {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasUpdatedAt() bool {
	if o != nil && o.UpdatedAt != nil {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given string and assigns it to the UpdatedAt field.
func (o *CreateRegistryProviderResponseDataAttributes) SetUpdatedAt(v string) {
	o.UpdatedAt = &v
}

// GetRegistryName returns the RegistryName field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetRegistryName() string {
	if o == nil || o.RegistryName == nil {
		var ret string
		return ret
	}
	return *o.RegistryName
}

// GetRegistryNameOk returns a tuple with the RegistryName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetRegistryNameOk() (*string, bool) {
	if o == nil || o.RegistryName == nil {
		return nil, false
	}
	return o.RegistryName, true
}

// HasRegistryName returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasRegistryName() bool {
	if o != nil && o.RegistryName != nil {
		return true
	}

	return false
}

// SetRegistryName gets a reference to the given string and assigns it to the RegistryName field.
func (o *CreateRegistryProviderResponseDataAttributes) SetRegistryName(v string) {
	o.RegistryName = &v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *CreateRegistryProviderResponseDataAttributes) GetPermissions() CreateRegistryProviderResponseDataAttributesPermissions {
	if o == nil || o.Permissions == nil {
		var ret CreateRegistryProviderResponseDataAttributesPermissions
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateRegistryProviderResponseDataAttributes) GetPermissionsOk() (*CreateRegistryProviderResponseDataAttributesPermissions, bool) {
	if o == nil || o.Permissions == nil {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *CreateRegistryProviderResponseDataAttributes) HasPermissions() bool {
	if o != nil && o.Permissions != nil {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given CreateRegistryProviderResponseDataAttributesPermissions and assigns it to the Permissions field.
func (o *CreateRegistryProviderResponseDataAttributes) SetPermissions(v CreateRegistryProviderResponseDataAttributesPermissions) {
	o.Permissions = &v
}

func (o CreateRegistryProviderResponseDataAttributes) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Namespace != nil {
		toSerialize["namespace"] = o.Namespace
	}
	if o.CreatedAt != nil {
		toSerialize["created-at"] = o.CreatedAt
	}
	if o.UpdatedAt != nil {
		toSerialize["updated-at"] = o.UpdatedAt
	}
	if o.RegistryName != nil {
		toSerialize["registry-name"] = o.RegistryName
	}
	if o.Permissions != nil {
		toSerialize["permissions"] = o.Permissions
	}
	return json.Marshal(toSerialize)
}

type NullableCreateRegistryProviderResponseDataAttributes struct {
	value *CreateRegistryProviderResponseDataAttributes
	isSet bool
}

func (v NullableCreateRegistryProviderResponseDataAttributes) Get() *CreateRegistryProviderResponseDataAttributes {
	return v.value
}

func (v *NullableCreateRegistryProviderResponseDataAttributes) Set(val *CreateRegistryProviderResponseDataAttributes) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateRegistryProviderResponseDataAttributes) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateRegistryProviderResponseDataAttributes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateRegistryProviderResponseDataAttributes(val *CreateRegistryProviderResponseDataAttributes) *NullableCreateRegistryProviderResponseDataAttributes {
	return &NullableCreateRegistryProviderResponseDataAttributes{value: val, isSet: true}
}

func (v NullableCreateRegistryProviderResponseDataAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateRegistryProviderResponseDataAttributes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}