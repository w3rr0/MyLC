-- Schema for users table with ENUM group_type

-- Create ENUM
CREATE TYPE group_type AS ENUM
    ( 'IT', 'HR', 'PR', 'IO', 'JFR', 'Grafika' );

-- Create Table
CREATE TABLE users
    ( id SERIAL PRIMARY KEY, first_name TEXT NOT NULL, last_name TEXT NOT NULL,
      email TEXT NOT NULL, "group" group_type NOT NULL );

-- Example Insert
INSERT INTO users
    (first_name, last_name, email, "group")
VALUES
    ('Jan', 'Kowalski', 'jan.kowalski@iaeste.pl', 'IT'),
    ('Alicja', 'Nowak', 'alicja.nowak@iaeste.pl', 'Grafika');


-- Schema for user login and password

-- Create Table
CREATE TABLE accounts
( email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  verification_token VARCHAR(64),
  is_verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schema for event_manager and event tables

-- Create event_manager
CREATE TABLE event_manager (id SERIAL PRIMARY KEY, name TEXT NOT NULL DEFAULT 'Untitled Event');

-- Add first event (create event table and update event_manager)
INSERT INTO event_manager DEFAULT VALUES;

CREATE TABLE table_1 (name TEXT NOT NULL, "13:00-13:30" INT[], "13:30-14:00" INT[]);

INSERT INTO table_1
    (name, "13:00-13:30", "13:30-14:00")
VALUES
    ('yes', '{}', '{}'),
    ('maybe', '{}', '{}'),
    ('no', '{}', '{}');

-- Update table
UPDATE table_1
SET "13:00-13:30" = '{26}'
WHERE name='yes';