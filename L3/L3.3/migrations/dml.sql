INSERT INTO users (username, email) VALUES
('admin', 'admin@example.com'),
('alice', 'alice@example.com'),
('bob', 'bob@example.com'),
('charlie', 'charlie@example.com');

-- корневые комменты
INSERT INTO comments (user_id, content) VALUES
(1, 'Всем привет! Это первый комментарий!'),
(2, 'Классная тема, согласен с автором.'),
(3, 'А мне кажется, что тут можно улучшить логику.');

-- ответы
INSERT INTO comments (user_id, content) VALUES
(4, 'Полностью согласен с тобой, Alice!'),
(2, 'Bob, можешь подробнее объяснить свою идею?'),
(1, 'Charlie, интересная мысль, расскажи подробнее.');

-- корневые комменты
INSERT INTO comment_paths (ancestor_id, descendant_id, depth) VALUES
(1, 1, 0),
(2, 2, 0),
(3, 3, 0),
(4, 4, 0),
(5, 5, 0),
(6, 6, 0);

-- пути для 1 → 4
INSERT INTO comment_paths (ancestor_id, descendant_id, depth) VALUES
(1, 4, 1);

-- пути для 3 → 5
INSERT INTO comment_paths (ancestor_id, descendant_id, depth) VALUES
(3, 5, 1);

-- пути для 3 → 6
INSERT INTO comment_paths (ancestor_id, descendant_id, depth) VALUES
(3, 6, 1);
