// Code generated by go-swagger; DO NOT EDIT.

package bootstrap

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/RafaySystems/rcloud-base/components/common/api/def/clients/sentry/models"
)

// BootstrapDeleteRelayNetworkReader is a Reader for the BootstrapDeleteRelayNetwork structure.
type BootstrapDeleteRelayNetworkReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *BootstrapDeleteRelayNetworkReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewBootstrapDeleteRelayNetworkOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 403:
		result := NewBootstrapDeleteRelayNetworkForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewBootstrapDeleteRelayNetworkNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewBootstrapDeleteRelayNetworkInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewBootstrapDeleteRelayNetworkDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewBootstrapDeleteRelayNetworkOK creates a BootstrapDeleteRelayNetworkOK with default headers values
func NewBootstrapDeleteRelayNetworkOK() *BootstrapDeleteRelayNetworkOK {
	return &BootstrapDeleteRelayNetworkOK{}
}

/* BootstrapDeleteRelayNetworkOK describes a response with status code 200, with default header values.

A successful response.
*/
type BootstrapDeleteRelayNetworkOK struct {
	Payload models.RPCDeleteRelayNetworkResponse
}

func (o *BootstrapDeleteRelayNetworkOK) Error() string {
	return fmt.Sprintf("[DELETE /v2/sentry/relaynetwork/{metadata.name}][%d] bootstrapDeleteRelayNetworkOK  %+v", 200, o.Payload)
}
func (o *BootstrapDeleteRelayNetworkOK) GetPayload() models.RPCDeleteRelayNetworkResponse {
	return o.Payload
}

func (o *BootstrapDeleteRelayNetworkOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBootstrapDeleteRelayNetworkForbidden creates a BootstrapDeleteRelayNetworkForbidden with default headers values
func NewBootstrapDeleteRelayNetworkForbidden() *BootstrapDeleteRelayNetworkForbidden {
	return &BootstrapDeleteRelayNetworkForbidden{}
}

/* BootstrapDeleteRelayNetworkForbidden describes a response with status code 403, with default header values.

Returned when the user does not have permission to access the resource.
*/
type BootstrapDeleteRelayNetworkForbidden struct {
	Payload interface{}
}

func (o *BootstrapDeleteRelayNetworkForbidden) Error() string {
	return fmt.Sprintf("[DELETE /v2/sentry/relaynetwork/{metadata.name}][%d] bootstrapDeleteRelayNetworkForbidden  %+v", 403, o.Payload)
}
func (o *BootstrapDeleteRelayNetworkForbidden) GetPayload() interface{} {
	return o.Payload
}

func (o *BootstrapDeleteRelayNetworkForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBootstrapDeleteRelayNetworkNotFound creates a BootstrapDeleteRelayNetworkNotFound with default headers values
func NewBootstrapDeleteRelayNetworkNotFound() *BootstrapDeleteRelayNetworkNotFound {
	return &BootstrapDeleteRelayNetworkNotFound{}
}

/* BootstrapDeleteRelayNetworkNotFound describes a response with status code 404, with default header values.

Returned when the resource does not exist.
*/
type BootstrapDeleteRelayNetworkNotFound struct {
	Payload interface{}
}

func (o *BootstrapDeleteRelayNetworkNotFound) Error() string {
	return fmt.Sprintf("[DELETE /v2/sentry/relaynetwork/{metadata.name}][%d] bootstrapDeleteRelayNetworkNotFound  %+v", 404, o.Payload)
}
func (o *BootstrapDeleteRelayNetworkNotFound) GetPayload() interface{} {
	return o.Payload
}

func (o *BootstrapDeleteRelayNetworkNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBootstrapDeleteRelayNetworkInternalServerError creates a BootstrapDeleteRelayNetworkInternalServerError with default headers values
func NewBootstrapDeleteRelayNetworkInternalServerError() *BootstrapDeleteRelayNetworkInternalServerError {
	return &BootstrapDeleteRelayNetworkInternalServerError{}
}

/* BootstrapDeleteRelayNetworkInternalServerError describes a response with status code 500, with default header values.

Returned for internal server error
*/
type BootstrapDeleteRelayNetworkInternalServerError struct {
	Payload interface{}
}

func (o *BootstrapDeleteRelayNetworkInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /v2/sentry/relaynetwork/{metadata.name}][%d] bootstrapDeleteRelayNetworkInternalServerError  %+v", 500, o.Payload)
}
func (o *BootstrapDeleteRelayNetworkInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *BootstrapDeleteRelayNetworkInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewBootstrapDeleteRelayNetworkDefault creates a BootstrapDeleteRelayNetworkDefault with default headers values
func NewBootstrapDeleteRelayNetworkDefault(code int) *BootstrapDeleteRelayNetworkDefault {
	return &BootstrapDeleteRelayNetworkDefault{
		_statusCode: code,
	}
}

/* BootstrapDeleteRelayNetworkDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type BootstrapDeleteRelayNetworkDefault struct {
	_statusCode int

	Payload *models.GooglerpcStatus
}

// Code gets the status code for the bootstrap delete relay network default response
func (o *BootstrapDeleteRelayNetworkDefault) Code() int {
	return o._statusCode
}

func (o *BootstrapDeleteRelayNetworkDefault) Error() string {
	return fmt.Sprintf("[DELETE /v2/sentry/relaynetwork/{metadata.name}][%d] Bootstrap_DeleteRelayNetwork default  %+v", o._statusCode, o.Payload)
}
func (o *BootstrapDeleteRelayNetworkDefault) GetPayload() *models.GooglerpcStatus {
	return o.Payload
}

func (o *BootstrapDeleteRelayNetworkDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GooglerpcStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}