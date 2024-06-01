--
-- PostgreSQL database dump
--
CREATE TABLE public.menus (
    id character varying(20) NOT NULL,
    code character varying(32) NOT NULL,
    name character varying(128) NOT NULL,
    description character varying(1024) NOT NULL,
    sequence integer DEFAULT NULL, 
    parent_id character varying(20) NOT NULL,
    type character varying(20) DEFAULT NULL,
    path character varying(255) DEFAULT NULL,
    parent_path character varying(255) NOT NULL,
    properties text DEFAULT NULL,
    status character varying(20) DEFAULT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.loggers (
    id character varying(20) NOT NULL,
    level character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    trace_id character varying(64) NOT NULL,
    tag character varying(32) NOT NULL,
    message character varying(1024) DEFAULT NULL,
    stack text NOT NULL,
    data text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
);

CREATE TABLE public.permissions (
    id character varying(20) NOT NULL,
    name character varying(128) NOT NULL,
    http_method character varying(255),
    http_path text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_menus (
    id character varying(20) NOT NULL,
    role_id character varying(20) NOT NULL,
    menu_id character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_permissions (
    id character varying(20) NOT NULL,
    role_id character varying(20) NOT NULL,
    permission_id character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.role_users (
    id character varying(20) NOT NULL,
    role_id character varying(20) NOT NULL,
    user_id character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.roles (
    id character varying(20) NOT NULL,
    name character varying(255) NOT NULL,
    description text DEFAULT NULL,
    status character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.user_permissions (
    id character varying(20) NOT NULL
    user_id character varying(20) NOT NULL,
    permission_id character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.users (
    id character varying(20) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    full_name character varying(255) NOT NULL,
    phone character varying(32) NOT NULL,
    remark character varying(1024) NOT NULL,
    avatar character varying(255) DEFAULT NULL,
    status character varying(20) NOT NULL
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
