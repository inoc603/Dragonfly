// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PeerCreateRequest PeerCreateRequest is used to create a peer instance in supernode.
// Usually, when dfget is going to register in supernode as a peer,
// it will send PeerCreateRequest to supernode.
//
// swagger:model PeerCreateRequest
type PeerCreateRequest struct {

	// Peer ID of dfget client. Every peer has a unique ID among peer network.
	ID string `json:"ID,omitempty"`

	// IP address which peer client carries
	// Format: ipv4
	IP strfmt.IPv4 `json:"IP,omitempty"`

	// host name of peer client node, as a valid RFC 1123 hostname.
	// Min Length: 1
	// Format: hostname
	HostName strfmt.Hostname `json:"hostName,omitempty"`

	// when registering, dfget will setup one uploader process.
	// This one acts as a server for peer pulling tasks.
	// This port is which this server listens on.
	//
	// Maximum: 65000
	// Minimum: 15000
	Port int32 `json:"port,omitempty"`

	// version number of dfget binary.
	Version string `json:"version,omitempty"`
}

// Validate validates this peer create request
func (m *PeerCreateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateIP(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHostName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePort(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PeerCreateRequest) validateIP(formats strfmt.Registry) error {

	if swag.IsZero(m.IP) { // not required
		return nil
	}

	if err := validate.FormatOf("IP", "body", "ipv4", m.IP.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PeerCreateRequest) validateHostName(formats strfmt.Registry) error {

	if swag.IsZero(m.HostName) { // not required
		return nil
	}

	if err := validate.MinLength("hostName", "body", string(m.HostName), 1); err != nil {
		return err
	}

	if err := validate.FormatOf("hostName", "body", "hostname", m.HostName.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PeerCreateRequest) validatePort(formats strfmt.Registry) error {

	if swag.IsZero(m.Port) { // not required
		return nil
	}

	if err := validate.MinimumInt("port", "body", int64(m.Port), 15000, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("port", "body", int64(m.Port), 65000, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PeerCreateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PeerCreateRequest) UnmarshalBinary(b []byte) error {
	var res PeerCreateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
