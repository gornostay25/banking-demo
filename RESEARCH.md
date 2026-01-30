# Research

App: Mini Banking Platform

- handle financial transactions
- proper accounting principles
- user-friendly interface


Stack:
- Backend: Node.js with NestJS or Go
- DB: PostgreSQL
- API: RESTful API
- Authentication: JWT (maybe betterauth)
- Frontend: React/Next.js
- UI Framework: Tailwind CSS (shadcn/ui)
- State Management: Zustand


Database Schema:
**double-entry ledger system** - double-entry accounting is a method of recording financial transactions in a way that ensures the balance of the accounting equation is always maintained.

- Ledger entries serve as the authoritative audit trail for all transactions
- Each transaction creates balanced entries (sum of amounts must equal zero)
- Account balances should be maintained for performance
- !! Ensure consistency between ledger entries and account balances

Tables
- users: User information
- accounts: User currency accounts with balances (each user has USD and EUR accounts)
- ledger: Double-entry bookkeeping records (using positive/negative amounts)
- transactions: High-level transaction records for user-facing history

User Creation: (Option B): Create at least 3 pre-seeded test users on initialization

Account Structure: Each user automatically has:
1 USD account (initial balance: $1000.00)
1 EUR account (initial balance: €500.00)


## Api Endpoints

### AUTH
- POST /auth/login - User login
- GET /auth/me - Get current user info

### Account Operations

- GET /accounts - List user's accounts with balances
- GET /accounts/:id/balance - Get specific account balance

### Transaction Operations
- POST /transactions/transfer - Transfer between users (same currency)
- POST /transactions/exchange - Currency exchange within user's accounts

### Transaction History
- GET /transactions - List transactions with filters:
    - type : transfer, exchange
    - page & limit : Pagination

## Frontend Features
- Dashboard 
    - Display current balance for each wallet (USD and EUR)
    - last 5 transactions
    - Forms for: Transfer/Exchange

- Transactions Pages
    - Transfer
        - Select recipient
        - Select currency
        - Enter amount
        - Validation 
    - Exchange
        - Select source currency
        - Enter amount
        - Show converted amount **(1 USD = 0.92 EUR)**
        - Display exchange rate

- Transaction History Page
    Table/list of all transactions
    Filter by transaction type
    Pagination controls

## Business Requirements
1. Insufficient Funds
    Transfers and exchanges must not exceed available balance **(check if the account has enough funds)**
    System must handle concurrent transaction attempts **(use a queue to handle the transactions)**
2. Transaction Integrity
    All financial operations must maintain data consistency **(use a transaction to ensure the integrity of the data)**
    Ledger and balances must remain synchronized 
    Partial transaction completion is not acceptable 
3. Currency Precision
    All monetary amounts must maintain 2 decimal place precision 
    No rounding errors should affect user balances 
4. Exchange Operations
    Fixed exchange rate: 1 USD = 0.92 EUR
    Exchange calculations must be transparent to users