CREATE TABLE "song_favorites" (
  "user_id" int NOT NULL,
  "song_id" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("user_id", "song_id")
);


CREATE INDEX ON "song_favorites" ("user_id");

CREATE INDEX ON "song_favorites" ("song_id");

ALTER TABLE "song_favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "song_favorites" ADD FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;