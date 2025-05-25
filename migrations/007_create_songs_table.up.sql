CREATE TABLE "songs" (
  "id" serial,
  "album_id" int NOT NULL,
  "audio" varchar(255),
  "title" varchar(100),
  "duration" int,
  "image" jsonb,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE INDEX ON "songs" ("album_id");

ALTER TABLE "songs" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;