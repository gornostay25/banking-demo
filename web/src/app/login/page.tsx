"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import * as React from "react";
import { useForm } from "react-hook-form";
import * as z from "zod/mini";

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
import { useAuth } from "@/hooks/use-auth";
import { login } from "@/lib/api/auth";

const loginSchema = z.object({
  email: z.email("Invalid email address"),
  password: z.string().check(z.minLength(1, "Password is required")),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const { isAuthenticated, refetchUser } = useAuth();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const loginMutation = useMutation({
    mutationFn: login,
    onSuccess: async () => {
      await refetchUser();
      router.replace("/");
    },
  });

  React.useEffect(() => {
    if (isAuthenticated) {
      router.replace("/");
    }
  }, [isAuthenticated, router]);

  const onSubmit = (values: LoginFormValues) => {
    loginMutation.mutate(values);
  };

  return (
    <div className="flex min-h-[70vh] w-full items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Welcome back</CardTitle>
          <CardDescription>
            Sign in to manage your wallets and transactions.
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit(onSubmit)}>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@mail.com"
                aria-invalid={Boolean(errors.email)}
                {...register("email")}
              />
              {errors.email ? (
                <p className="text-xs text-destructive">
                  {errors.email.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                aria-invalid={Boolean(errors.password)}
                {...register("password")}
              />
              {errors.password ? (
                <p className="text-xs text-destructive">
                  {errors.password.message}
                </p>
              ) : null}
            </div>
            {loginMutation.isError ? (
              <p className="text-sm text-destructive">
                {loginMutation.error instanceof Error
                  ? loginMutation.error.message
                  : "Login failed"}
              </p>
            ) : null}
          </CardContent>
          <CardFooter className="flex flex-col gap-3 mt-4">
            <Button
              className="w-full"
              type="submit"
              disabled={loginMutation.isPending}
            >
              {loginMutation.isPending ? (
                <>
                  <Loader2 className="size-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                "Sign in"
              )}
            </Button>
            <p className="bg-muted/60 text-muted-foreground text-xs leading-relaxed rounded-md border px-3 py-2 w-full">
              <span className="font-medium text-foreground">Emails:</span>{" "}
              <span className="select-all">user1@test.com</span>
              &nbsp;&bull;&nbsp;
              <span className="select-all">user2@test.com</span>
              &nbsp;&bull;&nbsp;
              <span className="select-all">user3@test.com</span>&nbsp;
              <br />
              <span className="font-medium text-foreground">Password:</span>{" "}
              <span className="select-all">password</span>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
