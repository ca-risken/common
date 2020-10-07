set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- OSINT ------------------------------------------------
INSERT INTO osint(osint_id, project_id, name) VALUES
  (1, 1001, 'project-a-osint'),
  (2, 1001, 'project-a-osint_2');

INSERT INTO osint_result(osint_result_id, osint_data_source_id, osint_id, project_id, resource_type, resource_name) VALUES
  (1, 1, 1, 1001, 'Domain', 'security-hub.jp'),
  (2, 1, 1, 1001, 'DnsRecord', 'security-hub.jp');

commit;
