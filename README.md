<!--
[x] Setup instructions
[x] Specify which user management approach was chosen (registration or pre-seeded)
[ ] Design decisions and trade-offs
[ ] Known limitations
[ ] Any incomplete features due to time constraints
 -->

# Mini Banking Platform
_* Uses Go Blueprint for the project structure_

User Management Approach: Pre-seeded users
- 3 pre-seeded test users on initialization
- Each user has:
    - 1 USD account (initial balance: $1000.00)
    - 1 EUR account (initial balance: €500.00)

## Setup

Run the application
```bash
make run
```

More: See [banking/README.md](banking/README.md) for more details.

## Design decisions and trade-offs

- Using `Serial` type for id's instead of `UUID` for simplicity

## Questions to Consider

<details>
<summary>How do you ensure transaction atomicity?</summary>
<p>
    Using transactions
</p>
</details>

<details>
<summary>How do you prevent double-spending?</summary>
<p>
    Using transactions
</p>
</details>

<details>
<summary>How do you maintain consistency between ledger entries and account balances?</summary>
<p>
    Using transactions
</p>
</details>

<details>
<summary>How would you handle decimal precision for different currencies?</summary>
<p></p>
</details>

<details>
<summary>What indexing strategy would you use for the ledger table?</summary>
<p></p>
</details>

<details>
<summary>How would you verify that balances are correctly synchronized?</summary>
<p></p>
</details>

<details>
<summary>How would you scale this system for millions of users?</summary>
<p></p>
</details>
