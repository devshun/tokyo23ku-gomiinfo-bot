CREATE TABLE IF NOT EXISTS garbage_collection_schedules (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  region VARCHAR(255),
  burnable_garbage_days VARCHAR(255),
  non_burnable_garbage_days VARCHAR(255),
  recyclable_garbage_days VARCHAR(255)
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
);
