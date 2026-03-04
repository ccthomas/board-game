-- Table: game.ability

-- DROP TABLE IF EXISTS game.ability;

CREATE TABLE IF NOT EXISTS game.ability
(
    id UUID NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    pattern character varying(20) COLLATE pg_catalog."default" NOT NULL,
    range integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT ability_pkey PRIMARY KEY (id)
);

-- Table: game.ability_effect

-- DROP TABLE IF EXISTS game.ability_effect;

CREATE TABLE IF NOT EXISTS game.ability_effect
(
    id UUID NOT NULL,
    ability_id UUID NOT NULL,
    expression jsonb NOT NULL,
    effect_type character varying(20) COLLATE pg_catalog."default" NOT NULL,
    damage_type_id UUID,
    alignment character varying(20) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT ability_effect_pkey PRIMARY KEY (id)
);