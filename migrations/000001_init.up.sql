-- Enum Types
CREATE TYPE user_role_enum AS ENUM ('ORGANISATEUR', 'TENEUR_STAND', 'PARENT', 'ENFANT');
CREATE TYPE stand_type_enum AS ENUM ('VENTE', 'ACTIVITE');
CREATE TYPE statut_enum AS ENUM ('STARTED', 'ENDED');
CREATE TYPE interaction_type_enum AS ENUM ('TRANSACTION', 'ACTIVITE');


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

-- Table: kermesses 
CREATE TABLE "kermesses" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT DEFAULT '',
  "statut" statut_enum NOT NULL DEFAULT 'STARTED'
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

-- Table de liaison entre les kermesses et les utilisateurs
CREATE TABLE "kermesses_users" (
  "id" SERIAL PRIMARY KEY,
  "kermesse_id" INTEGER NOT NULL REFERENCES "kermesses"("id"),
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  UNIQUE ("kermesse_id", "user_id")
);

-- Table de liaison entre les kermesses et les stands
CREATE TABLE "kermesses_stands" (
  "id" SERIAL PRIMARY KEY,
  "kermesse_id" INTEGER NOT NULL REFERENCES "kermesses"("id"),
  "stand_id" INTEGER NOT NULL REFERENCES "stands"("id"),
  UNIQUE ("kermesse_id", "stand_id")
);

--- Table: Interactions
CREATE TABLE "interactions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "kermesse_id" INTEGER NOT NULL REFERENCES "kermesses"("id"),
  "stand_id" INTEGER NOT NULL REFERENCES "stands"("id"),
  "type" interaction_type_enum NOT NULL,
  "statut" statut_enum NOT NULL DEFAULT 'STARTED',
  "jetons" INTEGER NOT NULL DEFAULT 0,
  "points" INTEGER NOT NULL DEFAULT 0
);

--- Table: Tombolas
CREATE TABLE "tombolas" (
  "id" SERIAL PRIMARY KEY,
  "kermesse_id" INTEGER NOT NULL UNIQUE REFERENCES "kermesses"("id"),
  "name" VARCHAR(255) NOT NULL,
  "statut" statut_enum NOT NULL DEFAULT 'STARTED',
  "price" INTEGER NOT NULL DEFAULT 0,
  "lot" VARCHAR(255) NOT NULL
);

--- Table: Tickets
CREATE TABLE "tickets" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER NOT NULL REFERENCES "users"("id"),
  "tombola_id" INTEGER NOT NULL REFERENCES "tombolas"("id"),
  "gagnant" BOOLEAN NOT NULL DEFAULT FALSE
);
