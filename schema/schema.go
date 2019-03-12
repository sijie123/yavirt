package schema

import (
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/db"
)

var stmts = []string{
	`
CREATE TABLE IF NOT EXISTS guest_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  image_id INT NOT NULL,
  host_id INT NOT NULL,
  cpu INT NOT NULL,
  mem BIGINT NOT NULL,
  state VARCHAR(16) NOT NULL,
  transit_status VARCHAR(16) NOT NULL DEFAULT '',
  create_time INT NOT NULL,
  transit_time INT NOT NULL NOT NULL DEFAULT 0,
  update_time INT NOT NULL DEFAULT 0
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
	`
CREATE TABLE IF NOT EXISTS image_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  parent_id INT NOT NULL DEFAULT 0,
  size INT NOT NULL DEFAULT 0,
  image_name VARCHAR(64) NOT NULL DEFAULT '',
  host_id INT NOT NULL,
  state VARCHAR(16) NOT NULL,
  transit_status VARCHAR(16) NOT NULL DEFAULT '',
  create_time INT NOT NULL,
  transit_time INT NOT NULL DEFAULT 0,
  update_time INT NOT NULL DEFAULT 0,
  UNIQUE KEY idx_image_name (image_name)
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
	`
CREATE TABLE IF NOT EXISTS host_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  hostname VARCHAR(15) NOT NULL,
  host_type CHAR(4) NOT NULL,
  state VARCHAR(16) NOT NULL,
  cpu INT NOT NULL DEFAULT 1,
  mem BIGINT not NULL DEFAULT 1073741824,
  UNIQUE KEY idx_hostname (hostname)
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
	`
CREATE TABLE IF NOT EXISTS volume_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  format VARCHAR(8) NOT NULL NOT NULL DEFAULT '',
  capacity INT NOT NULL,
  volume_type VARCHAR(6) NOT NULL,
  host_id INT NOT NULL,
  state VARCHAR(16) NOT NULL,
  transit_status VARCHAR(16) NOT NULL DEFAULT '',
  create_time INT NOT NULL,
  transit_time INT NOT NULL DEFAULT 0,
  update_time INT NOT NULL DEFAULT 0
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
	`
CREATE TABLE IF NOT EXISTS guest_volume_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  guest_id INT NOT NULL,
  volume_id INT NOT NULL,
  KEY idx_guest (guest_id),
  KEY idx_vlume (volume_id)
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
	`
CREATE TABLE IF NOT EXISTS addr_tab (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  high_value BIGINT NOT NULL DEFAULT 0,
  low_value BIGINT NOT NULL,
  prefix INT NOT NULL,
  gateway VARCHAR(16) NOT NULL DEFAULT '',
  guest_id INT NOT NULL,
  addr_type VARCHAR(16) NOT NULL,
  state VARCHAR(16) NOT NULL,
  host_id INT NOT NULL,
  KEY idx_state (state),
  UNIQUE KEY idx_value (low_value, high_value, addr_type)
) ENGINE=InnoDB COLLATE utf8mb4_unicode_ci`,
}

func InitSchema() error {
	for _, st := range stmts {
		if _, err := db.Exec(st); err != nil {
			return errors.Annotatef(err, st)
		}
	}
	return nil
}
