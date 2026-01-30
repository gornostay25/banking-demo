-- +goose Up

-- Insert 3 test users with password hash for "password"
INSERT INTO users (id, email, password_hash, created_at) VALUES
  (gen_random_uuid(), 'user1@test.com', '$2a$10$pTNVC3e60pVrAa0YBDOYSOrqJhDKBCezcD5Z.ceu0fAb.zrn8mI8W', NOW()),
  (gen_random_uuid(), 'user2@test.com', '$2a$10$pTNVC3e60pVrAa0YBDOYSOrqJhDKBCezcD5Z.ceu0fAb.zrn8mI8W', NOW()),
  (gen_random_uuid(), 'user3@test.com', '$2a$10$pTNVC3e60pVrAa0YBDOYSOrqJhDKBCezcD5Z.ceu0fAb.zrn8mI8W', NOW());

-- Create USD and EUR accounts for each user with initial balances
-- User 1 accounts
INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'USD',
  1000.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user1@test.com';

INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'EUR',
  500.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user1@test.com';

-- User 2 accounts
INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'USD',
  1000.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user2@test.com';

INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'EUR',
  500.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user2@test.com';

-- User 3 accounts
INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'USD',
  1000.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user3@test.com';

INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at)
SELECT 
  gen_random_uuid(),
  u.id,
  'EUR',
  500.00,
  NOW(),
  NOW()
FROM users u WHERE u.email = 'user3@test.com';

-- +goose Down

-- Delete test accounts
DELETE FROM accounts WHERE user_id IN (
  SELECT id FROM users WHERE email IN ('user1@test.com', 'user2@test.com', 'user3@test.com')
);

-- Delete test users
DELETE FROM users WHERE email IN ('user1@test.com', 'user2@test.com', 'user3@test.com');
