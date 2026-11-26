ALTER TABLE report_finding
  ADD INDEX idx_report_finding_project_date (project_id, report_date),
  ALGORITHM=INPLACE,
  LOCK=NONE;
