INSERT INTO shortener_analytics.click_analytics (short_code, client_id, user_agent, ip, timestamp)
VALUES
(
  'abc123',
  generateUUIDv4(),
  'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
  '192.168.0.1',
  now()
);
