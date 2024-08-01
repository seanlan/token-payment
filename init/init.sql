# ************************************************************
# Sequel Ace SQL dump
# 版本号： 20058
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主机: 127.0.0.1 (MySQL 8.3.0)
# 数据库: tokenpay
# 生成时间: 2024-08-01 03:35:20 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 转储表 admin_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_group`;

CREATE TABLE `admin_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '组名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分组表';



# 转储表 admin_group_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_group_permission`;

CREATE TABLE `admin_group_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL COMMENT '组ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分组权限表';



# 转储表 admin_group_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_group_user`;

CREATE TABLE `admin_group_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '用户ID',
  `group_id` int NOT NULL COMMENT '组ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分组成员表';



# 转储表 admin_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_log`;

CREATE TABLE `admin_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `asset_id` bigint NOT NULL COMMENT '资源ID',
  `asset_type` int NOT NULL COMMENT '资源类型',
  `remarks` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志',
  `create_at` bigint NOT NULL COMMENT '记录时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='日志记录';



# 转储表 admin_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_permission`;

CREATE TABLE `admin_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名称',
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限代码',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限类型',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';



# 转储表 admin_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_user`;

CREATE TABLE `admin_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `account` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '电子邮箱',
  `is_super` int NOT NULL COMMENT '是否超级用户',
  `is_frozen` int NOT NULL COMMENT '是否被封禁',
  `is_delete` int NOT NULL COMMENT '是否已删除',
  `hide_article` int NOT NULL COMMENT '是否隐藏稿件',
  `create_at` bigint NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';



# 转储表 admin_user_message
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_user_message`;

CREATE TABLE `admin_user_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `message` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息',
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源路径',
  `is_read` int NOT NULL COMMENT '是否处理',
  `create_at` bigint NOT NULL COMMENT '消息时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 admin_user_permission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_user_permission`;

CREATE TABLE `admin_user_permission` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户权限';



# 转储表 admin_user_token
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin_user_token`;

CREATE TABLE `admin_user_token` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户id',
  `token` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token',
  `expire_at` bigint NOT NULL DEFAULT '0' COMMENT 'token过期时间',
  `refresh_token` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'refresh token',
  `refresh_expire_at` bigint NOT NULL DEFAULT '0' COMMENT 'refresh token过期时间',
  `create_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户Token信息表';



# 转储表 application
# ------------------------------------------------------------

DROP TABLE IF EXISTS `application`;

CREATE TABLE `application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'app key',
  `app_secret` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'app secret',
  `app_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'app name',
  `hook_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知url',
  `create_at` bigint NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 application_chain
# ------------------------------------------------------------

DROP TABLE IF EXISTS `application_chain`;

CREATE TABLE `application_chain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `application_id` bigint NOT NULL COMMENT '应用id',
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `hot_wallet` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '热钱包地址',
  `cold_wallet` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '冷钱包地址',
  `fee_wallet` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '零钱整理费用钱包',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 application_withdraw_order
# ------------------------------------------------------------

DROP TABLE IF EXISTS `application_withdraw_order`;

CREATE TABLE `application_withdraw_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `application_id` bigint NOT NULL COMMENT '应用ID',
  `serial_no` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单序列号',
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `contract_address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币合约地址，如果是空表示是主币',
  `symbol` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币符号',
  `to_address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收款地址',
  `value` decimal(65,30) NOT NULL COMMENT '数量',
  `gas_price` bigint NOT NULL COMMENT 'gas费用',
  `token_id` bigint NOT NULL COMMENT 'tokenid （NFT）',
  `tx_hash` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交易hash值',
  `nonce` bigint NOT NULL DEFAULT '-1' COMMENT '交易nonce',
  `hook` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '到账变动通知url',
  `create_at` bigint NOT NULL COMMENT '申请时间',
  `transfer_at` bigint NOT NULL COMMENT '转账时间',
  `generated` int NOT NULL COMMENT '是否生成',
  `transfer_success` tinyint NOT NULL COMMENT '是否转账成功',
  `transfer_failed_times` int NOT NULL COMMENT '转账失败次数',
  `transfer_next_time` bigint NOT NULL COMMENT '下次转账时间',
  `received` tinyint NOT NULL COMMENT '是否到账',
  `receive_at` bigint NOT NULL COMMENT '到账时间',
  PRIMARY KEY (`id`),
  KEY `chain_symbol` (`chain_symbol`,`contract_address`,`symbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain`;

CREATE TABLE `chain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain_symbol` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链名称',
  `chain_id` bigint NOT NULL COMMENT '链ID',
  `currency` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '货币',
  `chain_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'evm' COMMENT '链类型 默认evm',
  `confirm` int NOT NULL COMMENT '确认区块数量',
  `gas` bigint NOT NULL COMMENT 'gas费用配置',
  `gas_price` bigint NOT NULL COMMENT 'gas price 配置',
  `latest_block` bigint NOT NULL COMMENT '最新区块',
  `rebase_block` bigint NOT NULL COMMENT '重新构建区块',
  `has_branch` tinyint NOT NULL COMMENT '是否出现分叉',
  `concurrent` int NOT NULL COMMENT '并发量',
  `address_pool` int NOT NULL COMMENT '地址池',
  `watch` tinyint NOT NULL COMMENT '是否监听',
  PRIMARY KEY (`id`),
  UNIQUE KEY `chain_symbol` (`chain_symbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain_address
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain_address`;

CREATE TABLE `chain_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `application_id` bigint NOT NULL COMMENT '应用ID',
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '地址',
  `enc_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '加密后的私钥',
  `hook` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '到账变动通知url',
  `watch` int NOT NULL COMMENT '是否监听',
  `create_at` bigint NOT NULL COMMENT '创建时间',
  `used` int NOT NULL COMMENT '是否已使用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `chain_symbol` (`chain_symbol`,`address`),
  KEY `application_id` (`application_id`),
  KEY `watch` (`watch`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain_block
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain_block`;

CREATE TABLE `chain_block` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `block_number` bigint NOT NULL COMMENT '区块高度',
  `block_hash` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '区块hash值',
  `parent_hash` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '上一个区块hash值',
  `checked` int NOT NULL COMMENT '是否检测完成',
  PRIMARY KEY (`id`),
  UNIQUE KEY `chain_symbol_2` (`chain_symbol`,`block_number`),
  KEY `chain_symbol` (`chain_symbol`),
  KEY `block_number` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain_rpc
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain_rpc`;

CREATE TABLE `chain_rpc` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `rpc_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'rpc地址',
  `disable` int NOT NULL COMMENT '是否可用',
  PRIMARY KEY (`id`),
  KEY `chain_symbol` (`chain_symbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain_token
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain_token`;

CREATE TABLE `chain_token` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `contract_address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币合约地址，如果是空表示是主币',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '币种名称',
  `symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '币种符号',
  `decimals` int NOT NULL COMMENT '小数位',
  `threshold` decimal(65,30) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# 转储表 chain_tx
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chain_tx`;

CREATE TABLE `chain_tx` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `application_id` bigint NOT NULL COMMENT '应用ID',
  `chain_symbol` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链的符号',
  `block_number` bigint NOT NULL COMMENT '区块高度',
  `block_hash` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '区块hash值',
  `tx_hash` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交易hash值',
  `from_address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '支付地址',
  `to_address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收款地址',
  `contract_address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币合约地址，如果是空表示是主币',
  `symbol` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代币符号',
  `value` decimal(65,30) NOT NULL COMMENT '数量',
  `token_id` bigint NOT NULL COMMENT 'tokenid （NFT）',
  `tx_index` bigint NOT NULL COMMENT '交易序号',
  `batch_index` bigint NOT NULL COMMENT '交易批次号',
  `confirm` int NOT NULL COMMENT '确认次数',
  `removed` int NOT NULL COMMENT '是否已移除',
  `transfer_type` int NOT NULL COMMENT '交易类型 1到账 2提币',
  `arranged` int NOT NULL COMMENT '是否整理过',
  `create_at` bigint NOT NULL COMMENT '交易时间',
  `serial_no` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单序列号',
  `notify_success` int NOT NULL COMMENT '是否通知成功',
  `notify_failed_times` int NOT NULL COMMENT '通知失败次数',
  `notify_next_time` bigint NOT NULL COMMENT '下次通知时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `chain_symbol` (`chain_symbol`,`block_hash`,`tx_hash`,`tx_index`,`batch_index`,`transfer_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
