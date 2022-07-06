set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO code_gitleaks(gitleaks_id, code_data_source_id, name, project_id, type, target_resource, repository_pattern, github_user, personal_access_token, gitleaks_config, status, status_detail, scan_at) VALUES
  (1001, 1001, "gitleaksk-name", 1001, 'USER', 'gitleakstest', 'gronit', '', '', '', 'CONFIGURED', '', null);

INSERT INTO code_github_setting(code_github_setting_id, project_id, name, github_user, personal_access_token, type, base_url, target_resource) VALUES
  (1001, 1001, 'code-github-setting-name', 'github-user', 'github-personal-access-token', 'USER', '', 'username');

INSERT INTO code_gitleaks_setting(code_github_setting_id, project_id, code_data_source_id, repository_pattern, scan_public, scan_internal, scan_private, status, status_detail, scan_at) VALUES
  (1001, 1001, 1001, '', 'true', 'true', 'true', 'CONFIGURED', '', null);

commit;
