set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- CORE ------------------------------------------------
INSERT INTO project(project_id, name) VALUES
  (1001, 'project-a'),
  (1002, 'project-b');

INSERT INTO user(user_id, sub, name, activated) VALUES
  (1001, 'alice', 'alice', 'true'),
  (1002, 'bob', 'bob', 'true'),
  (1003, 'john', 'john', 'true');

INSERT INTO user_role(user_id, role_id, project_id) VALUES
  (1001, 1, null),
  (1002, 1001, 1001),
  (1003, 1002, 1001);

INSERT INTO role(role_id, name, project_id) VALUES
  (1001, 'project-a-admin', 1001),
  (1002, 'project-a-viewer', 1001);

INSERT INTO role_policy(role_id, policy_id, project_id) VALUES
  (1001, 1001, 1001),
  (1002, 1002, 1001);

INSERT INTO policy(policy_id, name, project_id, action_ptn, resource_ptn) VALUES
  (1001, 'admin-policy', 1001, '.*', '.*'),
  (1002, 'viewer-policy', 1001, '^.*/(get-|list-|describe-|show-)', '.*');

INSERT INTO access_token(access_token_id, token_hash, name, description, project_id, last_updated_user_id) VALUES
  -- plain-text: `token`
  (1001, '2265daba0872fc3aef169d079365e590f0cbc8ed46c2a7984c8a642803cfd96cb47804a63cf22a79f6ca469268c29ee9e72a5059b62d0a598fe42dfc8dcc51bc', 'token-name', 'test', 1001, 1001);

INSERT INTO access_token_role(access_token_id, role_id) VALUES
  (1001, 1001);

INSERT INTO finding(finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data) VALUES
  (1001, 'desc-1001', 'aws:guard-duty', 'guard-duty-0001', 'arn:aws:s3:::example-bucket',      1001, 100.00, 1.00, '{"data":{"key":"value"}}'),
  (1002, 'desc-1002', 'aws:guard-duty', 'guard-duty-0002', 'arn:aws:s3:::example-bucket',      1001, 99.05, 0.99,  '{"data":{"key":"value"}}'),
  (1003, 'desc-1003', 'osint:website',  'website-0001',    'example-tool.',                    1001, 100.00, 1.00, '{"data":{"key":"value"}}'),
  (1004, 'desc-1004', 'aws:guard-duty', 'guard-duty-1004', 'arn:aws:s3:::example-bucket-1004', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1005, 'desc-1005', 'aws:guard-duty', 'guard-duty-1005', 'arn:aws:s3:::example-bucket-1005', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1006, 'desc-1006', 'aws:guard-duty', 'guard-duty-1006', 'arn:aws:s3:::example-bucket-1006', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1007, 'desc-1007', 'aws:guard-duty', 'guard-duty-1007', 'arn:aws:s3:::example-bucket-1007', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1008, 'desc-1008', 'aws:guard-duty', 'guard-duty-1008', 'arn:aws:s3:::example-bucket-1008', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1009, 'desc-1009', 'aws:guard-duty', 'guard-duty-1009', 'arn:aws:s3:::example-bucket-1009', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1010, 'desc-1010', 'aws:guard-duty', 'guard-duty-1010', 'arn:aws:s3:::example-bucket-1010', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1011, 'desc-1011', 'aws:guard-duty', 'guard-duty-0011', 'arn:aws:s3:::example-bucket-0011', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1012, 'desc-1012', 'aws:guard-duty', 'guard-duty-0012', 'arn:aws:s3:::example-bucket-0012', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1013, 'desc-1013', 'aws:guard-duty', 'guard-duty-1013', 'arn:aws:s3:::example-bucket-1013', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1014, 'desc-1014', 'aws:guard-duty', 'guard-duty-1014', 'arn:aws:s3:::example-bucket-1014', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1015, 'desc-1015', 'aws:guard-duty', 'guard-duty-1015', 'arn:aws:s3:::example-bucket-1015', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1016, 'desc-1016', 'aws:guard-duty', 'guard-duty-1016', 'arn:aws:s3:::example-bucket-1016', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1017, 'desc-1017', 'aws:guard-duty', 'guard-duty-1017', 'arn:aws:s3:::example-bucket-1017', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1018, 'desc-1018', 'aws:guard-duty', 'guard-duty-1018', 'arn:aws:s3:::example-bucket-1018', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1019, 'desc-1019', 'aws:guard-duty', 'guard-duty-1019', 'arn:aws:s3:::example-bucket-1019', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1020, 'desc-1020', 'aws:guard-duty', 'guard-duty-1020', 'arn:aws:s3:::example-bucket-1020', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1030, 'desc-1030', 'aws:guard-duty', 'guard-duty-1030', 'arn:aws:s3:::example-bucket-1030', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1031, 'desc-1031', 'aws:guard-duty', 'guard-duty-0031', 'arn:aws:s3:::example-bucket-0031', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1032, 'desc-1032', 'aws:guard-duty', 'guard-duty-0032', 'arn:aws:s3:::example-bucket-0032', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1033, 'desc-1033', 'aws:guard-duty', 'guard-duty-1033', 'arn:aws:s3:::example-bucket-1033', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1034, 'desc-1034', 'aws:guard-duty', 'guard-duty-1034', 'arn:aws:s3:::example-bucket-1034', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1035, 'desc-1035', 'aws:guard-duty', 'guard-duty-1035', 'arn:aws:s3:::example-bucket-1035', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1036, 'desc-1036', 'aws:guard-duty', 'guard-duty-1036', 'arn:aws:s3:::example-bucket-1036', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1037, 'desc-1037', 'aws:guard-duty', 'guard-duty-1037', 'arn:aws:s3:::example-bucket-1037', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1038, 'desc-1038', 'aws:guard-duty', 'guard-duty-1038', 'arn:aws:s3:::example-bucket-1038', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1039, 'desc-1039', 'aws:guard-duty', 'guard-duty-1039', 'arn:aws:s3:::example-bucket-1039', 1001, 10.00, 0.5, '{"data":{"key":"value"}}'),
  (1040, 'desc-1040', 'aws:guard-duty', 'guard-duty-1040', 'arn:aws:s3:::example-bucket-1040', 1001, 10.00, 0.5, '{"data":{"key":"value"}}');

INSERT INTO finding_tag(finding_tag_id, finding_id, project_id, tag) VALUES
  (1001, 1001, 1001, "tag"),
  (1002, 1001, 1001, "tag:key"),
  (1003, 1002, 1001, "tag:key");

INSERT INTO resource(resource_id, resource_name, project_id) VALUES
  (1001, 'arn:aws:s3:::example-bucket',          1001),
  (1002, 'arn:aws:iam::123456789012:user/alice', 1001),
  (1003, 'Cross Site Sciripting',                1001),
  (1004, 'arn:aws:s3:::example-bucket-1004',     1001),
  (1005, 'arn:aws:s3:::example-bucket-1005',     1001),
  (1006, 'arn:aws:s3:::example-bucket-1006',     1001),
  (1007, 'arn:aws:s3:::example-bucket-1007',     1001),
  (1008, 'arn:aws:s3:::example-bucket-1008',     1001),
  (1009, 'arn:aws:s3:::example-bucket-1009',     1001),
  (1010, 'arn:aws:s3:::example-bucket-1010',     1001),
  (1011, 'arn:aws:s3:::example-bucket-0011',     1001),
  (1012, 'arn:aws:s3:::example-bucket-0012',     1001),
  (1013, 'arn:aws:s3:::example-bucket-1013',     1001),
  (1014, 'arn:aws:s3:::example-bucket-1014',     1001),
  (1015, 'arn:aws:s3:::example-bucket-1015',     1001),
  (1016, 'arn:aws:s3:::example-bucket-1016',     1001),
  (1017, 'arn:aws:s3:::example-bucket-1017',     1001),
  (1018, 'arn:aws:s3:::example-bucket-1018',     1001),
  (1019, 'arn:aws:s3:::example-bucket-1019',     1001),
  (1020, 'arn:aws:s3:::example-bucket-1020',     1001),
  (1030, 'arn:aws:s3:::example-bucket-1030',     1001),
  (1031, 'arn:aws:s3:::example-bucket-0031',     1001),
  (1032, 'arn:aws:s3:::example-bucket-0032',     1001),
  (1033, 'arn:aws:s3:::example-bucket-1033',     1001),
  (1034, 'arn:aws:s3:::example-bucket-1034',     1001),
  (1035, 'arn:aws:s3:::example-bucket-1035',     1001),
  (1036, 'arn:aws:s3:::example-bucket-1036',     1001),
  (1037, 'arn:aws:s3:::example-bucket-1037',     1001),
  (1038, 'arn:aws:s3:::example-bucket-1038',     1001),
  (1039, 'arn:aws:s3:::example-bucket-1039',     1001),
  (1040, 'arn:aws:s3:::example-bucket-1040',     1001);

INSERT INTO resource_tag(resource_tag_id, resource_id, project_id, tag) VALUES
  (1001, 1001, 1001, 'tag1'),
  (1002, 1001, 1001, "tag:key"),
  (1003, 1002, 1001, "tag:key");

INSERT INTO alert_condition(alert_condition_id, description, severity, project_id, and_or, enabled) VALUES
  (1001, 'test_alert_condition', 'high', 1001, 'and', true),
  (1002, 'test_alert_condition_2', 'medium', 1001, 'or', false);

INSERT INTO alert(alert_id, alert_condition_id, description, severity, project_id, status) VALUES
  (1001, 1001, 'test_alert', 'high', 1001, 'ACTIVE'),
  (1002, 1001, 'test_alert_2', 'medium', 1001, 'ACTIVE'),
  (1003, 1001, 'test_alert_3', 'low', 1001, 'ACTIVE');

INSERT INTO alert_history(alert_history_id, history_type, alert_id, description, severity,finding_history, project_id) VALUES
  (1001, 'created', 1001, 'test_alert_history', 'high','{"finding_id":[1001,1002,1003,1004,1005,1006]}', 1001),
  (1002, 'deleted', 1001, 'test_alert_history_2', 'high','{"finding_id":[]}', 1001),
  (1003, 'created', 1002, 'test_alert_history', 'low','{"finding_id":[1001,1002]}', 1001),
  (1009, 'updated', 1002, 'test_alert_history_2', 'low','{"finding_id":[1001]}', 1001),
  (1010, 'deleted', 1002, 'test_alert_history_2', 'high','{"finding_id":[]}', 1001);

INSERT INTO rel_alert_finding(alert_id, finding_id, project_id) VALUES
  (1001, 1001, 1001),
  (1001, 1002, 1001),
  (1002, 1002, 1001);

INSERT INTO alert_rule(alert_rule_id, name, project_id, score, resource_name, tag, finding_cnt) VALUES
  (1001, 'test_alert_rule', 1001, 1.0, '', '', 1),
  (1002, 'test_alert_rule_2', 1001, 1.0, 'test', '', 1);

INSERT INTO alert_cond_rule(alert_condition_id, alert_rule_id, project_id) VALUES
  (1001, 1001, 1001),
  (1002, 1002, 1001);

INSERT INTO notification(notification_id, name, project_id, type, notify_setting) VALUES
  (1001, 'test_notification', 1001, 'slack', '{"webhook_url":"http://hogehoge.com/fuga"}'),
  (1002, 'test_notification_2', 1001, 'slack', '{"webhook_url":"http://hogehoge2.com/fuga2"}');

INSERT INTO alert_cond_notification(alert_condition_id, notification_id, project_id, cache_second, notified_at) VALUES
  (1001, 1001, 1001, 1800, '2020-09-01 16:00:00'),
  (1002, 1002, 1001, 1800, '2020-09-02 16:00:00');

commit;
