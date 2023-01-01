-- Active: 1672157646835@@127.0.0.1@5432@block@public;

CREATE TABLE users (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name    VARCHAR(250),
    email   VARCHAR(250),
    image   VARCHAR(1000)
);

CREATE TABLE folder (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name            VARCHAR,
    addres          BYTEA,
    -- 0 private
    -- 1 public
    access          BIT(1),
);

CREATE TABLE file (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    time_stamp      DATE,
    folder_addres   BYTEA,
    file_addres     BYTEA,
    hash            BYTEA,
    provHash        BYTEA,
    
    -- 0 private
    -- 1 public
    access          BIT(1),

    title           BYTEA,
    data            BYTEA
);

CREATE TABLE history (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    hash        BYTEA,
    provHash    BYTEA,
    

    user        INT,
    title       VARCHAR,
    description VARCHAR,
)