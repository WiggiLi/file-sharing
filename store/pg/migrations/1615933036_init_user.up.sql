CREATE TABLE users (
	id uuid NOT NULL,
	name text NOT NULL,
	email text NOT NULL,
	password text NOT NULL,
	token text,
	created_at timestamp default current_timestamp,
	CONSTRAINT "pk_user_id" PRIMARY KEY (id)
);

CREATE TABLE files_of_users (
	id_user uuid NOT NULL,
	id_file uuid NOT NULL
);
