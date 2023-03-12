CREATE TABLE events (
	id SERIAL PRIMARY KEY NOT NULL,
	title CHAR(150) NOT NULL,
	content TEXT NOT NULL,
	created TIMESTAMP NOT NULL,
	setting TIMESTAMP NOT NULL
);

CREATE INDEX idx_events_created ON events (created);

CREATE USER webevents;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.events TO webevents;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.sessions TO webevents;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.users TO webevents;


CREATE TABLE sessions (
	token CHAR(43) PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users (
	id SERIAL NOT NULL PRIMARY KEY,
	name CHAR(255) NOT NULL,
	email CHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created TIMESTAMP NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);



INSERT INTO events (title, content, created, setting) VALUES (
	'An old silent pond',
	E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
	NOW()::timestamp(0),
	NOW()::timestamp(0)+INTERVAL '1 YEAR'
);
INSERT INTO events (title, content, created, setting) VALUES (
	'Over the wintry forest',
	E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
	NOW()::timestamp(0),
	NOW()::timestamp(0)+INTERVAL '1 YEAR'
);
INSERT INTO events (title, content, created, setting) VALUES (
	'First autumn morning',
	E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
	NOW()::timestamp(0),
	NOW()::timestamp(0)+ INTERVAL '7 DAYS'
);
