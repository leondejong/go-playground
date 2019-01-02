-- psql

DROP TABLE "item";
DROP DATABASE "list";
DROP USER "user";
CREATE DATABASE "list";
CREATE ROLE "user" WITH LOGIN PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE "list" TO "user";

-- psql list -U user

CREATE TABLE "item" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(255) NOT NULL,
  "content" text NOT NULL,
  "date" timestamp DEFAULT NOW() NOT NULL
);

INSERT INTO "item" ("name", "content") VALUES
('name01', 'content01'),
('name02', 'content02'),
('name03', 'content03');

SELECT * FROM "item";
