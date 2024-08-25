CREATE TABLE Categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE Videos (
    id serial PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    category_id INT REFERENCES Categories(id)
);