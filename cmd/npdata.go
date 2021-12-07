package cmd

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/internal/service"
	"github.com/kmhebb/serverExample/web"
)

type NPDataService interface {
	// Methods implemented in API and in service:
	ImportGridData(ctx cloud.Context, req service.ImportGridDataRequest) (interface{}, *cloud.Error)
	DeleteBatchOfGridData(ctx cloud.Context, req service.ProcessBatchGridDataRequest) (interface{}, *cloud.Error)
	GetListOfGridBatches(ctx cloud.Context, req service.GridBatchListRequest) (interface{}, *cloud.Error)
	GetBillingDataList(ctx cloud.Context, req service.GetBillingDataListRequest) (interface{}, *cloud.Error)
	GetInvoiceDataForDisplay(ctx cloud.Context, req service.GetInvoiceDataRequest) (interface{}, *cloud.Error)
	GetListOfInvoices(ctx cloud.Context, req service.InvoiceListRequest) (interface{}, *cloud.Error)
	PullMeterDataFromUB(ctx cloud.Context, req service.UBRequest) (interface{}, *cloud.Error)
	ListCustomers(ctx cloud.Context, req service.UBRequest) (interface{}, *cloud.Error)
	GetCustomerDetails(ctx cloud.Context, req service.InvoiceListRequest) (interface{}, *cloud.Error)
	UpdateCustomerDetails(ctx cloud.Context, req service.UpdateCustomerDetailRequest) (interface{}, *cloud.Error)

	// Private methods, mostly for operational use.
	InitializeInvoiceDataFromUtilibill(ctx cloud.Context, req service.InitializeUtilibillRequest) (interface{}, *cloud.Error)
	InitializeMeterData(ctx cloud.Context, req service.InitializeMeterDataRequest) (interface{}, *cloud.Error)

	// Not yet implemented in the service, but has API handling
	ProcessBatchGridData(ctx cloud.Context, req service.ProcessBatchGridDataRequest) (interface{}, *cloud.Error)
	GetBillingDataCSV(ctx cloud.Context, req service.GetBillingDataRequest) (interface{}, *cloud.Error)

	// Does not have an API implementation, either a package procedure or may be for some kind of automated chronjob operation.

	SyncInvoiceDataFromUB(ctx cloud.Context) (interface{}, *cloud.Error)
	ValidateUserAccess(ctx cloud.Context) *cloud.Error

	// Unimplemented at this time
	PullCustomerDataFromUB(ctx cloud.Context) (interface{}, *cloud.Error)
}

func RegisterDataServiceRoutes(srv *web.Server, svc service.NPDataService) {

	routes := map[string]web.HandlerOpts{
		// GRID ENDPOINTS
		"/data/ImportGridData": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.ImportGridDataRequest
				cr := csv.NewReader(r.Body)
				for {
					var gd service.GridData
					var err error
					gd, err = cr.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						return nil, cloud.NewError(cloud.ErrOpts{
							Kind:    cloud.ErrKindBadRequest,
							Message: "failed to read csv data",
							Cause:   err,
						}) //fmt.Errorf("decode import request: %w", err)
					}
					request.GD = append(request.GD, gd)
				}
				//fmt.Printf("data from reader: %+v", request.GD)
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.ImportGridDataRequest)
				return svc.ImportGridData(ctx, req)
			},
		},
		"/data/ListGridBatches": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.GridBatchListRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode batch list request",
						Cause:   err,
					}) //fmt.Errorf("decode batch list request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.GridBatchListRequest)
				return svc.GetListOfGridBatches(ctx, req)
			},
		},
		"/data/DeleteGridBatch": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.ProcessBatchGridDataRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode delete batch request",
						Cause:   err,
					}) //fmt.Errorf("decode delete batch request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.ProcessBatchGridDataRequest)
				return svc.DeleteBatchOfGridData(ctx, req)
			},
		},

		// CUSTOMER ENDPOINTS
		"/data/InitializeCustomers": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UBRequest
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode init ub invoice request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UBRequest)
				return svc.InitCustomersFromUB(ctx, req)
			},
		},
		"/data/UpdateAllCustomerData": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UBRequest
				// There is no request data to process, but we can add to this in the future.
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode get invoice data request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UBRequest)
				return svc.PullCustomerDataFromUB(ctx, req)
			},
		},
		"/data/ListCustomers": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UBRequest
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode init ub invoice request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UBRequest)
				return svc.ListCustomers(ctx, req)
			},
		},
		"/data/GetCustomerDetail": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.InvoiceListRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode get customer detail request",
						Cause:   err,
					}) //fmt.Errorf("decode get customer detail request failed: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.InvoiceListRequest)
				return svc.GetCustomerDetails(ctx, req)
			},
		},
		"/data/UpdateCustomerDetail": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UpdateCustomerDetailRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode update customer detail request",
						Cause:   err,
					}) //fmt.Errorf("decode update customer detail request failed: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UpdateCustomerDetailRequest)
				return svc.UpdateCustomerDetails(ctx, req)
			},
		},

		// INVOICE ENDPOINTS
		"/data/InitializeInvoices": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.InitializeUtilibillRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode init ub invoice request",
						Cause:   err,
					}) //fmt.Errorf("decode init ub invoice request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.InitializeUtilibillRequest)
				return svc.InitializeInvoiceDataFromUtilibill(ctx, req)
			},
		},
		"/data/SyncInvoiceDataFromUB": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UBRequest
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode init ub invoice request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UBRequest)
				return svc.SyncInvoiceDataFromUB(ctx, req)
			},
		},
		"/data/listCustomerInvoices": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.InvoiceListRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode invoice list request",
						Cause:   err,
					}) //fmt.Errorf("decode batch list request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.InvoiceListRequest)
				return svc.GetListOfInvoices(ctx, req)
			},
		},
		"/data/invoice": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.GetInvoiceDataRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode invoice request",
						Cause:   err,
					}) //fmt.Errorf("decode get invoice data request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.GetInvoiceDataRequest)
				return svc.GetInvoiceDataForDisplay(ctx, req)
			},
		},

		// BILLING ENDPOINTS
		"/data/ProcessGridBatchForBilling": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.ProcessBatchGridDataRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode process batch request",
						Cause:   err,
					}) //fmt.Errorf("decode process request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.ProcessBatchGridDataRequest)
				return svc.ProcessBatchGridData(ctx, req)
			},
		},
		"/data/ListBillingDataBatches": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.GetBillingDataListRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode get billing data list request",
						Cause:   err,
					}) //fmt.Errorf("decode get billing data list request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.GetBillingDataListRequest)
				return svc.GetBillingDataList(ctx, req)
			},
		},
		"/data/GetBillingDataCSV": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.GetBillingDataRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode get billing data request",
						Cause:   err,
					}) //fmt.Errorf("decode billing data csv request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.GetBillingDataRequest)
				return svc.GetBillingDataCSV(ctx, req)
			},
		},
		"/data/SyncMeterData": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.UBRequest
				// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				// 	return nil, fmt.Errorf("decode init ub invoice request: %w", err)
				// }
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.UBRequest)
				return svc.PullMeterDataFromUB(ctx, req)

			},
		},
		"/data/InitMeterData": {
			Decoder: func(ctx *cloud.Context, r *http.Request) (interface{}, *cloud.Error) {
				ctx.TokenRequired = true
				var request service.InitializeMeterDataRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					return nil, cloud.NewError(cloud.ErrOpts{
						Kind:    cloud.ErrKindBadRequest,
						Message: "failed to decode init meter data request",
						Cause:   err,
					}) //fmt.Errorf("decode init ub invoice request: %w", err)
				}
				return request, nil
			},
			Endpoint: func(ctx cloud.Context, request interface{}) (interface{}, *cloud.Error) {
				req := request.(service.InitializeMeterDataRequest)
				return svc.InitializeMeterData(ctx, req)

			},
		},
	}

	for path, opts := range routes {
		h := web.NewHandler(opts)
		h.Use(web.LoggingMiddleware)
		srv.Handle(path, h)
	}
}
