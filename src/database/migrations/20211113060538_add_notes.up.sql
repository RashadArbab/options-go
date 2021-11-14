CREATE TABLE notes (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    ticker_id UUID REFERENCES positions(id) NOT NULL, 
    title VARCHAR(55), 
    body VARCHAR(1024)
);