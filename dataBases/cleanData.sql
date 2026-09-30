CREATE OR REPLACE FUNCTION cleanup_old_data()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM forecasts WHERE start_time < now() - interval '7 days';

    DELETE FROM tle_history WHERE fetched_at < now() - interval '1 year';

    DELETE FROM notification_recipients
    WHERE is_sent = true AND sent_at < now() - interval '30 days';
END;
$$;