set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO jira_setting(jira_setting_id, name, diagnosis_data_source_id, project_id, identity_field, identity_value, jira_id, jira_key, status, status_detail, scan_at) VALUES
  (1001, "test-a", 1001, 1001, '10035', '1393', '', '', 'CONFIGURED', '', null),
  (1002, "test-b", 1001, 1001, '', '', '10241', '', 'CONFIGURED', '', null),
  (1003, "test-c", 1001, 1001, '', '', '', 'OC202001', 'CONFIGURED', '', null);

INSERT INTO wpscan_setting(wpscan_setting_id, diagnosis_data_source_id,project_id, target_url, status, status_detail, scan_at) VALUES
  (1001, 1002, 1001, 'http://example.com', 'CONFIGURED', '', null);

commit;
