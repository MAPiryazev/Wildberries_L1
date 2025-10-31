create table if not exists shortcuts(
    id serial primary key,
    original_url text unique not null,
    short_code varchar(12) unique not null,
    client_id UUID null,
    expires_at timestamp null,
    created_at timestamp default now() not null  
);

create index if not exists idx_shortcuts_short_code on shortcuts(short_code);

create index if not exists idx_shortcuts_client_id on shortcuts(client_id);