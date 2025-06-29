```mermaid

erDiagram
direction LR
    USERS {
        int user_id PK
        string email UK
        string phone UK
        string full_name
        string password_hash
        string pin_hash
        decimal balance
        string registration_status
        boolean is_verified
        datetime created_at
        datetime updated_at
        datetime last_login
        boolean is_active
        string email_verification_token
        string phone_verification_token
        datetime email_verified_at
        datetime phone_verified_at
    }

    CONTACTS {
        int contact_id PK
        int user_id FK
        int contact_user_id FK
        string contact_name
        string contact_phone
        boolean is_favorite
        datetime created_at
    }

    TRANSACTIONS {
        int transaction_id PK
        int sender_id FK
        int receiver_id FK
        string transaction_type
        decimal amount
        decimal fee
        string description
        string reference_number UK
        string status
        datetime created_at
        datetime completed_at
        string category
    }

    TRANSACTION_HISTORY {
        int history_id PK
        int user_id FK
        int transaction_id FK
        string transaction_summary
        decimal balance_before
        decimal balance_after
        datetime recorded_at
    }

    TOPUP_METHODS {
        int method_id PK
        string method_name
        string method_type
        boolean is_active
        decimal min_amount
        decimal max_amount
        decimal fee_percentage
    }

    TOPUP_TRANSACTIONS {
        int topup_id PK
        int user_id FK
        int method_id FK
        decimal amount
        decimal fee
        string status
        string external_reference
        datetime created_at
        datetime completed_at
    }

    PASSWORD_RESETS {
        int reset_id PK
        int user_id FK
        string reset_token
        datetime expires_at
        datetime created_at
        boolean is_used
    }

    %% Relationships
    USERS ||--o{ CONTACTS : owns
    USERS ||--o{ TRANSACTIONS : sends
    USERS ||--o{ TRANSACTIONS : receives
    USERS ||--o{ TRANSACTION_HISTORY : has
    USERS ||--o{ TOPUP_TRANSACTIONS : makes
    USERS ||--o{ PASSWORD_RESETS : requests

    CONTACTS }o--|| USERS : refers_to
    TRANSACTIONS ||--|| TRANSACTION_HISTORY : generates
    TOPUP_METHODS ||--o{ TOPUP_TRANSACTIONS : used_in