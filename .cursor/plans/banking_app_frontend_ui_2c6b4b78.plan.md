---
name: Banking App Frontend UI
overview: Create a modern, minimalistic banking app frontend with Dashboard, Transfer, Exchange, and Transaction History pages using Tailwind CSS, shadcn/ui components, and React Query setup (UI only, no API calls yet).
todos:
  - id: setup-shadcn
    content: Create components.json and install required shadcn components (button, card, input, label, select, table, separator, skeleton, form)
    status: completed
  - id: setup-utils
    content: Create src/lib/utils.ts with cn() utility function
    status: completed
  - id: setup-react-query
    content: Add React Query provider to layout.tsx (setup only, no API calls)
    status: completed
  - id: create-navigation
    content: Create navigation component with links to all pages
    status: completed
  - id: create-dashboard
    content: Build dashboard page with wallet balances, last 5 transactions, and quick action forms
    status: completed
  - id: create-transfer
    content: Create transfer page with recipient, currency, and amount form fields
    status: completed
  - id: create-exchange
    content: Create exchange page with source/target currency, amount, and exchange rate display
    status: completed
  - id: create-transactions
    content: Create transaction history page with table, filters, and pagination UI
    status: completed
isProject: false
---

# Banking App Frontend UI Implementation Plan

## Overview

Build a responsive banking app frontend with Dashboard, Transfer, Exchange, and Transaction History pages. Use Tailwind CSS v4, shadcn/ui components, and set up React Query provider (UI only - no API integration yet).

## Project Structure

```
web/
├── components.json          # shadcn configuration
├── src/
│   ├── app/
│   │   ├── layout.tsx       # Root layout with React Query provider
│   │   ├── page.tsx         # Dashboard page
│   │   ├── transfer/
│   │   │   └── page.tsx     # Transfer form page
│   │   ├── exchange/
│   │   │   └── page.tsx     # Exchange form page
│   │   └── transactions/
│   │       └── page.tsx     # Transaction history page
│   ├── components/
│   │   ├── ui/              # shadcn components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── input.tsx
│   │   │   ├── label.tsx
│   │   │   ├── select.tsx
│   │   │   ├── table.tsx
│   │   │   ├── separator.tsx
│   │   │   ├── skeleton.tsx
│   │   │   └── form.tsx
│   │   └── layout/
│   │       └── navigation.tsx  # Navigation component
│   └── lib/
│       └── utils.ts         # cn() utility for shadcn
```

## Implementation Steps

### 1. Setup shadcn Configuration

- Create `web/components.json` with proper configuration for Next.js 16, Tailwind v4, and TypeScript
- Configure component path: `src/components/ui`
- Set up style configuration

### 2. Install shadcn Components

Install the following components using shadcn MCP:

- `button` - For form submissions and actions
- `card` - For wallet balance display and form containers
- `input` - For amount and recipient inputs
- `label` - For form labels
- `select` - For currency and recipient dropdowns
- `table` - For transaction history
- `separator` - For visual separation
- `skeleton` - For loading states
- `form` - For form handling structure

### 3. Create Utility Files

- `src/lib/utils.ts` - Add `cn()` utility function for className merging (required by shadcn)

### 4. Setup React Query Provider

- Update `src/app/layout.tsx` to wrap children with `QueryClientProvider`
- Create QueryClient instance (no API calls yet, just setup)

### 5. Create Navigation Component

- `src/components/layout/navigation.tsx` - Header navigation with links to:
  - Dashboard (/)
  - Transfer (/transfer)
  - Exchange (/exchange)
  - Transactions (/transactions)

### 6. Dashboard Page (`src/app/page.tsx`)

Features:

- Display wallet balances (USD and EUR) in cards
- Show last 5 transactions in a list/table
- Transfer form section (inline or card)
- Exchange form section (inline or card)
- Responsive grid layout

UI Structure:

```
┌─────────────────────────────────────┐
│ Navigation                          │
├─────────────────────────────────────┤
│ Wallet Balances                     │
│ ┌──────────┐  ┌──────────┐        │
│ │ USD Card │  │ EUR Card │        │
│ └──────────┘  └──────────┘        │
├─────────────────────────────────────┤
│ Last 5 Transactions                 │
│ ┌──────────────────────────────┐    │
│ │ Transaction List            │    │
│ └──────────────────────────────┘    │
├─────────────────────────────────────┤
│ Quick Actions                       │
│ ┌──────────┐  ┌──────────┐        │
│ │ Transfer │  │ Exchange │        │
│ │  Form    │  │  Form    │        │
│ └──────────┘  └──────────┘        │
└─────────────────────────────────────┘
```

### 7. Transfer Page (`src/app/transfer/page.tsx`)

Features:

- Form with fields:
  - Recipient select/dropdown (placeholder for now)
  - Currency select (USD/EUR)
  - Amount input
- Validation display area (for errors)
- Submit button
- Responsive card layout

### 8. Exchange Page (`src/app/exchange/page.tsx`)

Features:

- Form with fields:
  - Source currency select (USD/EUR)
  - Amount input
  - Display converted amount (calculated: 1 USD = 0.92 EUR)
  - Target currency select (USD/EUR)
- Exchange rate display (1 USD = 0.92 EUR)
- Submit button
- Responsive card layout

### 9. Transaction History Page (`src/app/transactions/page.tsx`)

Features:

- Table displaying transactions with columns:
  - Date
  - Type (transfer/exchange)
  - Amount
  - Currency
  - Direction (send/receive)
- Filter dropdown by transaction type (transfer/exchange)
- Pagination controls (UI only)
- Empty state when no transactions

### 10. Styling & Design

- Use Tailwind CSS v4 with minimalistic design
- Modern color scheme (use shadcn default theme)
- Responsive breakpoints:
  - Mobile: single column
  - Tablet: 2 columns for balances
  - Desktop: full layout
- Consistent spacing and typography
- Subtle shadows and borders for cards
- Clean form styling with proper focus states

## Key Files to Create/Modify

1. `web/components.json` - shadcn configuration
2. `web/src/lib/utils.ts` - Utility functions
3. `web/src/app/layout.tsx` - Add React Query provider
4. `web/src/components/layout/navigation.tsx` - Navigation component
5. `web/src/app/page.tsx` - Dashboard page
6. `web/src/app/transfer/page.tsx` - Transfer form page
7. `web/src/app/exchange/page.tsx` - Exchange form page
8. `web/src/app/transactions/page.tsx` - Transaction history page
9. All shadcn UI components in `web/src/components/ui/`

## Design Principles

- Minimalistic: Clean, uncluttered interface
- Modern: Contemporary design patterns
- Responsive: Works on mobile, tablet, desktop
- Accessible: Proper labels, semantic HTML
- Consistent: Unified spacing, typography, colors

## Notes

- All forms use placeholder/mock data (no API calls)
- Exchange rate is hardcoded: 1 USD = 0.92 EUR
- Transaction data will be mocked/empty for now
- Focus on UI/UX, not functionality

