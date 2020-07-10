set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- CORE ------------------------------------------------

INSERT INTO project(project_id, name) VALUES
  (1001, 'project-a');

INSERT INTO user(user_id, sub, name, activated) VALUES
  (1001, 'alice', 'alice', 'true'),
  (1002, 'bob', 'bob', 'true'),
  (1003, 'john', 'john', 'true');

INSERT INTO user_role(user_id, role_id, project_id) VALUES
  (1001, 1001, 1001),
  (1002, 1002, 1001),
  (1003, 1003, 1001);

INSERT INTO role(role_id, name, project_id) VALUES
  (1, 'system-admin', null),
  (1001, 'admin-role', 1001),
  (1002, 'viewer-role', 1001),
  (1003, 'aws-guardduty-role', 1001);

INSERT INTO role_policy(role_id, policy_id, project_id) VALUES
  (1001, 1001, 1001),
  (1002, 1002, 1001),
  (1003, 1003, 1001);

INSERT INTO policy(policy_id, name, project_id, action_ptn, resource_ptn) VALUES
  (1, 'system-admin-policy', null, '.*', '.*'),
  (1001, 'admin-policy', 1001, '.*', '.*'),
  (1002, 'viewer-policy', 1001, '(/List|/Get|/Describe)', '.*'),
  (1003, 'aws-guardduty-policy', 1001, '^finding/', '^aws:guardduty/');

INSERT INTO finding(finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data) VALUES
  (1001, 'desc-1001', 'aws:guardduty',       'guardduty-0001',       'arn:aws:s3:::example-bucket',          1001, 100.00, 1.00, '{"data":{"key":"value"}}'),
  (1002, 'desc-1002', 'aws:access-analizer', 'access-analizer-0001', 'arn:aws:s3:::example-bucket',          1001, 99.05, 0.99,  '{"data":{"key":"value"}}'),
  (1003, 'desc-1003', 'aws:iam-checker',     'iam-checker-0001',     'arn:aws:iam::123456789012:user/alice', 1001, 100.00, 1.00, '{"data":{"key":"value"}}');

INSERT INTO finding_tag(finding_tag_id, finding_id, project_id, tag_key, tag_value) VALUES
  (1001, 1001, 1001, "key", "value");

INSERT INTO resource(resource_id, resource_name, project_id) VALUES
  (1001, 'arn:aws:s3:::example-bucket',          1001),
  (1002, 'arn:aws:iam::123456789012:user/alice', 1001);

INSERT INTO resource_tag(resource_tag_id, resource_id, project_id, tag_key, tag_value) VALUES
  (1001, 1001, 1001, 'key1', "value"),
  (1002, 1001, 1001, 'key2', "value");

-- AWS ------------------------------------------------

INSERT INTO aws(aws_id, name, project_id, aws_account_id) VALUES
  (1001, 'account-a', 1001, '123456789012');

INSERT INTO aws_data_source(aws_data_source_id, data_source, max_score) VALUES
  (1001, 'aws:guard-duty', 10.0),
  (1002, 'aws:access-analyzer', 1.0),
  (1003, 'aws:prowler', 1.0),
  (1004, 'aws:iam-activity', 1.0),
  (1005, 'aws:iam-admin', 1.0);

INSERT INTO aws_rel_data_source(aws_id, aws_data_source_id, project_id, assume_role_arn, external_id) VALUES
  (1001, 1001, 1001, 'arn:aws:iam::123456789012:role/role-name', ''),
  (1001, 1002, 1001, 'arn:aws:iam::123456789012:role/role-name', ''),
  (1001, 1003, 1001, 'arn:aws:iam::123456789012:role/role-name', '');


commit;
