CREATE TABLE "withdrawals" (
  "id" bigserial PRIMARY KEY,
  "from_account" bigserial NOT NULL,
  "amout" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

