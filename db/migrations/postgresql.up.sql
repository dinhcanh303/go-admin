--
-- PostgreSQL database dump
--
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

CREATE TABLE public.menus (
    id integer DEFAULT nextval('public.goadmin_menu_myid_seq'::regclass) NOT NULL,
    parent_id integer DEFAULT 0 NOT NULL,
    type integer DEFAULT 0,
    "order" integer DEFAULT 0 NOT NULL,
    title character varying(50) NOT NULL,
    header character varying(100),
    plugin_name character varying(100) NOT NULL,
    icon character varying(50) NOT NULL,
    uri character varying(3000) NOT NULL,
    uuid character varying(100),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.loggers (
    id integer DEFAULT nextval('public.goadmin_operation_log_myid_seq'::regclass) NOT NULL,
    level character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    tag character varying(255) NOT NULL,
    message character varying(255),
    stack text NOT NULL,
    data text NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
);

CREATE TABLE public.permissions (
    id integer DEFAULT nextval('public.goadmin_permissions_myid_seq'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    http_method character varying(255),
    http_path text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_menus (
    role_id integer NOT NULL,
    menu_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_users (
    role_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.roles (
    id integer DEFAULT nextval('public.goadmin_roles_myid_seq'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    description text DEFAULT NULL,
    status character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.user_permissions (
    user_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.users (
    id integer DEFAULT nextval('public.goadmin_users_myid_seq'::regclass) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    full_name character varying(255) NOT NULL,
    phone character varying(32) NOT NULL,
    remark character varying(1024) NOT NULL,
    avatar character varying(255) DEFAULT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
