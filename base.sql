CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    pin_hash VARCHAR(255) NOT NULL,
    balance DECIMAL(15, 2) DEFAULT 0.00 CHECK (balance >= 0),
    registration_status VARCHAR(50) DEFAULT 'pending' CHECK (
        registration_status IN (
            'pending',
            'completed',
            'suspended'
        )
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS contacts (
    contact_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    contact_user_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE,
    contact_name VARCHAR(100) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users (user_id) ON DELETE SET NULL,
    receiver_id INTEGER REFERENCES users (user_id) ON DELETE SET NULL,
    transaction_type VARCHAR(50) NOT NULL CHECK (
        transaction_type IN (
            'transfer',
            'topup'
        )
    ),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    fee DECIMAL(15, 2) DEFAULT 0.00 CHECK (fee >= 0),
    description TEXT,
    reference_number VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'processing',
            'completed',
            'failed',
            'cancelled'
        )
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS transaction_history (
    history_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    transaction_id INTEGER NOT NULL REFERENCES transactions (transaction_id) ON DELETE CASCADE,
    transaction_summary TEXT NOT NULL,
    balance_before DECIMAL(15, 2) NOT NULL,
    balance_after DECIMAL(15, 2) NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payment_methods (
    method_id SERIAL PRIMARY KEY,
    method_name VARCHAR(255) NOT NULL,
    method_type VARCHAR(100) NOT NULL CHECK (
        method_type IN (
            'bank_transfer',
            'e_wallet',
            'retail',
        )
    ),
    is_active BOOLEAN DEFAULT TRUE,
    min_amount DECIMAL(15, 2) DEFAULT 10000.00,
    max_amount DECIMAL(15, 2) DEFAULT 10000000.00,
    fee_percentage DECIMAL(5, 4) DEFAULT 0.0000
);

INSERT INTO
    topup_methods (
        method_name,
        method_type,
        min_amount,
        max_amount,
        fee_percentage
    )
VALUES (
        'Bank Transfer - BCA',
        'bank_transfer',
        10000.00,
        10000000.00,
        0.0000
    ),
    (
        'Bank Transfer - Mandiri',
        'bank_transfer',
        10000.00,
        10000000.00,
        0.0000
    ),
    (
        'Bank Transfer - BNI',
        'bank_transfer',
        10000.00,
        10000000.00,
        0.0000
    ),
    (
        'Gopay',
        'e_wallet',
        10000.00,
        2000000.00,
        0.0250
    ),
    (
        'Ovo',
        'e_wallet',
        10000.00,
        2000000.00,
        0.0150
    ),
    (
        'Indomaret',
        'retail',
        10000.00,
        1000000.00,
        0.0200
    ),
    (
        'Alfamart',
        'retail',
        10000.00,
        1000000.00,
        0.0200
    );