package cloud

import (
	"time"

	"github.com/gofrs/uuid"
)

type GridDataRecord struct {
	HostAcct         int       `json:"host_acct" db:"host_acct" csv:"host_acct"`
	SatAcct          int       `json:"sat_acct" db:"sat_acct" csv:"sat_acct"`
	SatelliteName    string    `json:"satellite_name" db:"satellite_name" csv:"satellite_name"`
	SatServClass     string    `json:"sat_serv_class" db:"sat_serv_class" csv:"sat_serv_class"`
	SatVDL           int       `json:"sat_vdl" db:"sat_vdl" csv:"sat_vdl"`
	SatStatus        string    `json:"sat_status" db:"sat_status" csv:"sat_status"`
	VderEnergy       float64   `json:"vder_energy" db:"vder_energy" csv:"vder_energy"`
	VderCap          float64   `json:"vder_cap" db:"vder_cap" csv:"vder_cap"`
	VderEnv          float64   `json:"vder_env" db:"vder_env" csv:"vder_env"`
	VderDrv          float64   `json:"vder_drv" db:"vdr_drv" csv:"vdr_drv"`
	VderLsrv         float64   `json:"vder_lsrv" db:"vder_lsrv" csv:"vder_lsrv"`
	VderMTC          float64   `json:"vder_mtc" db:"vder_mtc" csv:"vder_mtc"`
	VderCc           float64   `json:"vder_cc" db:"vder_cc" csv:"vder_cc"`
	VderTotal        float64   `json:"vder_total" db:"vder_total" csv:"vder_total"`
	TransKWH         float64   `json:"trans_kwh" db:"trans_kwh" csv:"trans_kwh"`
	Allocation       float64   `json:"allocation" db:"allocation" csv:"allocation"`
	HostBillPeriod   string    `json:"host_bill_period" db:"host_bill_period" csv:"host_bill_period"`
	TransferDate     time.Time `json:"transfer_date" db:"transfer_date" csv:"transfer_date"`
	SatBillDate      time.Time `json:"sat_bill_date" db:"sat_bill_date" csv:"sat_bill_date"`
	BankedPriorMonth float64   `json:"banked_prior_month" db:"banked_prior_month" csv:"banked_prior_month"`
	CurrentVDER      float64   `json:"current_vder" db:"current_vder" csv:"current_vder"`
	TotalAvailable   float64   `json:"total_available" db:"total_available" csv:"total_available"`
	SatBillAmt       float64   `json:"sat_bill_amt" db:"sat_bill_amt" csv:"sat_bill_amt"`
	Applied          float64   `json:"applied" db:"applied" csv:"applied"`
	BankedCarryOver  float64   `json:"banked_carry_over" db:"banked_carry_over" csv:"banked_carry_over"`
}

type BatchData struct {
	BatchID   uuid.UUID
	BatchDate time.Time
}

type BillingMetaData struct {
	BillingDataID   uuid.UUID
	BillingDataDate time.Time
}

type BillingData struct {
	CustomerNumber    int       `json:"customerNumber" mapstructure:"customerNumber" db:"customer_number" csv:"Customer Number"`
	MeterNumber       int       `json:"meterNumber" mapstructure:"meterNumber" db:"meter_number" csv:"Meter Number"`
	RollupDesc        string    `json:"rollupDescription" mapstructure:"rollupDescription" db:"rollup_description" csv:"Rollup Description"`
	ChargeDesc        string    `json:"chargeDesc" mapstructure:"chargeDesc" db:"charge_description" csv:"Charge Description"`
	StartDate         time.Time `json:"startDate" mapstructure:"startDate" db:"start_date" csv:"Start Date"`
	EndDate           time.Time `json:"endDate" mapstructure:"endDate" db:"end_date" csv:"End Date"`
	Units             int       `json:"units" mapstructure:"units" db:"units" csv:"Units"`
	ChargeAmount      float32   `json:"chargeAmount" mapstructure:"chargeAmount" db:"charge_amount" csv:"Charge Amount"`
	Rate              float32   `json:"rate" mapstructure:"rate" db:"rate" csv:"Rate"`
	TaxID             int       `json:"taxid" mapstructure:"taxid" db:"taxid" csv:"Tax Id"`
	SpecialChargeCode int       `json:"specialChargeCode" mapstructure:"specialChargeCode" db:"special_charge_code" csv:"Special Charge Code"`
	BillingDate       time.Time `json:"billingDate" mapstructure:"billingDate" db:"billing_date" csv:"-"`
	BillingBatchID    uuid.UUID `json:"billingBatchID" mapstructure:"billingBatchID" db:"billing_batch_id" csv:"-"`
}

type BillingBatchList struct {
	BillingList []BillingMetaData
}

type StatementData struct {
	CustomerNumber  string `json:"customerID"`
	Adjustments     string `json:"adjustments" mapstructure:"adjustments" db:"adjustments"`
	CarriedForward  string `json:"carriedForward" mapstructure:"carriedForward" db:"carried_forward"`
	CurrentBalance  string `json:"currentBalance" mapstructure:"currentBalance" db:"current_balance"`
	CurrentCharges  string `json:"currentCharges" mapstructure:"currentCharges" db:"current_charges"`
	DueDate         string `json:"dueDate" mapstructure:"dueDate" db:"due_date"`
	IssuedDate      string `json:"issuedDate" mapstructure:"issuedDate" db:"issued_date"`
	Payment         string `json:"payment" mapstructure:"payment" db:"payment"`
	PreviousBalance string `json:"previousBalance" mapstructure:"previousBalance" db:"previous_balance"`
	StatementNumber string `json:"statementNumber" mapstructure:"statementNumber" db:"statement_number"`
	StatementType   string `json:"statementType" mapstructure:"statementType" db:"statement_type"`
	Tax             string `json:"tax" mapstructure:"tax" db:"tax"`
}

type StatementSummary struct {
	StatementNumber string `json:"statementNumber" mapstructure:"statementNumber" db:"statement_number"`
	CurrentBalance  string `json:"currentBalance" mapstructure:"currentBalance" db:"current_balance"`
	DueDate         string `json:"dueDate" mapstructure:"dueDate" db:"due_date"`
	IssuedDate      string `json:"issuedDate" mapstructure:"issuedDate" db:"issued_date"`
}

type CustomerListInterface struct {
	Data []UBCustomerSummary `json:"success"`
}

type StatementInterface struct {
	Data []interface{} `json:"success"`
}

type CustomerStatements struct {
	CustomerID string
	Statements []StatementData
}

type UBCustomerNumberList struct {
	Customers []int
}

type UBCustomerSummary struct {
	CustomerNumber    string `json:"customerNumber"`
	AlternativeNumber string `json:"alternativeNumber"`
	CustType          string `json:"custType"`
	CurrentBalance    string `json:"currentBalance"`
	Status            string `json:"status"`
	Salutation        string `json:"salutation"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	MobileNumber      string `json:"mobileNumber"`
	Email             string `json:"emailAddress"`
	Company           string `json:"company"`
}

type UBCustomerDetail struct {
	Success struct {
		CustomerNumber             string `json:"customerNumber"`
		AlternativeNumber          string `json:"alternativeNumber"`
		CustType                   string `json:"custType"`
		CurrentBalance             string `json:"currentBalance"`
		Status                     string `json:"status"`
		Salutation                 string `json:"salutation"`
		FirstName                  string `json:"firstName"`
		LastName                   string `json:"lastName"`
		Abn                        string `json:"abn"`
		Acn                        string `json:"acn"`
		PhoneNumber                string `json:"phoneNumber"`
		PhoneNumberAh              string `json:"phoneNumberAh"`
		MobileNumber               string `json:"mobileNumber"`
		FaxNumber                  string `json:"faxNumber"`
		EmailAddress               string `json:"emailAddress"`
		Company                    string `json:"company"`
		BillingAddress             string `json:"billingAddress"`
		BillingAddress2nd          string `json:"billingAddress2nd"`
		BillingSuburb              string `json:"billingSuburb"`
		BillingState               string `json:"billingState"`
		BillingPostalCode          string `json:"billingPostalCode"`
		BillingCountry             string `json:"billingCountry"`
		BillingEmail               string `json:"billingEmail"`
		Category                   string `json:"category"`
		IsHomeAddressSameAsBilling string `json:"isHomeAddressSameAsBilling"`
		HomeAddress                string `json:"homeAddress"`
		HomeAddress2nd             string `json:"homeAddress2nd"`
		HomeSuburb                 string `json:"homeSuburb"`
		HomeState                  string `json:"homeState"`
		HomePostalCode             string `json:"homePostalCode"`
		HomeCountry                string `json:"homeCountry"`
		Feedbacks                  string `json:"feedbacks"`
		LegalEntity                string `json:"legalEntity"`
		AccountName                string `json:"accountName"`
		AuthenticationType         string `json:"authenticationType"`
		AuthenticationNumber       string `json:"authenticationNumber"`
		BirthDate                  string `json:"birthDate"`
		RefName                    string `json:"refName"`
		RefContactNumber           string `json:"refContactNumber"`
		RefRelationshipTenant      string `json:"refRelationshipTenant"`
		EnableInternetAccess       string `json:"enableInternetAccess"`
		InternetPassword           string `json:"internetPassword"`
		ThirdPartyAutoPay          string `json:"thirdPartyAutoPay"`
		EmailPdf                   string `json:"emailPdf"`
		PrintBill                  string `json:"printBill"`
		IsLifeSupport              string `json:"isLifeSupport"`
		SecondaryCustomerList      []struct {
			SecondaryCustomerName string `json:"secondaryCustomerName"`
			SecondaryPhoneNumber  string `json:"secondaryPhoneNumber"`
			SecondaryMobileNumber string `json:"secondaryMobileNumber"`
			SecondaryFaxNumber    string `json:"secondaryFaxNumber"`
			SecondaryEmail        string `json:"secondaryEmail"`
		} `json:"secondaryCustomerList"`
		InvoiceDefault       string `json:"invoiceDefault"`
		BillingDefault       string `json:"billingDefault"`
		GuarantorFirstName   string `json:"guarantorFirstName"`
		GuarantorLastName    string `json:"guarantorLastName"`
		GuarantorAddress1    string `json:"guarantorAddress1"`
		GuarantorAddress2    string `json:"guarantorAddress2"`
		GuarantorCity        string `json:"guarantorCity"`
		GuarantorState       string `json:"guarantorState"`
		GuarantorZipCode     string `json:"guarantorZipCode"`
		GuarantorHomePhone   string `json:"guarantorHomePhone"`
		GuarantorMobilePhone string `json:"guarantorMobilePhone"`
		GuarantorEmail       string `json:"guarantorEmail"`
		GuarantorAuthType    string `json:"guarantorAuthType"`
		GuarantorAuthNo      string `json:"guarantorAuthNo"`
		CycleNumber          int    `json:"cycleNumber"`
	} `json:"success"`
}

type SecondaryCustomerlist struct {
	SecondaryCustomerName string `json:"secondaryCustomerName"`
	SecondaryPhoneNumber  string `json:"secondaryPhoneNumber"`
	SecondaryMobileNumber string `json:"secondaryMobileNumber"`
	SecondaryFaxNumber    string `json:"secondaryFaxNumber"`
	SecondaryEmail        string `json:"secondaryEmail"`
}

type UBCustomerDirectDebit struct {
	CustomerNumber    string `json:"customerNumber"`
	DirectDebitStatus string `json:"directDebitStatus"`
	TokenCustomerID   string `json:"tokenCustomerId"`
	TokenDateCreated  string `json:"tokenDateCreated"`
	TokenLastModified string `json:"tokenLastModified"`
	CustomerDetails   struct {
		FirstName                   string `json:"firstName"`
		LastName                    string `json:"lastName"`
		MobileNumber                string `json:"mobileNumber"`
		EnableSmsPaymentReminder    string `json:"enableSmsPaymentReminder"`
		EnableSmsExpiredCard        string `json:"enableSmsExpiredCard"`
		EnableSmsFailedNotification string `json:"enableSmsFailedNotification"`
	} `json:"customerDetails"`
}

type NPCustomerSummary struct {
	CustomerNumber int64  `json:"customerNumber"`
	CustType       string `json:"custType"`
	Status         string `json:"status"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	MobileNumber   string `json:"mobileNumber"`
	Email          string `json:"emailAddress"`
	Company        string `json:"company"`
}

type NPCustomerDetail struct {
	CustomerNumber             int    `json:"customerNumber"`
	CustType                   string `json:"custType"`
	Status                     string `json:"status"`
	FirstName                  string `json:"firstName"`
	LastName                   string `json:"lastName"`
	PhoneNumber                string `json:"phoneNumber"`
	MobileNumber               string `json:"mobileNumber"`
	EmailAddress               string `json:"emailAddress"`
	Company                    string `json:"company"`
	BillingAddress             string `json:"billingAddress"`
	BillingAddress2nd          string `json:"billingAddress2nd"`
	BillingCity                string `json:"billingCity"`
	BillingState               string `json:"billingState"`
	BillingZip                 string `json:"billingZip"`
	IsHomeAddressSameAsBilling bool   `json:"isHomeAddressSameAsBilling"`
	HomeAddress                string `json:"homeAddress"`
	HomeAddress2nd             string `json:"homeAddress2nd"`
	HomeCity                   string `json:"homeCity"`
	HomeState                  string `json:"homeState"`
	HomeZip                    string `json:"homeZip"`
	EmailPdf                   bool   `json:"emailPdf"`
	PrintBill                  bool   `json:"printBill"`
	CycleNumber                int    `json:"cycleNumber"`
	DirectDebitStatus          string `json:"directDebitStatus"`
}

type CustomerMeterData struct {
	//MeterID         int     `json:"-"`
	GridAcctNumber  string `json:"grid_acct" db:"grid_acct"`
	CustomerNumber  string `json:"customer_number" db:"customer_number"`
	StreetAddress   string `json:"street" db:"street"`
	City            string `json:"city" db:"city"`
	State           string `json:"state" db:"state"`
	Zip             string `json:"zip" db:"zip"`
	GridBillGroup   string `json:"grid_bill_group" db:"grid_bill_group"`
	MeterClass      string `json:"meter_class" db:"meter_class"`
	PaperlessCredit string `json:"paperless_credit" db:"paperless_credit"`
	HostFacility    string `json:"host_facility" db:"host_facility"`
}
