-- +goose Up

-- users table for authentication
CREATE TABLE users (
    id            BIGSERIAL      PRIMARY KEY,
    username      VARCHAR(50)    NOT NULL UNIQUE,
    email         VARCHAR(100)   NOT NULL UNIQUE,
    password_hash VARCHAR(255)   NOT NULL,
    is_active     BOOLEAN        DEFAULT true,
    created_at    TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);

-- members table
CREATE TABLE members (
    id            BIGSERIAL      PRIMARY KEY,
    full_name     VARCHAR(100)   NOT NULL,
    is_active     BOOLEAN        DEFAULT true,
    benefit_limit DECIMAL(10, 2) NOT NULL,
    used_amount   DECIMAL(10, 2) DEFAULT 0.00,
    created_at    TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_members_is_active ON members (is_active);

-- procedures table
CREATE TABLE procedures (
    id           BIGSERIAL      PRIMARY KEY,
    code         VARCHAR(20)    NOT NULL UNIQUE,  -- UNIQUE required for FK reference from claims
    description  TEXT,
    average_cost DECIMAL(10, 2) NOT NULL,
    created_at   TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_procedures_code ON procedures (code);
CREATE INDEX idx_procedures_average_cost ON procedures (average_cost);

-- providers table
CREATE TABLE providers (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    location   VARCHAR(100),
    created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_providers_name ON providers (name);
CREATE INDEX idx_providers_location ON providers (location);

-- claims table
CREATE TYPE CLAIM_STATUS AS ENUM ('APPROVED', 'PARTIAL', 'REJECTED');

CREATE TABLE claims (
    id               BIGSERIAL      PRIMARY KEY,
    member_id        BIGINT         REFERENCES members(id),
    provider_id      BIGINT         REFERENCES providers(id),
    procedure_code   VARCHAR(20)    REFERENCES procedures(code),
    diagnosis_code   VARCHAR(20),
    requested_amount DECIMAL(10, 2) NOT NULL,
    approved_amount  DECIMAL(10, 2) DEFAULT 0.00,
    status           CLAIM_STATUS   NOT NULL,
    fraud_flag       BOOLEAN        DEFAULT false,
    rejection_reason TEXT,
    created_at       TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_claims_member_id      ON claims (member_id);
CREATE INDEX idx_claims_provider_id    ON claims (provider_id);
CREATE INDEX idx_claims_procedure_code ON claims (procedure_code);
CREATE INDEX idx_claims_diagnosis_code ON claims (diagnosis_code);
CREATE INDEX idx_claims_status         ON claims (status);

-- +goose Down

DROP INDEX IF EXISTS idx_claims_status;
DROP INDEX IF EXISTS idx_claims_diagnosis_code;
DROP INDEX IF EXISTS idx_claims_procedure_code;
DROP INDEX IF EXISTS idx_claims_provider_id;
DROP INDEX IF EXISTS idx_claims_member_id;
DROP TABLE IF EXISTS claims;
DROP TYPE  IF EXISTS CLAIM_STATUS;

DROP INDEX IF EXISTS idx_providers_location;
DROP INDEX IF EXISTS idx_providers_name;
DROP TABLE IF EXISTS providers;

DROP INDEX IF EXISTS idx_procedures_average_cost;
DROP INDEX IF EXISTS idx_procedures_code;
DROP TABLE IF EXISTS procedures;

DROP INDEX IF EXISTS idx_members_is_active;
DROP TABLE IF EXISTS members;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;