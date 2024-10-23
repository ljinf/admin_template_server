CREATE TABLE `sys_user`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id`      int(11) NOT NULL COMMENT '用户ID',
    `dept_id`      int(20) DEFAULT NULL COMMENT '部门ID',
    `username`     varchar(30) NOT NULL COMMENT '用户账号',
    `nick_name`    varchar(30) NOT NULL COMMENT '用户昵称',
    `user_type`    varchar(2)   DEFAULT '00' COMMENT '用户类型（00系统用户）',
    `email`        varchar(50)  DEFAULT '' COMMENT '用户邮箱',
    `phone_number` varchar(11)  DEFAULT '' COMMENT '手机号码',
    `sex`          char(1)      DEFAULT '0' COMMENT '用户性别（0男 1女 2未知）',
    `avatar`       varchar(100) DEFAULT '' COMMENT '头像地址',
    `password`     varchar(100) DEFAULT '' COMMENT '密码',
    `status`       char(1)      DEFAULT '0' COMMENT '帐号状态（0正常 1停用）',
    `login_ip`     varchar(128) DEFAULT '' COMMENT '最后登录IP',
    `login_date`   datetime     DEFAULT NULL COMMENT '最后登录时间',
    `create_by`    varchar(64)  DEFAULT '' COMMENT '创建者',
    `created_at`   datetime     DEFAULT NULL,
    `updated_at`   datetime     DEFAULT NULL,
    `deleted_at`   datetime     DEFAULT NULL,
    `remark`       varchar(500) DEFAULT NULL COMMENT '备注',
    PRIMARY KEY (`id`),
    key            user_idx(`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';


CREATE TABLE `sys_menu`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `menu_id`    int(11) NOT NULL COMMENT '菜单ID',
    `menu_name`  varchar(50) NOT NULL COMMENT '菜单名称',
    `parent_id`  bigint(20) DEFAULT '0' COMMENT '父菜单ID',
    `order_num`  int(4) DEFAULT '0' COMMENT '显示顺序',
    `path`       varchar(200) DEFAULT '' COMMENT '路由地址',
    `component`  varchar(255) DEFAULT NULL COMMENT '组件路径',
    `query`      varchar(255) DEFAULT NULL COMMENT '路由参数',
    `is_frame`   int(1) DEFAULT '1' COMMENT '是否为外链（0是 1否）',
    `is_cache`   int(1) DEFAULT '0' COMMENT '是否缓存（0缓存 1不缓存）',
    `menu_type`  char(1)      DEFAULT '' COMMENT '菜单类型（M目录 C菜单 F按钮）',
    `visible`    char(1)      DEFAULT '0' COMMENT '菜单状态（0显示 1隐藏）',
    `status`     char(1)      DEFAULT '0' COMMENT '菜单状态（0正常 1停用）',
    `perms`      varchar(100) DEFAULT NULL COMMENT '权限标识',
    `icon`       varchar(100) DEFAULT '#' COMMENT '菜单图标',
    `create_by`  varchar(64)  DEFAULT '' COMMENT '创建者',
    `created_at` datetime     DEFAULT NULL,
    `updated_at` datetime     DEFAULT NULL,
    `deleted_at` datetime     DEFAULT NULL,
    `remark`     varchar(500) DEFAULT '' COMMENT '备注',
    PRIMARY KEY (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单权限表';


-- 古典诗词四 个组成部份  全诗词  重点句 作者简介 赏析
-- 现代诗  全诗句 重点句
-- 句读两个部份就一段话 出处

DROP TABLE IF EXISTS `ps_poem_list`;
CREATE TABLE `ps_poem_list`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `poem_type`    INT(11) NOT NULL DEFAULT 1 COMMENT '类型，1古诗  3现代诗  3句读',
    `author`       VARCHAR(128) NOT NULL DEFAULT '0' COMMENT '作者\出处',
    `introduction` TEXT COMMENT '作者简介',
    `title`        VARCHAR(128) NOT NULL COMMENT '标题',
    `content`      TEXT COMMENT '内容',
    `extract`      TEXT COMMENT '摘录\重点句',
    `summary`      TEXT COMMENT '赏析',
    `created_at`   INT(11) DEFAULT NULL COMMENT '创建时间',
    `updated_at`   INT(11) DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='诗词句读';
