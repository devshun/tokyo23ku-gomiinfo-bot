CREATE TABLE garbage_days (
  id INT AUTO_INCREMENT PRIMARY KEY,
  region_id INT,
  garbage_type INT NOT NULL,
  day_of_week INT NOT NULL,
  week_number_of_month INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (region_id) REFERENCES regions(id)
  CONSTRAINT garbage_days_unique UNIQUE (region_id, garbage_type, day_of_week, week_number_of_month)
);
