package storage

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	// use gotm dialect
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Config stores the db configuration.
type Config struct {
	User             string
	Password         string
	Port             string
	Db               string
	Host             string
	AdditionalParams string
}

// DBHandler stores a db handler.
type DBHandler struct {
	Conn *gorm.DB
	Cfg  *Config
}

// Connect establishes a connection to the database.
func (db *DBHandler) Connect() (err error) {

	user := db.Cfg.User
	password := db.Cfg.Password
	host := db.Cfg.Host
	port := db.Cfg.Port
	dbName := db.Cfg.Db
	params := db.Cfg.AdditionalParams

	db.Conn, err = gorm.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+dbName+params)

	return err
}

// Close closes current db connection.
func (db *DBHandler) Close() error {
	return db.Conn.Close()
}

// Ping verifies a connection to the database is still alive establishing a connection if necessary.
func (db *DBHandler) Ping() error {
	return db.Conn.DB().Ping()
}

// SetLogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs.
func (db *DBHandler) SetLogMode(mode bool) {
	db.Conn.LogMode(mode)
}

// SetLogger replace default logger.
func (db *DBHandler) SetLogger(logger gorm.Logger) {
	db.Conn.SetLogger(logger)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (db *DBHandler) SetMaxOpenConns(maxOpened int) {
	db.Conn.DB().SetMaxOpenConns(maxOpened)
}

// SetMaxIdleConns sets the maximum number of connections in the idle pool.
func (db *DBHandler) SetMaxIdleConns(maxIddled int) {
	db.Conn.DB().SetMaxIdleConns(maxIddled)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (db *DBHandler) SetConnMaxLifetime(seconds time.Duration) {
	db.Conn.DB().SetConnMaxLifetime(time.Second * time.Duration(seconds))
}

// --------------------------------------------------------------------------------
// GORM requires functions that translate StructureNames to the actual table names
// --------------------------------------------------------------------------------

// TableName for beneficiary_party
func (BeneficiaryParty) TableName() string {
	return "beneficiary_party"
}

// TableName for charges_information
func (ChargesInformation) TableName() string {
	return "charges_information"
}

// TableName for debtor_party
func (DebtorParty) TableName() string {
	return "debtor_party"
}

// TableName for fx
func (Fx) TableName() string {
	return "fx"
}

// TableName for payments
func (Payments) TableName() string {
	return "payments"
}

// --------------------------------------------------------------------------------

// GetPayments returns a list of payments
func (db *DBHandler) GetPayments() ([]*Payments, error) {

	var payments []*Payments

	err := db.Conn.Find(&payments).Error
	if err != nil {
		return nil, err
	}

	for i, k := range payments {
		payments[i], err = db.GetPayment(*k.UUID)
		if err != nil {
			return nil, err
		}
	}

	return payments, nil
}

// GetPayment finds a single payment record
func (db *DBHandler) GetPayment(uuid string) (*Payments, error) {

	var payments Payments

	err := db.Conn.Where("uuid = ?", uuid).First(&payments).Error
	if err != nil {
		return nil, err
	}

	err = db.Conn.Where("payments_id = ?", payments.PaymentsID).First(&payments.DebtorParty).Error
	if err != nil {
		return nil, err
	}

	err = db.Conn.Where("payments_id = ?", payments.PaymentsID).First(&payments.ChargesInformation).Error
	if err != nil {
		return nil, err
	}

	err = db.Conn.Where("payments_id = ?", payments.PaymentsID).First(&payments.BeneficiaryParty).Error
	if err != nil {
		return nil, err
	}

	err = db.Conn.Where("payments_id = ?", payments.PaymentsID).First(&payments.Fx).Error
	if err != nil {
		return nil, err
	}

	return &payments, nil
}

// UpdatePayment updates a single payment record
func (db *DBHandler) UpdatePayment(uuid string, p *Payments) error {

	var payments Payments

	handler := db.Conn.Where("uuid = ?", uuid).First(&payments)
	if handler.RecordNotFound() {
		return errors.New("payment not found")
	}

	tx := db.Conn.Begin()
	defer tx.Close()

	var b BeneficiaryParty
	var c ChargesInformation
	var d DebtorParty
	var f Fx

	b = p.BeneficiaryParty
	c = p.ChargesInformation
	d = p.DebtorParty
	f = p.Fx

	err := tx.Model(&p).Updates(&p).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&d).Where("payments_id = ?", payments.PaymentsID).Update(&d).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&c).Where("payments_id = ?", payments.PaymentsID).Update(&c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&f).Where("payments_id = ?", payments.PaymentsID).Update(&f).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&b).Where("payments_id = ?", payments.PaymentsID).Update(&b).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// DeletePayment removes a single payment record
func (db *DBHandler) DeletePayment(uuid string) error {

	handler := db.Conn.Where("uuid = ?", uuid).Delete(Payments{})
	if handler.RowsAffected == 0 {
		return errors.New("0 rows affected")
	}
	return handler.Error
}

// CreatePayment adds a new payment record
func (db *DBHandler) CreatePayment(p *Payments) (err error) {

	var b BeneficiaryParty
	var c ChargesInformation
	var d DebtorParty
	var f Fx

	b = p.BeneficiaryParty
	c = p.ChargesInformation
	d = p.DebtorParty
	f = p.Fx

	tx := db.Conn.Begin()
	defer tx.Close()

	err = tx.Model(&p).Create(&p).Error
	if err != nil {
		return err
	}

	b.PaymentsID = p.PaymentsID
	err = tx.Model(&b).Create(&b).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	c.PaymentsID = p.PaymentsID
	err = tx.Model(&c).Create(&c).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	d.PaymentsID = p.PaymentsID
	err = tx.Model(&d).Create(&d).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	f.PaymentsID = p.PaymentsID
	err = tx.Model(&f).Create(&f).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// DBConn is a helper function returning a pointer to the current DBHandler struct
func (db *DBHandler) DBConn() *DBHandler {
	return &DBHandler{db.Conn, db.Cfg}
}

// InitMysql ...
func InitMysql(c Connector) Connector {
	return c
}
