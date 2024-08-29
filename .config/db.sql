create database if not exists insider;
use insider;

drop table if exists messages;
create table messages (
    id serial primary key,
    content varchar(1000) not null,
    recipient_phone varchar(20) not null,
    created_at timestamp default current_timestamp,
    sent_at timestamp,

    index sent_at_index (sent_at)
);