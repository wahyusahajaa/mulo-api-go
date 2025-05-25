CREATE TABLE "song_genres" (
  "song_id" int NOT NULL,
  "genre_id" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("song_id", "genre_id")
);

CREATE INDEX ON "song_genres" ("song_id");

CREATE INDEX ON "song_genres" ("genre_id");

ALTER TABLE "song_genres" ADD FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "song_genres" ADD FOREIGN KEY ("genre_id") REFERENCES "genres" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;