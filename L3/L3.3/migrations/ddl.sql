-- таблица пользователей 
create table users(
    id serial primary key,
    username varchar(50) not null unique,
    email varchar(100) unique,
    created_at timestamp default now()
);

-- таблица комментариев
create table comments(
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    content text not null,
    created_at timestamp default now(),
    updated_at timestamp
);

--  Closure Table (для того чтобы востанавливать пути до комметариев в дереве)
create table comment_paths(
    ancestor_id int not null references comments(id) on delete cascade,
    descendant_id int not null references comments(id) on delete cascade,
    depth int not null,
    primary key (ancestor_id, descendant_id)
);

-- индексы
-- для быстрого для поиска предка и последователя
create index idx_comment_paths_ancestor on comment_paths(ancestor_id);
create index idx_comment_paths_descendant on comment_paths(descendant_id);

-- для поиска по словам
create index idx_comments_content_fts on comments using GIN(to_tsvector('russian', content));
create index idx_comments_content_fts_en on comments using GIN(to_tsvector('english', content));
