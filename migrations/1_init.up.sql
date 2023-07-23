CREATE TABLE accounts (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    login VARCHAR UNIQUE NOT NULL ,
    hash_pass VARCHAR NOT NULL,
    points DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE TYPE order_status AS enum ('NEW','PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id VARCHAR PRIMARY KEY,
    account_id UUID REFERENCES accounts (id),
    status order_status NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    points DOUBLE PRECISION,

    CONSTRAINT uc_order_account UNIQUE (id, account_id)
);

CREATE TABLE withdraws (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    order_id VARCHAR NOT NULL,
    points DOUBLE PRECISION NOT NULL,
    account_id UUID NOT NULL,
    processed_at TIMESTAMP NOT NULL,

    CONSTRAINT withdraws_unique_order_id UNIQUE (order_id),
    CONSTRAINT withdraws_fk_account foreign key (account_id) references accounts (id)
)
