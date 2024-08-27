CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE videos (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL,
    description TEXT,
    category_id INT REFERENCES categories(id)
);