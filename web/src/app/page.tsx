import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import Link from "next/link"
import { ArrowLeftRight, Send, List } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const balances = [
  {
    currency: "USD",
    balance: "12,480.50",
    account: "1234567890",
    change: "+2.4% this week",
  },
  {
    currency: "EUR",
    balance: "8,120.00",
    account: "1234567890",
    change: "+1.1% this week",
  },
]

const recentTransactions = [
  {
    id: "tx-001",
    date: "Jan 30, 2026",
    description: "Transfer to Alex Morgan",
    type: "Transfer",
    amount: "-$240.00",
  },
  {
    id: "tx-002",
    date: "Jan 30, 2026",
    description: "Exchange USD to EUR",
    type: "Exchange",
    amount: "€460.00",
  },
  {
    id: "tx-003",
    date: "Jan 29, 2026",
    description: "Payment from Jamie Rivera",
    type: "Transfer",
    amount: "+$1,200.00",
  },
  {
    id: "tx-004",
    date: "Jan 28, 2026",
    description: "Exchange EUR to USD",
    type: "Exchange",
    amount: "$820.00",
  },
  {
    id: "tx-005",
    date: "Jan 27, 2026",
    description: "Transfer to Taylor Price",
    type: "Transfer",
    amount: "-$85.50",
  },
]

export default function Home() {
  return (
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
        {balances.map((wallet) => (
          <Card key={wallet.currency}>
            <CardHeader>
              <CardDescription>ID: <span className="font-bold select-all">{wallet.account}</span></CardDescription>
              <CardTitle className="text-2xl">
                {wallet.currency} {wallet.balance}
              </CardTitle>
            </CardHeader>
            <CardContent className="text-muted-foreground text-sm">
              {wallet.change}
            </CardContent>
          </Card>
        ))}
      </section>

      <section className="grid gap-6">
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

        <Card className="flex flex-col">
          <CardHeader>
            <CardTitle>Last 5 transactions</CardTitle>
            <CardDescription>Latest activity across wallets.</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Date</TableHead>
                  <TableHead>Details</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recentTransactions.map((tx) => (
                  <TableRow key={tx.id}>
                    <TableCell className="text-muted-foreground">
                      {tx.date}
                    </TableCell>
                    <TableCell className="font-medium">{tx.description}</TableCell>
                    <TableCell>{tx.type}</TableCell>
                    <TableCell className="text-right font-medium">
                      {tx.amount}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

      </section>
    </div>
  )
}
