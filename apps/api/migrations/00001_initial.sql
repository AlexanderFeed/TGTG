-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text NOT NULL,
    name text NOT NULL,
    city text NOT NULL DEFAULT 'Москва',
    role text NOT NULL DEFAULT 'customer',
    verified_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_email_normalized CHECK (email = lower(btrim(email))),
    CONSTRAINT users_email_length CHECK (char_length(email) BETWEEN 3 AND 320),
    CONSTRAINT users_name_length CHECK (char_length(name) BETWEEN 2 AND 80),
    CONSTRAINT users_city_length CHECK (char_length(city) BETWEEN 2 AND 80),
    CONSTRAINT users_role_allowed CHECK (role IN ('customer', 'merchant', 'admin'))
);

CREATE UNIQUE INDEX users_email_unique ON users (email);

CREATE TABLE email_challenges (
    id uuid PRIMARY KEY,
    email text NOT NULL,
    purpose text NOT NULL,
    pending_name text,
    code_hash text NOT NULL,
    attempts integer NOT NULL DEFAULT 0,
    max_attempts integer NOT NULL,
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT email_challenges_purpose_allowed CHECK (purpose IN ('register', 'login')),
    CONSTRAINT email_challenges_attempts_valid CHECK (attempts >= 0 AND max_attempts > 0),
    CONSTRAINT email_challenges_email_normalized CHECK (email = lower(btrim(email)))
);

CREATE INDEX email_challenges_lookup_idx
    ON email_challenges (email, purpose, created_at DESC);

CREATE TABLE user_sessions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX user_sessions_token_unique ON user_sessions (token_hash);
CREATE INDEX user_sessions_user_idx ON user_sessions (user_id, created_at DESC);

CREATE TABLE offers (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    merchant text NOT NULL,
    category text NOT NULL,
    description text NOT NULL DEFAULT '',
    contents text NOT NULL DEFAULT '',
    image_url text NOT NULL DEFAULT '',
    price_kopecks bigint NOT NULL,
    original_price_kopecks bigint NOT NULL,
    pickup_start timestamptz NOT NULL,
    pickup_end timestamptz NOT NULL,
    quantity integer NOT NULL,
    address text NOT NULL,
    district text NOT NULL DEFAULT '',
    latitude double precision,
    longitude double precision,
    delivery boolean NOT NULL DEFAULT false,
    status text NOT NULL DEFAULT 'active',
    created_by uuid REFERENCES users(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT offers_title_length CHECK (char_length(title) BETWEEN 2 AND 120),
    CONSTRAINT offers_merchant_length CHECK (char_length(merchant) BETWEEN 2 AND 120),
    CONSTRAINT offers_category_length CHECK (char_length(category) BETWEEN 2 AND 80),
    CONSTRAINT offers_prices_valid CHECK (price_kopecks > 0 AND original_price_kopecks >= price_kopecks),
    CONSTRAINT offers_pickup_valid CHECK (pickup_end > pickup_start),
    CONSTRAINT offers_quantity_valid CHECK (quantity >= 0),
    CONSTRAINT offers_status_allowed CHECK (status IN ('draft', 'active', 'paused', 'deleted')),
    CONSTRAINT offers_latitude_valid CHECK (latitude IS NULL OR latitude BETWEEN -90 AND 90),
    CONSTRAINT offers_longitude_valid CHECK (longitude IS NULL OR longitude BETWEEN -180 AND 180)
);

CREATE INDEX offers_public_list_idx ON offers (status, created_at DESC);
CREATE INDEX offers_category_idx ON offers (category) WHERE status = 'active';

INSERT INTO offers (
    id, title, merchant, category, description, contents, image_url,
    price_kopecks, original_price_kopecks, pickup_start, pickup_end,
    quantity, address, district, latitude, longitude, delivery, status
) VALUES
    ('11111111-1111-4111-8111-111111111111', 'Сюрприз-пакет из пекарни', 'Хлебная история', 'Выпечка', 'Свежая выпечка, оставшаяся к закрытию.', 'Хлеб, булочки или десерты — состав зависит от дня.', '/images/bakery-rescue.png', 29900, 85000, now() + interval '4 hours', now() + interval '7 hours', 5, 'ул. Покровка, 12', 'Басманный', 55.7598, 37.6452, false, 'active'),
    ('22222222-2222-4222-8222-222222222222', 'Кофе и десерт', 'Тёплый свет', 'Кафе', 'Набор из напитка и десерта дня.', 'Кофе и десерт на выбор заведения.', '/images/cafe-rescue.png', 24900, 62000, now() + interval '3 hours', now() + interval '6 hours', 4, 'Мясницкая ул., 24', 'Красносельский', 55.7646, 37.6388, false, 'active'),
    ('33333333-3333-4333-8333-333333333333', 'Овощи и фрукты', 'Рядом маркет', 'Продукты', 'Набор хороших продуктов с коротким сроком реализации.', 'Сезонные овощи и фрукты.', '/images/grocery-rescue.png', 39900, 99000, now() + interval '5 hours', now() + interval '9 hours', 7, 'Садовая-Каретная ул., 8', 'Тверской', 55.7712, 37.6084, true, 'active');

-- +goose Down
DROP TABLE IF EXISTS offers;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS email_challenges;
DROP TABLE IF EXISTS users;
