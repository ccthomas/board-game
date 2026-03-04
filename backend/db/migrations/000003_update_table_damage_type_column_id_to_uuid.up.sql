-- ============================================================
-- MIGRATION UP: Convert game.damage_type.id from VARCHAR(36) to UUID
-- Goal: Reduce column storage to native 16-byte UUID type
-- ============================================================

BEGIN;

-- Step 1: Add a new UUID column alongside the existing id
ALTER TABLE game.damage_type
    ADD COLUMN id_new UUID;

-- Step 2: Populate the new column by casting existing VARCHAR UUIDs
UPDATE game.damage_type
SET id_new = id::UUID;

-- Step 3: Ensure no NULLs exist before setting NOT NULL constraint
--         (This will error if any id values are not valid UUIDs)
ALTER TABLE game.damage_type
    ALTER COLUMN id_new SET NOT NULL;

-- Step 4: Drop the existing primary key constraint
ALTER TABLE game.damage_type
    DROP CONSTRAINT damage_type_pkey;

-- Step 5: Drop the old VARCHAR id column
ALTER TABLE game.damage_type
    DROP COLUMN id;

-- Step 6: Rename the new UUID column to id
ALTER TABLE game.damage_type
    RENAME COLUMN id_new TO id;

-- Step 7: Re-add the primary key on the new UUID column
ALTER TABLE game.damage_type
    ADD CONSTRAINT damage_type_pkey PRIMARY KEY (id);

COMMIT;