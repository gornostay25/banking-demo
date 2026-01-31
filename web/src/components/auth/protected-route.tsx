"use client";

import { useRouter } from "next/navigation";
import * as React from "react";

import { useAuth } from "@/hooks/use-auth";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  React.useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading || !isAuthenticated) {
    return (
      <div className="flex min-h-[50vh] w-full items-center justify-center">
        <div className="size-6 animate-spin rounded-full border-2 border-muted-foreground/20 border-t-muted-foreground" />
      </div>
    );
  }

  return <>{children}</>;
}
