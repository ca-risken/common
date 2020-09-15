CREATE DATABASE IF NOT EXISTS mimosa DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
use mimosa;

-- CORE ------------------------------------------------

CREATE TABLE user (
  user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  sub VARCHAR(255) NOT NULL,
  name VARCHAR(64) NOT NULL,
  activated ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id),
  UNIQUE KEY uidx_sub (sub)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE user_role (
  user_id INT UNSIGNED NOT NULL,
  role_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, role_id)
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

CREATE TABLE finding (
  finding_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  data_source VARCHAR(64) NOT NULL,
  data_source_id VARCHAR(255) NOT NULL,
  resource_name VARCHAR(255) NOT NULL,
  project_id INT UNSIGNED NULL,
  original_score FLOAT(5,2) UNSIGNED NOT NULL,
  score FLOAT(3,2) UNSIGNED NULL,
  data JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_id),
  UNIQUE KEY uidx_data_source (project_id, data_source, data_source_id),
  INDEX idx_resource_name(resource_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE finding_tag (
  finding_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  finding_id BIGINT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  tag_key VARCHAR(64) NOT NULL,
  tag_value VARCHAR(200) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(finding_tag_id),
  UNIQUE KEY uidx_finding_tag (finding_id, tag_key)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE resource (
  resource_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_name VARCHAR(255) NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(resource_id),
  UNIQUE KEY uidx_resource (project_id, resource_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE resource_tag (
  resource_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  resource_id BIGINT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  tag_key VARCHAR(64) NOT NULL,
  tag_value VARCHAR(200) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(resource_tag_id),
  UNIQUE KEY uidx_resource_tag (resource_id, tag_key),
  INDEX idx_tag_key(tag_key)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert (
  alert_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  alert_condition_id INT UNSIGNED NOT NULL,
  description VARCHAR(200) NULL,
  severity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NULL,
  activated ENUM('false', 'true') NOT NULL DEFAULT 'true',
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
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_history_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE rel_alert_finding (
  alert_id INT UNSIGNED NOT NULL,
  finding_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_id, finding_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_condition (
  alert_condition_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  severity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NULL,
  and_or ENUM('and', 'or') NOT NULL DEFAULT 'and',
  enabled ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_cond_rule (
  alert_condition_id INT UNSIGNED NOT NULL,
  alert_rule_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, alert_rule_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_rule (
  alert_rule_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NULL,
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
  project_id INT UNSIGNED NULL,
  cache_second INT UNSIGNED NOT NULL DEFAULT 1800,
  notified_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, notification_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE notification (
  notification_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NULL,
  type VARCHAR(64) NOT NULL,
  notify_setting JSON NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(notification_id)
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
  UNIQUE KEY uidx_aws_account_id (aws_account_id)
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
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(aws_id, aws_data_source_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;
