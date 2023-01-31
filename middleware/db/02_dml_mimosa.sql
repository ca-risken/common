set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- CORE ------------------------------------------------
INSERT INTO role(role_id, name, project_id) VALUES
  (1, 'system-admin-role', null),
  (2, 'system-readonly-role', null);

INSERT INTO policy(policy_id, name, project_id, action_ptn, resource_ptn) VALUES
  (1, 'system-admin-policy', null, '.*', '.*'),
  (2, 'system-readonly-policy', null, '^.*/(list|get|describe)-.*', '.*');

INSERT INTO role_policy(role_id, policy_id, project_id) VALUES
  (1, 1, null),
  (2, 2, null);

-- AWS ------------------------------------------------
INSERT INTO aws_data_source(aws_data_source_id, data_source, max_score) VALUES
  (1001, 'aws:guard-duty', 10.0),
  (1002, 'aws:access-analyzer', 1.0),
  (1003, 'aws:admin-checker', 1.0),
  (1004, 'aws:cloudsploit', 10.0),
  (1005, 'aws:portscan', 10.0);

-- OSINT ------------------------------------------------
INSERT INTO osint_data_source(osint_data_source_id, name, description, max_score) VALUES
  (1002, 'osint:subdomain', 'researcher about subdomain', 10.0),
  (1003, 'osint:website', 'researcher about website', 1.0);

-- DIAGNOSIS ------------------------------------------------
INSERT INTO diagnosis_data_source(diagnosis_data_source_id, name, description, max_score) VALUES
  (1002, 'diagnosis:wpscan', 'Vulnerability scan for WordPress', 10.0),
  (1003, 'diagnosis:portscan', 'Portscan for ip/fqdn', 10.0),
  (1004, 'diagnosis:application-scan', 'vulnerability scan for web application', 10.0);

-- CODE ------------------------------------------------
INSERT INTO code_data_source(code_data_source_id, name, description, max_score) VALUES
  (1001, 'code:gitleaks', 'Credential scanning for GitHub', 1.0),
  (1002, 'code:dependency', 'Dependency vulnerability scanning for GitHub', 1.0);

-- GOOGLE ------------------------------------------------
INSERT INTO google_data_source(google_data_source_id, name, description, max_score) VALUES
  (1001, 'google:asset', 'Cloud Asset Inventory', 1.0),
  (1002, 'google:cloudsploit', 'Aqua CloudSploit for GCP', 1.0),
  (1003, 'google:scc', 'Security Command Center', 1.0),
  (1004, 'google:portscan', 'Portscan for GCP', 10.0);

commit;
