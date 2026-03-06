-- Table: game.creature

-- DROP TABLE IF EXISTS game.creature;

CREATE TABLE IF NOT EXISTS game.creature
(
    id uuid NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    health_points smallint NOT NULL,
    defence jsonb,
    initiative smallint NOT NULL,
    movement smallint NOT NULL,
    action_count smallint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT creature_pkey PRIMARY KEY (id)
);

-- Table: game.creature_ability_slot

-- DROP TABLE IF EXISTS game.creature_ability_slot;

CREATE TABLE IF NOT EXISTS game.creature_ability_slot
(
    creature_id uuid NOT NULL,
    ability_id uuid NOT NULL,
    roll_threshold smallint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT creature_ability_slot_pkey PRIMARY KEY (creature_id, ability_id)
);

