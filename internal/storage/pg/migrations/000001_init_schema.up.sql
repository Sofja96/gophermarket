CREATE TABLE "users" (
   "id" bigserial PRIMARY KEY,
   "login" varchar NOT NULL,
   "balance" real NOT NULL default 0,
   "withdrawn" real NOT NULL default 0,
   "password" varchar NOT NULL
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES users(id),
  "number" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT ('NEW'),
  "accrual" real NOT NULL default 0,
  "uploaded_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "withdrawals" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES users(id),
  "number" varchar NOT NULL,
  "sum" real NOT NULL,
  "processed_at"  timestamptz NOT NULL DEFAULT (now())
);




