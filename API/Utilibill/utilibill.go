package utilibill

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/API/goxml"
	mps "github.com/mitchellh/mapstructure"
)

type payload struct {
	Params params `xml:"ser:get"`
}
type login struct {
	Password string `xml:"password"`
	UserName string `xml:"userName"`
}

type parameter struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type data struct {
	Code       string      `xml:"code"`
	Parameters []parameter `xml:"parameters"`
}
type params struct {
	Login login `xml:"login"`
	Data  data  `xml:"data"`
}

var statementURL string = "https://go.utilitybilling.com/billing/UtbServiceStatement?wsdl"
var customerURL string = "https://go.utilitybilling.com/billing/UtbServiceCustomer?wsdl"
var transactionURL string = "https://go.utilitybilling.com/billing/UtbServiceTransaction?wsdl"
var meterURL string = "https://go.utilitybilling.com/billing/UtbServiceMeter?wsdl"
var rateplanURL string = "https://go.utilitybilling.com/billing/UtbServiceRatePlan?wsdl"

var loginCredentials login

//:= cloud.CustomerInterface{
// 		Data: summaries,
// 		}
// 	}

func SetCredentials(username string, password string) {
	loginCredentials = login{
		UserName: username,
		Password: password,
	}
}

func GetCustomerListFromUB() ([]cloud.UBCustomerSummary, error) {
	var RawResponse cloud.CustomerListInterface

	d := data{
		Code:       "CUSTOMER_LIST",
		Parameters: []parameter{},
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(customerURL, "GET", pl)
	if err != nil {
		return []cloud.UBCustomerSummary{}, fmt.Errorf("error getting data from utilibill api: %w", err)
	}
	var details []cloud.UBCustomerSummary
	RawResponse = cloud.CustomerListInterface{
		Data: details,
	}

	jsonErr := json.Unmarshal([]byte(data), &RawResponse)
	if jsonErr != nil {
		return []cloud.UBCustomerSummary{}, fmt.Errorf("error decoding customer list json from utilibill api: %w", err)
	}

	return RawResponse.Data, nil
}

func GetCustomerDetailsFromUB(customerNumber int) (cloud.UBCustomerDetail, error) {
	customer := parameter{
		Key:   "customerNumber",
		Value: fmt.Sprint(customerNumber),
	}
	d := data{
		Code: "CUSTOMER_DETAILS",
		Parameters: []parameter{
			customer,
		},
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(customerURL, "GET", pl)
	if err != nil {
		return cloud.UBCustomerDetail{}, fmt.Errorf("error getting data from utilibill api: %w", err)
	}
	var details cloud.UBCustomerDetail

	jsonErr := json.Unmarshal([]byte(data), &details)
	if jsonErr != nil {
		return cloud.UBCustomerDetail{}, fmt.Errorf("error decoding customer detail json from utilibill api: %w", err)
	}

	return details, nil
}

func GetDirectDebitFromUtilibill(customerNumber int) (cloud.UBCustomerDirectDebit, error) {
	customer := parameter{
		Key:   "customerNumber",
		Value: fmt.Sprint(customerNumber),
	}
	d := data{
		Code: "CUSTOMER_DIRECT_DEBIT",
		Parameters: []parameter{
			customer,
		},
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(customerURL, "GET", pl)
	if err != nil {
		return cloud.UBCustomerDirectDebit{}, fmt.Errorf("error getting data from utilibill api: %w", err)
	}

	var details cloud.UBCustomerDirectDebit

	jsonErr := json.Unmarshal([]byte(data), &details)
	if jsonErr != nil {
		return cloud.UBCustomerDirectDebit{}, fmt.Errorf("error decoding customer direct debit json from utilibill api: %w", err)
	}

	return details, nil
}

func GetInvoiceDataFromUtilibill() ([]cloud.CustomerStatements, error) {
	yesterday := time.Now().Add(-time.Hour * 24)
	issueDateFrom := parameter{
		Key:   "issueDateFrom",
		Value: yesterday.Format("1/2/2006"),
	}
	issueDateTo := parameter{
		Key:   "issueDateTo",
		Value: time.Now().Format("1/2/2006"),
	}
	d := data{
		Code: "STATEMENT",
		Parameters: []parameter{
			issueDateFrom,
			issueDateTo,
		},
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(statementURL, "GET", pl)
	if err != nil {
		return nil, fmt.Errorf("error getting data from utilibill api: %w", err)
	}

	var RawResponse cloud.StatementInterface
	err = json.Unmarshal([]byte(data), &RawResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding invoice data json from utilibill api: %w", err)
	}

	var statementbook []cloud.CustomerStatements
	for _, v := range RawResponse.Data {
		fm := v.(map[string]interface{})
		for ii, vv := range fm {
			statements := cloud.CustomerStatements{}
			statements.CustomerID = ii
			fms := vv.([]interface{})
			for _, vvv := range fms {
				statement := cloud.StatementData{}
				fms2 := vvv.(map[string]interface{})
				err = mps.WeakDecode(fms2, &statement)
				if err != nil {
					return nil, fmt.Errorf("failed to decode map: %w", err)
				}
				statements.Statements = append(statements.Statements, statement)
			}
			statementbook = append(statementbook, statements)
		}
		//fmt.Printf("statementbook: %+v", statementbook)
	}

	return statementbook, nil
}

func InitializeInvoiceDataFromUtilibill(ctx context.Context) ([]cloud.CustomerStatements, error) {
	d := data{
		Code: "STATEMENT",
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(statementURL, "GET", pl)
	if err != nil {
		return nil, fmt.Errorf("error getting data from utilibill api: %w", err)
	}
	var RawResponse cloud.StatementInterface
	err = json.Unmarshal([]byte(data), &RawResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding invoice data json from utilibill api: %w", err)
	}

	var statementbook []cloud.CustomerStatements
	for _, v := range RawResponse.Data {
		fm := v.(map[string]interface{})
		for ii, vv := range fm {
			statements := cloud.CustomerStatements{}
			statements.CustomerID = ii
			fms := vv.([]interface{})
			for _, vvv := range fms {
				statement := cloud.StatementData{}
				fms2 := vvv.(map[string]interface{})
				err = mps.WeakDecode(fms2, &statement)
				if err != nil {
					return nil, fmt.Errorf("failed to decode map: %w", err)
				}
				statements.Statements = append(statements.Statements, statement)
			}
			statementbook = append(statementbook, statements)
		}
		//fmt.Printf("statementbook: %+v", statementbook)
	}

	return statementbook, nil
}

func UpdateCustomerDetailOnUtilibill(ctx cloud.Context, details cloud.NPCustomerDetail) (cloud.UBCustomerDetail, error) {
	customerNumber := parameter{
		Key:   "customerNumber",
		Value: strconv.Itoa(details.CustomerNumber),
	}
	custType := parameter{
		Key:   "custType",
		Value: details.CustType,
	}
	firstName := parameter{
		Key:   "firstName",
		Value: details.FirstName,
	}
	lastName := parameter{
		Key:   "lastName",
		Value: details.LastName,
	}
	phoneNumber := parameter{
		Key:   "phoneNumber",
		Value: details.PhoneNumber,
	}
	mobileNumber := parameter{
		Key:   "mobileNumber",
		Value: details.MobileNumber,
	}
	email := parameter{
		Key:   "emailAddress",
		Value: details.EmailAddress,
	}
	company := parameter{
		Key:   "company",
		Value: details.Company,
	}
	billingAddress := parameter{
		Key:   "billingAddress",
		Value: details.BillingAddress,
	}
	billingAddress2 := parameter{
		Key:   "billingAddress2nd",
		Value: details.BillingAddress2nd,
	}
	billingAddressCity := parameter{
		Key:   "billingSuburb",
		Value: details.BillingCity,
	}
	billingAddressState := parameter{
		Key:   "billingState",
		Value: details.BillingState,
	}
	billingAddressZip := parameter{
		Key:   "billingPostalCode",
		Value: details.BillingZip,
	}
	isHomeAddrSame := parameter{
		Key:   "isHomeAddressSameAsBilling",
		Value: evaluateParameter(details.IsHomeAddressSameAsBilling, "bool"),
	}
	homeAddress := parameter{
		Key:   "homeAddress",
		Value: details.HomeAddress,
	}
	homeAddress2 := parameter{
		Key:   "homeAddress2nd",
		Value: details.HomeAddress2nd,
	}
	homeAddressCity := parameter{
		Key:   "homeSuburb",
		Value: details.HomeCity,
	}
	homeAddressState := parameter{
		Key:   "homeState",
		Value: details.HomeState,
	}
	homeAddressZip := parameter{
		Key:   "homePostalCode",
		Value: details.HomeZip,
	}
	emailPdf := parameter{
		Key:   "emailPdf",
		Value: evaluateParameter(details.EmailPdf, "bool"),
	}
	printBill := parameter{
		Key:   "printBill",
		Value: evaluateParameter(details.PrintBill, "bool"),
	}

	d := data{
		Code: "CUSTOMER_DETAILS",
		Parameters: []parameter{
			customerNumber,
			custType,
			firstName,
			lastName,
			phoneNumber,
			mobileNumber,
			email,
			company,
			billingAddress,
			billingAddress2,
			billingAddressCity,
			billingAddressState,
			billingAddressZip,
			isHomeAddrSame,
			homeAddress,
			homeAddress2,
			homeAddressCity,
			homeAddressState,
			homeAddressZip,
			emailPdf,
			printBill,
		},
	}
	p := params{
		Login: loginCredentials,
		Data:  d,
	}
	pl := payload{
		Params: p,
	}

	data, err := goxml.UBsoapCall(customerURL, "UPDATE", pl)
	if err != nil {
		return cloud.UBCustomerDetail{}, fmt.Errorf("error getting data from utilibill api: %w", err)
	}
	var retDetail cloud.UBCustomerDetail
	err = json.Unmarshal([]byte(data), &retDetail)
	if err != nil {
		return cloud.UBCustomerDetail{}, fmt.Errorf("error decoding invoice data json from utilibill api: %w", err)
	}
	return retDetail, nil
}

func evaluateParameter(param interface{}, outputType string) string {
	// Utilibill has odd handling of bool type variables. This is the only evaluation at the moment, but its set up to handle others.
	switch outputType {
	case "bool":
		if param == true {
			return "Yes"
		}
		return "No"
	}
	return ""
}
