"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
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
import { transfer } from "@/lib/api/transactions";

const transferSchema = z.object({
  to_account_id: z.string().check(z.minLength(1, "Account ID is required")),
  amount: z
    .string()
    .check(z.minLength(1, "Amount is required"))
    .check(
      z.refine((val) => {
        const num = Number.parseFloat(val);
        return !Number.isNaN(num) && num > 0;
      }, "Amount must be greater than 0"),
    ),
  currency: z.enum(["USD", "EUR"], "Currency is required"),
});

type TransferFormValues = z.infer<typeof transferSchema>;

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

export default function TransferPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    watch,
  } = useForm<TransferFormValues>({
    resolver: zodResolver(transferSchema),
    defaultValues: {
      to_account_id: "",
      amount: "",
      currency: "USD",
    },
  });

  const selectedCurrency = watch("currency");

  const { data: accounts } = useQuery({
    queryKey: ["accounts"],
    queryFn: getAccounts,
  });

  const selectedAccount = accounts?.find(
    (acc) => acc.currency === selectedCurrency,
  );

  const transferMutation = useMutation({
    mutationFn: transfer,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      router.push(`/transactions?highlight=${data.id}`);
    },
  });

  const onSubmit = (values: TransferFormValues) => {
    transferMutation.mutate({
      to_account_id: values.to_account_id,
      amount: values.amount,
      currency: values.currency,
    });
  };

  return (
    <ProtectedRoute>
      <div className="flex flex-1 flex-col gap-8">
        <section className="space-y-1">
          <h1 className="text-3xl font-semibold tracking-tight">Transfer</h1>
          <p className="text-muted-foreground">
            Send money between accounts or to another user.
          </p>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>New transfer</CardTitle>
            <CardDescription>
              Enter recipient account ID, currency, and amount to send.
            </CardDescription>
          </CardHeader>
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="to_account_id">Recipient Account ID</Label>
                <Input
                  id="to_account_id"
                  placeholder="Enter account ID"
                  aria-invalid={Boolean(errors.to_account_id)}
                  {...register("to_account_id")}
                />
                {errors.to_account_id ? (
                  <p className="text-xs text-destructive">
                    {errors.to_account_id.message}
                  </p>
                ) : (
                  <p className="text-muted-foreground text-xs">
                    Enter the account ID of the recipient.
                  </p>
                )}
              </div>

              <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
                <div className="space-y-2">
                  <Label htmlFor="transfer-amount">Amount</Label>
                  <Input
                    id="transfer-amount"
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
                  ) : null}
                </div>
                <div className="space-y-2">
                  <Label>Currency</Label>
                  <Controller
                    name="currency"
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
                  {selectedAccount ? (
                    <p className="text-muted-foreground text-xs">
                      Available:{" "}
                      {formatBalance(
                        selectedAccount.balance,
                        selectedAccount.currency,
                      )}
                    </p>
                  ) : null}
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between mt-4">
              {transferMutation.isError ? (
                <p className="text-sm text-destructive">
                  {transferMutation.error instanceof Error
                    ? transferMutation.error.message
                    : "Transfer failed"}
                </p>
              ) : null}
              <Button type="submit" disabled={transferMutation.isPending}>
                {transferMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 size-4 animate-spin" />
                    Processing...
                  </>
                ) : (
                  "Submit transfer"
                )}
              </Button>
              <p className="text-muted-foreground text-xs">
                Transfers are processed instantly.
              </p>
            </CardFooter>
          </form>
        </Card>
      </div>
    </ProtectedRoute>
  );
}
