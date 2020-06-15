CREATE DATABASE IF NOT EXISTS mimosa DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
use mimosa;

CREATE TABLE user (
  user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  actevated ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id),
  INDEX idx_name(name)
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
  PRIMARY KEY(role_id)
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
  PRIMARY KEY(policy_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE project (
  project_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(project_id)
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
  serverity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NULL,
  actevated ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_condition (
  alert_condition_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  description VARCHAR(200) NULL,
  serverity ENUM('high', 'medium', 'low') NOT NULL DEFAULT 'low',
  project_id INT UNSIGNED NULL,
  enabled ENUM('false', 'true') NOT NULL DEFAULT 'true',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_cond_rule (
  alert_condition_id INT UNSIGNED NOT NULL,
  alert_rule_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  and_or ENUM('and', 'or') NOT NULL DEFAULT 'and',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, alert_rule_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_rule (
  alert_rule_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NULL,
  score FLOAT(3,2) UNSIGNED NULL,
  resource_name VARCHAR(255) NULL,
  tag_key VARCHAR(64) NOT NULL,
  tag_value VARCHAR(200) NOT NULL,
  finding_cnt INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_rule_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

CREATE TABLE alert_cond_action (
  alert_condition_id INT UNSIGNED NOT NULL,
  alert_action_id INT UNSIGNED NOT NULL,
  project_id INT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_condition_id, alert_action_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE alert_action (
  alert_action_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NULL,
  project_id INT UNSIGNED NULL,
  notification_id INT UNSIGNED NULL,
  casche_secound INT UNSIGNED NOT NULL DEFAULT 1800,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY(alert_action_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin AUTO_INCREMENT = 1001;

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
