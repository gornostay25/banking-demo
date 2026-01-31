"use client";

import { useQuery } from "@tanstack/react-query";
import { ArrowLeftRight, List, Send } from "lucide-react";
import Link from "next/link";
import * as React from "react";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getAccounts } from "@/lib/api/accounts";
import { getTransactions } from "@/lib/api/transactions";
import {
  capitalizeFirst,
  formatTransactionAmount,
  formatTransactionDate,
} from "@/lib/utils";

const BALANCE_SKELETON_COUNT = 2;
const RECENT_TRANSACTIONS_LIMIT = 5;
const REFRESH_INTERVAL_MS = 30000;

const formatBalance = (amount: string, currency: "USD" | "EUR") => {
  const numericAmount = Number(amount);
  const safeAmount = Number.isNaN(numericAmount) ? 0 : numericAmount;
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(safeAmount);
};

export default function Home() {
  const {
    data: accounts,
    isLoading: isAccountsLoading,
    isError: isAccountsError,
    error: accountsError,
  } = useQuery({
    queryKey: ["accounts"],
    queryFn: getAccounts,
    refetchInterval: REFRESH_INTERVAL_MS,
  });

  const {
    data: transactionsResponse,
    isLoading: isTransactionsLoading,
    isError: isTransactionsError,
    error: transactionsError,
  } = useQuery({
    queryKey: ["transactions", "recent", RECENT_TRANSACTIONS_LIMIT],
    queryFn: () => getTransactions({ limit: RECENT_TRANSACTIONS_LIMIT }),
    refetchInterval: REFRESH_INTERVAL_MS,
  });

  const balances = accounts ?? [];
  const recentTransactions = transactionsResponse?.transactions ?? [];

  const showBalancesSkeleton = isAccountsLoading;
  const showBalancesError = isAccountsError;
  const balancesErrorMessage =
    accountsError instanceof Error
      ? accountsError.message
      : "Failed to load balances.";

  const showTransactionsSkeleton = isTransactionsLoading;
  const showTransactionsError = isTransactionsError;
  const transactionsErrorMessage =
    transactionsError instanceof Error
      ? transactionsError.message
      : "Failed to load transactions.";

  const balanceSkeletons = React.useMemo(
    () => Array.from({ length: BALANCE_SKELETON_COUNT }, (_, index) => index),
    [],
  );
  const transactionSkeletons = React.useMemo(
    () =>
      Array.from({ length: RECENT_TRANSACTIONS_LIMIT }, (_, index) => index),
    [],
  );

  return (
    <ProtectedRoute>
      <div className="flex flex-1 flex-col gap-8">
        <section className="flex flex-wrap items-center justify-between gap-4">
          <div className="space-y-1">
            <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground">
              Track balances, transfers, and currency exchanges in one place.
            </p>
          </div>
        </section>

        <section className="grid gap-4 sm:grid-cols-2">
          {showBalancesSkeleton
            ? balanceSkeletons.map((item) => (
                <Card key={`balance-skeleton-${item}`}>
                  <CardHeader className="space-y-2">
                    <Skeleton className="h-4 w-36" />
                    <Skeleton className="h-7 w-40" />
                  </CardHeader>
                </Card>
              ))
            : null}
          {showBalancesError ? (
            <Card className="sm:col-span-2">
              <CardHeader>
                <CardTitle>Balances</CardTitle>
                <CardDescription className="text-destructive">
                  {balancesErrorMessage}
                </CardDescription>
              </CardHeader>
            </Card>
          ) : null}
          {!showBalancesSkeleton &&
          !showBalancesError &&
          balances.length === 0 ? (
            <Card className="sm:col-span-2">
              <CardHeader>
                <CardTitle>Balances</CardTitle>
                <CardDescription>No accounts available yet.</CardDescription>
              </CardHeader>
            </Card>
          ) : null}
          {!showBalancesSkeleton && !showBalancesError
            ? balances.map((wallet) => (
                <Card key={wallet.id}>
                  <CardHeader>
                    <CardDescription>
                      ID:{" "}
                      <span className="font-bold select-all">{wallet.id}</span>
                    </CardDescription>
                    <CardTitle className="text-2xl">
                      {wallet.currency}{" "}
                      {formatBalance(wallet.balance, wallet.currency)}
                    </CardTitle>
                  </CardHeader>
                </Card>
              ))
            : null}
        </section>

        <section>
          <Card>
            <CardHeader>
              <CardTitle>Quick actions</CardTitle>
              <CardDescription>
                Jump straight to transfers or currency exchanges.
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-3">
              <div className="flex flex-wrap gap-3">
                <Button asChild size="sm" className="flex-1">
                  <Link href="/transfer">
                    <Send className="size-4" />
                    Go to transfer
                  </Link>
                </Button>
                <Button asChild variant="outline" size="sm" className="flex-1">
                  <Link href="/exchange">
                    <ArrowLeftRight className="size-4" />
                    Go to exchange
                  </Link>
                </Button>
              </div>
            </CardContent>
            <CardFooter>
              <Button asChild variant="ghost" size="sm" className="w-full">
                <Link href="/transactions">
                  <List className="size-4" />
                  View all transactions
                </Link>
              </Button>
            </CardFooter>
          </Card>
        </section>

        <section>
          <Card>
            <CardHeader>
              <CardTitle>Last 5 transactions</CardTitle>
              <CardDescription>Latest activity across wallets.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="w-full overflow-x-auto">
                <Table className="min-w-[520px]">
                  <TableHeader>
                    <TableRow>
                      <TableHead>Date</TableHead>
                      <TableHead>Type</TableHead>
                      <TableHead>Direction</TableHead>
                      <TableHead className="text-right">Amount</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {showTransactionsSkeleton
                      ? transactionSkeletons.map((row) => (
                          <TableRow key={`transactions-skeleton-${row}`}>
                            <TableCell>
                              <Skeleton className="h-4 w-20" />
                            </TableCell>
                            <TableCell>
                              <Skeleton className="h-4 w-32" />
                            </TableCell>
                            <TableCell>
                              <Skeleton className="h-4 w-16" />
                            </TableCell>
                            <TableCell className="text-right">
                              <Skeleton className="ml-auto h-4 w-20" />
                            </TableCell>
                          </TableRow>
                        ))
                      : null}
                    {showTransactionsError ? (
                      <TableRow>
                        <TableCell
                          colSpan={4}
                          className="text-sm text-destructive"
                        >
                          {transactionsErrorMessage}
                        </TableCell>
                      </TableRow>
                    ) : null}
                    {!showTransactionsSkeleton &&
                    !showTransactionsError &&
                    recentTransactions.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={4}
                          className="text-sm text-muted-foreground"
                        >
                          No transactions yet.
                        </TableCell>
                      </TableRow>
                    ) : null}
                    {!showTransactionsSkeleton && !showTransactionsError
                      ? recentTransactions.map((tx) => (
                          <TableRow key={`${tx.id}-${tx.direction}`}>
                            <TableCell className="text-muted-foreground">
                              {formatTransactionDate(tx.created_at)}
                            </TableCell>
                            <TableCell>{capitalizeFirst(tx.type)}</TableCell>
                            <TableCell className="font-medium">
                              {capitalizeFirst(tx.direction)}
                            </TableCell>
                            <TableCell className="text-right font-medium">
                              {formatTransactionAmount(
                                tx.amount,
                                tx.currency,
                                tx.direction,
                              )}
                            </TableCell>
                          </TableRow>
                        ))
                      : null}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </section>
      </div>
    </ProtectedRoute>
  );
}
