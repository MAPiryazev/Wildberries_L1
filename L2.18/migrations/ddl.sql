create table Events(
    id serial primary key,
    user_id integer not null,
    date date not null,
    title text not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);