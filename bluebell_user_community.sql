-- 用户社区关系表
DROP TABLE IF EXISTS `user_community`;
CREATE TABLE `user_community` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `community_id` bigint(20) NOT NULL COMMENT '社区ID',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_community` (`user_id`, `community_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户社区关系表';
