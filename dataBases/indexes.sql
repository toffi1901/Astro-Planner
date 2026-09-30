CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_satellites_active ON satellites(is_active);

CREATE INDEX idx_forecasts_satellite_time ON forecasts(norad_id, start_time);
CREATE INDEX idx_forecasts_location ON forecasts(latitude, longitude);

CREATE INDEX idx_favorites_user ON favorites(user_id);

CREATE INDEX idx_observations_user_date ON observations(user_id, observation_date);
CREATE INDEX idx_observations_satellite ON observations(norad_id);

CREATE INDEX idx_reports_user_status ON reports(user_id, status);

CREATE INDEX idx_comments_satellite ON comments(norad_id);
CREATE INDEX idx_comments_user ON comments(user_id);

CREATE INDEX idx_tle_history_satellite ON tle_history(norad_id, fetched_at DESC);

CREATE INDEX idx_audit_log_admin ON audit_log(admin_id, created_at DESC);
CREATE INDEX idx_audit_log_target ON audit_log(target_type, target_id);

CREATE INDEX idx_notification_recipients_user ON notification_recipients(user_id, is_sent);
CREATE INDEX idx_notification_recipients_scheduled ON notification_recipients(scheduled_at) WHERE NOT is_sent;

CREATE INDEX idx_tags_user ON tags(user_id);
CREATE INDEX idx_observation_tags_tag ON observation_tags(tag_id);