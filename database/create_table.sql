CREATE DATABASE ir;

create user 'tester' identified by '123456';

USE ir;

SELECT DATABASE();

SHOW TABLES;

-- 建表
create table if not exists user(
	id int auto_increment comment '用户id，自增',
	name varchar(20) not null comment '用户名',
	password char(32) not null comment '密码的md5',
    create_time datetime default current_timestamp comment '用户注册时间',
    update_time datetime default current_timestamp on update current_timestamp comment '最后修改时间',
	primary key (id),
  -- 用户名唯一索引，确保name是唯一的。
	unique key idx_name (name)
)default charset=utf8mb4 comment '用户信息';
