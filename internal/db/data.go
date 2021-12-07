package db

import (
	"fmt"
	"strconv"
	"time"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/internal/log"
	"github.com/kmhebb/serverExample/pg"
	"github.com/pborman/uuid"
)

func ImportGridData(ctx cloud.Context, tx pg.Tx, data []cloud.GridDataRecord) error {
	logger := log.NewLogger()
	upload_id := uuid.NewUUID()
	upload_date := time.Now()
	var errCount int
	for _, row := range data {
		query := `INSERT INTO customer.utility_data (
		host_acct,
		sat_acct,
		satellite_name,
		sat_serv_class,
		sat_vdl,
		sat_status,
		vder_energy,
		vder_cap,
		vder_env,
		vder_drv,
		vder_lsrv,
		vder_mtc,
		vder_total,
		trans_kwh,
		allocation,
		host_bill_period,
		transfer_date,
		sat_bill_date,
		banked_prior_month,
		current_vder,
		total_available,
		sat_bill_amt,
		applied,
		banked_carry_over,
		upload_id,
		upload_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26);`

		err := tx.Exec(ctx.Ctx, query,
			row.HostAcct, row.SatAcct, row.SatelliteName, row.SatServClass, row.SatVDL, row.SatStatus, row.VderEnergy, row.VderCap,
			row.VderEnv, row.VderDrv, row.VderLsrv, row.VderMTC, row.VderTotal, row.TransKWH, row.Allocation, row.HostBillPeriod, row.TransferDate,
			row.SatBillDate, row.BankedPriorMonth, row.CurrentVDER, row.TotalAvailable, row.SatBillAmt, row.Applied, row.BankedCarryOver,
			upload_id, upload_date)
		if err != nil {
			logger.Debug(fmt.Sprintf("error saving to db: %+v", err))
			errCount += 1
		}
	}

	if errCount > 0 {
		// log the problem for debugging.
		logger.Debug(fmt.Sprintf("failed to add %d records in upload: %s", errCount, upload_id))
		return fmt.Errorf("there were %d insert errors in upload: %s", errCount, upload_id)
	}

	return nil
}

func DeleteGridData(ctx cloud.Context, tx pg.Tx, BatchID string) error {
	q := `DELETE from customer.utility_data WHERE upload_id = $1`
	err := tx.Exec(ctx.Ctx, q, BatchID)
	if err != nil {
		return fmt.Errorf("DeleteGridData failed to delete: %w", err)
	}
	return nil
}

func GetGridBatchList(ctx cloud.Context, tx pg.Tx, ListType string) ([]cloud.BatchData, error) {
	var q string
	switch ListType {
	case "all":
		q = `SELECT distinct(upload_id), upload_date FROM customer.utility_data`
	default:
		q = `SELECT distinct(upload_id), upload_date FROM customer.utility_data`
	}

	var batchList []cloud.BatchData
	rows, err := tx.Query(ctx.Ctx, q)
	if err != nil {
		return nil, fmt.Errorf("batch list query failed: %w", err)
	}
	for rows.Next() {
		var b cloud.BatchData
		err = rows.Scan(&b.BatchID, &b.BatchDate)
		if err != nil {
			return nil, fmt.Errorf("batch list assignment failed: %w", err)
		}
		batchList = append(batchList, b)
	}
	return batchList, nil

}

func GetBillingDataList(ctx cloud.Context, tx pg.Tx, ListType string) ([]cloud.BillingMetaData, error) {
	var q string
	switch ListType {
	case "all":
		q = `SELECT distinct(billing_batch_id), billing_date FROM customer.billing_data`
	default:
		q = `SELECT distinct(billing_batch_id), billing_date FROM customer.billing_data`
	}

	var batchList []cloud.BillingMetaData
	rows, err := tx.Query(ctx.Ctx, q)
	if err != nil {
		return nil, fmt.Errorf("billing data list query failed: %w", err)
	}
	for rows.Next() {
		var b cloud.BillingMetaData
		err = rows.Scan(&b.BillingDataID, &b.BillingDataDate)
		if err != nil {
			return nil, fmt.Errorf("billing data list assignment failed: %w", err)
		}
		batchList = append(batchList, b)
	}
	return batchList, nil
}

func SyncronizeUtilibillStatementData(ctx cloud.Context, tx pg.Tx, data []cloud.CustomerStatements) error {

	query := `INSERT INTO customer.statement_data (customerid, adjustments, carried_forward, current_charges, current_balance, due_date, issued_date, payment, previous_balance, statement_number, statement_type, tax) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, v := range data {
		for _, vv := range v.Statements {
			err := tx.Exec(ctx.Ctx, query, v.CustomerID, vv.Adjustments, vv.CarriedForward, vv.CurrentCharges, vv.CurrentBalance, vv.DueDate, vv.IssuedDate, vv.Payment, vv.PreviousBalance, vv.StatementNumber, vv.StatementType, vv.Tax)
			if err != nil {
				return fmt.Errorf("statement data insert failed: %w", err)
			}
		}
	}

	return nil
}

func ListInvoiceDataByCustomerID(ctx cloud.Context, tx pg.Tx, customerid int) ([]cloud.StatementSummary, error) {
	query := `SELECT statement_number, issued_date, due_date, current_balance FROM customer.statement_data WHERE customerid = $1`

	var data []cloud.StatementSummary
	rows, err := tx.Query(ctx.Ctx, query, fmt.Sprint(customerid))
	if err != nil {
		return []cloud.StatementSummary{}, fmt.Errorf("list invoice query failed: %w", err)
	}

	for rows.Next() {
		var st cloud.StatementSummary
		err = rows.Scan(&st.StatementNumber, &st.IssuedDate, &st.DueDate, &st.CurrentBalance)
		if err != nil {
			return []cloud.StatementSummary{}, fmt.Errorf("list invoice data assignment failed: %w", err)
		}
		data = append(data, st)
	}

	return data, nil
}

func GetInvoiceData(ctx cloud.Context, tx pg.Tx, customerID int, statementID int) (cloud.StatementData, error) {
	query := `SELECT * from customer.statement_data WHERE customerid = $1 and statement_number = $2`

	var data cloud.StatementData
	rows, err := tx.Query(ctx.Ctx, query, fmt.Sprint(customerID), fmt.Sprint(statementID))
	if err != nil {
		return cloud.StatementData{}, fmt.Errorf("statement data query failed: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&data.CustomerNumber, &data.Adjustments, &data.CarriedForward, &data.CurrentCharges, &data.CurrentBalance, &data.DueDate, &data.IssuedDate, &data.Payment, &data.PreviousBalance, &data.StatementNumber, &data.StatementType, &data.Tax)
		if err != nil {
			return cloud.StatementData{}, fmt.Errorf("failed to load query into data: %w", err)
		}
	}

	return data, nil
}

func UserAccess(ctx cloud.Context, tx pg.Tx) error {
	query := `SELECT CAST(id AS varchar), firstname, lastname, email, passhash, mustchange, lastactivity, datecreated, datemodified FROM users.profile WHERE id = $1;`
	var u cloud.User

	rows, err := tx.Query(ctx.Ctx, query, ctx.UserKey)
	if err != nil {
		return fmt.Errorf("pg/Tx.UserFindByIDQuery: %w", err)
	}

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.MustChange, &u.LastActivity, &u.DateModified, &u.DateCreated); err != nil {
			return fmt.Errorf("pg/Tx.UserFindByIDAssignment: %w", err)
		}
	}

	return nil
}

func SynchronizeUtilibillCustomerData(ctx cloud.Context, tx pg.Tx, data cloud.UBCustomerDetail) error {
	query := `INSERT INTO customer.customers (
		customer_number, 
		alternative_number, 
		cust_type, 
		status, 
		firstname, 
		lastname, 
		abn, 
		acn, 
		phonenumber, 
		phonenumber_ah, 
		mobilenumber, 
		email, 
		company, 
		billing_address, 
		billing_address2, 
		billing_city, 
		billing_state, 
		billing_zip, 
		is_billing_addr_same, 
		home_address, 
		home_address2, 
		home_city, 
		home_state, 
		home_zip,
		birthdate, 
		emailpdf, 
		printbill,
		is_life_support,
		cycle_number
		) VALUES ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
		$29
		) ON CONFLICT (customer_number) DO UPDATE SET (
		customer_number, 
		alternative_number, 
		cust_type, 
		status, 
		firstname, 
		lastname, 
		abn, 
		acn, 
		phonenumber, 
		phonenumber_ah, 
		mobilenumber, 
		email, 
		company, 
		billing_address, 
		billing_address2, 
		billing_city, 
		billing_state, 
		billing_zip, 
		is_billing_addr_same, 
		home_address, 
		home_address2, 
		home_city, 
		home_state, 
		home_zip,
		birthdate, 
		emailpdf, 
		printbill,
		is_life_support,
		cycle_number
		) = ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
		$29
		)`

	CustomerInt, intErr := stringToInt(data.Success.CustomerNumber, "int")
	if intErr != nil {
		fmt.Printf("error converting customer number: %+v", intErr)
	}
	IsHome, intErr := stringToInt(data.Success.IsHomeAddressSameAsBilling, "bool")
	if intErr != nil {
		fmt.Printf("error converting IsHomeAddress: %+v", intErr)
	}
	emailPDF, intErr := stringToInt(data.Success.EmailPdf, "bool")
	if intErr != nil {
		fmt.Printf("error converting emailpdf: %+v", intErr)
	}
	printBill, intErr := stringToInt(data.Success.PrintBill, "bool")
	if intErr != nil {
		fmt.Printf("error converting printbill: %+v", intErr)
	}
	isLifeSupport, intErr := stringToInt(data.Success.IsLifeSupport, "bool")
	if intErr != nil {
		fmt.Printf("error converting isLifeSupport: %+v", intErr)
	}
	err := tx.Exec(ctx.Ctx, query,
		CustomerInt,
		data.Success.AlternativeNumber,
		data.Success.CustType,
		data.Success.Status,
		data.Success.FirstName,
		data.Success.LastName,
		data.Success.Abn,
		data.Success.Acn,
		data.Success.PhoneNumber,
		data.Success.PhoneNumberAh,
		data.Success.MobileNumber,
		data.Success.EmailAddress,
		data.Success.Company,
		data.Success.BillingAddress,
		data.Success.BillingAddress2nd,
		data.Success.BillingSuburb,
		data.Success.BillingState,
		data.Success.BillingPostalCode,
		IsHome,
		data.Success.HomeAddress,
		data.Success.HomeAddress2nd,
		data.Success.HomeSuburb,
		data.Success.HomeState,
		data.Success.HomePostalCode,
		data.Success.BirthDate,
		emailPDF,
		printBill,
		isLifeSupport,
		data.Success.CycleNumber)
	if err != nil {
		return fmt.Errorf("there were errors inserting customer data: %w", err)
	}

	return nil
}

func UpdateUtilibillDirectDebitData(ctx cloud.Context, tx pg.Tx, data cloud.UBCustomerDirectDebit) error {
	query := `UPDATE customer.customers SET direct_debit_status = $1 WHERE customer_number = $2`
	id, convErr := strconv.Atoi(data.CustomerNumber)
	if convErr != nil {
		return convErr
	}
	err := tx.Exec(ctx.Ctx, query, data.DirectDebitStatus, id)
	if err != nil {
		return fmt.Errorf("failed to update customer direct debit status: %w", err)
	}

	return nil
}

func SynchronizeUtilibillCustomerList(ctx cloud.Context, tx pg.Tx, data []cloud.UBCustomerSummary) error {
	query := `INSERT INTO customer.customers (
		customer_number, 
		alternative_number, 
		cust_type, 
		status, 
		firstname, 
		lastname, 
		mobilenumber, 
		email, 
		company
		) VALUES ($1, $2,$3, $4, $5, $6, $7, $8, $9) 
		ON CONFLICT (customer_number) DO UPDATE SET (
		customer_number, 
		alternative_number, 
		cust_type, 
		status, 
		firstname, 
		lastname, 
		mobilenumber, 
		email, 
		company
		) = ($1, $2,$3, $4, $5, $6, $7, $8, $9)  `

	var insertErrors int
	var captureError error
	for _, v := range data {
		CustomerInt, intErr := stringToInt(v.CustomerNumber, "int")
		if intErr != nil {
			fmt.Printf("error converting customer number: %+v", intErr)
		}

		err := tx.Exec(ctx.Ctx, query,
			CustomerInt,
			v.AlternativeNumber,
			v.CustType,
			v.Status,
			v.FirstName,
			v.LastName,
			v.MobileNumber,
			v.Email,
			v.Company,
		)
		if err != nil {
			insertErrors = +1
			captureError = err
		}
	}
	if insertErrors != 0 {
		return fmt.Errorf("there were %d errors inserting customer list data: %w", insertErrors, captureError)
	}

	return nil
}

func GetCustomerNumberList(ctx cloud.Context, tx pg.Tx) (cloud.UBCustomerNumberList, error) {
	query := `SELECT customer_number FROM customer.customers`

	rows, err := tx.Query(ctx.Ctx, query)
	if err != nil {
		return cloud.UBCustomerNumberList{}, fmt.Errorf("failed to query for customer number list: %w", err)
	}

	var data cloud.UBCustomerNumberList
	for rows.Next() {
		var customer int
		err := rows.Scan(&customer)
		if err != nil {
			return cloud.UBCustomerNumberList{}, fmt.Errorf("customer list assignment failed: %w", err)
		}
		data.Customers = append(data.Customers, customer)
	}
	return data, nil
}

func stringToInt(st string, t string) (interface{}, error) {
	switch t {
	case "int":
		if st == "" {
			return 0, nil
		}
		number, err := strconv.Atoi(st)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %v to int: %v", st, err)
		}
		return number, nil
	case "float":
		if st == "" {
			return 0.00, nil
		}
		number, err := strconv.ParseFloat(st, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %v to float: %v", st, err)
		}
		return number, nil
	case "bool":
		if st == "" {
			return false, nil
		}
		if st == "Yes" {
			return true, nil
		}
		if st == "No" {
			return false, nil
		}
		number, err := strconv.ParseBool(st)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %v to bool: %v", st, err)
		}
		return number, nil
	}

	return nil, fmt.Errorf("failed to convert %v", st)
}

func ListNPCustomers(ctx cloud.Context, tx pg.Tx) ([]cloud.NPCustomerSummary, error) {
	query := `SELECT customer_number,
	cust_type,
	status,
	firstname,
	lastname,
	email,
	mobilenumber,
	company FROM customer.customers`
	var list []cloud.NPCustomerSummary
	rows, err := tx.Query(ctx.Ctx, query)
	if err != nil {
		return []cloud.NPCustomerSummary{}, fmt.Errorf("customer list query failed: %w", err)
	}
	for rows.Next() {
		var c cloud.NPCustomerSummary
		err = rows.Scan(&c.CustomerNumber, &c.CustType, &c.Status, &c.FirstName, &c.LastName, &c.Email, &c.MobileNumber, &c.Company)
		if err != nil {
			return []cloud.NPCustomerSummary{}, fmt.Errorf("customer list assignment failed: %w", err)
		}
		list = append(list, c)
	}
	return list, nil
}

func GetNPCustomerDetail(ctx cloud.Context, tx pg.Tx, customerid int) (cloud.NPCustomerDetail, error) {
	query := `SELECT 
	customer_number,
	cust_type,
	status,
	firstname,
	lastname,
	phonenumber,
	mobilenumber,
	email,
	company,
	billing_address,
	billing_address2,
	billing_city,
	billing_state,
	billing_zip,
	is_billing_addr_same,
	home_address,
	home_address2,
	home_city,
	home_state,
	home_zip,
	emailpdf,
	printbill,
	cycle_number,
	direct_debit_status FROM customer.customers WHERE customer_number = $1`

	var c cloud.NPCustomerDetail
	row, err := tx.Query(ctx.Ctx, query, customerid)
	if err != nil {
		return cloud.NPCustomerDetail{}, fmt.Errorf("customer detail query error: %w", err)
	}

	for row.Next() {
		err = row.Scan(&c.CustomerNumber, &c.CustType, &c.Status, &c.FirstName, &c.LastName, &c.PhoneNumber, &c.MobileNumber, &c.EmailAddress, &c.Company, &c.BillingAddress, &c.BillingAddress2nd, &c.BillingCity, &c.BillingState, &c.BillingZip, &c.IsHomeAddressSameAsBilling, &c.HomeAddress, &c.HomeAddress2nd, &c.HomeCity, &c.BillingState, &c.HomeZip, &c.EmailPdf, &c.PrintBill, &c.CycleNumber, &c.DirectDebitStatus)
		if err != nil {
			return cloud.NPCustomerDetail{}, fmt.Errorf("customer detail assignment error: %w", err)
		}
	}

	return c, nil
}

func UpdateNPCustomerDetail(ctx cloud.Context, tx pg.Tx, data cloud.NPCustomerDetail) error {
	query := `UPDATE customer.customers SET (
		cust_type, 
		status, 
		firstname, 
		lastname, 
		phonenumber, 
		mobilenumber, 
		email, 
		company, 
		billing_address, 
		billing_address2, 
		billing_city, 
		billing_state, 
		billing_zip, 
		is_billing_addr_same, 
		home_address, 
		home_address2, 
		home_city, 
		home_state, 
		home_zip, 
		emailpdf, 
		printbill, 
		cycle_number) = ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
			$21, $22) WHERE customer_number = $23`

	err := tx.Exec(ctx.Ctx, query,
		data.CustType,
		data.Status,
		data.FirstName,
		data.LastName,
		data.PhoneNumber,
		data.MobileNumber,
		data.EmailAddress,
		data.Company,
		data.BillingAddress,
		data.BillingAddress2nd,
		data.BillingCity,
		data.BillingState,
		data.BillingZip,
		data.IsHomeAddressSameAsBilling,
		data.HomeAddress,
		data.HomeAddress2nd,
		data.HomeCity,
		data.HomeState,
		data.HomeZip,
		data.EmailPdf,
		data.PrintBill,
		data.CycleNumber,
		data.CustomerNumber)
	if err != nil {
		return fmt.Errorf("update customer detail failed: %w", err)
	}

	return nil
}

func InitMeterData(ctx cloud.Context, tx pg.Tx, data []cloud.CustomerMeterData) error {
	query := `INSERT INTO customer.meters (
		meter_id,
		grid_acct,
		customer_number,
		street,
		city,
		state,
		zip,
		grid_bill_group,
		meter_class,
		paperless_credit,
		host_facility) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) `

	var errCount int
	var errRecord error
	for _, v := range data {
		// we need to manipulate some of the data coming in. its all in json as strings.
		CustomerInt, intErr := stringToInt(v.CustomerNumber, "int")
		if intErr != nil {
			fmt.Printf("error converting customer number: %+v", intErr)
		}
		var GridBillGr interface{}
		var GBErr error
		if v.GridBillGroup == "" {
			GridBillGr = 0
		}
		GridBillGr, GBErr = stringToInt(v.GridBillGroup, "int")
		if GBErr != nil {
			fmt.Printf("error converting grid bill group number: %+v", GBErr)
		}
		var PaperlessCredit float64
		var convErr error
		if v.PaperlessCredit == "" {
			PaperlessCredit = 0.00
		} else {
			PaperlessCredit, convErr = strconv.ParseFloat(v.PaperlessCredit, 32)
			if convErr != nil {
				fmt.Printf("error converting paperless credit: %+v", convErr)
				PaperlessCredit = 0.00
			}
		}

		// now we will run the query with the data
		err := tx.Exec(ctx.Ctx, query,
			uuid.New(),
			v.GridAcctNumber,
			CustomerInt,
			v.StreetAddress,
			v.City,
			v.State,
			v.Zip,
			GridBillGr,
			v.MeterClass,
			PaperlessCredit,
			v.HostFacility,
		)
		if err != nil {
			errCount += 1
			errRecord = err
		}
	}

	// because we are ranging over all the data, we are going to count the issues and return them in bulk.
	if errCount != 0 {
		return fmt.Errorf("there were %v errors inserting meter data: %v", errCount, errRecord)
	}

	return nil

}
