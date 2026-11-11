DROP TRIGGER IF EXISTS observation_update ON observations;
DROP TRIGGER IF EXISTS trg_report_update ON reports;
DROP TRIGGER IF EXISTS trg_comment_update ON comments;
DROP TRIGGER IF EXISTS trg_tle_history_insert ON tle_history;
DROP FUNCTION IF EXISTS trg_set_updated_date();
DROP FUNCTION IF EXISTS update_satellite_tle_timestamp();