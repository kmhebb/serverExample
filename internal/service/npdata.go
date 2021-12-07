package service

import (
	"fmt"
	"strconv"
	"time"

	cloud "github.com/kmhebb/serverExample"

	utilibill "github.com/kmhebb/serverExample/API/Utilibill"
	"github.com/kmhebb/serverExample/internal/db"
	"github.com/kmhebb/serverExample/lib/email"
	"github.com/kmhebb/serverExample/log"
	"github.com/kmhebb/serverExample/pg"
	uuid "github.com/pborman/uuid"
)

type NPDataService struct {
	DB Database
	L  log.Logger
	Em email.Service
}

type GridData []string

type ImportGridDataRequest struct {
	GD []GridData
}

type ImportGridDataResponse struct {
	BatchID uuid.UUID
}

type GridBatchListRequest struct {
	ListType string
}

type GridBatchListResponse struct {
	BatchList []cloud.BatchData
}

type InitializeUtilibillRequest struct {
	Init string `json:"init"`
}

type InvoiceListRequest struct {
	CustomerID int `json:"customerid"`
}

type InvoiceListResponse struct {
	Invoices []cloud.StatementSummary
}

type ProcessBatchGridDataRequest struct {
	GridDataID string `json:"id"`
}

type ProcessBatchGridDataResponse struct {
}

type UBRequest struct {
}

type GetBillingDataListRequest struct {
	ListType string
}

type GetBillingDataListResponse struct {
	BD cloud.BillingBatchList
}

type GetBillingDataRequest struct {
	BillingDataID string
}

type GetBillingDataCSVResponse struct {
	BD []cloud.BillingData
}

type GetInvoiceDataRequest struct {
	CustomerID      int
	StatementNumber int
}

type GetInvoiceDataResponse struct {
	I cloud.StatementData
}

type UpdateCustomerDetailRequest struct {
	Data cloud.NPCustomerDetail `json:"data"`
}

type InitializeMeterDataRequest struct {
	Data []cloud.CustomerMeterData `json:"data"`
}

// This method will import data from a csv into the database with a unique identifier.
func (svc NPDataService) ImportGridData(ctx cloud.Context, req ImportGridDataRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	// Here we are receiving a .csv of grid data that we need to save into the data base.
	// At this point, the csv has been loaded into a slice contained in the req parameter, lets make sure it has values:
	if len(req.GD) == 0 {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "import csv did not contain data",
			Cause:   err,
		}) //fmt.Errorf("import csv did not contain data")
	}

	// We are going to load the data into a golang data structure for handling.
	var GridData []cloud.GridDataRecord
	for i, col := range req.GD {
		if i == 0 {
			continue
		}
		var gd cloud.GridDataRecord
		gd.HostAcct, err = strconv.Atoi(col[0])
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: HostAcct", log.Fields{"error": err})
		}
		gd.SatAcct, err = strconv.Atoi(col[1])
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: SatAcct", log.Fields{"error": err})
		}
		gd.SatelliteName = col[2]
		gd.SatServClass = col[3]
		gd.SatVDL, err = strconv.Atoi(col[4])
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: SatVDL", log.Fields{"error": err})
		}
		gd.SatStatus = col[5]
		gd.VderEnergy, err = strconv.ParseFloat(col[6], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: VderEnergy", log.Fields{"error": err})
		}
		gd.VderCap, err = strconv.ParseFloat(col[7], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: VderCap", log.Fields{"error": err})
		}
		gd.VderEnv, err = strconv.ParseFloat(col[8], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: VderEnv", log.Fields{"error": err})
		}
		gd.VderDrv, err = strconv.ParseFloat(col[9], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: VderDrv", log.Fields{"error": err})
		}
		gd.VderLsrv, err = strconv.ParseFloat(col[10], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: VderLsrv", log.Fields{"error": err})
		}
		gd.VderMTC, err = strconv.ParseFloat(col[11], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: vdermtc", log.Fields{"error": err})
		}
		gd.VderCc, err = strconv.ParseFloat(col[12], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: vdercc", log.Fields{"error": err})
		}
		gd.VderTotal, err = strconv.ParseFloat(col[13], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: vdertotal", log.Fields{"error": err})
		}
		gd.TransKWH, err = strconv.ParseFloat(col[14], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: transkwh", log.Fields{"error": err})
		}
		gd.Allocation, err = strconv.ParseFloat(col[15], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: allocation", log.Fields{"error": err})
		}
		gd.HostBillPeriod = col[16]
		gd.TransferDate, err = time.Parse("2006-01-02", col[17])
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: transferdate", log.Fields{"error": err})
		}
		gd.SatBillDate, err = time.Parse("2006-01-02", col[18])
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: satbilldate", log.Fields{"error": err})
		}
		gd.BankedPriorMonth, err = strconv.ParseFloat(col[19], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: bankedpriormonth", log.Fields{"error": err})
		}
		gd.CurrentVDER, err = strconv.ParseFloat(col[20], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: currentvder", log.Fields{"error": err})
		}
		gd.TotalAvailable, err = strconv.ParseFloat(col[21], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: totalavailable", log.Fields{"error": err})
		}
		gd.SatBillAmt, err = strconv.ParseFloat(col[22], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: satbillamount", log.Fields{"error": err})
		}
		gd.Applied, err = strconv.ParseFloat(col[23], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: applied", log.Fields{"error": err})
		}
		gd.BankedCarryOver, err = strconv.ParseFloat(col[24], 64)
		if err != nil {
			svc.L.Debug(ctx.Ctx, "decode csv: bankedcarryover", log.Fields{"error": err})
		}
		GridData = append(GridData, gd)
	}

	// Now we are going to pass this over to the database.
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.ImportGridData(ctx, tx, GridData)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice import transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/DataService.ImportGridData.RunInTransaction failed: %w", err)
	}

	return nil, nil
}

// This method will delete a batch of data that was uploaded via csv.
func (svc NPDataService) DeleteBatchOfGridData(ctx cloud.Context, req ProcessBatchGridDataRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	if req.GridDataID == "" {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindBadRequest,
			Message: "batch id required",
			Cause:   err,
		}) //fmt.Errorf("batch id required")
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.DeleteGridData(ctx, tx, req.GridDataID)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice delete batch transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/Dataservice.DeleteBatchofGridData.RunInTransaction failed: %w", err)
	}

	return nil, nil
}

// This method will list the batches of data that are in the database.
func (svc NPDataService) GetListOfGridBatches(ctx cloud.Context, req GridBatchListRequest) (GridBatchListResponse, *cloud.Error) {
	var err error
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err = svc.ValidateUserAccess(ctx)
	if err != nil {
		return GridBatchListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var resp GridBatchListResponse
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		resp.BatchList, err = db.GetGridBatchList(ctx, tx, req.ListType)
		if err != nil {
			return cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindInternal,
				Message: "db action failed",
				Cause:   err,
			}) //fmt.Errorf("service/db.GetGridBatchList failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return GridBatchListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get batch list transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/Dataservice.DeleteBatchofGridData.RunInTransaction failed: %w", err)
	}

	return resp, nil
}

// This method will take the grid data that has been uploaded and turn it into billing information.
func (svc NPDataService) ProcessBatchGridData(ctx cloud.Context, req ProcessBatchGridDataRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}
	return nil, nil
}

// This method will list batches of billing information in the database.
func (svc NPDataService) GetBillingDataList(ctx cloud.Context, req GetBillingDataListRequest) (interface{}, *cloud.Error) {
	var err error
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err = svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}
	var resp GetBillingDataListResponse

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		resp.BD.BillingList, err = db.GetBillingDataList(ctx, tx, req.ListType)
		if err != nil {
			return cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindInternal,
				Message: "db action failed",
				Cause:   err,
			}) //fmt.Errorf("service/db.GetGridBatchList failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return GridBatchListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get billing data list transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/Dataservice.DeleteBatchofGridData.RunInTransaction failed: %w", err)
	}

	return resp, nil
}

// This method will return a set of billing data as csv.
func (svc NPDataService) GetBillingDataCSV(ctx cloud.Context, req GetBillingDataRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}
	return nil, nil
}

// This method will pull all statement data from utilibill and synchronize the statement data in the database.
func (svc NPDataService) SyncInvoiceDataFromUB(ctx cloud.Context, req UBRequest) (interface{}, *cloud.Error) {

	var sd []cloud.CustomerStatements

	sd, err := utilibill.GetInvoiceDataFromUtilibill()
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindExternal,
			Message: "failed to get invoice data from utilibill api",
			Cause:   err,
		}) //fmt.Errorf("failed to get invoice data from utilibill api: %w", err)
	}

	// Now we are going to pass this over to the database.
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.SyncronizeUtilibillStatementData(ctx, tx, sd)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice sync statement data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/DataService.SyncInvoiceData.RunInTransaction failed: %w", err)
	}

	return nil, nil
}

func (svc NPDataService) InitializeInvoiceDataFromUtilibill(ctx cloud.Context, req InitializeUtilibillRequest) (interface{}, *cloud.Error) {
	var err error
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err = svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var sd []cloud.CustomerStatements
	sd, err = utilibill.InitializeInvoiceDataFromUtilibill(ctx.Ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindExternal,
			Message: "failed to initialize invoice data from utilibill api",
			Cause:   err,
		}) //fmt.Errorf("failed to initialize invoice data from utilibill: %w", err)
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.SyncronizeUtilibillStatementData(ctx, tx, sd)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice initialize statement data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/DataService.InitUBInvoiceData.RunInTransaction failed: %w", err)
	}

	return nil, nil
}

func (svc NPDataService) InitCustomersFromUB(ctx cloud.Context, req interface{}) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	data, err := utilibill.GetCustomerListFromUB()
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindExternal,
			Message: "failed to customer list from utilibill api",
			Cause:   err,
		}) //fmt.Errorf("failed to get customer list from utilibill api: %w", err)
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.SynchronizeUtilibillCustomerList(ctx, tx, data)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice init customer data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/DataService.InitUBCustomerunInTransaction failed: %w", err)
	}

	return nil, nil
}

// This method will pull all customer data from utilibill and synchronize the customer data in the database.
func (svc NPDataService) PullCustomerDataFromUB(ctx cloud.Context, req UBRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var List cloud.UBCustomerNumberList
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		List, dbErr = db.GetCustomerNumberList(ctx, tx)
		return dbErr
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get customer list transaction failed",
			Cause:   err,
		}) //fmt.Errorf("failed to get customer number list in transaction: %w", err)
	}

	for _, v := range List.Customers {
		// utilibill will only allow one request per second to these endpoints, so we will start with a wait.
		time.Sleep(1 * time.Second)
		fmt.Printf("printing: %v\n", v)
		// this section retrieves the customer details.
		data, err := utilibill.GetCustomerDetailsFromUB(v)
		if err != nil {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindExternal,
				Message: "failed to get customer data from utilibill api",
				Cause:   err,
			}) //fmt.Errorf("failed to get customer data from utilibill api: %w", err)
		}

		err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
			return db.SynchronizeUtilibillCustomerData(ctx, tx, data)
		})
		if err != nil {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindInternal,
				Message: "dataservice sync utilibill customer data transaction failed",
				Cause:   err,
			}) //fmt.Errorf("service/DataService.SyncUBCustomerunInTransaction failed: %w", err)
		}

		// this section will get the direct_debit_status from the customer's payment info
		data2, err := utilibill.GetDirectDebitFromUtilibill(v)
		if err != nil {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindExternal,
				Message: "failed to get direct debit data from utilibill api",
				Cause:   err,
			}) //fmt.Errorf("failed to get customer data from utilibill api: %w", err)
		}

		err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
			return db.UpdateUtilibillDirectDebitData(ctx, tx, data2)
		})
		if err != nil {
			return nil, cloud.NewError(cloud.ErrOpts{
				Kind:    cloud.ErrKindInternal,
				Message: "dataservice sync utilibill direct debit data transaction failed",
				Cause:   err,
			}) //fmt.Errorf("service/DataService.SyncUBCustomerDDSInTransaction failed: %w", err)
		}

	}

	return nil, nil
}

func (svc NPDataService) PullMeterDataFromUB(ctx cloud.Context, req UBRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return InvoiceListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	return nil, nil
}

func (svc NPDataService) GetListOfInvoices(ctx cloud.Context, req InvoiceListRequest) (InvoiceListResponse, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return InvoiceListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var resp InvoiceListResponse
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		resp.Invoices, dbErr = db.ListInvoiceDataByCustomerID(ctx, tx, req.CustomerID)
		return dbErr
	})
	if err != nil {
		return InvoiceListResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get list of invoice data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("service/DataService.GetListOfInvoices.RunInTransaction failed: %w", err)
	}

	return resp, nil
}

// This method will return a set of customer statement data.
func (svc NPDataService) GetInvoiceDataForDisplay(ctx cloud.Context, req GetInvoiceDataRequest) (GetInvoiceDataResponse, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}
	var resp GetInvoiceDataResponse
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		resp.I, dbErr = db.GetInvoiceData(ctx, tx, req.CustomerID, req.StatementNumber)
		if dbErr != nil {
			return err //fmt.Errorf("service/DataService.GetInvoiceData.RunInTransaction failed: %w", dbErr)
		}
		return nil
	})
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get invoice data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("failed to get data: %w", err)
	}

	return resp, nil
}

func (svc NPDataService) ListCustomers(ctx cloud.Context, req UBRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var list []cloud.NPCustomerSummary
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		list, dbErr = db.ListNPCustomers(ctx, tx)
		if dbErr != nil {
			return dbErr
		}
		return nil
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice list customer data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("listCustomers/RunInTransaction failed: %w", err)
	}

	return list, nil
}

func (svc NPDataService) GetCustomerDetails(ctx cloud.Context, req InvoiceListRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	var customer cloud.NPCustomerDetail
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		var dbErr error
		customer, dbErr = db.GetNPCustomerDetail(ctx, tx, req.CustomerID)
		if dbErr != nil {
			return dbErr
		}
		return nil
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice get customer detail transaction failed",
			Cause:   err,
		}) //fmt.Errorf("getCustomerDetail/RunInTransaction failed: %w", err)
	}
	return customer, nil
}

func (svc NPDataService) UpdateCustomerDetails(ctx cloud.Context, req UpdateCustomerDetailRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	// First, lets save these values to the database.
	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.UpdateNPCustomerDetail(ctx, tx, req.Data)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice update customer data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("updateCustomerDetail/RunInTransaction failed: %w", err)
	}

	// Now, lets send these values to Utilibill to update that record.
	retDetail, err := utilibill.UpdateCustomerDetailOnUtilibill(ctx, req.Data)
	if err != nil || retDetail.Success.CustomerNumber != strconv.Itoa(req.Data.CustomerNumber) {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindExternal,
			Message: "failed to update customer details on utilibill api",
			Cause:   err,
		}) //fmt.Errorf("error updating utilibill: %w", err)
	}

	return nil, nil
}

func (svc NPDataService) InitializeMeterData(ctx cloud.Context, req InitializeMeterDataRequest) (interface{}, *cloud.Error) {
	// We checked that the token existed and was valid in the handler, but now we need to make sure the user is allowed to do this action.
	err := svc.ValidateUserAccess(ctx)
	if err != nil {
		return GetInvoiceDataResponse{}, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindForbidden,
			Message: "user not allowed to perform this action",
			Cause:   err,
		}) //fmt.Errorf("user not allowed to perform this action: %w", err)
	}

	err = svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		return db.InitMeterData(ctx, tx, req.Data)
	})
	if err != nil {
		return nil, cloud.NewError(cloud.ErrOpts{
			Kind:    cloud.ErrKindInternal,
			Message: "dataservice initialize meter data transaction failed",
			Cause:   err,
		}) //fmt.Errorf("updateCustomerDetail/RunInTransaction failed: %w", err)
	}

	return nil, nil
}

func (svc NPDataService) ValidateUserAccess(ctx cloud.Context) error {
	err := svc.DB.RunInTransaction(ctx, func(ctx cloud.Context, tx pg.Tx) error {
		dbErr := db.UserAccess(ctx, tx)
		if dbErr != nil {
			return dbErr
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return nil
}
