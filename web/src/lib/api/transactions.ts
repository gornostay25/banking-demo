import { apiRequest } from "@/lib/api/client";

export type TransactionType = "transfer" | "exchange";
export type TransactionDirection = "send" | "receive";
export type TransactionCurrency = "USD" | "EUR";

export type TransactionResponse = {
  id: string;
  type: TransactionType;
  amount: string;
  currency: TransactionCurrency;
  direction: TransactionDirection;
  created_at: string;
};

export type TransactionListResponse = {
  transactions: TransactionResponse[];
  total: number;
  page: number;
  limit: number;
};

export type TransactionListParams = {
  type?: TransactionType;
  page?: number;
  limit?: number;
};

export const getTransactions = ({
  type,
  page,
  limit,
}: TransactionListParams = {}) => {
  const params = new URLSearchParams();

  if (type) {
    params.set("type", type);
  }
  if (page) {
    params.set("page", String(page));
  }
  if (limit) {
    params.set("limit", String(limit));
  }

  const query = params.toString();
  const path = query ? `/transactions?${query}` : "/transactions";

  return apiRequest<TransactionListResponse>(path);
};

export type TransferRequest = {
  amount: string;
  currency: TransactionCurrency;
  to_account_id: string;
};

export type ExchangeRequest = {
  amount: string;
  from_currency: TransactionCurrency;
  to_currency: TransactionCurrency;
};

export const transfer = (payload: TransferRequest) =>
  apiRequest<TransactionResponse>("/transactions/transfer", {
    method: "POST",
    body: JSON.stringify(payload),
  });

export const exchange = (payload: ExchangeRequest) =>
  apiRequest<TransactionResponse>("/transactions/exchange", {
    method: "POST",
    body: JSON.stringify(payload),
  });
