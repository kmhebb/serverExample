package service

import (
	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/pg"
)

type Tx interface {
	// User DB methods
	FindByToken(ctx cloud.Context, token string) (*cloud.User, error)
	FindByID(ctx cloud.Context, token string) (*cloud.User, error)
	UpdateUserRecord(ctx cloud.Context, u *cloud.User) error
	FindByEmail(ctx cloud.Context, email string) (*cloud.User, error)
	CreateUser(ctx cloud.Context, u *cloud.User) error
	SaveValidationCode(ctx cloud.Context, email string, code string) error
	GetValidationCode(ctx cloud.Context, e string) (string, error)
	DeleteValidationCode(ctx cloud.Context, email string) error
	UpdateLastActivity(ctx cloud.Context, uid string) error
	GetUserList(ctx cloud.Context, listType string) ([]cloud.User, error)

	// Data Service DB methods
	ImportGridData(ctx cloud.Context, data *[]cloud.GridDataRecord) error
	DeleteGridData(ctx cloud.Context, BatchID string) error
	GetGridBatchList(ctx cloud.Context, ListType string) error
	GetBillingDataList(ctx cloud.Context, ListType string) ([]cloud.BillingData, error)
	SyncronizeUtilibillStatementData(ctx cloud.Context, data []cloud.StatementData) error
	ListInvoiceDataByCustomerID(ctx cloud.Context, customerid int) ([]cloud.StatementData, error)
	GetInvoiceData(ctx cloud.Context, customerID int, statementID int) (cloud.StatementData, error)
	UserAccess(ctx cloud.Context) error
	SynchronizeUtilibillCustomerList(ctx cloud.Context, data []cloud.UBCustomerSummary) error
	GetCustomerNumberList(ctx cloud.Context) (cloud.UBCustomerNumberList, error)
	SynchronizeUtilibillCustomerData(ctx cloud.Context, data []cloud.UBCustomerDetail) error
	UpdateUtilibillDirectDebitData(ctx cloud.Context, data cloud.UBCustomerDirectDebit) error
	ListNPCustomers(ctx cloud.Context) ([]cloud.NPCustomerSummary, error)
	GetNPCustomerDetail(ctx cloud.Context, customerid int) (cloud.NPCustomerDetail, error)
	UpdateNPCustomerDetail(ctx cloud.Context, data cloud.NPCustomerDetail) error
	InitMeterData(ctx cloud.Context, data cloud.CustomerMeterData) error
}

type Database interface {
	RunInTransaction(cloud.Context, func(cloud.Context, pg.Tx) error) error
}
