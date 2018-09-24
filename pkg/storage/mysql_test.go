package storage

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	storage *DBHandler
	dbGorm  *gorm.DB
	db      *sql.DB
}

var (
	pt     = "Payment"
	uuid   = "test_uuid"
	ver    = 0
	oID    = "test_organisation_id"
	amount = "12300.21"
	pc     = "GBP"
	pp     = "test"
	ptt    = "test"
	ppd    = "2017-01-18"

	sampleBPData = &BeneficiaryParty{PaymentsID: 1}
	sampleCIData = &ChargesInformation{PaymentsID: 1}
	sampleDPData = &DebtorParty{PaymentsID: 1}
	sampleFXData = &Fx{PaymentsID: 1}

	samplePaymentData = &Payments{
		PaymentsType:            &pt,
		UUID:                    &uuid,
		PaymentsVersion:         &ver,
		PaymentsOrganisationID:  &oID,
		PaymentsAmount:          &amount,
		BeneficiaryParty:        *sampleBPData,
		ChargesInformation:      *sampleCIData,
		PaymentsCurrency:        &pc,
		DebtorParty:             *sampleDPData,
		Fx:                      *sampleFXData,
		PaymentsPurpose:         &pp,
		PaymentsTransactionType: &ptt,
		PaymentsProcessingDate:  &ppd,
	}
)

var testDB = &DBHandler{
	Cfg: &Config{
		User:             "rafal",
		Password:         "password",
		Port:             "3306",
		Db:               "payments_demo_test",
		Host:             "localhost",
		AdditionalParams: "?charset=utf8&parseTime=true&loc=UTC",
	},
}

func setupGormDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql",
		testDB.Cfg.User+":"+
			testDB.Cfg.Password+
			"@tcp("+testDB.Cfg.Host+":"+
			testDB.Cfg.Port+")/"+
			testDB.Cfg.Db+
			testDB.Cfg.AdditionalParams)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("mysql",
		testDB.Cfg.User+":"+
			testDB.Cfg.Password+
			"@tcp("+testDB.Cfg.Host+":"+
			testDB.Cfg.Port+")/"+
			testDB.Cfg.Db+
			testDB.Cfg.AdditionalParams)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *StorageSuite) SetupSuite() {

	var err error
	s.dbGorm, err = setupGormDB()
	if err != nil {
		s.T().Fatal(err)
	}
	// for comparison, it's worth to have another driver used for tests
	s.db, err = setupDB()
	if err != nil {
		s.T().Fatal(err)
	}

	s.storage = &DBHandler{Conn: s.dbGorm, Cfg: testDB.Cfg}
}

// SetupTest will ensure we have a consistent state before tests run.
func (s *StorageSuite) SetupTest() {
	s.db.Exec("DELETE FROM payments")
	s.db.Exec("ALTER TABLE payments AUTO_INCREMENT = 1")
	s.db.Exec("DELETE FROM fx")
	s.db.Exec("ALTER TABLE fx AUTO_INCREMENT = 1")
	s.db.Exec("DELETE FROM beneficiary_party")
	s.db.Exec("ALTER TABLE beneficiary_party AUTO_INCREMENT = 1")
	s.db.Exec("DELETE FROM charges_information")
	s.db.Exec("ALTER TABLE charges_information AUTO_INCREMENT = 1")
	s.db.Exec("DELETE FROM debtor_party")
	s.db.Exec("ALTER TABLE debtor_party AUTO_INCREMENT = 1")

	// First let's create some data
	s.PrepareData()
}

func (s *StorageSuite) TearDownSuite() {
	// Close the connection after all tests in the suite finish
	s.storage.Close()
	s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StorageSuite)
	suite.Run(t, s)
}

func (s *StorageSuite) TestCreatePayment() {

	// Create a payment through the CreatePayment method
	err := s.storage.CreatePayment(samplePaymentData)
	if err != nil {
		s.T().Fatal(err)
	}

	// Query the database for the entry we just created
	res, err := s.db.Query(`SELECT COUNT(*) FROM payments WHERE uuid='test_uuid'`)

	// Get the count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 2 {
		s.T().Errorf("incorrect count, wanted 2, got %d", count)
	}
}

func (s *StorageSuite) TestUpdatePayment() {
	var err error

	// Firstly check if such row exists
	res, err := s.db.Query(`SELECT COUNT(*) FROM payments WHERE uuid='test_uuid'`)

	// Get the count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}

	// Update a payment through the UpdatePayment method
	err = s.storage.UpdatePayment("test_uuid", samplePaymentData)
	if err != nil {
		s.T().Fatal(err)
	}

	// Let's see if the record got updated
	res, err = s.db.Query(`SELECT COUNT(*) FROM payments`)

	// Get the count result
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StorageSuite) TestDeletePayment() {

	var err error
	s.dbGorm, err = setupGormDB()
	if err != nil {
		s.T().Fatal(err)
	}
	s.storage = &DBHandler{Conn: s.dbGorm, Cfg: testDB.Cfg}

	err = s.storage.DeletePayment("test_uuid")
	if err != nil {
		s.T().Fatal(err)
	}

	res, err := s.db.Query(`SELECT COUNT(*) FROM payments WHERE uuid='test_uuid'`)

	// Get the count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 0 {
		s.T().Errorf("incorrect count, wanted 0, got %d", count)
	}

	res, err = s.db.Query(`SELECT COUNT(*) FROM beneficiary_party WHERE payments_id='1'`)

	// Get the count result
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 0 {
		s.T().Errorf("incorrect count, wanted 0, got %d", count)
	}
}

func (s *StorageSuite) TestGetPayments() {

	var err error

	s.dbGorm, err = setupGormDB()
	if err != nil {
		s.T().Fatal(err)
	}
	s.storage = &DBHandler{Conn: s.dbGorm, Cfg: testDB.Cfg}

	payments, err := s.storage.GetPayments()
	if err != nil {
		s.T().Fatal(err)
	}

	count := len(payments)
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}

	expectedPayment := samplePaymentData
	if *payments[0].UUID != *expectedPayment.UUID {
		s.T().Errorf("incorrect details, expected %v, got %v", *expectedPayment.UUID, *payments[0].UUID)
	}
}

func (s *StorageSuite) PrepareData() {
	var err error
	_, err = s.db.Query(`INSERT INTO payments (payments_type, uuid, payments_version, payments_organisation_id,
		payments_amount, payments_currency, payments_purpose, payments_transaction_type, payments_processing_date)
		VALUES('Payment','test_uuid', '0', 'test_organisation_id2', '500.20', 'GBP', 'test', 'Credit', '2017-01-18')`)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = s.db.Query(`INSERT INTO beneficiary_party (payments_id, account_name, account_number, account_number_code, 
		beneficiary_party_address, bank_id, bank_id_code, beneficiary_party_name)
		VALUES('1', 'test', '0001', '0002', 'tesing address', '0003', '0004', 'blablabla')`)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = s.db.Query(`INSERT INTO charges_information (payments_id, bearer_code, sender_charges_amount, 
		sender_charges_currency, receiver_charges_amount, receiver_charges_currency) 
		VALUES('1', 'SHAR', '100.00', 'EUR', '2.00', 'USD')`)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = s.db.Query(`INSERT INTO debtor_party (payments_id, account_name, account_number, account_number_code, 
		debtor_party_address, bank_id, bank_id_code, debtor_party_name) 
		VALUES('1', 'test', '0001', '0002', 'tesing address', '0003', '0004', 'blablabla')`)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = s.db.Query(`INSERT INTO fx (payments_id, contract_reference, exchange_rate, original_amount, original_currency) 
		VALUES('1', 'FX123', '5.00000', '2.00', 'USD')`)
	if err != nil {
		s.T().Fatal(err)
	}
}
