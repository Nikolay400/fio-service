CREATE TABLE people (
    id SERIAL PRIMARY KEY NOT NULL, 
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100) NOT NULL,
    age SMALLINT NOT NULL,
    gender VARCHAR(100) NOT NULL,
    country VARCHAR(2) NOT NULL
);