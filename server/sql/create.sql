-- Active: 1672157646835@@127.0.0.1@5432@block@public;

CREATE TABLE users (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name    VARCHAR(250),
    email   VARCHAR(250),
    image   VARCHAR(1000)
);

CREATE TABLE folder (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    address         BYTEA,
    name            VARCHAR,
    userID          INT,
    file            BYTEA,
    -- 0 private
    -- 1 public
    access          BOOLEAN
);

CREATE TABLE file (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    timeStamp       INT,
    hash            BYTEA,
    prevHash        BYTEA,
    -- 0 private
    -- 1 public
    access          BOOLEAN,
    
    title           BYTEA,
    type            BYTEA,
    file            BYTEA
);

CREATE TABLE history (
    id int NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    hash        BYTEA,
    provHash    BYTEA,
    

    user        INT,
    title       VARCHAR,
    description VARCHAR,
)