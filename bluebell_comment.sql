-- 评论表
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `post_id` bigint(20) NOT NULL COMMENT '帖子ID',
    `author_id` bigint(20) NOT NULL COMMENT '评论作者ID',
    `content` text NOT NULL COMMENT '评论内容',
    `parent_id` bigint(20) DEFAULT 0 COMMENT '父评论ID，0表示顶级评论',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1正常 0删除',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='评论表';
