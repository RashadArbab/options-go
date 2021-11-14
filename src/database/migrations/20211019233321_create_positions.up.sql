CREATE TABLE positions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id)  NOT NULL,
    ticker VARCHAR(55), 
    amount INT
)