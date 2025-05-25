CREATE TABLE "albums" (
  "id" serial,
  "artist_id" int NOT NULL,
  "name" varchar(100),
  "slug" varchar(100) UNIQUE,
  "image" jsonb,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE INDEX ON "albums" ("artist_id");

CREATE UNIQUE INDEX ON "albums" ("slug");

ALTER TABLE "albums" ADD FOREIGN KEY ("artist_id") REFERENCES "artists" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;