CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    phone TEXT,
    password TEXT
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    price DOUBLE PRECISION,
    status TEXT,
    seller_id INT REFERENCES users(id)
);