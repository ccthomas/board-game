-- Table: game.damage_type

-- DROP TABLE IF EXISTS game.damage_type;

CREATE TABLE IF NOT EXISTS game.damage_type
(
    id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT damage_type_pkey PRIMARY KEY (id)
)