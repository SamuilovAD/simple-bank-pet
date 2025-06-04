-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-06-04T18:16:47.848Z

CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY NOT NULL,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigint PRIMARY KEY NOT NULL,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY NOT NULL,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar UNIQUE NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigint PRIMARY KEY NOT NULL,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00+00'),
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "owner_currency_key" ON "accounts" ("owner", "currency");

CREATE INDEX "accounts_owner_idx" ON "accounts" USING BTREE ("owner");

CREATE INDEX "entries_account_id_idx" ON "entries" USING BTREE ("account_id");

CREATE INDEX "transfers_from_account_id_idx" ON "transfers" USING BTREE ("from_account_id");

CREATE INDEX "transfers_from_account_id_to_account_id_idx" ON "transfers" USING BTREE ("from_account_id", "to_account_id");

CREATE INDEX "transfers_to_account_id_idx" ON "transfers" USING BTREE ("to_account_id");

ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_fkey" FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD CONSTRAINT "entries_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "sessions" ADD CONSTRAINT "sessions_username_fkey" FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_from_account_id_fkey" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_to_account_id_fkey" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
