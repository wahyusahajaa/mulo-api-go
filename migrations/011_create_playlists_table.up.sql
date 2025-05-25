CREATE TABLE "playlists" (
  "id" serial,
  "user_id" int NOT NULL,
  "name" varchar(100),
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

ALTER TABLE "playlists" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;