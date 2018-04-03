CREATE TABLE setting (
	key	  TEXT NOT NULL PRIMARY KEY UNIQUE,
	value TEXT
);

CREATE TABLE location (
	id   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	name TEXT,
	time TEXT NOT NULL,
	lat  TEXT NOT NULL,
	lng  TEXT NOT NULL
);

CREATE TABLE person (
	id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	first_name TEXT NOT NULL,
	last_name  TEXT NOT NULL
);

CREATE TABLE location_person (
	location_id INTEGER NOT NULL,
	person_id   INTEGER NOT NULL,
	PRIMARY KEY(location_id, person_id),
	FOREIGN KEY(location_id) REFERENCES location(id) ON DELETE CASCADE,
	FOREIGN KEY(person_id) REFERENCES person(id) ON DELETE CASCADE
);

INSERT INTO setting (key, value) VALUES ('db_version', '1');