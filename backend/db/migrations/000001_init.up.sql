CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ══════════════ ref address ══════════════
CREATE TABLE ref_provinces (
  code    varchar(2)   PRIMARY KEY,
  name_th varchar(100) NOT NULL
);

CREATE TABLE ref_districts (
  code          varchar(4)   PRIMARY KEY,
  province_code varchar(2)   NOT NULL REFERENCES ref_provinces(code),
  name_th       varchar(100) NOT NULL,
  kind          varchar(10)  NOT NULL DEFAULT 'อำเภอ'  
);

CREATE TABLE ref_subdistricts (
  code          varchar(6)   PRIMARY KEY,
  district_code varchar(4)   NOT NULL REFERENCES ref_districts(code),
  name_th       varchar(100) NOT NULL,
  kind          varchar(10)  NOT NULL DEFAULT 'ตำบล'  
);

-- ตำบล 1 แห่งมีได้หลายรหัสไปรษณีย์ และรหัสเดียวครอบหลายตำบล
CREATE TABLE ref_subdistrict_postal (
  subdistrict_code varchar(6) NOT NULL REFERENCES ref_subdistricts(code),
  postal_code      varchar(5) NOT NULL,
  is_primary       boolean    NOT NULL DEFAULT false,
  PRIMARY KEY (subdistrict_code, postal_code)
);

CREATE INDEX idx_district_province    ON ref_districts(province_code);
CREATE INDEX idx_subdistrict_district ON ref_subdistricts(district_code);

-- ══════════════ org / user / client ══════════════
CREATE TABLE organizations (
  id            uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  code          varchar(20)  NOT NULL UNIQUE,
  name_th       varchar(200) NOT NULL,
  org_type      varchar(30)  NOT NULL,
    -- BANK_BRANCH | DISTRICT_OFFICE | TREASURY_PROVINCIAL
    -- | TREASURY_DISTRICT | SPECIAL_AREA | CENTRAL
  parent_id     uuid         REFERENCES organizations(id),
  province_code varchar(2),        
  district_code varchar(4),
  scope_policy  varchar(20)  NOT NULL DEFAULT 'OWN_UNIT',
    -- OWN_UNIT | DISTRICT | PROVINCE | NATIONWIDE
  is_active     boolean      NOT NULL DEFAULT true,
  created_at    timestamptz  NOT NULL DEFAULT now()
);

CREATE TABLE users (
  id                uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  username          varchar(100) NOT NULL UNIQUE,
  password_hash     varchar(100),           
  role              varchar(20)  NOT NULL,
  organization_id   uuid         NOT NULL REFERENCES organizations(id),
  identity_provider varchar(20)  NOT NULL DEFAULT 'LOCAL', 
  is_active         boolean      NOT NULL DEFAULT true,
  created_at        timestamptz  NOT NULL DEFAULT now()
);

-- ระบบพาร์ตเนอร์ที่ยิง API เข้ามา (MVP1 ยังไม่ใช้ แต่ schema พร้อม)
CREATE TABLE api_clients (
  id              uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id       varchar(64)  NOT NULL UNIQUE,
  name_th         varchar(200) NOT NULL,
  organization_id uuid         NOT NULL REFERENCES organizations(id),
  secret_hash     varchar(100),
  is_active       boolean      NOT NULL DEFAULT true,
  created_at      timestamptz  NOT NULL DEFAULT now()
);

-- ══════════════ citizen ══════════════
CREATE TABLE addresses (
  id               uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  house_no         varchar(50)  NOT NULL,
  moo              varchar(10),
  road             varchar(100),
  subdistrict_code varchar(6)   NOT NULL REFERENCES ref_subdistricts(code),
  district_code    varchar(4)   NOT NULL REFERENCES ref_districts(code),
  province_code    varchar(2)   NOT NULL REFERENCES ref_provinces(code),
  postal_code      varchar(5)   NOT NULL,
  subdistrict_name varchar(100) NOT NULL,
  district_name    varchar(100) NOT NULL,
  province_name    varchar(100) NOT NULL
);

CREATE TABLE citizens (
  id               uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  national_id_hash varchar(64)  NOT NULL UNIQUE,   
  national_id_enc  bytea        NOT NULL,
  title            varchar(20)  NOT NULL,
  first_name       varchar(100) NOT NULL,
  last_name        varchar(100) NOT NULL,
  birth_year       int          NOT NULL,         
  birth_month      int,                           
  birth_day        int,                           
  birth_precision  varchar(12)  NOT NULL,          
  phone            varchar(20)  NOT NULL,
  address_id       uuid         NOT NULL REFERENCES addresses(id),
  created_at       timestamptz  NOT NULL DEFAULT now(),
  updated_at       timestamptz  NOT NULL DEFAULT now()
);

CREATE TABLE identity_verifications (
  id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id    uuid        NOT NULL REFERENCES citizens(id),
  method        varchar(24) NOT NULL,
  verified      boolean     NOT NULL,
  note          text,                              
  verified_by   uuid        REFERENCES users(id),  
  provider_code varchar(20),
  reference_no  varchar(50),
  verified_at   timestamptz NOT NULL DEFAULT now()
);

-- ══════════════ application ══════════════
CREATE TABLE applications (
  id                   uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  application_no       varchar(20) NOT NULL UNIQUE,
  citizen_id           uuid        NOT NULL REFERENCES citizens(id),
  fiscal_year          int         NOT NULL,
  status               varchar(24) NOT NULL DEFAULT 'SUBMITTED',
  is_farmer            boolean     NOT NULL DEFAULT false,
  marital_status       varchar(12),
  expense_to_others    bigint      NOT NULL DEFAULT 0,
  has_credit_card      boolean     NOT NULL DEFAULT false,
  applicant_snapshot   jsonb       NOT NULL,       
  registration_unit_id uuid        NOT NULL REFERENCES organizations(id),
  created_by_user_id   uuid        REFERENCES users(id),        
  created_by_client_id uuid        REFERENCES api_clients(id),  
  submission_channel   varchar(20) NOT NULL,
  submitted_at         timestamptz NOT NULL DEFAULT now(),
  created_at           timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT chk_app_creator CHECK (
    (created_by_user_id IS NOT NULL) <> (created_by_client_id IS NOT NULL)
  )
);

-- ★ ชั้นกันซ้ำที่ 2 — 1 คน 1 ใบต่อปีงบประมาณ
CREATE UNIQUE INDEX uq_app_active_per_year
  ON applications (citizen_id, fiscal_year)
  WHERE status <> 'CANCELLED';

CREATE INDEX idx_app_unit ON applications(registration_unit_id, submitted_at DESC);

CREATE SEQUENCE app_no_seq START 1;

CREATE TABLE application_status_history (
  id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id uuid        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  from_status    varchar(24),
  to_status      varchar(24) NOT NULL,
  actor_type     varchar(10) NOT NULL,  
  actor_id       uuid        NOT NULL,  
  actor_role     varchar(20) NOT NULL,
  reason         text,
  created_at     timestamptz NOT NULL DEFAULT now()
);

-- ══════════════ children ══════════════
CREATE TABLE household_members (
  id               uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id   uuid         NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  relation         varchar(12)  NOT NULL,
  national_id_hash varchar(64),
  full_name        varchar(200) NOT NULL,
  birth_year       int,
  annual_income    bigint CHECK (annual_income >= 0)
);

CREATE TABLE income_sources (
  id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id uuid        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  source_type    varchar(12) NOT NULL,
  annual_amount  bigint      NOT NULL CHECK (annual_amount >= 0)
);

CREATE TABLE assets (
  id                    uuid          PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id        uuid          NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  asset_type            varchar(24)   NOT NULL,
  amount                numeric(14,2) NOT NULL CHECK (amount >= 0),
  unit                  varchar(8)    NOT NULL,
  joint_account_holders int           NOT NULL DEFAULT 1 CHECK (joint_account_holders >= 1),
  is_minor_account      boolean       NOT NULL DEFAULT false
);

CREATE TABLE liabilities (
  id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  application_id uuid        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
  liability_type varchar(16) NOT NULL,
  credit_limit   bigint      NOT NULL CHECK (credit_limit >= 0)
);

CREATE INDEX idx_hh_app     ON household_members(application_id);
CREATE INDEX idx_income_app ON income_sources(application_id);
CREATE INDEX idx_asset_app  ON assets(application_id);
CREATE INDEX idx_liab_app   ON liabilities(application_id);
