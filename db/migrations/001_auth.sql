








SET search_path TO public;




GRANT USAGE, CREATE ON SCHEMA public TO bro;

-- Таблица одноразовых кодов
CREATE TABLE IF NOT EXISTS auth_codes (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    used BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_codes_email_code
ON auth_codes(email, code)
WHERE used = false;


-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,               -- уникальный идентификатор пользователя
    email VARCHAR(255) UNIQUE NOT NULL,-- email, обязательно уникальный и не может быть NULL
    name VARCHAR(255) NOT NULL,        -- имя пользователя, теперь без DEFAULT
    nick VARCHAR(255) UNIQUE NOT NULL,  -- ник пользователя, сразу создаём столбец
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- статус пользователя
    device_id VARCHAR(255),            -- идентификатор устройства
    created_at TIMESTAMPTZ DEFAULT NOW() -- дата и время создания
);


-- Таблица сессий
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at
ON sessions(expires_at);