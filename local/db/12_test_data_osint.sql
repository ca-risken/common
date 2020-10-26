set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

-- OSINT ------------------------------------------------
INSERT INTO osint(osint_id, project_id, resource_type, resource_name) VALUES
  (1001, 1001, 'Domain', 'security-hub.jp'),
  (1002, 1001, 'Domain', 'cyberagent.co.jp');

INSERT INTO rel_osint_data_source(rel_osint_data_source_id, osint_data_source_id, osint_id, project_id, status) VALUES
  (1001, 1001, 1001, 1001, 'CONFIGURED'),
  (1002, 1001, 1001, 1001, 'CONFIGURED');

INSERT INTO rel_osint_detect_word(rel_osint_detect_word_id,rel_osint_data_source_id,osint_detect_word_id, project_id) VALUES
  (1001, 1001, 1001, 1001),
  (1002, 1001, 1002, 1001);

INSERT INTO osint_detect_word(osint_detect_word_id,word,project_id) VALUES
  (1001, 'jenkins', 1001),
  (1002, 'stg', 1001),
  (1003, 'admin', 1001);

commit;
