set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- AWS ------------------------------------------------
INSERT INTO aws_data_source(aws_data_source_id, data_source, max_score) VALUES
  (1001, 'aws:guard-duty', 10.0),
  (1002, 'aws:access-analyzer', 1.0),
  (1003, 'aws:prowler', 1.0),
  (1004, 'aws:iam-activity', 1.0),
  (1005, 'aws:iam-admin', 1.0);

INSERT INTO aws_rel_data_source(aws_id, aws_data_source_id, project_id, assume_role_arn, external_id, status, status_detail, scan_at) VALUES
  (1001, 1001, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1002, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1003, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null);

-- OSINT ------------------------------------------------
INSERT INTO osint_data_source(osint_data_source_id, name, description, max_score) VALUES
  (1, 'intrigue', 'intrigue', 10.0),
  (2, 'tmp_osint', 'this_datasource_does_not exist.', 10.0);

-- DIAGNOSIS ------------------------------------------------
INSERT INTO diagnosis_data_source(diagnosis_data_source_id, name, description, max_score) VALUES
  (1, 'diagnosis:jira', 'jira', 10.0);

commit;
