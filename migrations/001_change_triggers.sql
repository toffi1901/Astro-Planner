DROP TRIGGER IF EXISTS trg_comment_update ON comments;

CREATE TRIGGER trg_comment_update
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION trg_set_updated_date();

ALTER TABLE users
ADD CONSTRAINT chk_users_latitude
CHECK (latitude IS NULL OR latitude BETWEEN -90 AND 90);

ALTER TABLE users
ADD CONSTRAINT chk_users_longitude
CHECK (longitude IS NULL OR longitude BETWEEN -180 AND 180);

ALTER TABLE forecasts
ADD CONSTRAINT chk_forecast_time
CHECK (end_time > start_time);