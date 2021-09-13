// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewV2DownloadClusterLogsParams creates a new V2DownloadClusterLogsParams object
// with the default values initialized.
func NewV2DownloadClusterLogsParams() *V2DownloadClusterLogsParams {
	var ()
	return &V2DownloadClusterLogsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewV2DownloadClusterLogsParamsWithTimeout creates a new V2DownloadClusterLogsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewV2DownloadClusterLogsParamsWithTimeout(timeout time.Duration) *V2DownloadClusterLogsParams {
	var ()
	return &V2DownloadClusterLogsParams{

		timeout: timeout,
	}
}

// NewV2DownloadClusterLogsParamsWithContext creates a new V2DownloadClusterLogsParams object
// with the default values initialized, and the ability to set a context for a request
func NewV2DownloadClusterLogsParamsWithContext(ctx context.Context) *V2DownloadClusterLogsParams {
	var ()
	return &V2DownloadClusterLogsParams{

		Context: ctx,
	}
}

// NewV2DownloadClusterLogsParamsWithHTTPClient creates a new V2DownloadClusterLogsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewV2DownloadClusterLogsParamsWithHTTPClient(client *http.Client) *V2DownloadClusterLogsParams {
	var ()
	return &V2DownloadClusterLogsParams{
		HTTPClient: client,
	}
}

/*V2DownloadClusterLogsParams contains all the parameters to send to the API endpoint
for the v2 download cluster logs operation typically these are written to a http.Request
*/
type V2DownloadClusterLogsParams struct {

	/*ClusterID
	  The cluster whose logs should be downloaded.

	*/
	ClusterID strfmt.UUID
	/*HostID
	  A specific host in the cluster whose logs should be downloaded.

	*/
	HostID *strfmt.UUID
	/*LogsType
	  The type of logs to be downloaded.

	*/
	LogsType *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithTimeout(timeout time.Duration) *V2DownloadClusterLogsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithContext(ctx context.Context) *V2DownloadClusterLogsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithHTTPClient(client *http.Client) *V2DownloadClusterLogsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithClusterID(clusterID strfmt.UUID) *V2DownloadClusterLogsParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetClusterID(clusterID strfmt.UUID) {
	o.ClusterID = clusterID
}

// WithHostID adds the hostID to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithHostID(hostID *strfmt.UUID) *V2DownloadClusterLogsParams {
	o.SetHostID(hostID)
	return o
}

// SetHostID adds the hostId to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetHostID(hostID *strfmt.UUID) {
	o.HostID = hostID
}

// WithLogsType adds the logsType to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) WithLogsType(logsType *string) *V2DownloadClusterLogsParams {
	o.SetLogsType(logsType)
	return o
}

// SetLogsType adds the logsType to the v2 download cluster logs params
func (o *V2DownloadClusterLogsParams) SetLogsType(logsType *string) {
	o.LogsType = logsType
}

// WriteToRequest writes these params to a swagger request
func (o *V2DownloadClusterLogsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID.String()); err != nil {
		return err
	}

	if o.HostID != nil {

		// query param host_id
		var qrHostID strfmt.UUID
		if o.HostID != nil {
			qrHostID = *o.HostID
		}
		qHostID := qrHostID.String()
		if qHostID != "" {
			if err := r.SetQueryParam("host_id", qHostID); err != nil {
				return err
			}
		}

	}

	if o.LogsType != nil {

		// query param logs_type
		var qrLogsType string
		if o.LogsType != nil {
			qrLogsType = *o.LogsType
		}
		qLogsType := qrLogsType
		if qLogsType != "" {
			if err := r.SetQueryParam("logs_type", qLogsType); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
