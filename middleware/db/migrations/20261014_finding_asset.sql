-- Run once before deploying asset-aware Core. Equality uses the full binary key.
ALTER TABLE finding
 ADD COLUMN asset_type VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
 ADD COLUMN asset_key VARBINARY(8192) NOT NULL DEFAULT '',
 ADD COLUMN asset_name TEXT NOT NULL DEFAULT (''),
 ADD INDEX idx_asset(project_id, asset_type, asset_key(191));
