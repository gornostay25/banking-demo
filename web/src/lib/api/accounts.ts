import { apiRequest } from "@/lib/api/client";

export type AccountResponse = {
  id: string;
  currency: "USD" | "EUR";
  balance: string;
  created_at: string;
  updated_at: string;
};

export const getAccounts = () => apiRequest<AccountResponse[]>("/accounts");
