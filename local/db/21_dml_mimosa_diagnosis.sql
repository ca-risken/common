set character_set_connection = utf8mb4;
set character_set_client = utf8mb4;

use mimosa;

INSERT INTO diagnosis(diagnosis_id, project_id, name) VALUES
  (1, 1001, 'project-a-diagnosis'),
  (2, 1001, 'project-a-diagnosis_2'),
  (3, 1001, 'project-a-diagnosis_3');

INSERT INTO diagnosis_data_source(diagnosis_data_source_id, name, description, max_score) VALUES
  (1, 'jira', 'jira', 10.0);

INSERT INTO rel_diagnosis_data_source(rel_diagnosis_data_source_id, diagnosis_data_source_id, diagnosis_id, project_id, record_id, jira_id, jira_key) VALUES
  (1, 1, 1, 1001, '1353', '', ''),
  (2, 1, 1, 1001, '', '10241', ''),
  (3, 1, 1, 1001, '', '', 'OC202001');

commit;
