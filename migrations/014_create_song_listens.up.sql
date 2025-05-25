CREATE TABLE "song_listens" (
  "id" serial,
  "user_id" int NOT NULL,
  "song_id" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE INDEX ON "song_listens" ("user_id");

CREATE INDEX ON "song_listens" ("song_id");

ALTER TABLE "song_listens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "song_listens" ADD FOREIGN KEY ("song_id") REFERENCES "songs" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
