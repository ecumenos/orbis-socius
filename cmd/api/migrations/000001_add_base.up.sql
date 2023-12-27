BEGIN;

CREATE TABLE public.accounts
(
  id              BIGINT PRIMARY KEY,
  created_at      TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at      TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at      TIMESTAMP(0) WITH TIME ZONE,
  unique_name     VARCHAR(255) NOT NULL,
  domain          VARCHAR(255) NOT NULL,
  civitas         BIGINT NOT NULL,
  display_name    VARCHAR(255) NOT NULL,
  tombstoned      BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX accounts_unique_name_uindex ON accounts (unique_name);
CREATE INDEX accounts_civitas_index ON accounts (civitas);

COMMIT;
