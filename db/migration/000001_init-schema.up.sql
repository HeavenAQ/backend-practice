CREATE TABLE "accounts" (
  id bigserial PRIMARY KEY,
  owner VARCHAR NOT NULL,
  balance BIGINT NOT NULL,
  currency VARCHAR NOT NULL,
  created TIMESTAMPTZ NOT NULL DEFAULT (now()) 
);

CREATE TABLE "entries" (
  id bigserial  PRIMARY KEY,
  account_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created TIMESTAMPTZ NOT NULL DEFAULT (now()),

  CONSTRAINT fk_account_id
    FOREIGN KEY(account_id)
      REFERENCES "accounts"(id)
);

CREATE TABLE "transfers" (
  id bigserial PRIMARY KEY,
  from_account_id BIGINT NOT NULL,
  to_account_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created TIMESTAMPTZ NOT NULL DEFAULT (now()),

  CONSTRAINT fk_from_account_id
    FOREIGN KEY(from_account_id)
      REFERENCES "accounts"(id),

  CONSTRAINT fk_to_account_id
    FOREIGN KEY(to_account_id)
      REFERENCES "accounts"(id)
);


CREATE INDEX ON "accounts" (owner);
CREATE INDEX ON "entries" (account_id);
CREATE INDEX ON "transfers" (from_account_id);
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries".amount IS 'can be negative or positive';
COMMENT ON COLUMN "transfers".amount IS 'can only be positive';

