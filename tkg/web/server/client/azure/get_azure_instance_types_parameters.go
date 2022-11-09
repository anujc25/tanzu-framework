// Code generated by go-swagger; DO NOT EDIT.

package azure

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetAzureInstanceTypesParams creates a new GetAzureInstanceTypesParams object
// with the default values initialized.
func NewGetAzureInstanceTypesParams() *GetAzureInstanceTypesParams {
	var ()
	return &GetAzureInstanceTypesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAzureInstanceTypesParamsWithTimeout creates a new GetAzureInstanceTypesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAzureInstanceTypesParamsWithTimeout(timeout time.Duration) *GetAzureInstanceTypesParams {
	var ()
	return &GetAzureInstanceTypesParams{

		timeout: timeout,
	}
}

// NewGetAzureInstanceTypesParamsWithContext creates a new GetAzureInstanceTypesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAzureInstanceTypesParamsWithContext(ctx context.Context) *GetAzureInstanceTypesParams {
	var ()
	return &GetAzureInstanceTypesParams{

		Context: ctx,
	}
}

// NewGetAzureInstanceTypesParamsWithHTTPClient creates a new GetAzureInstanceTypesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAzureInstanceTypesParamsWithHTTPClient(client *http.Client) *GetAzureInstanceTypesParams {
	var ()
	return &GetAzureInstanceTypesParams{
		HTTPClient: client,
	}
}

/*
GetAzureInstanceTypesParams contains all the parameters to send to the API endpoint
for the get azure instance types operation typically these are written to a http.Request
*/
type GetAzureInstanceTypesParams struct {

	/*Location
	  Azure region name

	*/
	Location string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get azure instance types params
func (o *GetAzureInstanceTypesParams) WithTimeout(timeout time.Duration) *GetAzureInstanceTypesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get azure instance types params
func (o *GetAzureInstanceTypesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get azure instance types params
func (o *GetAzureInstanceTypesParams) WithContext(ctx context.Context) *GetAzureInstanceTypesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get azure instance types params
func (o *GetAzureInstanceTypesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get azure instance types params
func (o *GetAzureInstanceTypesParams) WithHTTPClient(client *http.Client) *GetAzureInstanceTypesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get azure instance types params
func (o *GetAzureInstanceTypesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithLocation adds the location to the get azure instance types params
func (o *GetAzureInstanceTypesParams) WithLocation(location string) *GetAzureInstanceTypesParams {
	o.SetLocation(location)
	return o
}

// SetLocation adds the location to the get azure instance types params
func (o *GetAzureInstanceTypesParams) SetLocation(location string) {
	o.Location = location
}

// WriteToRequest writes these params to a swagger request
func (o *GetAzureInstanceTypesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param location
	if err := r.SetPathParam("location", o.Location); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}