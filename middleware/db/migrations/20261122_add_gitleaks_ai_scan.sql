ALTER TABLE code_gitleaks_setting
  ADD COLUMN ai_scan BOOLEAN NOT NULL DEFAULT FALSE AFTER scan_private;
