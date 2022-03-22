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

// ProviderPlatformResponseDataAttributes struct for ProviderPlatformResponseDataAttributes
type ProviderPlatformResponseDataAttributes struct {
	Os                     *string      `json:"os,omitempty"`
	Arch                   *string      `json:"arch,omitempty"`
	Filename               *string      `json:"filename,omitempty"`
	Shasum                 *string      `json:"shasum,omitempty"`
	Permissions            *Permissions `json:"permissions,omitempty"`
	ProviderBinaryUploaded *bool        `json:"provider-binary-uploaded,omitempty"`
}

// NewProviderPlatformResponseDataAttributes instantiates a new ProviderPlatformResponseDataAttributes object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProviderPlatformResponseDataAttributes() *ProviderPlatformResponseDataAttributes {
	this := ProviderPlatformResponseDataAttributes{}
	return &this
}

// NewProviderPlatformResponseDataAttributesWithDefaults instantiates a new ProviderPlatformResponseDataAttributes object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProviderPlatformResponseDataAttributesWithDefaults() *ProviderPlatformResponseDataAttributes {
	this := ProviderPlatformResponseDataAttributes{}
	return &this
}

// GetOs returns the Os field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetOs() string {
	if o == nil || o.Os == nil {
		var ret string
		return ret
	}
	return *o.Os
}

// GetOsOk returns a tuple with the Os field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetOsOk() (*string, bool) {
	if o == nil || o.Os == nil {
		return nil, false
	}
	return o.Os, true
}

// HasOs returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasOs() bool {
	if o != nil && o.Os != nil {
		return true
	}

	return false
}

// SetOs gets a reference to the given string and assigns it to the Os field.
func (o *ProviderPlatformResponseDataAttributes) SetOs(v string) {
	o.Os = &v
}

// GetArch returns the Arch field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetArch() string {
	if o == nil || o.Arch == nil {
		var ret string
		return ret
	}
	return *o.Arch
}

// GetArchOk returns a tuple with the Arch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetArchOk() (*string, bool) {
	if o == nil || o.Arch == nil {
		return nil, false
	}
	return o.Arch, true
}

// HasArch returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasArch() bool {
	if o != nil && o.Arch != nil {
		return true
	}

	return false
}

// SetArch gets a reference to the given string and assigns it to the Arch field.
func (o *ProviderPlatformResponseDataAttributes) SetArch(v string) {
	o.Arch = &v
}

// GetFilename returns the Filename field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetFilename() string {
	if o == nil || o.Filename == nil {
		var ret string
		return ret
	}
	return *o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetFilenameOk() (*string, bool) {
	if o == nil || o.Filename == nil {
		return nil, false
	}
	return o.Filename, true
}

// HasFilename returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasFilename() bool {
	if o != nil && o.Filename != nil {
		return true
	}

	return false
}

// SetFilename gets a reference to the given string and assigns it to the Filename field.
func (o *ProviderPlatformResponseDataAttributes) SetFilename(v string) {
	o.Filename = &v
}

// GetShasum returns the Shasum field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetShasum() string {
	if o == nil || o.Shasum == nil {
		var ret string
		return ret
	}
	return *o.Shasum
}

// GetShasumOk returns a tuple with the Shasum field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetShasumOk() (*string, bool) {
	if o == nil || o.Shasum == nil {
		return nil, false
	}
	return o.Shasum, true
}

// HasShasum returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasShasum() bool {
	if o != nil && o.Shasum != nil {
		return true
	}

	return false
}

// SetShasum gets a reference to the given string and assigns it to the Shasum field.
func (o *ProviderPlatformResponseDataAttributes) SetShasum(v string) {
	o.Shasum = &v
}

// GetPermissions returns the Permissions field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetPermissions() Permissions {
	if o == nil || o.Permissions == nil {
		var ret Permissions
		return ret
	}
	return *o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetPermissionsOk() (*Permissions, bool) {
	if o == nil || o.Permissions == nil {
		return nil, false
	}
	return o.Permissions, true
}

// HasPermissions returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasPermissions() bool {
	if o != nil && o.Permissions != nil {
		return true
	}

	return false
}

// SetPermissions gets a reference to the given Permissions and assigns it to the Permissions field.
func (o *ProviderPlatformResponseDataAttributes) SetPermissions(v Permissions) {
	o.Permissions = &v
}

// GetProviderBinaryUploaded returns the ProviderBinaryUploaded field value if set, zero value otherwise.
func (o *ProviderPlatformResponseDataAttributes) GetProviderBinaryUploaded() bool {
	if o == nil || o.ProviderBinaryUploaded == nil {
		var ret bool
		return ret
	}
	return *o.ProviderBinaryUploaded
}

// GetProviderBinaryUploadedOk returns a tuple with the ProviderBinaryUploaded field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProviderPlatformResponseDataAttributes) GetProviderBinaryUploadedOk() (*bool, bool) {
	if o == nil || o.ProviderBinaryUploaded == nil {
		return nil, false
	}
	return o.ProviderBinaryUploaded, true
}

// HasProviderBinaryUploaded returns a boolean if a field has been set.
func (o *ProviderPlatformResponseDataAttributes) HasProviderBinaryUploaded() bool {
	if o != nil && o.ProviderBinaryUploaded != nil {
		return true
	}

	return false
}

// SetProviderBinaryUploaded gets a reference to the given bool and assigns it to the ProviderBinaryUploaded field.
func (o *ProviderPlatformResponseDataAttributes) SetProviderBinaryUploaded(v bool) {
	o.ProviderBinaryUploaded = &v
}

func (o ProviderPlatformResponseDataAttributes) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Os != nil {
		toSerialize["os"] = o.Os
	}
	if o.Arch != nil {
		toSerialize["arch"] = o.Arch
	}
	if o.Filename != nil {
		toSerialize["filename"] = o.Filename
	}
	if o.Shasum != nil {
		toSerialize["shasum"] = o.Shasum
	}
	if o.Permissions != nil {
		toSerialize["permissions"] = o.Permissions
	}
	if o.ProviderBinaryUploaded != nil {
		toSerialize["provider-binary-uploaded"] = o.ProviderBinaryUploaded
	}
	return json.Marshal(toSerialize)
}

type NullableProviderPlatformResponseDataAttributes struct {
	value *ProviderPlatformResponseDataAttributes
	isSet bool
}

func (v NullableProviderPlatformResponseDataAttributes) Get() *ProviderPlatformResponseDataAttributes {
	return v.value
}

func (v *NullableProviderPlatformResponseDataAttributes) Set(val *ProviderPlatformResponseDataAttributes) {
	v.value = val
	v.isSet = true
}

func (v NullableProviderPlatformResponseDataAttributes) IsSet() bool {
	return v.isSet
}

func (v *NullableProviderPlatformResponseDataAttributes) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProviderPlatformResponseDataAttributes(val *ProviderPlatformResponseDataAttributes) *NullableProviderPlatformResponseDataAttributes {
	return &NullableProviderPlatformResponseDataAttributes{value: val, isSet: true}
}

func (v NullableProviderPlatformResponseDataAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProviderPlatformResponseDataAttributes) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
