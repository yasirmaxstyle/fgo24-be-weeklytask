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
        datetime created_at
        datetime updated_at
        datetime last_login
        boolean is_active
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
        int sender_id FK "nullable"
        int receiver_id FK
        string transaction_type "TRANSFER, TOPUP"
        int payment_method_id FK "nullable"
        decimal amount
        decimal fee
        string description
        string reference_number UK
        string payment_reference "nullable"
        string status
        datetime created_at
        datetime completed_at
        string category
    }

    PAYMENT_METHODS {
        int method_id PK
        string method_name
        string method_type
        boolean is_active
        decimal min_amount
        decimal max_amount
        decimal fee_percentage
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
    USERS ||--o{ TRANSACTIONS : "initiates (sender_id)"
    USERS ||--o{ TRANSACTIONS : "benefits (receiver_id)"
    USERS ||--o{ TRANSACTION_HISTORY : has
    USERS ||--o{ PASSWORD_RESETS : requests

    CONTACTS }o--|| USERS : refers_to
    TRANSACTIONS ||--o{ TRANSACTION_HISTORY : generates
    PAYMENT_METHODS ||--o{ TRANSACTIONS : used_in