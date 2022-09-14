-- 建表
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) CHARACTER SET utf8mb4  DEFAULT '' COMMENT '用户名',
  `password` varchar(60) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '密码',
  `profile_picture` varchar(100) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '头像',
  PRIMARY KEY (`id`),
  UNIQUE INDEX index_username(`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 初始化
INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (1, 'KHighness', 'cea13832a6a48e6b83c472a50ca55934', 'http://127.0.0.1:10000/avatar/show/khighness.jpg');
INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (2, 'FlowerK', 'cea13832a6a48e6b83c472a50ca55934', 'http://127.0.0.1:10000/avatar/show/flowerk.jpeg');

-- 存储过程：从起始start_id开始，插入max_num条数据
DELIMITER $$
CREATE PROCEDURE insert_user(IN start_id INT(10),IN max_num INT(10))
BEGIN
    DECLARE i INT DEFAULT 0;
    SET @user_prefix = 'user_';
    SET autocommit = 0;
    REPEAT
        INSERT INTO entry_task.user(id,username,password,profile_picture) VALUES((start_id+i),CONCAT(@user_prefix,start_id+i),'e10adc3949ba59abbe56e057f20f883e','http://127.0.0.1:10000/avatar/default.jpg');
        SET i = i + 1;
    UNTIL i = max_num
END REPEAT;
COMMIT;
END $$

-- 插入1000,0000条数据
call insert_user(3, 10000000)

-- 函数：产生n位随机名字
-- DELIMITER $$
-- CREATE FUNCTION rand_name(n INT) RETURNS VARCHAR(255)
-- BEGIN
--     DECLARE chars_str VARCHAR(100) DEFAULT '@0123456789abcdefghijklmnopqrstuvwsyzABCDEFGHIJKLMNOPQRSTUVWXYZ=';
--     DECLARE return_str VARCHAR(255) DEFAULT '';
--     DECLARE i INT DEFAULT 0;
--     WHILE i < n DO
--             SET return_str = CONCAT(return_str,SUBSTRING(chars_str,FLOOR(1+RAND()*64),1));
--             SET i = i + 1;
-- END WHILE;
-- RETURN return_str;
-- END $$

-- 最大连接数
-- show variables like '%max_connections%';
-- 服务器响应的最大连接数
-- show global status like 'Max_used_connections';
-- 设置最大连接数
-- set global max_connections = 10000;
-- 客户端连接
-- show processlist;
-- 客户端连接ip数
-- select SUBSTRING_INDEX(host,':',1) as ip , count(*) from information_schema.processlist group by ip
