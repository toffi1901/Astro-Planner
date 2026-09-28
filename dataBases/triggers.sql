CREATE OR REPLACE FUNCTION trg_set_updated_date()
RETURNS trigger 
LANGUAGE plpgsql
AS $$
BEGIN 
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

CREATE TRIGGER observation_update()
AFTER UPDATE ON observations
FOR EACH ROW 
EXECUTE FUCNTION trg_set_updated_date();

CREATE TRIGGER trg_report_update()
AFTER UPDATE ON reports
FOR EACH ROW
EXECUTE FUNCTION trg_set_updated_date();

CREATE TRIGGER trg_comment_update()
AFTER UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION trg_set_updated_date();

CREATE OR REPLACE FUNCTION update_satellite_tle_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE satellites 
    SET last_updated_tle = NEW.fetched_at 
    WHERE norad_id = NEW.norad_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_tle_history_insert
    AFTER INSERT ON tle_history
    FOR EACH ROW EXECUTE FUNCTION update_satellite_tle_timestamp();