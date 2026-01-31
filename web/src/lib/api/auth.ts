import { apiRequest } from "@/lib/api/client";

export type LoginRequest = {
  email: string;
  password: string;
};

export type TokenPairResponse = {
  code: number;
  expire: string;
  message: string;
  token: string;
};

export type UserResponse = {
  id: string;
  email: string;
  created_at: string;
};

export const login = (payload: LoginRequest) =>
  apiRequest<TokenPairResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload),
  });

export const logout = () =>
  apiRequest<Record<string, unknown>>("/auth/logout", {
    method: "POST",
  });

export const refreshToken = () =>
  apiRequest<TokenPairResponse>(
    "/auth/refresh",
    {
      method: "POST",
      body: JSON.stringify({}),
    },
    false,
  );

export const getCurrentUser = () => apiRequest<UserResponse>("/auth/me");
