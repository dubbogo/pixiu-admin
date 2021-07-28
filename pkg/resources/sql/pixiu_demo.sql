/*pixiu_admin 数据库表设计demo示例*/

CREATE DATABASE  IF NOT EXISTS  `pixiu`  DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ;

USE `pixiu`;
-- ----------------------------
-- Table structure for pixiu_user
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_user`;
CREATE TABLE `pixiu_user` (
      `id` integer NOT NULL PRIMARY KEY AUTO_INCREMENT,
      `username` varchar(255) NOT NULL,
      `password` varchar(255) NOT NULL,
      `role` integer(5) NOT NULL DEFAULT 1,
      `enabled` integer(5) NOT NULL DEFAULT 0 COMMENT 'delete or not',
      `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
      `date_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);

-- ----------------------------
-- Table structure for pixiu_role
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_role`;
CREATE TABLE `pixiu_role` (
                              `id` integer NOT NULL PRIMARY KEY AUTO_INCREMENT,
                              `role_name` varchar(255) NOT NULL,
                              `description` varchar(255) NOT NULL,
                              `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
                              `date_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);

-- ----------------------------
-- Table structure for pixiu_user_role
-- ----------------------------
DROP TABLE IF EXISTS `pixiu_user_role`;
CREATE TABLE `pixiu_role` (
      `id` integer NOT NULL PRIMARY KEY AUTO_INCREMENT,
      `user_id` integer NOT NULL,
      `role_id` integer NOT NULL,
      `date_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
      `date_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time'
);