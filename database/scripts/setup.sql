BEGIN;

DO $$
BEGIN
	CREATE USER db_app WITH PASSWORD 'dev';
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating role db_app -- it already exists';
END
$$;

DO $$
BEGIN
	CREATE SCHEMA private;
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating schema private -- it already exists';
END
$$;

DO $$
BEGIN
	CREATE TYPE private.jwt_token as (
		role TEXT,
		user_id INTEGER,
		authentication_id INTEGER
	);
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating type private.jwt_token -- it already exists';
END
$$;

DO $$
BEGIN
	CREATE ROLE db_app_player;
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating role db_app_player -- it already exists';

	
END
$$;

DO $$
BEGIN
	CREATE ROLE db_app_guest;
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating role db_app_guest -- it already exists';

	
END
$$;

DO $$
BEGIN
	CREATE ROLE db_app_admin;
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating role db_app_admin -- it already exists';

	
END
$$;

DO $$
BEGIN
	CREATE ROLE data_analyst;
	EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'not creating role data_analyst -- it already exists';
END
$$;

COMMIT;


BEGIN;
DO $$
BEGIN
	GRANT db_app to postgres;
	GRANT db_app_player TO db_app;
	GRANT db_app_guest TO db_app;
	GRANT db_app_admin TO db_app;
	GRANT data_analyst TO db_app;
END
$$;
COMMIT;
