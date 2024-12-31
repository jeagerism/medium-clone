-- down.sql

-- Truncate the data in the tables in reverse order (to avoid foreign key constraint issues)
TRUNCATE TABLE "article_tags" CASCADE;
TRUNCATE TABLE "tags" CASCADE;
TRUNCATE TABLE "tokens" CASCADE;
TRUNCATE TABLE "images" CASCADE;
TRUNCATE TABLE "follows" CASCADE;
TRUNCATE TABLE "comments" CASCADE;
TRUNCATE TABLE "articles" CASCADE;
TRUNCATE TABLE "payment_records" CASCADE;
TRUNCATE TABLE "user_subscriptions" CASCADE;
TRUNCATE TABLE "membership_plans" CASCADE;
TRUNCATE TABLE "users" CASCADE;
TRUNCATE TABLE "likes" CASCADE;
