"use client";

import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "next/navigation";
import * as React from "react";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getTransactions, type TransactionType } from "@/lib/api/transactions";
import {
  capitalizeFirst,
  formatTransactionAmount,
  formatTransactionDate,
} from "@/lib/utils";

const DEFAULT_PAGE = 1;
const DEFAULT_LIMIT = 10;

function TransactionsContent() {
  const searchParams = useSearchParams();
  const [highlightId, setHighlightId] = React.useState<string | null>(null);
  const [page, setPage] = React.useState(DEFAULT_PAGE);
  const [typeFilter, setTypeFilter] = React.useState<"all" | TransactionType>(
    "all",
  );

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["transactions", page, DEFAULT_LIMIT, typeFilter],
    queryFn: () =>
      getTransactions({
        page,
        limit: DEFAULT_LIMIT,
        type: typeFilter === "all" ? undefined : typeFilter,
      }),
  });

  const total = data?.total ?? 0;
  const transactions = data?.transactions ?? [];
  const totalPages = Math.max(1, Math.ceil(total / DEFAULT_LIMIT));
  const safePage = Math.min(page, totalPages);

  React.useEffect(() => {
    if (page !== safePage) {
      setPage(safePage);
    }
  }, [page, safePage]);

  React.useEffect(() => {
    const targetId = searchParams.get("highlight");
    if (!targetId) {
      return;
    }
    setHighlightId(targetId);
    setPage(DEFAULT_PAGE);
    setTypeFilter("all");
  }, [searchParams]);

  const pageButtons = React.useMemo(() => {
    if (totalPages <= 5) {
      return Array.from({ length: totalPages }, (_, index) => index + 1);
    }

    let start = Math.max(1, page - 2);
    let end = Math.min(totalPages, page + 2);

    if (end - start < 4) {
      if (start === 1) {
        end = Math.min(totalPages, start + 4);
      } else if (end === totalPages) {
        start = Math.max(1, end - 4);
      }
    }

    const pages = [];
    for (let current = start; current <= end; current += 1) {
      pages.push(current);
    }
    return pages;
  }, [page, totalPages]);

  const showSkeleton = isLoading;
  const showError = isError;
  const errorMessage =
    error instanceof Error ? error.message : "Failed to load transactions.";

  const handleTypeChange = (value: "all" | TransactionType) => {
    setTypeFilter(value);
    setPage(DEFAULT_PAGE);
  };

  const skeletonRows = React.useMemo(
    () => Array.from({ length: DEFAULT_LIMIT }, (_, index) => index),
    [],
  );

  return (
    <ProtectedRoute>
      <div className="flex flex-1 flex-col gap-8">
        <section className="flex flex-wrap items-center justify-between gap-4">
          <div className="space-y-1">
            <h1 className="text-3xl font-semibold tracking-tight">
              Transaction history
            </h1>
            <p className="text-muted-foreground">
              Browse all transfers and exchanges across wallets.
            </p>
          </div>
          <Select value={typeFilter} onValueChange={handleTypeChange}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All types</SelectItem>
              <SelectItem value="transfer">Transfers</SelectItem>
              <SelectItem value="exchange">Exchanges</SelectItem>
            </SelectContent>
          </Select>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>All transactions</CardTitle>
            <CardDescription>
              {total > 0
                ? `Showing ${transactions.length} of ${total} results.`
                : "No transactions available yet."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="w-full overflow-x-auto">
              <Table className="min-w-[720px]">
                <TableHeader>
                  <TableRow>
                    <TableHead>Date</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Direction</TableHead>
                    <TableHead>Currency</TableHead>
                    <TableHead className="text-right">Amount</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {showSkeleton
                    ? skeletonRows.map((row) => (
                        <TableRow key={`skeleton-${row}`}>
                          <TableCell>
                            <Skeleton className="h-4 w-20" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-20" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-16" />
                          </TableCell>
                          <TableCell>
                            <Skeleton className="h-4 w-12" />
                          </TableCell>
                          <TableCell className="text-right">
                            <Skeleton className="ml-auto h-4 w-20" />
                          </TableCell>
                        </TableRow>
                      ))
                    : null}
                  {showError ? (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-sm text-destructive"
                      >
                        {errorMessage}
                      </TableCell>
                    </TableRow>
                  ) : null}
                  {!showSkeleton && !showError && transactions.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-sm text-muted-foreground"
                      >
                        No transactions found.
                      </TableCell>
                    </TableRow>
                  ) : null}
                  {!showSkeleton && !showError
                    ? transactions.map((tx) => (
                        <TableRow
                          key={`${tx.id}-${tx.direction}`}
                          className={
                            highlightId === tx.id
                              ? "bg-accent/80 ring-1 ring-accent/80"
                              : undefined
                          }
                        >
                          <TableCell className="text-muted-foreground whitespace-nowrap">
                            {formatTransactionDate(tx.created_at)}
                          </TableCell>
                          <TableCell className="font-medium">
                            {capitalizeFirst(tx.type)}
                          </TableCell>
                          <TableCell>{capitalizeFirst(tx.direction)}</TableCell>
                          <TableCell>{tx.currency}</TableCell>
                          <TableCell className="text-right font-medium whitespace-nowrap">
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

        <div className="flex flex-wrap items-center justify-between gap-4">
          <p className="text-muted-foreground text-sm">
            Page {safePage} of {totalPages} • {DEFAULT_LIMIT} items per page
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((current) => Math.max(1, current - 1))}
              disabled={safePage <= 1 || showSkeleton}
            >
              Previous
            </Button>
            {pageButtons.map((pageNumber) => (
              <Button
                key={`page-${pageNumber}`}
                size="sm"
                variant={pageNumber === safePage ? "default" : "outline"}
                onClick={() => setPage(pageNumber)}
                disabled={pageNumber === safePage || showSkeleton}
              >
                {pageNumber}
              </Button>
            ))}
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                setPage((current) => Math.min(totalPages, current + 1))
              }
              disabled={safePage >= totalPages || showSkeleton}
            >
              Next
            </Button>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
// Fix useSearchParams() should be wrapped in a suspense boundary at page "/transactions". Read more: https://nextjs.org/docs/messages/missing-suspense-with-csr-bailout
export default function TransactionsPage() {
  return (
    <React.Suspense fallback={null}>
      <TransactionsContent />
    </React.Suspense>
  );
}
