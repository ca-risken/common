set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- AWS ------------------------------------------------
INSERT INTO aws(aws_id, name, project_id, aws_account_id) VALUES
  (1001, 'account-a', 1001, '123456789012');

INSERT INTO aws_rel_data_source(aws_id, aws_data_source_id, project_id, assume_role_arn, external_id, status, status_detail, scan_at) VALUES
  (1001, 1001, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1002, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1003, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1004, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null),
  (1001, 1005, 1001, 'arn:aws:iam::123456789012:role/role-name', '', 'CONFIGURED', '', null);
commit;
