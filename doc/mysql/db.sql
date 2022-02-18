CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '用户名',
  `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '密码',
  `profile_picture` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '头像',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (1, 'Khighness', '$2a$12$7uFkvyQfIMIub/CE89yy.eQqBDhyRYHNFeyWJUhGef597XdFJD/bC', 'http://127.0.0.1:10000/avatar/default.jpg');
INSERT INTO entry_task.user (id, username, password, profile_picture) VALUES (2, 'FlowerK', '$2a$12$Mr1AmXijhEZ1IgM9HjkvPepQgP/TorE/migfzLsCE5Mh6Y84ysdsq', 'http://127.0.0.1:10000/avatar/youyu.jpeg');
