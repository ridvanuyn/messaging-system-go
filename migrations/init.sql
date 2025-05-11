CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "to" VARCHAR(20) NOT NULL,
    content VARCHAR(500) NOT NULL,
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP,
    message_id VARCHAR(100)
);

-- Test data
INSERT INTO messages ("to", content) VALUES 
    ('+905551111111', 'Insider - Project Test Message 1'),
    ('+905552222222', 'Insider - Project Test Message 2'),
    ('+905553333333', 'Insider - Project Test Message 3'),
    ('+905554444444', 'Insider - Project Test Message 4');
