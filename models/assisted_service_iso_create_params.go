// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AssistedServiceIsoCreateParams assisted service iso create params
//
// swagger:model assisted-service-iso-create-params
type AssistedServiceIsoCreateParams struct {

	// Version of the OpenShift cluster.
	OpenshiftVersion string `json:"openshift_version,omitempty"`

	// The pull secret obtained from Red Hat OpenShift Cluster Manager at console.redhat.com/openshift/install/pull-secret.
	PullSecret string `json:"pull_secret,omitempty"`

	// SSH public key for debugging the installation.
	SSHPublicKey string `json:"ssh_public_key,omitempty"`
}

// Validate validates this assisted service iso create params
func (m *AssistedServiceIsoCreateParams) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this assisted service iso create params based on context it is used
func (m *AssistedServiceIsoCreateParams) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AssistedServiceIsoCreateParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AssistedServiceIsoCreateParams) UnmarshalBinary(b []byte) error {
	var res AssistedServiceIsoCreateParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
