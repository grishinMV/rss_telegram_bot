create table feeds
(
    id         INTEGER       not null
        primary key autoincrement,
    link       varchar(256)  not null,
    last_new   INT default 0 not null,
    next_parse INT default 0
);

create table users
(
    id           INTEGER
        primary key autoincrement,
    telegram_id  INT not null
        constraint unique_telegram_id
        unique,
    chat_id      INT not null,
    last_message varchar(32)
);

create table users_feeds
(
    id      INTEGER not null
        primary key autoincrement
        unique,
    user_id INTEGER not null,
    feed_id INTEGER not null
);

create unique index users_feeds_user_id_feed_id_uindex
    on users_feeds (user_id, feed_id);

