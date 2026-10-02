-- MySQL dump 10.13  Distrib 9.2.0, for Linux (x86_64)
--
-- Host: localhost    Database: blog_db
-- ------------------------------------------------------
-- Server version	9.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `advertisements`
--

DROP TABLE IF EXISTS `advertisements`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `advertisements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `ad_image` varchar(255) DEFAULT NULL,
  `link` longtext,
  `title` longtext,
  `content` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_advertisements_delete_at` (`delete_at`),
  KEY `fk_advertisements_image` (`ad_image`),
  CONSTRAINT `fk_advertisements_image` FOREIGN KEY (`ad_image`) REFERENCES `images` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advertisements`
--

LOCK TABLES `advertisements` WRITE;
/*!40000 ALTER TABLE `advertisements` DISABLE KEYS */;
/*!40000 ALTER TABLE `advertisements` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_categories`
--

DROP TABLE IF EXISTS `article_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_categories` (
  `category` varchar(191) NOT NULL,
  `number` bigint DEFAULT NULL,
  PRIMARY KEY (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_categories`
--

LOCK TABLES `article_categories` WRITE;
/*!40000 ALTER TABLE `article_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `article_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_likes`
--

DROP TABLE IF EXISTS `article_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `article_id` longtext,
  `user_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_article_likes_delete_at` (`delete_at`),
  KEY `fk_article_likes_user` (`user_id`),
  CONSTRAINT `fk_article_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_likes`
--

LOCK TABLES `article_likes` WRITE;
/*!40000 ALTER TABLE `article_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `article_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_tags`
--

DROP TABLE IF EXISTS `article_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_tags` (
  `tag` varchar(191) NOT NULL,
  `number` bigint DEFAULT NULL,
  PRIMARY KEY (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_tags`
--

LOCK TABLES `article_tags` WRITE;
/*!40000 ALTER TABLE `article_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `article_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `article_id` longtext,
  `p_id` bigint unsigned DEFAULT NULL,
  `user_uuid` char(36) DEFAULT NULL,
  `content` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_comments_delete_at` (`delete_at`),
  KEY `fk_comments_children` (`p_id`),
  KEY `fk_comments_user` (`user_uuid`),
  CONSTRAINT `fk_comments_children` FOREIGN KEY (`p_id`) REFERENCES `comments` (`id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_uuid`) REFERENCES `users` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `feedbacks`
--

DROP TABLE IF EXISTS `feedbacks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `feedbacks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `user_uuid` char(36) DEFAULT NULL,
  `content` longtext,
  `reply` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_feedbacks_delete_at` (`delete_at`),
  KEY `fk_feedbacks_user` (`user_uuid`),
  CONSTRAINT `fk_feedbacks_user` FOREIGN KEY (`user_uuid`) REFERENCES `users` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `feedbacks`
--

LOCK TABLES `feedbacks` WRITE;
/*!40000 ALTER TABLE `feedbacks` DISABLE KEYS */;
/*!40000 ALTER TABLE `feedbacks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `footer_links`
--

DROP TABLE IF EXISTS `footer_links`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `footer_links` (
  `title` varchar(191) NOT NULL,
  `link` longtext,
  PRIMARY KEY (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `footer_links`
--

LOCK TABLES `footer_links` WRITE;
/*!40000 ALTER TABLE `footer_links` DISABLE KEYS */;
/*!40000 ALTER TABLE `footer_links` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friend_links`
--

DROP TABLE IF EXISTS `friend_links`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friend_links` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `logo` varchar(255) DEFAULT NULL,
  `link` longtext,
  `name` longtext,
  `description` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_friend_links_delete_at` (`delete_at`),
  KEY `fk_friend_links_image` (`logo`),
  CONSTRAINT `fk_friend_links_image` FOREIGN KEY (`logo`) REFERENCES `images` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friend_links`
--

LOCK TABLES `friend_links` WRITE;
/*!40000 ALTER TABLE `friend_links` DISABLE KEYS */;
/*!40000 ALTER TABLE `friend_links` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `images`
--

DROP TABLE IF EXISTS `images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `url` varchar(255) DEFAULT NULL,
  `category` bigint DEFAULT NULL,
  `storage` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_images_url` (`url`),
  KEY `idx_images_delete_at` (`delete_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `images`
--

LOCK TABLES `images` WRITE;
/*!40000 ALTER TABLE `images` DISABLE KEYS */;
/*!40000 ALTER TABLE `images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jwt_blacklists`
--

DROP TABLE IF EXISTS `jwt_blacklists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jwt_blacklists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `jwt` text,
  PRIMARY KEY (`id`),
  KEY `idx_jwt_blacklists_delete_at` (`delete_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jwt_blacklists`
--

LOCK TABLES `jwt_blacklists` WRITE;
/*!40000 ALTER TABLE `jwt_blacklists` DISABLE KEYS */;
/*!40000 ALTER TABLE `jwt_blacklists` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `logins`
--

DROP TABLE IF EXISTS `logins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `logins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `login_method` longtext,
  `ip` longtext,
  `address` longtext,
  `os` longtext,
  `device_info` longtext,
  `browser_info` longtext,
  `status` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_logins_delete_at` (`delete_at`),
  KEY `fk_logins_user` (`user_id`),
  CONSTRAINT `fk_logins_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `logins`
--

LOCK TABLES `logins` WRITE;
/*!40000 ALTER TABLE `logins` DISABLE KEYS */;
/*!40000 ALTER TABLE `logins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `uuid` char(36) DEFAULT NULL,
  `username` longtext,
  `password` longtext,
  `email` longtext,
  `openid` longtext,
  `avatar` varchar(255) DEFAULT NULL,
  `address` longtext,
  `signature` varchar(191) DEFAULT '签名是空白的，这位用户似乎比较低调。',
  `role_id` bigint DEFAULT NULL,
  `register` bigint DEFAULT NULL,
  `freeze` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_uuid` (`uuid`),
  KEY `idx_users_delete_at` (`delete_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-02-18 11:21:41

--
-- ===== 论坛模块 & 固定标签库（2026-09-28 追加）=====
--

--
-- Table structure for table `blog_tags`（固定标签库）
--
CREATE TABLE IF NOT EXISTS `blog_tags` (
  `tag` varchar(191) NOT NULL,
  `group` varchar(20) NOT NULL DEFAULT 'tech',
  `number` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Table structure for table `forum_posts`（论坛帖子）
--
CREATE TABLE IF NOT EXISTS `forum_posts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `category` varchar(50) NOT NULL DEFAULT '技术',
  `tags` varchar(500) DEFAULT NULL,
  `images` json DEFAULT NULL,
  `like_count` int NOT NULL DEFAULT 0,
  `comment_count` int NOT NULL DEFAULT 0,
  `view_count` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_forum_posts_delete_at` (`delete_at`),
  KEY `idx_forum_posts_deleted_at` (`deleted_at`),
  KEY `idx_forum_posts_user_id` (`user_id`),
  KEY `idx_forum_posts_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Table structure for table `forum_comments`（论坛评论，二级嵌套）
--
CREATE TABLE IF NOT EXISTS `forum_comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned NOT NULL,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `user_id` bigint unsigned NOT NULL,
  `content` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_forum_comments_delete_at` (`delete_at`),
  KEY `idx_forum_comments_deleted_at` (`deleted_at`),
  KEY `idx_forum_comments_post_id` (`post_id`),
  KEY `idx_forum_comments_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Table structure for table `forum_likes`（论坛点赞，防重复）
--
CREATE TABLE IF NOT EXISTS `forum_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_forum_likes_post_user` (`post_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 预置固定标签库 30 条
--
INSERT IGNORE INTO `blog_tags` (`tag`, `group`, `number`) VALUES
('Go', 'tech', 0),
('Vue', 'tech', 0),
('前端开发', 'tech', 0),
('后端开发', 'tech', 0),
('数据库', 'tech', 0),
('MySQL', 'tech', 0),
('Redis', 'tech', 0),
('Elasticsearch', 'tech', 0),
('微服务', 'tech', 0),
('云原生', 'tech', 0),
('Docker', 'tech', 0),
('Kubernetes', 'tech', 0),
('架构设计', 'tech', 0),
('性能优化', 'tech', 0),
('网络安全', 'tech', 0),
('算法', 'tech', 0),
('人工智能', 'tech', 0),
('开源', 'tech', 0),
('编程语言', 'tech', 0),
('开发工具', 'tech', 0),
('生活随笔', 'life', 0),
('读书笔记', 'life', 0),
('影视', 'life', 0),
('音乐', 'life', 0),
('旅行', 'life', 0),
('美食', 'life', 0),
('运动健身', 'life', 0),
('数码', 'life', 0),
('工作日常', 'life', 0),
('情感杂谈', 'life', 0);

--
-- 追加：对齐现有 mock 文章标签（避免"标签不存在"校验错误）
--
INSERT IGNORE INTO `blog_tags` (`tag`, `group`, `number`) VALUES
('后端', 'tech', 0),
('架构', 'tech', 0),
('运维', 'tech', 0),
('Nginx', 'tech', 0),
('前端', 'tech', 0),
('Vite', 'tech', 0),
('搜索', 'tech', 0),
('索引', 'tech', 0),
('缓存', 'tech', 0),
('认证', 'tech', 0),
('负载均衡', 'tech', 0),
('部署', 'tech', 0),
('工程化', 'tech', 0),
('安全', 'tech', 0),
('代码规范', 'tech', 0),
('中间件', 'tech', 0),
('Webpack', 'tech', 0),
('TypeScript', 'tech', 0),
('JWT', 'tech', 0),
('GORM', 'tech', 0),
('Gin', 'tech', 0),
('随笔', 'life', 0),
('博客', 'life', 0),
('经验', 'life', 0);
