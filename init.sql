CREATE DATABASE midasinvestment;

\c midasinvestment

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR (150) PRIMARY KEY
    );

CREATE TABLE IF NOT EXISTS debank_users_assets (
    id SERIAL PRIMARY KEY,
    body JSONB NOT NULL,
    created_at TIMESTAMP,
    user_id VARCHAR (100) NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(user_id)
);


INSERT INTO users (user_id) VALUES
            ('0xba8a8f39b2315d4bc725c026ce3898c2c7e74f57'),
            ('0x2bd4284509bf6626d5def7ef20d4ca38ce71792e'),
            ('0x3ea91c76b176779d10cc2a27fd2687888886f0c2'),
            ('0xe8e94110e568fd45c8eb578bef0f36b5f154b794'),
            ('0x21bce0768110b9a8c50942be257637a843a7eac6'),
            ('0x9429614ccabfb2b24f444f33ede29d4575ebcdd1'),
            ('0x12244c23101f66741dae553c8836a9b2fd4e413a'),
            ('0x8c2753ee27ba890fbb60653d156d92e1c334f528');
