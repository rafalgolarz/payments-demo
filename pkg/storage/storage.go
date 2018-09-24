package storage

// Connector lists the methods that matter when connecting to the data layer.
type Connector interface {
	Connect() (err error)
	Close() error
	GetPayments() ([]*Payments, error)
	GetPayment(uuid string) (*Payments, error)
	UpdatePayment(uuid string, p *Payments) error
	DeletePayment(uuid string) error
	CreatePayment(p *Payments) (err error)
}

// Apart json tags, it also has gorm tags used by Gorm package

// BeneficiaryParty identifies the beneficient of transaction
type BeneficiaryParty struct {
	BeneficiaryPartyID      uint    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	PaymentsID              uint    `gorm:"ForeignKey:PaymentsID" json:"-"`
	AccountName             *string `gorm:"type:varchar(100)" json:"account_name"`
	AccountNumber           *string `gorm:"type:varchar(100)" json:"account_number"`
	AccountNumberCode       *string `gorm:"type:varchar(10)" json:"account_number_code"`
	BeneficiaryPartyAddress *string `gorm:"type:varchar(255)" json:"address"`
	BankID                  *string `gorm:"type:varchar(50)" json:"bank_id"`
	BankIDCode              *string `gorm:"type:varchar(20)" json:"bank_id_code"`
	BeneficiaryPartyName    *string `gorm:"type:varchar(20)" json:"name"`
}

// ChargesInformation shows who and how much was charged for transaction
type ChargesInformation struct {
	ChargesInformationID    uint    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	PaymentsID              uint    `gorm:"ForeignKey:PaymentsID" json:"-"`
	BearerCode              *string `gorm:"type:varchar(10)" json:"bearer_code"`
	SenderChargesAmount     *string `gorm:"type:decimal(15,2)" json:"sender_charges_amount"`
	SenderChargesCurrency   *string `gorm:"type:varchar(3)" json:"sender_charges_currency"`
	ReceiverChargesAmount   *string `gorm:"type:decimal(15,2)" json:"receiver_charges_amount"`
	ReceiverChargesCurrency *string `gorm:"type:varchar(3)" json:"receiver_charges_currency"`
}

// DebtorParty identifies the debtor's side of transaction
type DebtorParty struct {
	DebtorPartyID      uint    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	PaymentsID         uint    `gorm:"ForeignKey:PaymentsID" json:"-"`
	AccountName        *string `gorm:"type:varchar(100)" json:"account_name"`
	AccountNumber      *string `gorm:"type:varchar(100)" json:"account_number"`
	AccountNumberCode  *string `gorm:"type:varchar(10)" json:"account_number_code"`
	DebtorPartyAddress *string `gorm:"type:varchar(255)" json:"address"`
	BankID             *string `gorm:"type:varchar(50)" json:"bank_id"`
	BankIDCode         *string `gorm:"type:varchar(20)" json:"bank_id_code"`
	DebtorPartyName    *string `gorm:"type:varchar(20)" json:"name"`
}

// Fx is all about currency exchange
type Fx struct {
	FxID              uint    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	PaymentsID        uint    `gorm:"ForeignKey:PaymentsID" json:"-"`
	ContractReference *string `gorm:"type:varchar(100)" json:"contract_reference"`
	ExchangeRate      *string `gorm:"type:decimal(13,5)" json:"exchange_rate"`
	OriginalAmount    *string `gorm:"type:decimal(15,2)" json:"original_amount"`
	OriginalCurrency  *string `gorm:"type:varchar(3)" json:"original_currency"`
}

// Payments is the main structure holding all data about transaction
type Payments struct {
	PaymentsID              uint               `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	PaymentsType            *string            `gorm:"type:varchar(20)" json:"type"`
	UUID                    *string            `gorm:"type:varchar(50)" json:"id"`
	PaymentsVersion         *int               `gorm:"type:int" json:"version"`
	PaymentsOrganisationID  *string            `gorm:"type:varchar(50)" json:"organisation_id"`
	PaymentsAmount          *string            `gorm:"type:decimal(15,2)" json:"amount"`
	BeneficiaryParty        BeneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation      ChargesInformation `json:"charges_information"`
	PaymentsCurrency        *string            `gorm:"type:varchar(3)" json:"currency"`
	DebtorParty             DebtorParty        `json:"debtor_party"`
	Fx                      Fx                 `json:"fx"`
	PaymentsPurpose         *string            `gorm:"type:varchar(50)" json:"payment_purpose"`
	PaymentsTransactionType *string            `gorm:"type:varchar(50)" json:"payment_type"`
	PaymentsProcessingDate  *string            `gorm:"type:date" json:"processing_date"`
}
