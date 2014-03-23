CREATE TABLE IF NOT EXISTS urls ( id serial, "when" timestamp, url text, title varchar(255) );
ALTER TABLE urls ADD COLUMN shared_by varchar(255);
