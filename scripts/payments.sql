CREATE USER 'rafal'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON *.* TO 'rafal'@'localhost' WITH GRANT OPTION;
CREATE USER 'rafal'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON *.* TO 'rafal'@'%' WITH GRANT OPTION;

ALTER USER 'rafal'@'localhost' IDENTIFIED WITH mysql_native_password BY 'password';
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'password';

FLUSH PRIVILEGES;

DROP DATABASE IF EXISTS payments_demo;
CREATE DATABASE payments_demo;
USE payments_demo;

DROP TABLE IF EXISTS payments;

CREATE TABLE payments (
    payments_id INT AUTO_INCREMENT PRIMARY KEY,
    payments_type VARCHAR(20) NOT NULL DEFAULT 'Payment',
    uuid VARCHAR(50) NOT NULL,
    payments_version INT,
    payments_organisation_id VARCHAR(50) NOT NULL,
    payments_amount DECIMAL(15,2),
    payments_currency VARCHAR(3),
    payments_purpose VARCHAR(255),
    payments_transaction_type VARCHAR(20),
    payments_processing_date DATE
)ENGINE=InnoDB;

DROP TABLE IF EXISTS beneficiary_party;

CREATE TABLE beneficiary_party (
    beneficiary_party_id INT AUTO_INCREMENT PRIMARY KEY,
    payments_id INT,
    account_name VARCHAR(100),
    account_number VARCHAR(100),
    account_number_code VARCHAR(10),
    beneficiary_party_address VARCHAR(255),
    bank_id VARCHAR(50),
    bank_id_code VARCHAR(20),
    beneficiary_party_name VARCHAR(255),
    FOREIGN KEY fk_payments4(payments_id) REFERENCES payments(payments_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
)ENGINE=InnoDB;


DROP TABLE IF EXISTS charges_information;

CREATE TABLE charges_information (
    charges_information_id INT AUTO_INCREMENT PRIMARY KEY,
    payments_id INT,
    bearer_code VARCHAR(10),
    sender_charges_amount DECIMAL(15,2),
    sender_charges_currency VARCHAR(3),
    receiver_charges_amount DECIMAL(15,2),
    receiver_charges_currency VARCHAR(3),
    FOREIGN KEY fk_payments3(payments_id) REFERENCES payments(payments_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
)ENGINE=InnoDB;

DROP TABLE IF EXISTS debtor_party;

CREATE TABLE debtor_party (
    debtor_party_id INT AUTO_INCREMENT PRIMARY KEY,
    payments_id INT,
    account_name VARCHAR(100),
    account_number VARCHAR(100),
    account_number_code VARCHAR(10),
    debtor_party_address VARCHAR(255),
    bank_id VARCHAR(50),
    bank_id_code VARCHAR(20),
    debtor_party_name VARCHAR(255),
    FOREIGN KEY fk_payments2(payments_id) REFERENCES payments(payments_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
)ENGINE=InnoDB;

DROP TABLE IF EXISTS fx;

CREATE TABLE fx (
    fx_id INT AUTO_INCREMENT PRIMARY KEY,
    payments_id INT,
    contract_reference VARCHAR(100),
    exchange_rate DECIMAL(13,5),
    original_amount DECIMAL(15,2),
    original_currency VARCHAR(3),
    FOREIGN KEY fk_payments(payments_id) REFERENCES payments(payments_id) 
        ON UPDATE CASCADE ON DELETE CASCADE
)ENGINE=InnoDB;
