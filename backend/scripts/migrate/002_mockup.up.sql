-- Insert mock data into users table
INSERT INTO "users" ("name", "email", "password_hash", "bio", "profile_image", "role") VALUES
('Alice', 'alice@example.com', 'hashed_password_1', 'Love writing tech articles.', 'https://example.com/profiles/alice.jpg', 'admin'),
('Bob', 'bob@example.com', 'hashed_password_2', 'Enthusiastic about science.', 'https://example.com/profiles/bob.jpg', 'user'),
('Charlie', 'charlie@example.com', 'hashed_password_3', 'Passionate about photography.', 'https://example.com/profiles/charlie.jpg', 'user');

-- Insert mock data into membership_plans table
INSERT INTO "membership_plans" ("name", "price", "duration_days") VALUES
('Basic Plan', 9.99, 30),
('Premium Plan', 19.99, 90),
('Pro Plan', 49.99, 365);

-- Insert mock data into user_subscriptions table
INSERT INTO "user_subscriptions" ("user_id", "plan_id", "subscription_id", "status", "start_date", "end_date") VALUES
(1, 1, 'sub_001', 'active', NOW() - INTERVAL '15 days', NOW() + INTERVAL '15 days'),
(2, 2, 'sub_002', 'active', NOW() - INTERVAL '45 days', NOW() + INTERVAL '45 days'),
(3, 3, 'sub_003', 'active', NOW() - INTERVAL '150 days', NOW() + INTERVAL '215 days');

-- Insert mock data into payment_records table
INSERT INTO "payment_records" ("user_id", "subscription_id", "payment_gateway", "payment_id", "amount", "currency") VALUES
(1, 1, 'Stripe', 'pay_001', 9.99, 'USD'),
(2, 2, 'PayPal', 'pay_002', 19.99, 'USD'),
(3, 3, 'Stripe', 'pay_003', 49.99, 'USD');

-- Insert mock data into articles table
INSERT INTO "articles" ("title", "content", "cover_image", "user_id") VALUES
('The Future of AI', 'Content about AI...', 'https://example.com/images/ai.jpg', 1),
('Exploring the Cosmos', 'Content about space...', 'https://example.com/images/space.jpg', 2),
('Photography Tips', 'Content about photography...', 'https://example.com/images/photography.jpg', 3);

-- Insert mock data into comments table
INSERT INTO "comments" ("article_id", "user_id", "content") VALUES
(1, 2, 'Great article on AI!'),
(2, 1, 'Space is so fascinating!'),
(3, 2, 'Thanks for the tips on photography.');

-- Insert mock data into follows table
INSERT INTO "follows" ("follower_id", "following_id") VALUES
(1, 2),
(2, 3),
(3, 1);

-- Insert mock data into images table
INSERT INTO "images" ("article_id", "image_url", "caption") VALUES
(1, 'https://example.com/images/ai_diagram.jpg', 'AI diagram'),
(2, 'https://example.com/images/galaxy.jpg', 'Galaxy image'),
(3, 'https://example.com/images/camera.jpg', 'Camera close-up');

-- Insert mock data into tokens table
INSERT INTO "tokens" ("user_id", "refresh_token", "expires_at", "status") VALUES
(1, 'refresh_token_1', NOW() + INTERVAL '30 days', 'active'),
(2, 'refresh_token_2', NOW() + INTERVAL '30 days', 'active'),
(3, 'refresh_token_3', NOW() + INTERVAL '30 days', 'active');

-- Insert mock data into tags table
INSERT INTO "tags" ("name") VALUES
('Science'),
('Technology'),
('Photography');

-- Insert mock data into article_tags table
INSERT INTO "article_tags" ("article_id", "tag_id") VALUES
(1, 2), -- Technology tag for AI article
(2, 1), -- Science tag for Cosmos article
(3, 3); -- Photography tag for Tips article

-- Insert mock data into likes table
INSERT INTO "likes" ("article_id", "user_id") VALUES
(1, 2),
(1, 3),
(2, 1),
(3, 2),
(3, 1);
