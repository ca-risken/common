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
  (1001, 1001, 1001),
  (1001, 1004, 1002),
  (1002, 1002, 1001),
  (1003, 1003, 1001);

INSERT INTO role(role_id, name, project_id) VALUES
  (1, 'system-admin', null),
  (1001, 'project-a-admin', 1001),
  (1002, 'project-a-viewer', 1001),
  (1003, 'project-a-aws-guardduty', 1001),
  (1004, 'project-b-admin', 1002);

INSERT INTO role_policy(role_id, policy_id, project_id) VALUES
  (1001, 1001, 1001),
  (1002, 1002, 1001),
  (1003, 1003, 1001),
  (1004, 1004, 1002);

INSERT INTO policy(policy_id, name, project_id, action_ptn, resource_ptn) VALUES
  (1, 'system-admin-policy', null, '.*', '.*'),
  (1001, 'admin-policy', 1001, '.*', '.*'),
  (1002, 'viewer-policy', 1001, '(/List|/Get|/Describe)', '.*'),
  (1003, 'aws-guardduty-policy', 1001, '^finding/', '^aws:guardduty/'),
  (1004, 'admin-policy', 1002, '.*', '.*');

INSERT INTO finding(finding_id, description, data_source, data_source_id, resource_name, project_id, original_score, score, data) VALUES
  (1001, 'desc-1001', 'aws:guardduty',       'guardduty-0001',       'arn:aws:s3:::example-bucket',          1001, 100.00, 1.00, '{"data":{"key":"value"}}'),
  (1002, 'desc-1002', 'aws:access-analizer', 'access-analizer-0001', 'arn:aws:s3:::example-bucket',          1001, 99.05, 0.99,  '{"data":{"key":"value"}}'),
  (1003, 'desc-1003', 'aws:iam-checker',     'iam-checker-0001',     'arn:aws:iam::123456789012:user/alice', 1001, 100.00, 1.00, '{"data":{"key":"value"}}');

INSERT INTO finding_tag(finding_tag_id, finding_id, project_id, tag) VALUES
  (1001, 1001, 1001, "tag"),
  (1002, 1001, 1001, "tag:key"),
  (1003, 1002, 1001, "tag:key");

INSERT INTO resource(resource_id, resource_name, project_id) VALUES
  (1001, 'arn:aws:s3:::example-bucket',          1001),
  (1002, 'arn:aws:iam::123456789012:user/alice', 1001);

INSERT INTO resource_tag(resource_tag_id, resource_id, project_id, tag) VALUES
  (1001, 1001, 1001, 'tag1'),
  (1002, 1001, 1001, "tag:key"),
  (1003, 1002, 1001, "tag:key");

INSERT INTO alert_condition(alert_condition_id, description, severity, project_id, and_or, enabled) VALUES
  (1001, 'test_alert_condition', 'high', 1001, 'and', true),
  (1002, 'test_alert_condition_2', 'medium', 1001, 'or', false);

INSERT INTO alert(alert_id, alert_condition_id, description, severity, project_id, activated) VALUES
  (1001, 1001, 'test_alert', 'high', 1001, true),
  (1002, 1001, 'test_alert_2', 'medium', 1001, true),
  (1003, 1001, 'test_alert_3', 'low', 1001, true);

INSERT INTO alert_history(alert_history_id, history_type, alert_id, description, severity, project_id) VALUES
  (1001, 'created', 1001, 'test_alert_history', 'high', 1001),
  (1002, 'deleted', 1001, 'test_alert_history_2', 'high', 1001),
  (1003, 'created', 1002, 'test_alert_history', 'low', 1001),
  (1004, 'updated', 1002, 'test_alert_history_2', 'low', 1001),
  (1005, 'updated', 1002, 'test_alert_history_2', 'medium', 1001),
  (1006, 'updated', 1002, 'test_alert_history_2', 'medium', 1001),
  (1007, 'updated', 1002, 'test_alert_history_2', 'high', 1001),
  (1008, 'updated', 1002, 'test_alert_history_2', 'medium', 1001),
  (1009, 'updated', 1002, 'test_alert_history_2', 'low', 1001),
  (1010, 'deleted', 1002, 'test_alert_history_2', 'high', 1001);

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
