-- Drop tables
DROP TABLE IF EXISTS "tickets";
DROP TABLE IF EXISTS "tombolas";
DROP TABLE IF EXISTS "interactions";
DROP TABLE IF EXISTS "kermesses_stands";
DROP TABLE IF EXISTS "kermesses_users";
DROP TABLE IF EXISTS "kermesses";
DROP TABLE IF EXISTS "stands";
DROP TABLE IF EXISTS "users";

-- Drop types
DROP TYPE IF EXISTS user_role_enum;
DROP TYPE IF EXISTS stand_type_enum;
DROP TYPE IF EXISTS statut_enum;
DROP TYPE IF EXISTS interaction_type_enum;