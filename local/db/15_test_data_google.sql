set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO gcp(gcp_id, name, project_id, gcp_project_id) VALUES
  (1001, "my-gcp",  1001, 'my-project');

INSERT INTO gcp_data_source(gcp_id, google_data_source_id, project_id) VALUES
  (1001, 1001, 1001);

commit;
