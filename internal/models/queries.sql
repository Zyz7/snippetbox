-- Para instalar y configurar mysql

--sudo apt install mysql-server
--sudo mysql_secure_installation
--sudo mysql

-- Dentro de mysql
SELECT user,authentication_string,plugin,host FROM mysql.user;
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'password';
ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'password';
GRANT ALL PRIVILEGES ON * . * TO 'newuser'@'localhost';
FLUSH PRIVILEGES;
SELECT user,authentication_string,plugin,host FROM mysql.user;

-------
-- mysql -u root -p

show databases;
create database snippetbox character set utf8mb4 collate utf8mb4_unicode_ci;

use snippetbox;

create table snippets (
    id integer not null primary key auto_increment,
    title varchar(100) not null,
    content text not null,
    created datetime not null,
    expires datetime not null
);

create index idx_snippets_created on snippets(created);

insert into snippets (title, content, created, expires) values (
    'An old silent pond',
    'An old silent pond...\nAfrog jumps into the pond,\nsplash! Silence again',
    utc_timestamp(),
    date_add(utc_timestamp(), interval 365 day)
);

insert into snippets (title, content, created, expires) values (
    'Over the wintry forest',
    'Over the wintry...\nforest winds howl in a rage\nwith no leaves to blow',
    utc_timestamp(),
    date_add(utc_timestamp(), interval 365 day)
);

create table sessions (
    token char(43) primary key,
    data blob not null,
    expiry timestamp(6) not null
);

create index sessions_expiry_idx on sessions (expiry);

select * from snippets;
select * from sesiones;

------
create table users (
    id integer not null primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created datetime not null 
);

alter table users add constraint users_uc_email unique (email);

------
-- Creamos un nuevo usuario

create user 'web'@'localhost' identified by 'password';
grant select, insert, update, delete on snippetbox.* to 'web'@'localhost';

-- mysql -D snippetbox -u web -p


