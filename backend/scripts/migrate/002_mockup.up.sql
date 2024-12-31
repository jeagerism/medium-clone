-- mockup_up.sql
-- Insert a user
INSERT INTO "users" ("name", "email", "password_hash") 
VALUES ('John Doe', 'john.doe@example.com', 'hashedpassword123');

-- Insert sample membership plans
INSERT INTO "membership_plans" ("name", "price", "duration_days") 
VALUES ('Basic Plan', 9.99, 30), 
       ('Premium Plan', 19.99, 30);

-- Insert a subscription
INSERT INTO "user_subscriptions" ("user_id", "plan_id", "start_date", "end_date") 
VALUES (1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '30 days');

-- Insert a payment record
INSERT INTO "payment_records" ("user_id", "subscription_id", "payment_gateway", "payment_id", "amount", "currency") 
VALUES (1, 1, 'PayPal', 'PAYID12345', 9.99, 'USD');

-- Insert an article
INSERT INTO "articles" ("user_id", "title", "content") 
VALUES (1, 'Sample Article', 'This is a sample article content.');

-- Insert a comment on the article
INSERT INTO "comments" ("article_id", "user_id", "content") 
VALUES (1, 1, 'Great article!');

-- Insert a follow relationship
INSERT INTO "follows" ("follower_id", "following_id") 
VALUES (1, 1);

-- Insert a tag and link it to the article
INSERT INTO "tags" ("name") 
VALUES ('Tech'), ('Science');

INSERT INTO "article_tags" ("article_id", "tag_id") 
VALUES (1, 1), (1, 2);

INSERT INTO "likes" ("article_id", "user_id")
VALUES (1,1);

-- New 

INSERT INTO "users" ("name", "email", "password_hash") 
VALUES ('John Doe2', 'john2.doe@example.com', 'hashedpassword123');

INSERT INTO "likes" ("article_id", "user_id")
VALUES (1,2);

INSERT INTO "comments" ("article_id", "user_id", "content") 
VALUES (1, 2, 'Wow it's'Great article!');

INSERT INTO "images" ("article_id", "image_url", "caption") 
VALUES (1, 'image1', 'image caption');