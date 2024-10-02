-- Type: user_role_enum
CREATE TYPE user_role_enum AS ENUM ('ORGANISATEUR', 'TENEUR_STAND', 'PARENT', 'ENFANT');

-- Table: users
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "role" users_role_enum NOT NULL,
  "jetons" INTEGER NOT NULL DEFAULT 0,
  "parent_id" INTEGER REFERENCES "users"("id") DEFAULT NULL
);
