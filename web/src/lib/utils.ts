import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function capitalizeFirst(value: string) {
  if (!value) {
    return value;
  }
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export function formatTransactionDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  const formatter = new Intl.DateTimeFormat("en-US", {
    month: "numeric",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
  return formatter.format(date);
}

export function formatTransactionAmount(
  amount: string,
  currency: "USD" | "EUR",
  direction: "send" | "receive",
) {
  const numericAmount = Number(amount);
  const safeAmount = Number.isNaN(numericAmount) ? 0 : Math.abs(numericAmount);
  const formatted = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(safeAmount);

  if (direction === "send") {
    return `- ${formatted}`;
  }

  return formatted;
}
