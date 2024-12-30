-- up.sql

-- Create membership_plans table
CREATE TABLE "membership_plans" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(100) NOT NULL,
  "price" DECIMAL(10,2) NOT NULL,
  "duration_days" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create user_subscriptions table
CREATE TABLE "user_subscriptions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "plan_id" INT NOT NULL,
  "subscription_id" VARCHAR(255),
  "status" VARCHAR(50) NOT NULL DEFAULT 'active',
  "start_date" TIMESTAMP NOT NULL,
  "end_date" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create payment_records table
CREATE TABLE "payment_records" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "subscription_id" INT NOT NULL,
  "payment_gateway" VARCHAR(50) NOT NULL,
  "payment_id" VARCHAR(255) NOT NULL,
  "amount" DECIMAL(10,2) NOT NULL,
  "currency" VARCHAR(10) NOT NULL,
  "payment_date" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create users table
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(100),
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "password_hash" TEXT NOT NULL,
  "bio" TEXT,
  "profile_image" TEXT,
  "role" VARCHAR(50) NOT NULL DEFAULT 'user',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create articles table
CREATE TABLE "articles" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "content" TEXT NOT NULL,
  "user_id" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create comments table
CREATE TABLE "comments" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "user_id" INT NOT NULL,
  "content" TEXT NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create follows table
CREATE TABLE "follows" (
  "follower_id" INT NOT NULL,
  "following_id" INT NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("follower_id", "following_id")
);

-- Create images table
CREATE TABLE "images" (
  "id" SERIAL PRIMARY KEY,
  "article_id" INT NOT NULL,
  "image_url" TEXT NOT NULL,
  "caption" VARCHAR(255),
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create tokens table
CREATE TABLE "tokens" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT,
  "refresh_token" TEXT,
  "expires_at" TIMESTAMP NOT NULL,
  "status" VARCHAR(20) DEFAULT 'active',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

-- Create tags table
CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(50) UNIQUE NOT NULL
);

-- Create article_tags table
CREATE TABLE "article_tags" (
  "article_id" INT NOT NULL,
  "tag_id" INT NOT NULL,
  PRIMARY KEY ("article_id", "tag_id")
);

-- Add foreign key constraints
ALTER TABLE "user_subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "user_subscriptions" ADD FOREIGN KEY ("plan_id") REFERENCES "membership_plans" ("id") ON DELETE CASCADE;

ALTER TABLE "payment_records" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "payment_records" ADD FOREIGN KEY ("subscription_id") REFERENCES "user_subscriptions" ("id") ON DELETE CASCADE;

ALTER TABLE "articles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "follows" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "follows" ADD FOREIGN KEY ("following_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "images" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;

ALTER TABLE "article_tags" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
ALTER TABLE "article_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE;

ALTER TABLE "tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
