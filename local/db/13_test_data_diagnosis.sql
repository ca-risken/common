set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO jira_setting(jira_setting_id, name, diagnosis_data_source_id, project_id, identity_field, identity_value, jira_id, jira_key, status, status_detail, scan_at) VALUES
  (1001, "test-a", 1001, 1001, '10035', '1393', '', '', 'CONFIGURED', '', null),
  (1002, "test-b", 1001, 1001, '', '', '10241', '', 'CONFIGURED', '', null),
  (1003, "test-c", 1001, 1001, '', '', '', 'OC202001', 'CONFIGURED', '', null);

INSERT INTO wpscan_setting(wpscan_setting_id, diagnosis_data_source_id,project_id, target_url, status, status_detail, scan_at) VALUES
  (1001, 1002, 1001, 'http://example.com', 'CONFIGURED', '', null);

INSERT INTO portscan_setting(portscan_setting_id, diagnosis_data_source_id, project_id, name) VALUES
  (1001, 1003, 1001, 'test_target');

INSERT INTO portscan_target(portscan_target_id, portscan_setting_id, project_id,target, status, status_detail, scan_at) VALUES
  (1001, 1001, 1001, '127.0.0.1', 'CONFIGURED', '', null);

INSERT INTO application_scan(application_scan_id, diagnosis_data_source_id, project_id, name, scan_type, status, status_detail, scan_at) VALUES
  (1001, 1004, 1001, 'test_target','BASIC', 'CONFIGURED','',null);

INSERT INTO application_scan_basic_setting(application_scan_basic_setting_id, application_scan_id, project_id,target,max_depth,max_children) VALUES
  (1001, 1001, 1001, 'http://exanmple.com', 10, 10);

commit;
