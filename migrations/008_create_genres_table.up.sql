CREATE TABLE "genres" (
  "id" serial,
  "name" varchar(100),
  "image" jsonb,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);