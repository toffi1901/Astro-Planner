CREATE TYPE role_enum AS ENUM ('guest', 'amateur', 'student', 'scientist', 'admin');

CREATE TYPE satellite_enum AS ENUM ('iss', 'starlink', 'scientific', 'weather', 'communication', 'other');

CREATE TYPE report_status_enum AS ENUM ('in_process', 'finished', 'frozen');

CREATE TYPE celestial_enum AS ENUM ('planet', 'star', 'constellation', 'galaxy', 'other');

CREATE TYPE notification_type_enum AS ENUM ('system', 'reminder');

CREATE TYPE target_enum AS ENUM ('user', 'satellite', 'observation', 'comment', 'session');

CREATE TYPE channel_enum AS ENUM ('email', 'telegram');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    user_id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email          varchar(255) NOT NULL UNIQUE,
    password_hash  varchar(255) NOT NULL,
    role           role_enum NOT NULL DEFAULT 'guest',
    latitude       double precision,
    longitude      double precision,
    university     varchar(255),
    organization   varchar(255),
    is_verified    boolean NOT NULL DEFAULT false,
    is_blocked     boolean NOT NULL DEFAULT false,
    register_date  timestamptz NOT NULL DEFAULT now(),
    last_login     timestamptz
);

CREATE TABLE satellites (
    norad_id         bigint PRIMARY KEY,
    satellite_name   varchar(100) NOT NULL,
    satellite_type   satellite_enum,
    tle_line1        varchar(80),
    tle_line2        varchar(80),
    last_updated_tle timestamptz,
    is_active        boolean NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE celestial_objects (
    object_id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         varchar(100) NOT NULL UNIQUE,
    object_type  celestial_enum NOT NULL,
    description  text,
    url_image    varchar(500)
);

CREATE TABLE forecasts (
    forecast_id  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    norad_id     bigint NOT NULL REFERENCES satellites(norad_id),
    latitude     double precision NOT NULL,
    longitude    double precision NOT NULL,
    start_time   timestamptz NOT NULL,
    end_time     timestamptz NOT NULL,
    max_altitude double precision,
    magnitude    double precision,
    azimuth      double precision,
    created_at   timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_forecast UNIQUE (norad_id, latitude, longitude, start_time)
);

CREATE TABLE favorites (
    favorite_id  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    norad_id     bigint REFERENCES satellites(norad_id),
    latitude     double precision NOT NULL,
    longitude    double precision NOT NULL,
    point_name   varchar(100),
    created_at   timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_favorite UNIQUE (user_id, norad_id, latitude, longitude)
);

CREATE TABLE observations (
    observation_id      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id             uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    norad_id            bigint REFERENCES satellites(norad_id),
    celestial_object_id bigint REFERENCES celestial_objects(object_id),
    observation_date    timestamptz NOT NULL,
    latitude            double precision,
    longitude           double precision,
    note                text,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CHECK (
        (norad_id IS NOT NULL AND celestial_object_id IS NULL) OR
        (norad_id IS NULL AND celestial_object_id IS NOT NULL)
    )
);

CREATE TABLE reports (
    report_id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id           uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    title             varchar(255) NOT NULL,
    content           text NOT NULL,
    status            report_status_enum NOT NULL DEFAULT 'in_process',
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE comments (
    comment_id        bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    norad_id          bigint NOT NULL REFERENCES satellites(norad_id) ON DELETE CASCADE,
    user_id           uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content           text NOT NULL,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tle_history (
    tle_note_id  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    norad_id     bigint NOT NULL REFERENCES satellites(norad_id) ON DELETE CASCADE,
    tle_line1    varchar(80) NOT NULL,
    tle_line2    varchar(80) NOT NULL,
    fetched_at   timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_tle_history UNIQUE (norad_id, fetched_at)
);

CREATE TABLE audit_log (
    log_id      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    admin_id    uuid NOT NULL REFERENCES users(user_id),
    action      varchar(50) NOT NULL,
    target_type target_enum NOT NULL,
    target_id   text NOT NULL,
    details     jsonb,
    ip_address  varchar(45),
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE notifications (
    notification_id   bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type              notification_type_enum NOT NULL,
    payload           jsonb NOT NULL DEFAULT '{}',
    created_at        timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE notification_recipients (
    recipient_id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    notification_id bigint NOT NULL REFERENCES notifications(notification_id) ON DELETE CASCADE,
    user_id         uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    channel         channel_enum NOT NULL DEFAULT 'email',
    scheduled_at    timestamptz,
    sent_at         timestamptz,
    is_sent         boolean NOT NULL DEFAULT false,
    CONSTRAINT uq_recipient UNIQUE (notification_id, user_id, channel)
);

CREATE TABLE tags (
    tag_id   bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id  uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name     varchar(50) NOT NULL,
    CONSTRAINT uq_tag UNIQUE (user_id, name)
);

CREATE TABLE observation_tags (
    observation_id bigint NOT NULL REFERENCES observations(observation_id) ON DELETE CASCADE,
    tag_id         bigint NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (observation_id, tag_id)
);