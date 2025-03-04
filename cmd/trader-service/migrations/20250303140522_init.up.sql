

CREATE TABLE quotes (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    quantity DECIMAL NOT NULL,
    price DECIMAL NOT NULL,
    side VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE trades (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    quote_id UUID NOT NULL REFERENCES quotes(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);