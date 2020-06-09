set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO user(name, actevated) VALUES
  (1001, 'alice', 'true'),
  (1002, 'bob', 'false');

INSERT INTO project(project_id, name) VALUES
  (1001, 'project-a');

INSERT INTO finding(finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data) VALUES
  (1001, 'desc-1001', 'aws:guardduty',        'guardduty-0001',       'arn:aws:s3:::example-bucket',          1001, 100.00, 1.00, '{"data":{"key":"value"}}'),
  (1002, 'desc-1002', 'aws:access-analizer',  'access-analizer-0001', 'arn:aws:s3:::example-bucket',          1001, 99.05, 0.99,  '{"data":{"key":"value"}}'),
  (1003, 'desc-1003', 'aws:iam-checker',      'iam-checker-0001',     'arn:aws:iam::123456789012:user/alice', 1001, 100.00, 1.00, '{"data":{"key":"value"}}');

INSERT INTO finding_tag(finding_tag_id, finding_id, tag_key, tag_value) VALUES
  (1001, 1001, "key", "value");

INSERT INTO resource(resource_id, resource_name, project_id) VALUES
  (1001, 'arn:aws:s3:::example-bucket',          1001),
  (1002, 'arn:aws:iam::123456789012:user/alice', 1001);

INSERT INTO resource_tag(resource_tag_id, resource_id, tag_key, tag_value) VALUES
  (1001, 1001, 'key1', "value"),
  (1002, 1001, 'key2', "value");

commit;
