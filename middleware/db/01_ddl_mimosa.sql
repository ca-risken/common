CREATE DATABASE IF NOT EXISTS mimosa DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
use mimosa;

-- CORE ------------------------------------------------
CREATE TABLE user (
  user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  sub VARCHAR(255) NOT NULL,
  name VARCHAR(64) NOT NULL,
  user_idp_key VARCHAR(255) NULL,
  activated ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id),
  UNIQUE KEY uidx_sub (sub),
  UNIQUE KEY uidx_user_idp_key (user_idp_key)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE user_reserved (
  reserved_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_idp_key VARCHAR(255) NOT NULL,
  role_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(reserved_id),
  UNIQUE KEY uidx_user_reserved (user_idp_key, role_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE user_role (
  user_id INT UNSIGNED NOT NULL,
  role_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, role_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE access_token (
  access_token_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  token_hash VARCHAR(255) NOT NULL,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(255) NULL,
  project_id INT UNSIGNED NOT NULL,
  expired_at DATETIME NOT NULL DEFAULT '9999-12-31 23:59:59',
  last_updated_user_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(access_token_id),
  UNIQUE KEY uidx_access_token (project_id, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE access_token_role (
  access_token_id INT UNSIGNED NOT NULL,
  role_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(access_token_id, role_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE role (
  role_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(role_id),
  UNIQUE KEY uidx_role (project_id, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE role_policy (
  role_id INT UNSIGNED NOT NULL,
  policy_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(role_id, policy_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE policy (
  policy_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  project_id INT UNSIGNED NULL,
  action_ptn TEXT NOT NULL,
  resource_ptn TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(policy_id),
  UNIQUE KEY uidx_policy (project_id, name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE project (
  project_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(project_id),
  UNIQUE KEY uidx_project (name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE project_tag (
  project_id INT UNSIGNED NOT NULL,
  tag VARCHAR(512) NOT NULL,
  color VARCHAR(32) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(project_id, tag)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE finding (
  finding_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  data_source VARCHAR(64) NOT NULL,
  data_source_id VARCHAR(255) NOT NULL,
  resource_name VARCHAR(255) NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  original_score FLOAT(5,2) UNSIGNED NOT NULL,
  score FLOAT(3,2) UNSIGNED NULL,
  data JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_id),
  UNIQUE KEY uidx_data_source (project_id, data_source, data_source_id),
  INDEX idx_score(project_id, score, updated_at)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE finding_tag (
  finding_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  finding_id BIGINT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  tag VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_tag_id),
  UNIQUE KEY uidx_finding_tag (finding_id, tag),
  INDEX idx_finding_tag(tag)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE INDEX idx_project_id ON finding_tag (project_id, tag);

CREATE TABLE resource (
  resource_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_name VARCHAR(255) NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(resource_id),
  UNIQUE KEY uidx_resource (project_id, resource_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE resource_tag (
  resource_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_id BIGINT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  tag VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(resource_tag_id),
  UNIQUE KEY uidx_resource_tag (resource_id, tag),
  INDEX idx_resource_tag(project_id, tag)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE pend_finding (
  finding_id BIGINT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  note VARCHAR(128) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE finding_setting (
  finding_setting_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id INT UNSIGNED NOT NULL,
  resource_name VARCHAR(255) NOT NULL,
  setting JSON NOT NULL,
  status ENUM('UNKNOWN', 'ACTIVE', 'DEACTIVE') NOT NULL DEFAULT 'UNKNOWN',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_setting_id),
  UNIQUE KEY uidx_finding_setting (project_id, resource_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE recommend (
  recommend_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  data_source VARCHAR(64) NOT NULL,
  type VARCHAR(128) NOT NULL,
  risk TEXT NULL,
  recommendation TEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(recommend_id),
  UNIQUE KEY uidx_recommend (data_source, type)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE recommend_finding (
  finding_id BIGINT UNSIGNED NOT NULL,
  recommend_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert (
  alert_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  alert_condition_id INT UNSIGNED NOT NULL,
  description VARCHAR(200) NULL,
  severity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NOT NULL,
  status ENUM('ACTIVE', 'PENDING', 'DEACTIVE') NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_history (
  alert_history_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  history_type ENUM('created', 'updated', 'deleted') NOT NULL,
  alert_id INT UNSIGNED NOT NULL,
  description VARCHAR(200) NULL,
  severity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  finding_history JSON NULL,
  project_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_history_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE rel_alert_finding (
  alert_id INT UNSIGNED NOT NULL,
  finding_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_id, finding_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_condition (
  alert_condition_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  severity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NOT NULL,
  and_or ENUM('and', 'or') NOT NULL DEFAULT 'and',
  enabled boolean NOT NULL DEFAULT true,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_cond_rule (
  alert_condition_id INT UNSIGNED NOT NULL,
  alert_rule_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, alert_rule_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_rule (
  alert_rule_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NOT NULL,
  score FLOAT(3,2) UNSIGNED NOT NULL,
  resource_name VARCHAR(255) NULL,
  tag VARCHAR(64) NULL,
  finding_cnt INT UNSIGNED NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_rule_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_cond_notification (
  alert_condition_id INT UNSIGNED NOT NULL,
  notification_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  cache_second INT UNSIGNED NOT NULL DEFAULT 1800,
  notified_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, notification_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE notification (
  notification_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NOT NULL,
  type VARCHAR(64) NOT NULL,
  notify_setting JSON NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(notification_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE report_finding (
  report_finding_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  report_date DATE NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  data_source VARCHAR(64) NOT NULL,
  score FLOAT(3,2) UNSIGNED NOT NULL,
  count INT UNSIGNED NOT NULL,
  PRIMARY KEY(report_finding_id),
  UNIQUE KEY uidx_report_finding (report_date, project_id, data_source, score)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

-- AWS ------------------------------------------------
CREATE TABLE aws (
  aws_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NOT NULL,
  aws_account_id VARCHAR(12) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(aws_id),
  UNIQUE KEY uidx_aws (project_id, aws_account_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE aws_data_source (
  aws_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  data_source VARCHAR(64) NOT NULL,
  max_score FLOAT(5,2) UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(aws_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE aws_rel_data_source (
  aws_id INT UNSIGNED NOT NULL,
  aws_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  assume_role_arn VARCHAR(255) NOT NULL,
  external_id VARCHAR(255) NULL,
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(aws_id, aws_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

-- OSINT ------------------------------------------------
CREATE TABLE osint (
  osint_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id INT UNSIGNED NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_name VARCHAR(200) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(osint_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE osint_data_source (
  osint_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL,
  description VARCHAR(200) NOT NULL,
  max_score FLOAT(5,2) UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(osint_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE rel_osint_data_source (
  rel_osint_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  osint_data_source_id INT UNSIGNED NOT NULL,
  osint_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(rel_osint_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE osint_detect_word (
  osint_detect_word_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  rel_osint_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  word VARCHAR(50) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(osint_detect_word_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

-- DIAGNOSIS ------------------------------------------------
CREATE TABLE diagnosis_data_source (
  diagnosis_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL,
  description VARCHAR(200) NOT NULL,
  max_score FLOAT(5,2) UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(diagnosis_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE wpscan_setting (
  wpscan_setting_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  diagnosis_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  target_url VARCHAR(200),
  options JSON NULL,
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(wpscan_setting_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE portscan_setting (
  portscan_setting_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  diagnosis_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  name VARCHAR(200),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(portscan_setting_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE portscan_target (
  portscan_target_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  portscan_setting_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  target VARCHAR(255),
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(portscan_target_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE application_scan (
  application_scan_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  diagnosis_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  name VARCHAR(200),
  scan_type ENUM('NOT_CONFIGURED', 'BASIC') NOT NULL DEFAULT 'BASIC',
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(application_scan_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE application_scan_basic_setting (
  application_scan_basic_setting_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  application_scan_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  target VARCHAR(255),
  max_depth INT UNSIGNED NOT NULL,
  max_children INT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(application_scan_basic_setting_id),
  UNIQUE KEY uidx_diagnosis_application_scan_basic_setting (application_scan_id, project_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;


-- CODE ------------------------------------------------
CREATE TABLE code_data_source (
  code_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(128) NOT NULL,
  max_score FLOAT(5,2) UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(code_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE code_github_setting (
  code_github_setting_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  project_id INT UNSIGNED NOT NULL,
  name VARCHAR(64) NULL,
  github_user VARCHAR(64) NULL,
  personal_access_token VARCHAR(255) NULL,
  type ENUM('UNKNOWN_TYPE', 'ORGANIZATION', 'USER') NOT NULL DEFAULT 'UNKNOWN_TYPE',
  base_url VARCHAR(128) NULL,
  target_resource VARCHAR(128) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(code_github_setting_id),
  UNIQUE KEY uidx_code_github_setting (name, project_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE code_gitleaks_setting (
  code_github_setting_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  code_data_source_id INT UNSIGNED NOT NULL,
  repository_pattern VARCHAR(128) NULL,
  scan_public ENUM('false', 'true') NOT NULL DEFAULT 'true',
  scan_internal ENUM('false', 'true') NOT NULL DEFAULT 'true',
  scan_private ENUM('false', 'true') NOT NULL DEFAULT 'false',
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(code_github_setting_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE code_gitleaks_cache (
  code_github_setting_id INT UNSIGNED NOT NULL,
  repository_full_name VARCHAR(255) NOT NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(code_github_setting_id, repository_full_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE code_dependency_setting (
  code_github_setting_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  code_data_source_id INT UNSIGNED NOT NULL,
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL,
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(code_github_setting_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

-- GOOGLE ------------------------------------------------
CREATE TABLE google_data_source (
  google_data_source_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(128) NOT NULL,
  max_score FLOAT(5,2) UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(google_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE gcp (
  gcp_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NULL,
  project_id INT UNSIGNED NOT NULL,
  gcp_organization_id VARCHAR(128) NULL,
  gcp_project_id VARCHAR(128) NOT NULL,
  verification_code VARCHAR(128) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(gcp_id),
  UNIQUE KEY uidx_gcp (project_id, gcp_project_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE gcp_data_source (
  gcp_id INT UNSIGNED NOT NULL,
  google_data_source_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NOT NULL,
  status ENUM('UNKNOWN', 'OK' ,'CONFIGURED', 'IN_PROGRESS', 'ERROR') NOT NULL DEFAULT 'UNKNOWN',
  status_detail VARCHAR(255) NULL,
  scan_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(gcp_id, google_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
