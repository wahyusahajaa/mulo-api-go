CREATE TABLE "artist_genres" (
  "artist_id" int NOT NULL,
  "genre_id" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("artist_id", "genre_id")
);

CREATE INDEX ON "artist_genres" ("artist_id");

CREATE INDEX ON "artist_genres" ("genre_id");

ALTER TABLE "artist_genres" ADD FOREIGN KEY ("artist_id") REFERENCES "artists" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "artist_genres" ADD FOREIGN KEY ("genre_id") REFERENCES "genres" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;