CREATE TABLE "oauth_accounts" (
  "id" serial,
  "user_id" int NOT NULL,
  "provider" varchar(20) NOT NULL,
  "provider_user_id" varchar(100) UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE INDEX ON "oauth_accounts" ("user_id");

CREATE INDEX ON "oauth_accounts" ("provider");

CREATE UNIQUE INDEX ON "oauth_accounts" ("provider_user_id");

ALTER TABLE "oauth_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;