-- table create karne ke liye db me, agar na ho toh

CREATE TABLE IF NOT EXISTS test_data (
    id TEXT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);