-- 建表
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '用户名',
  `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '密码',
  `profile_picture` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '头像',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 初始化
INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (1, 'Khighness', '$2a$12$7uFkvyQfIMIub/CE89yy.eQqBDhyRYHNFeyWJUhGef597XdFJD/bC', 'http://127.0.0.1:10000/avatar/khighness.jpg');
INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (2, 'FlowerK', '$2a$12$Mr1AmXijhEZ1IgM9HjkvPepQgP/TorE/migfzLsCE5Mh6Y84ysdsq', 'http://127.0.0.1:10000/avatar/flowerk.jpeg');

-- 函数：产生n位随机名字
DELIMITER $$
CREATE FUNCTION rand_name(n INT) RETURNS VARCHAR(255)
BEGIN
    DECLARE chars_str VARCHAR(100) DEFAULT '@0123456789abcdefghijklmnopqrstuvwsyzABCDEFGHIJKLMNOPQRSTUVWXYZ=';
    DECLARE return_str VARCHAR(255) DEFAULT '';
    DECLARE i INT DEFAULT 0;
    WHILE i < n DO
    SET return_str = CONCAT(return_str,SUBSTRING(chars_str,FLOOR(1+RAND()*64),1));
    SET i = i + 1;
    END WHILE;
    RETURN return_str;
END $$


-- 存储过程：从起始start_id+1开始，插入max_num条数据
DELIMITER $$
CREATE PROCEDURE insert_user(IN start_id INT(10),IN max_num INT(10))
BEGIN
DECLARE i INT DEFAULT 0;
    SET autocommit = 0;
    REPEAT
    SET i = i + 1;
    INSERT INTO user(id,username,password,profile_picture) VALUES((start_id+i),rand_name(12),'$2a$12$DOrlHGuWsNMs4GdY/E3d9edISCpZfX1PnycD6iP.1P4kAJLIwmtZu','http://127.0.0.1:10000/avatar/khighness.jpg');
    UNTIL i = max_num
    END REPEAT;
    COMMIT;
END $$

-- 插入1000,0000条数据
call insert_user(10, 10000000)