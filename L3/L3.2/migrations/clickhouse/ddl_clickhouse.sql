CREATE DATABASE IF NOT EXISTS shortener_analytics;

CREATE TABLE IF NOT EXISTS shortener_analytics.click_analytics(
    short_code String,
    client_id UUID,
    user_agent String,
    ip String,
    timestamp DateTime default now()
)
ENGINE = MergeTree()
ORDER BY (short_code, timestamp);
