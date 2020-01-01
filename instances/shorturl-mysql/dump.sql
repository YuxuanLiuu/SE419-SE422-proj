create database shorturl;
use shorturl;
create table short(
`id` int(32) primary key auto_increment, 
`long` varchar(100)
);
