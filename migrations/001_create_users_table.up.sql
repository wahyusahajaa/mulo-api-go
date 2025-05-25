CREATE SCHEMA "u";

CREATE TYPE "u"."role" AS ENUM (
  'admin',
  'member'
);

CREATE TABLE "users" (
  "id" serial,
  "full_name" varchar(100),
  "username" varchar(100) UNIQUE,
  "email" varchar(100) UNIQUE,
  "password" varchar(100),
  "image" jsonb,
  "role" u.role,
  "email_verified_at" timestamp,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX ON "users" ("username");
CREATE UNIQUE INDEX ON "users" ("email");