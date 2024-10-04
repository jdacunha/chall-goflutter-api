-- Enum Types
CREATE TYPE user_role_enum AS ENUM ('ORGANISATEUR', 'TENEUR_STAND', 'PARENT', 'ENFANT');
CREATE TYPE stand_type_enum AS ENUM ('VENTE', 'ACTIVITE');
CREATE TYPE statut_enum AS ENUM ('STARTED', 'ENDED');


-- Table: users
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "role" user_role_enum NOT NULL,
  "jetons" INTEGER NOT NULL DEFAULT 0,
  "parent_id" INTEGER REFERENCES "users"("id") DEFAULT NULL
);

-- Table: stands
CREATE TABLE "stands" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '',
  "type" stand_type_enum NOT NULL,
  "price" INTEGER NOT NULL DEFAULT 0,
  "stock" INTEGER NOT NULL DEFAULT 0
);

-- Table: kermesses
CREATE TABLE "kermesses" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '',
  "statut" statut_enum NOT NULL DEFAULT 'STARTED'
);