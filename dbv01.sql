CREATE DATABASE mathrix2023;
\c mathrix2023;

CREATE TABLE events (
	eventid SERIAL PRIMARY KEY NOT NULL,
	eventname VARCHAR(150) NOT NULL,
	content TEXT NOT NULL,
	venue VARCHAR(150) NOT NULL,
	technicality boolean NOT NULL,	
	eventdate DATE NOT NULL,
	UNIQUE(eventname)
);
CREATE TABLE sponsers (
	sponserid SERIAL PRIMARY KEY NOT NULL,
	sponsername VARCHAR(150) NOT NULL,
	url VARCHAR(150) NOT NULL,
	UNIQUE(sponsername)
);
CREATE TABLE sessions (
	token VARCHAR(43) PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);
CREATE TABLE users (
	userid SERIAL NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	urn VARCHAR(20) NOT NULL,
	phone bigint NOT NULL check (phone between 6000000000 and 9999999999),
	college VARCHAR(255) NOT NULL,
	dept VARCHAR(255) NOT NULL,
	year int check (year between 1 and 10),
	degree VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created TIMESTAMP NOT NULL,
	authority BOOLEAN NOT NULL DEFAULT false,
	UNIQUE(urn, phone, email)
);
CREATE TABLE eventlist (
	id SERIAL NOT NULL,
	userid INT NOT NULL,
		CONSTRAINT referenceuser
			FOREIGN KEY(userid)
				REFERENCES users(userid) ON DELETE CASCADE,
	eventid INT NOT NULL,
		CONSTRAINT referenceevent
			FOREIGN KEY(eventid)
				REFERENCES events(eventid)ON DELETE CASCADE,
	CONSTRAINT eventlistpk
		PRIMARY KEY(eventid, userid)
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
CREATE USER webmathrix;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.events TO webmathrix;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.eventlist TO webmathrix;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.sessions TO webmathrix;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.users TO webmathrix;
GRANT SELECT, INSERT, UPDATE, DELETE ON public.sponsers TO webmathrix;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO webmathrix;
ALTER USER webmathrix WITH PASSWORD 'neomathrix';
INSERT INTO users (name, urn, phone, college,dept, year, degree, email, hashed_password, authority, created) VALUES('Suryaa','2019242025', 8667219826, 'College of Engineering, Guindy', 'Mathematics', 4, 'Information Technology', 'suryaamana@gmail.com', '$2a$12$L9r2pVOv0U.UR6cmfBmgS.DqEms0DVT06hjRIponOIEWnNQj68QxS', true, NOW()::timestamp(0));
