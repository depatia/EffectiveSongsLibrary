CREATE TABLE service.songs (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"group" varchar NOT NULL,
	song varchar NOT NULL,
	release_date date NOT NULL,
	link varchar NOT NULL,
	"text" varchar NOT NULL,
	CONSTRAINT songs_pk PRIMARY KEY (id),
	CONSTRAINT songs_unique UNIQUE (link)
);
CREATE UNIQUE INDEX songs_id_idx ON service.songs USING btree (id);