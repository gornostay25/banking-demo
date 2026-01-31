"use client";

import {
  isServer,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import type * as React from "react";
import { AuthProvider } from "@/contexts/auth-context";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined;

function getQueryClient() {
  if (isServer) {
    return makeQueryClient();
  }

  if (!browserQueryClient) {
    browserQueryClient = makeQueryClient();
  }

  return browserQueryClient;
}

export default function Providers({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>{children}</AuthProvider>
    </QueryClientProvider>
  );
}
