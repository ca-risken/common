set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO google_gcp(gcp_id, google_data_source_id, name, project_id, gcp_organization_id, gcp_project_id, scan_at) VALUES
  (1001, 1001, "CloudAsset",  1001, 'my-org', 'my-project', null),
  (1002, 1002, "CloudSploit", 1001, 'my-org', 'my-project', null);

commit;
