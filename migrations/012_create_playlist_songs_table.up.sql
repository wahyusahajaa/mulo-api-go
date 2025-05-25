CREATE TABLE "playlist_songs" (
  "playlist_id" int NOT NULL,
  "song_id" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("playlist_id", "song_id")
);

CREATE INDEX ON "playlist_songs" ("playlist_id");

CREATE INDEX ON "playlist_songs" ("song_id");

ALTER TABLE "playlist_songs" ADD FOREIGN KEY ("playlist_id") REFERENCES "playlists" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "playlist_songs" ADD FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;