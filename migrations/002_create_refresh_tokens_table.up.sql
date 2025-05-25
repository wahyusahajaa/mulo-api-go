CREATE TABLE "refresh_tokens" (
  "id" serial,
  "user_id" int NOT NULL,
  "token" TEXT UNIQUE NOT NULL,
  "revoked" bool DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "revoked_at" timestamp,
  "expires_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);

CREATE INDEX ON "refresh_tokens" ("user_id");
CREATE UNIQUE INDEX ON "refresh_tokens" ("token");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;