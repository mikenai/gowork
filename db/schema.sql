CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(255) primary key);
CREATE TABLE users (
    id string PRIMARY key,
    name string
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20221202132622');
