CREATE TABLE accounts (
                          id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
                          login VARCHAR UNIQUE NOT NULL ,
                          hash_pass VARCHAR NOT NULL,
                          points INT NOT NULL DEFAULT 0
);

CREATE TABLE orders (
                        id BIGINT PRIMARY KEY,
                        account_id UUID REFERENCES accounts (id),
                        status INT NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        points INT,

                        CONSTRAINT uc_order_account UNIQUE (id, account_id)
);