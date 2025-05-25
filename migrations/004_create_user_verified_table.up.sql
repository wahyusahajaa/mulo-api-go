CREATE TABLE "user_verified" (
  "id" serial,
  "user_id" int NOT NULL,
  "code" varchar(5) UNIQUE,
  "expired_at" timestamp DEFAULT (now() + interval '5 minutes'),
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);


CREATE INDEX ON "user_verified" ("user_id");

CREATE UNIQUE INDEX ON "user_verified" ("code");

ALTER TABLE "user_verified" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;