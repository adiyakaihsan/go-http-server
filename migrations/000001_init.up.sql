CREATE TABLE public.users (
    id serial PRIMARY KEY,
    username varchar(50) NOT NULL UNIQUE,
    password varchar(100) NOT NULL
);