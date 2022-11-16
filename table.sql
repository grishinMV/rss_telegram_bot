create table feeds
(
    id         int auto_increment
        primary key,
    link       varchar(256)  not null,
    last_new   int default 0 not null,
    next_parse int default 0 null
);

create table users
(
    id           int auto_increment
        primary key,
    telegram_id  int         not null,
    chat_id      int         not null,
    last_message varchar(32) null,
    constraint unique_telegram_id
        unique (telegram_id)
);

create table users_feeds
(
    id      int auto_increment
        primary key,
    user_id int not null,
    feed_id int not null,
    constraint user_feed
        unique (user_id, feed_id),
    constraint feed
        foreign key (feed_id) references feeds (id)
            on update cascade on delete cascade,
    constraint user
        foreign key (user_id) references users (id)
            on update cascade on delete cascade
);

