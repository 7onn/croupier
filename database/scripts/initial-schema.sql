--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

CREATE TABLE public.test (
    id integer NOT NULL,
    name character varying,
    email character varying,
    password character varying,
    login character varying,
    active boolean,
    phone character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    password_digest character varying,
    role character varying DEFAULT 'user'::character varying,
    description text,
    confirmation_code character varying,
    confirmation_date timestamp without time zone,
    birth date
);
