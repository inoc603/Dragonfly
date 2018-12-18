// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PieceUpdateRequest request peer uses to update its status of downloading piece in supernode.
// swagger:model PieceUpdateRequest
type PieceUpdateRequest struct {

	// contains the peer ID.
	//
	PeerID string `json:"PeerID,omitempty"`
}

// Validate validates this piece update request
func (m *PieceUpdateRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PieceUpdateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PieceUpdateRequest) UnmarshalBinary(b []byte) error {
	var res PieceUpdateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
