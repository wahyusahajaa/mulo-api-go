CREATE TABLE "artists" (
  "id" serial,
  "name" varchar(100),
  "slug" varchar(100) UNIQUE,
  "image" jsonb,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX ON "artists" ("slug");