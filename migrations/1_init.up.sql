CREATE TABLE IF NOT EXISTS users (
    	uid SERIAL PRIMARY KEY,
    	name TEXT,
    	login TEXT,
    	password TEXT
);

CREATE TABLE IF NOT EXISTS books (
		b_id SERIAL PRIMARY KEY,
		author TEXT,
		title TEXT,
		uid INTEGER
);