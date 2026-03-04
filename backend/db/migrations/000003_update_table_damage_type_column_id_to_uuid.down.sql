-- ============================================================
-- MIGRATION DOWN: Revert game.damage_type.id from UUID to VARCHAR(36)
-- ============================================================

BEGIN;

-- Step 1: Drop the UUID default (if set during UP migration)
ALTER TABLE game.damage_type
    ALTER COLUMN id DROP DEFAULT;

-- Step 2: Add a new VARCHAR(36) column alongside the UUID id
ALTER TABLE game.damage_type
    ADD COLUMN id_old CHARACTER VARYING(36);

-- Step 3: Populate the old-style column by casting UUID back to text
UPDATE game.damage_type
SET id_old = id::TEXT;

-- Step 4: Set NOT NULL on the reverted column
ALTER TABLE game.damage_type
    ALTER COLUMN id_old SET NOT NULL;

-- Step 5: Drop the existing primary key constraint
ALTER TABLE game.damage_type
    DROP CONSTRAINT damage_type_pkey;

-- Step 6: Drop the UUID id column
ALTER TABLE game.damage_type
    DROP COLUMN id;

-- Step 7: Rename the VARCHAR column back to id
ALTER TABLE game.damage_type
    RENAME COLUMN id_old TO id;

-- Step 8: Re-add the primary key on the VARCHAR column
ALTER TABLE game.damage_type
    ADD CONSTRAINT damage_type_pkey PRIMARY KEY (id);

COMMIT;