-- Migration: 001_create_users
-- Description: Create the user schema

CREATE TABLE IF NOT EXISTS users (
    id      SERIAL  PRIMARY KEY,
    name    TEXT    NOT NULL,
    dob     DATE    NOT NULL 
);