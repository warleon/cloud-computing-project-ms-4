# MS4 – Compliance & Risk Service (golang, MySQL)

This is an example microservice for compliance & risk checks using Gin and GORM.

## Run

1. Copy `.env.example` to `.env` and edit values.
2. `docker-compose up --build`
3. The service will be available at `http://localhost:8080/api/v1`.

## Endpoints

- `POST /api/v1/validateTransaction` - validate a transaction
- `POST /api/v1/rules` - create a compliance rule
- `GET /api/v1/rules` - list rules

## ER Diagram

```mermaid
erDiagram
    Rule {
        uint ID PK
        string Name
        string Description
        string Type
        string Account
        float Threshold
        time CreatedAt
        time UpdatedAt
        time DeletedAt
    }

    Sanction {
        uint ID PK
        string AccID
        time CreatedAt
        time UpdatedAt
        time DeletedAt
    }

    Decision {
        uint ID PK
        bool Approved
        string Reason
        time CreatedAt
        time UpdatedAt
        time DeletedAt
    }

        AuditLog {
        uint ID PK
        string TransactionID
        string CustomerID
        uint DecisionID FK
        time CreatedAt
        time UpdatedAt
        time DeletedAt
    }

    %% Relationships (assumed)
    %% You didn’t define explicit foreign keys, so these are logical guesses.
    Rule ||--o{ Decision : "generates"
    Rule ||--o{ Sanction : "uses"
    Decision ||--o{ AuditLog : "referenced by"
    %% Notes
    %% Threshold is optional in RuleExtras
    %% db *gorm.DB is not persisted, so it's ignored

```
