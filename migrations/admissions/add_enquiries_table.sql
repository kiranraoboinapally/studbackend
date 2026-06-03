-- Table: admissions.enquiries

-- DROP TABLE IF EXISTS admissions.enquiries;

-- Create sequence if it doesn't exist
CREATE SEQUENCE IF NOT EXISTS admissions.enquiries_id_seq;

CREATE TABLE IF NOT EXISTS admissions.enquiries
(
    id bigint NOT NULL DEFAULT nextval('admissions.enquiries_id_seq'),
    full_name text COLLATE pg_catalog."default" NOT NULL,
    mobile_number text COLLATE pg_catalog."default" NOT NULL,
    email_address text COLLATE pg_catalog."default" NOT NULL,
    country text COLLATE pg_catalog."default",
    state text COLLATE pg_catalog."default",
    district text COLLATE pg_catalog."default",
    preferred_campus bigint,
    qualification_type text COLLATE pg_catalog."default",
    program_id bigint,
    status text COLLATE pg_catalog."default" DEFAULT 'pending'::text,
    otp_verified boolean DEFAULT false,
    otp_token text,
    otp_sent_at timestamp with time zone,
    otp_expires_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT enquiries_pkey PRIMARY KEY (id),
    CONSTRAINT enquiries_mobile_number_key UNIQUE (mobile_number),
    CONSTRAINT enquiries_email_address_key UNIQUE (email_address)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS admissions.enquiries
    OWNER to postgres;

-- Index: idx_admissions_enquiries_status

-- DROP INDEX IF EXISTS admissions.idx_admissions_enquiries_status;

CREATE INDEX IF NOT EXISTS idx_admissions_enquiries_status
    ON admissions.enquiries USING btree
    (status COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Index: idx_admissions_enquiries_created_at

-- DROP INDEX IF EXISTS admissions.idx_admissions_enquiries_created_at;

CREATE INDEX IF NOT EXISTS idx_admissions_enquiries_created_at
    ON admissions.enquiries USING btree
    (created_at DESC NULLS LAST)
    TABLESPACE pg_default;
