set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO code_gitleaks(gitleaks_id, code_data_source_id, name, project_id, type, target_resource, repository_pattern, github_user, personal_access_token, gitleaks_config, status, status_detail, scan_at) VALUES
  (1001, 1001, "gitleaksk-name", 1001, 'USER', 'gitleakstest', 'gronit', '', '', '', 'CONFIGURED', '', null);

commit;
