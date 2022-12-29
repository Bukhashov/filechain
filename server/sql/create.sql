-- Active: 1672157646835@@127.0.0.1@5432@block@public;

CREATE TABLE users (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name    VARCHAR(250),
    email   VARCHAR(250),
    image   VARCHAR(1000)
);