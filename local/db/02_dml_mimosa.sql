set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- CORE ------------------------------------------------
INSERT INTO role(role_id, name, project_id) VALUES
  (1, 'system-admin-role', null);

INSERT INTO policy(policy_id, name, project_id, action_ptn, resource_ptn) VALUES
  (1, 'system-admin-policy', null, '.*', '.*');

INSERT INTO role_policy(role_id, policy_id, project_id) VALUES
  (1, 1, null);

-- AWS ------------------------------------------------
INSERT INTO aws_data_source(aws_data_source_id, data_source, max_score) VALUES
  (1001, 'aws:guard-duty', 10.0),
  (1002, 'aws:access-analyzer', 1.0),
  (1003, 'aws:admin-checker', 1.0);

-- OSINT ------------------------------------------------
INSERT INTO osint_data_source(osint_data_source_id, name, description, max_score) VALUES
  (1001, 'osint:subdomainsearch', 'subdomain researcher', 10.0);

-- DIAGNOSIS ------------------------------------------------
INSERT INTO diagnosis_data_source(diagnosis_data_source_id, name, description, max_score) VALUES
  (1001, 'diagnosis:jira', 'jira', 10.0);

commit;
