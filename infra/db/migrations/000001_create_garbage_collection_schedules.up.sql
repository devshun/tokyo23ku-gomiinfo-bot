CREATE TABLE IF NOT EXISTS garbage_collection_schedules (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  region VARCHAR(255),
  burnable_garbage_days VARCHAR(255),
  non_burnable_garbage_days VARCHAR(255),
  recyclable_garbage_days VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);
