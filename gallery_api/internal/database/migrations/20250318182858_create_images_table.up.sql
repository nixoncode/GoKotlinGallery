CREATE TABLE images
(
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP
    WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);