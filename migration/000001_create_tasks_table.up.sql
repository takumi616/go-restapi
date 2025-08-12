CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(30) NOT NULL UNIQUE,
    description TEXT,
    status BOOLEAN NOT NULL
);