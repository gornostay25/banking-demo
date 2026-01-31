"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import * as React from "react";
import { Controller, useForm } from "react-hook-form";
import * as z from "zod/mini";
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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { getAccounts } from "@/lib/api/accounts";
import { exchange } from "@/lib/api/transactions";

const EXCHANGE_RATE = 0.92; // 1 USD = 0.92 EUR
const MIN_EXCHANGE_AMOUNT = 10.0;

const exchangeSchema = z
  .object({
    amount: z
      .string()
      .check(z.minLength(1, "Amount is required"))
      .check(
        z.refine((val) => {
          const num = Number.parseFloat(val);
          return !Number.isNaN(num) && num >= MIN_EXCHANGE_AMOUNT;
        }, `Amount must be at least ${MIN_EXCHANGE_AMOUNT}`),
      ),
    from_currency: z.enum(["USD", "EUR"], "From currency is required"),
    to_currency: z.enum(["USD", "EUR"], "To currency is required"),
  })
  .check(
    z.superRefine((data, ctx) => {
      if (data.from_currency === data.to_currency) {
        ctx.addIssue({
          code: "custom",
          message: "From and to currencies must be different",
          path: ["to_currency"],
        });
      }
    }),
  );

type ExchangeFormValues = z.infer<typeof exchangeSchema>;

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

const calculateConvertedAmount = (
  amount: string,
  fromCurrency: "USD" | "EUR",
  toCurrency: "USD" | "EUR",
): number => {
  const numAmount = Number.parseFloat(amount);
  if (Number.isNaN(numAmount)) return 0;

  if (fromCurrency === "USD" && toCurrency === "EUR") {
    return numAmount * EXCHANGE_RATE;
  }
  if (fromCurrency === "EUR" && toCurrency === "USD") {
    return numAmount / EXCHANGE_RATE;
  }
  return numAmount;
};

export default function ExchangePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    watch,
    setValue,
  } = useForm<ExchangeFormValues>({
    resolver: zodResolver(exchangeSchema),
    defaultValues: {
      amount: "",
      from_currency: "USD",
      to_currency: "EUR",
    },
  });

  const amount = watch("amount");
  const fromCurrency = watch("from_currency");
  const toCurrency = watch("to_currency");

  // Automatically switch to_currency when from_currency changes
  React.useEffect(() => {
    if (fromCurrency === "USD") {
      setValue("to_currency", "EUR");
    } else if (fromCurrency === "EUR") {
      setValue("to_currency", "USD");
    }
  }, [fromCurrency, setValue]);

  const convertedAmount = React.useMemo(() => {
    if (!amount) return 0;
    return calculateConvertedAmount(amount, fromCurrency, toCurrency);
  }, [amount, fromCurrency, toCurrency]);

  const { data: accounts } = useQuery({
    queryKey: ["accounts"],
    queryFn: getAccounts,
  });

  const fromAccount = accounts?.find((acc) => acc.currency === fromCurrency);

  const exchangeMutation = useMutation({
    mutationFn: exchange,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      router.push(`/transactions?highlight=${data.id}`);
    },
  });

  const onSubmit = (values: ExchangeFormValues) => {
    exchangeMutation.mutate({
      amount: values.amount,
      from_currency: values.from_currency,
      to_currency: values.to_currency,
    });
  };

  return (
    <ProtectedRoute>
      <div className="flex flex-1 flex-col gap-8">
        <section className="space-y-1">
          <h1 className="text-3xl font-semibold tracking-tight">Exchange</h1>
          <p className="text-muted-foreground">
            Convert between USD and EUR at a fixed rate.
          </p>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>Currency exchange</CardTitle>
            <CardDescription>
              Review the rate and confirm the transfer between wallets.
            </CardDescription>
          </CardHeader>
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
                <div className="space-y-2">
                  <Label htmlFor="exchange-amount">Amount</Label>
                  <Input
                    id="exchange-amount"
                    type="number"
                    step="0.01"
                    placeholder="0.00"
                    aria-invalid={Boolean(errors.amount)}
                    {...register("amount")}
                  />
                  {errors.amount ? (
                    <p className="text-xs text-destructive">
                      {errors.amount.message}
                    </p>
                  ) : (
                    <p className="text-muted-foreground text-xs">
                      Minimum exchange amount: {MIN_EXCHANGE_AMOUNT.toFixed(2)}
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label>From</Label>
                  <Controller
                    name="from_currency"
                    control={control}
                    render={({ field }) => (
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="USD">USD</SelectItem>
                          <SelectItem value="EUR">EUR</SelectItem>
                        </SelectContent>
                      </Select>
                    )}
                  />
                  {fromAccount ? (
                    <p className="text-muted-foreground text-xs">
                      Available:{" "}
                      {formatBalance(fromAccount.balance, fromAccount.currency)}
                    </p>
                  ) : null}
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
                <div className="space-y-2">
                  <Label>Converted amount</Label>
                  <Input
                    value={convertedAmount.toFixed(2)}
                    readOnly
                    className="bg-muted"
                  />
                  <p className="text-muted-foreground text-xs">
                    You receive{" "}
                    {formatBalance(convertedAmount.toFixed(2), toCurrency)} in
                    the target wallet.
                  </p>
                </div>
                <div className="space-y-2">
                  <Label>To</Label>
                  <Controller
                    name="to_currency"
                    control={control}
                    render={({ field }) => (
                      <Select
                        value={field.value}
                        onValueChange={(value) => {
                          // Prevent selecting the same currency
                          if (value !== fromCurrency) {
                            field.onChange(value);
                          }
                        }}
                        disabled
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="EUR">EUR</SelectItem>
                          <SelectItem value="USD">USD</SelectItem>
                        </SelectContent>
                      </Select>
                    )}
                  />
                  {errors.to_currency ? (
                    <p className="text-xs text-destructive">
                      {errors.to_currency.message}
                    </p>
                  ) : null}
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
              <p className="text-muted-foreground text-xs">
                Exchange rate: 1 USD = {EXCHANGE_RATE} EUR
              </p>
              {exchangeMutation.isError ? (
                <p className="text-sm text-destructive">
                  {exchangeMutation.error instanceof Error
                    ? exchangeMutation.error.message
                    : "Exchange failed"}
                </p>
              ) : null}
              <Button type="submit" disabled={exchangeMutation.isPending}>
                {exchangeMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 size-4 animate-spin" />
                    Processing...
                  </>
                ) : (
                  "Exchange funds"
                )}
              </Button>
            </CardFooter>
          </form>
        </Card>
      </div>
    </ProtectedRoute>
  );
}
