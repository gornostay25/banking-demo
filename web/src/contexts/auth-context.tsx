"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import * as React from "react";

import { getCurrentUser, logout as logoutRequest } from "@/lib/api/auth";
import { setUnauthorizedHandler } from "@/lib/api/client";

type AuthContextValue = {
  user: Awaited<ReturnType<typeof getCurrentUser>> | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLoggingOut: boolean;
  logout: () => Promise<void>;
  refetchUser: () => void;
};

const AuthContext = React.createContext<AuthContextValue | undefined>(
  undefined,
);

const AUTH_QUERY_KEY = ["auth", "me"] as const;

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data, isLoading, refetch } = useQuery({
    queryKey: AUTH_QUERY_KEY,
    queryFn: getCurrentUser,
    retry: false,
    refetchOnWindowFocus: false,
  });

  const logoutMutation = useMutation({
    mutationFn: logoutRequest,
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: AUTH_QUERY_KEY });
      router.replace("/login");
    },
  });

  React.useEffect(() => {
    setUnauthorizedHandler(() => {
      queryClient.removeQueries({ queryKey: AUTH_QUERY_KEY });
      if (
        typeof window !== "undefined" &&
        window.location.pathname !== "/login"
      ) {
        router.replace("/login");
      }
    });

    return () => setUnauthorizedHandler(null);
  }, [queryClient, router]);

  const value = React.useMemo<AuthContextValue>(
    () => ({
      user: data ?? null,
      isAuthenticated: Boolean(data),
      isLoading,
      isLoggingOut: logoutMutation.isPending,
      logout: async () => {
        await logoutMutation.mutateAsync();
      },
      refetchUser: () => {
        void refetch();
      },
    }),
    [data, isLoading, logoutMutation, refetch],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export const useAuthContext = () => {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error("useAuthContext must be used within AuthProvider");
  }
  return context;
};
