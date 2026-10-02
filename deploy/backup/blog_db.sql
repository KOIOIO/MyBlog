-- MySQL dump 10.13  Distrib 8.0.46, for Linux (aarch64)
--
-- Host: localhost    Database: blog_db
-- ------------------------------------------------------
-- Server version	8.0.46

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
-- Current Database: `blog_db`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `blog_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `blog_db`;

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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_advertisements_delete_at` (`delete_at`),
  KEY `fk_advertisements_image` (`ad_image`),
  KEY `idx_advertisements_deleted_at` (`deleted_at`),
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
INSERT INTO `article_categories` VALUES ('技术',13),('生活',3);
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_article_likes_delete_at` (`delete_at`),
  KEY `fk_article_likes_user` (`user_id`),
  KEY `idx_article_likes_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_article_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=105 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_likes`
--

LOCK TABLES `article_likes` WRITE;
/*!40000 ALTER TABLE `article_likes` DISABLE KEYS */;
INSERT INTO `article_likes` VALUES (104,'2026-09-28 19:14:41.223','2026-09-28 19:14:41.223',NULL,'15',1,NULL);
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
INSERT INTO `article_tags` VALUES ('Docker',3),('Elasticsearch',1),('Gin',2),('Go',5),('GORM',1),('JWT',2),('Kubernetes',1),('MySQL',2),('Nginx',2),('Redis',2),('TypeScript',1),('Vite',2),('Vue',1),('Webpack',1),('中间件',2),('云原生',2),('人工智能',2),('代码规范',1),('前端',3),('博客',1),('后端',4),('安全',1),('工作日常',2),('工程化',2),('开发工具',1),('开源',1),('微服务',4),('性能优化',3),('搜索',1),('数据库',3),('数码',1),('旅行',1),('架构',4),('生活随笔',1),('索引',1),('经验',1),('缓存',1),('编程语言',1),('网络安全',1),('美食',1),('认证',2),('读书笔记',1),('负载均衡',1),('运维',2),('部署',3),('随笔',2);
/*!40000 ALTER TABLE `article_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `blog_tags`
--

DROP TABLE IF EXISTS `blog_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `blog_tags` (
  `tag` varchar(191) NOT NULL,
  `group` varchar(20) NOT NULL DEFAULT 'tech',
  `number` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `blog_tags`
--

LOCK TABLES `blog_tags` WRITE;
/*!40000 ALTER TABLE `blog_tags` DISABLE KEYS */;
INSERT INTO `blog_tags` VALUES ('Docker','tech',-1),('Elasticsearch','tech',0),('Gin','tech',0),('Go','tech',0),('GORM','tech',2),('JWT','tech',0),('Kubernetes','tech',0),('MySQL','tech',0),('Nginx','tech',0),('Redis','tech',1),('TypeScript','tech',0),('Vite','tech',0),('Vue','tech',0),('Webpack','tech',0),('中间件','tech',-1),('云原生','tech',-1),('人工智能','tech',0),('代码规范','tech',0),('前端','tech',0),('前端开发','tech',0),('博客','life',0),('后端','tech',1),('后端开发','tech',1),('安全','tech',0),('工作日常','life',0),('工程化','tech',0),('开发工具','tech',0),('开源','tech',0),('影视','life',0),('微服务','tech',-1),('性能优化','tech',-1),('情感杂谈','life',0),('搜索','tech',0),('数据库','tech',-1),('数码','life',0),('旅行','life',0),('架构','tech',-1),('架构设计','tech',0),('生活随笔','life',0),('算法','tech',0),('索引','tech',0),('经验','life',0),('缓存','tech',0),('编程语言','tech',0),('网络安全','tech',0),('美食','life',0),('认证','tech',0),('读书笔记','life',0),('负载均衡','tech',0),('运动健身','life',0),('运维','tech',0),('部署','tech',-1),('随笔','life',0),('音乐','life',0);
/*!40000 ALTER TABLE `blog_tags` ENABLE KEYS */;
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comments_delete_at` (`delete_at`),
  KEY `fk_comments_children` (`p_id`),
  KEY `fk_comments_user` (`user_uuid`),
  KEY `idx_comments_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_comments_children` FOREIGN KEY (`p_id`) REFERENCES `comments` (`id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_uuid`) REFERENCES `users` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (53,'2026-09-28 20:11:35.878','2026-09-28 20:11:35.878',NULL,'1',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','你好','2026-09-28 20:11:59.656'),(54,'2026-09-28 20:11:43.626','2026-09-28 20:11:43.626',NULL,'1',53,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈','2026-09-28 20:11:59.650'),(55,'2026-09-29 08:36:26.780','2026-09-29 08:36:26.780',NULL,'15',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','真的可以',NULL),(56,'2026-09-29 08:36:34.777','2026-09-29 08:36:34.777',NULL,'15',55,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈\n',NULL),(57,'2026-09-29 08:38:31.233','2026-09-29 08:38:31.233',NULL,'15',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','![](/emoji/xiaochun_emoji_13.png)![](/emoji/xiaochun_emoji_14.png)','2026-09-29 11:07:23.465'),(58,'2026-09-29 08:45:36.970','2026-09-29 08:45:36.970',NULL,'3',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈😉','2026-09-29 08:45:45.420'),(59,'2026-09-29 08:45:42.255','2026-09-29 08:45:42.255',NULL,'3',58,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈','2026-09-29 08:45:44.499'),(60,'2026-09-29 11:04:09.349','2026-09-29 11:04:09.349',NULL,'16',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','真的很有用',NULL),(61,'2026-09-29 11:07:31.303','2026-09-29 11:07:31.303',NULL,'15',56,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈\n',NULL),(62,'2026-09-29 11:11:34.320','2026-09-29 11:11:34.320',NULL,'12',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','你好',NULL),(63,'2026-09-29 11:11:43.839','2026-09-29 11:11:43.839',NULL,'12',62,'fbd5364d-bb15-11f1-b230-16b7a2303b52','哈哈',NULL);
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_feedbacks_delete_at` (`delete_at`),
  KEY `fk_feedbacks_user` (`user_uuid`),
  KEY `idx_feedbacks_deleted_at` (`deleted_at`),
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
-- Table structure for table `forum_comments`
--

DROP TABLE IF EXISTS `forum_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `forum_comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned NOT NULL,
  `parent_id` bigint unsigned NOT NULL DEFAULT '0',
  `user_id` bigint unsigned NOT NULL,
  `content` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_forum_comments_delete_at` (`delete_at`),
  KEY `idx_forum_comments_deleted_at` (`deleted_at`),
  KEY `idx_forum_comments_post_id` (`post_id`),
  KEY `idx_forum_comments_parent_id` (`parent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `forum_comments`
--

LOCK TABLES `forum_comments` WRITE;
/*!40000 ALTER TABLE `forum_comments` DISABLE KEYS */;
INSERT INTO `forum_comments` VALUES (1,'2026-09-28 17:57:16.713','2026-09-28 17:57:16.713',NULL,'2026-09-28 18:49:56.435',2,0,1,'你好');
/*!40000 ALTER TABLE `forum_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `forum_likes`
--

DROP TABLE IF EXISTS `forum_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `forum_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `post_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_forum_likes_post_user` (`post_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `forum_likes`
--

LOCK TABLES `forum_likes` WRITE;
/*!40000 ALTER TABLE `forum_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `forum_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `forum_posts`
--

DROP TABLE IF EXISTS `forum_posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `forum_posts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `category` varchar(50) NOT NULL DEFAULT 'æŠ€æœ¯',
  `tags` varchar(500) DEFAULT NULL,
  `images` json DEFAULT NULL,
  `like_count` int NOT NULL DEFAULT '0',
  `comment_count` int NOT NULL DEFAULT '0',
  `view_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_forum_posts_delete_at` (`delete_at`),
  KEY `idx_forum_posts_deleted_at` (`deleted_at`),
  KEY `idx_forum_posts_user_id` (`user_id`),
  KEY `idx_forum_posts_category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `forum_posts`
--

LOCK TABLES `forum_posts` WRITE;
/*!40000 ALTER TABLE `forum_posts` DISABLE KEYS */;
INSERT INTO `forum_posts` VALUES (1,'2026-09-28 17:43:05.966','2026-09-28 17:57:22.772',NULL,'2026-09-28 18:49:58.819',1,'Docker是什么','哈哈','技术','[\"JWT\",\"GORM\",\"Go\",\"Redis\"]','[\"/uploads/forum/92cc16f485a51463c9f30e92be435374-20260928174305.jpg\"]',0,0,9),(2,'2026-09-28 17:53:13.593','2026-09-28 17:57:25.557',NULL,'2026-09-28 18:49:56.436',1,'哈哈','你好','生活','[\"工作日常\",\"影视\"]','[\"/uploads/forum/92cc16f485a51463c9f30e92be435374-20260928175308.jpg\"]',1,1,6),(5,'2026-09-28 18:54:34.304','2026-09-28 18:54:50.550',NULL,'2026-09-28 18:54:55.917',1,'哈哈','哈哈','技术','[\"Docker\",\"Gin\"]','[]',0,0,2),(6,'2026-09-28 18:55:29.151','2026-09-28 18:55:29.183',NULL,'2026-09-28 18:55:44.428',1,'你好','其实我','技术','[\"Gin\"]','[]',0,0,1);
/*!40000 ALTER TABLE `forum_posts` ENABLE KEYS */;
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_friend_links_delete_at` (`delete_at`),
  KEY `fk_friend_links_image` (`logo`),
  KEY `idx_friend_links_deleted_at` (`deleted_at`),
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_images_url` (`url`),
  KEY `idx_images_delete_at` (`delete_at`),
  KEY `idx_images_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `images`
--

LOCK TABLES `images` WRITE;
/*!40000 ALTER TABLE `images` DISABLE KEYS */;
INSERT INTO `images` VALUES (25,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'01.jpg','/uploads/image/01.jpg',3,0,NULL),(26,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'02.jpg','/uploads/image/02.jpg',3,0,NULL),(27,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'03.jpg','/uploads/image/03.jpg',3,0,NULL),(28,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'04.jpg','/uploads/image/04.jpg',3,0,NULL),(29,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'05.jpg','/uploads/image/05.jpg',3,0,NULL),(30,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'06.jpg','/uploads/image/06.jpg',3,0,NULL),(31,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'07.jpg','/uploads/image/07.jpg',3,0,NULL),(32,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'08.jpg','/uploads/image/08.jpg',3,0,NULL),(33,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'09.jpg','/uploads/image/09.jpg',3,0,NULL),(34,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'10.jpg','/uploads/image/10.jpg',3,0,NULL),(35,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'11.jpg','/uploads/image/11.jpg',3,0,NULL),(36,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'12.jpg','/uploads/image/12.jpg',3,0,NULL),(37,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'13.jpg','/uploads/image/13.jpg',3,0,NULL),(38,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'14.jpg','/uploads/image/14.jpg',3,0,NULL),(39,'2026-09-28 11:10:18.000','2026-09-28 11:10:18.000',NULL,'15.jpg','/uploads/image/15.jpg',3,0,NULL),(40,'2026-09-28 11:19:37.000','2026-09-28 11:19:37.000',NULL,'16.jpg','/uploads/image/16.jpg',3,0,NULL);
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_jwt_blacklists_delete_at` (`delete_at`),
  KEY `idx_jwt_blacklists_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jwt_blacklists`
--

LOCK TABLES `jwt_blacklists` WRITE;
/*!40000 ALTER TABLE `jwt_blacklists` DISABLE KEYS */;
INSERT INTO `jwt_blacklists` VALUES (1,'2026-09-28 16:27:47.275','2026-09-28 16:27:47.275',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE4ODgwMn0.wnvjvNaS2sWUJCZyFc7YFfysxPGy0rDgHSDENoip9Wk',NULL),(2,'2026-09-28 17:42:19.869','2026-09-28 17:42:19.869',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE4ODg2N30.MSD5mKFMNpSTl9X0sJVilJ77yXkwIyGE-P62CaAcZ2k',NULL),(3,'2026-09-28 17:58:03.002','2026-09-28 17:58:03.002',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE5MzMzOX0.hcDete-I4wrdtD8gCmjklL1ydWkNulVkvA4w64YmzI0',NULL),(4,'2026-09-28 18:18:54.522','2026-09-28 18:18:54.522',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjIsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE5NTM3Mn0.kcRkuwg9E642Tl8yuzDEkuyZ-y-5bsNZDW3yidZAVeo',NULL),(5,'2026-09-28 18:21:17.027','2026-09-28 18:21:17.027',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjExLCJpc3MiOiJnb19ibG9nIiwiYXVkIjpbIlRBUCJdLCJleHAiOjE3OTExOTQ0MzJ9.kAa47DKj6jEnDJonAwm1pBBqGEuEd_DSOL6uM3R8Wj4',NULL),(6,'2026-09-28 18:27:33.384','2026-09-28 18:27:33.384',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE5NTYwN30.QIH0UbYwBGFL-WlFJUG3yXLwZCZxNnxkZoZQwzEsVuA',NULL),(7,'2026-09-29 08:45:22.399','2026-09-29 08:45:22.399',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTE5NjA1M30.GLezGnd0rtFsKvgF8yaoJeq2WZtvcIA6cqOrGm3jyrk',NULL),(8,'2026-09-29 11:03:57.755','2026-09-29 11:03:57.755',NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsImlzcyI6ImdvX2Jsb2ciLCJhdWQiOlsiVEFQIl0sImV4cCI6MTc5MTI0NzUyMn0.7Yo3nqjRPIFFdTouW5aztBGJvswmKEzMNtk0EuFMiDY',NULL);
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_logins_delete_at` (`delete_at`),
  KEY `fk_logins_user` (`user_id`),
  KEY `idx_logins_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_logins_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `logins`
--

LOCK TABLES `logins` WRITE;
/*!40000 ALTER TABLE `logins` DISABLE KEYS */;
INSERT INTO `logins` VALUES (1,'2026-09-28 16:26:42.788','2026-09-28 16:26:42.788',NULL,1,'email','127.0.0.1','未知','Mac OS X','Mac','Chrome',200,NULL),(2,'2026-09-28 16:27:47.748','2026-09-28 16:27:47.748',NULL,1,'email','127.0.0.1','未知','Mac OS X','Mac','Chrome',200,NULL),(3,'2026-09-28 17:42:20.347','2026-09-28 17:42:20.347',NULL,1,'email','127.0.0.1','河南省-郑州市','Mac OS X','Mac','Chrome',200,NULL),(9,'2026-09-28 18:20:07.944','2026-09-28 18:20:07.944',NULL,1,'email','127.0.0.1','河南省-郑州市','Mac OS X','Mac','Chrome',200,NULL),(10,'2026-09-28 18:27:33.838','2026-09-28 18:27:33.838',NULL,1,'email','127.0.0.1','河南省-郑州市','Mac OS X','Mac','Chrome',200,NULL),(11,'2026-09-29 08:45:22.535','2026-09-29 08:45:22.535',NULL,1,'email','127.0.0.1','未知','Mac OS X','Mac','Chrome',200,NULL),(12,'2026-09-29 11:03:57.892','2026-09-29 11:03:57.892',NULL,1,'email','127.0.0.1','未知','Mac OS X','Mac','Chrome',200,NULL);
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
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_uuid` (`uuid`),
  KEY `idx_users_delete_at` (`delete_at`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2026-09-28 08:24:11.895','2026-09-28 17:13:43.472',NULL,'fbd5364d-bb15-11f1-b230-16b7a2303b52','Xiaoyu_Wang','$2a$10$lKSsQ6V3aSCpx05.2FY7w.oapNLUPCz12vjyKfdRprQ1AyOOT8LPi','2652777599@qq.com',NULL,'/uploads/avatar/92cc16f485a51463c9f30e92be435374-20260928171343.jpg','河南省郑州市','没有理想的人不伤心',2,0,0,NULL);
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

-- Dump completed on 2026-09-29  4:49:25
