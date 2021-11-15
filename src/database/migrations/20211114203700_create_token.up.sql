CREATE TABLE tokens (
    id UUID PRIMARY KEY,
    valid BOOLEAN DEFAULT true NOT NULL, 
    refresh_token UUID NOT NULL, 
    email VARCHAR(255) REFERENCES users(email) NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL 
);

CREATE  INDEX valid_index
ON tokens (email, valid) ;

CREATE INDEX refresh_tokens 
ON tokens (email, refresh_token);
 