-- Users
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  email VARCHAR(190) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  INDEX idx_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Wallets
CREATE TABLE IF NOT EXISTS wallets (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(120) NOT NULL,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  INDEX idx_wallet_user (user_id),
  CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Currencies
CREATE TABLE IF NOT EXISTS currencies (
  code CHAR(3) PRIMARY KEY,
  name VARCHAR(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Wallet Balances
CREATE TABLE IF NOT EXISTS wallet_balances (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  wallet_id BIGINT UNSIGNED NOT NULL,
  currency_code CHAR(3) NOT NULL,
  balance DECIMAL(38,18) NOT NULL DEFAULT 0,
  updated_at DATETIME NULL,
  UNIQUE KEY uniq_wallet_currency (wallet_id, currency_code),
  INDEX idx_wb_wallet (wallet_id),
  INDEX idx_wb_wallet_cur (wallet_id, currency_code),
  CONSTRAINT fk_wb_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_wb_currency FOREIGN KEY (currency_code) REFERENCES currencies(code)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Transactions
CREATE TABLE IF NOT EXISTS transactions (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  wallet_id BIGINT UNSIGNED NOT NULL,
  to_wallet_id BIGINT UNSIGNED NULL,
  type VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL,
  amount DECIMAL(38,18) NOT NULL,
  currency_code CHAR(3) NOT NULL,
  exchange_rate DECIMAL(38,18) NOT NULL DEFAULT 1,
  converted_amount DECIMAL(38,18) NOT NULL DEFAULT 0,
  reference VARCHAR(120) NULL,
  metadata TEXT NULL,
  created_at DATETIME NULL,
  INDEX idx_trx_wallet (wallet_id),
  INDEX idx_trx_type (type),
  INDEX idx_trx_currency (currency_code),
  INDEX idx_trx_status (status),
  INDEX idx_trx_wallet_created (wallet_id, created_at),
  INDEX idx_trx_reference (reference),
  CONSTRAINT fk_trx_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO currencies (code, name) VALUES ('USD','US Dollar'), ('EUR','Euro'), ('JPY','Japanese Yen');
