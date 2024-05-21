-- Drop the index on the 'coords' column of the 'locations' table
DROP INDEX IF EXISTS locations_coords_idx;

-- Drop the 'resources' table, automatically removing the foreign keys and unique constraints
DROP TABLE IF EXISTS resources;

-- Drop the 'locations' table
DROP TABLE IF EXISTS locations;
