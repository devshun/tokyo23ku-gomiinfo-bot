CREATE TABLE regions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  ward_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (ward_id) REFERENCES wards(id)
  UNIQUE KEY uc_regions (ward_id, name)
);
