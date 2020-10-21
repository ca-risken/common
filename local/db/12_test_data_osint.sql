set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- OSINT ------------------------------------------------
INSERT INTO osint(osint_id, project_id, name) VALUES
  (1001, 1001, 'project-a-osint'),
  (1002, 1001, 'project-a-osint_2');

INSERT INTO rel_osint_data_source(rel_osint_data_source_id, osint_data_source_id, osint_id, project_id, resource_type, resource_name, status) VALUES
  (1001, 1001, 1001, 1001, 'Domain', 'security-hub.jp', 'CONFIGURED'),
  (1002, 1001, 1001, 1001, 'DnsRecord', 'security-hub.jp', 'CONFIGURED');

commit;
