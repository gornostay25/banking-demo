import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const transactions = [
  {
    id: "tx-3101",
    date: "Jan 30, 2026",
    type: "Transfer",
    direction: "Send",
    currency: "USD",
    amount: "-240.00",
    status: "Completed",
  },
  {
    id: "tx-3102",
    date: "Jan 30, 2026",
    type: "Exchange",
    direction: "Receive",
    currency: "EUR",
    amount: "460.00",
    status: "Completed",
  },
  {
    id: "tx-3103",
    date: "Jan 29, 2026",
    type: "Transfer",
    direction: "Receive",
    currency: "USD",
    amount: "1,200.00",
    status: "Completed",
  },
  {
    id: "tx-3104",
    date: "Jan 28, 2026",
    type: "Exchange",
    direction: "Send",
    currency: "EUR",
    amount: "-320.00",
    status: "Completed",
  },
  {
    id: "tx-3105",
    date: "Jan 27, 2026",
    type: "Transfer",
    direction: "Send",
    currency: "USD",
    amount: "-85.50",
    status: "Completed",
  },
  {
    id: "tx-3106",
    date: "Jan 27, 2026",
    type: "Transfer",
    direction: "Receive",
    currency: "EUR",
    amount: "640.00",
    status: "Completed",
  },
  {
    id: "tx-3107",
    date: "Jan 26, 2026",
    type: "Exchange",
    direction: "Receive",
    currency: "USD",
    amount: "900.00",
    status: "Completed",
  },
  {
    id: "tx-3108",
    date: "Jan 25, 2026",
    type: "Transfer",
    direction: "Send",
    currency: "EUR",
    amount: "-220.00",
    status: "Completed",
  },
]

export default function TransactionsPage() {
  return (
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
          <Select defaultValue="all">
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
            Showing 8 of 25 results for January 2026.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Direction</TableHead>
                <TableHead>Currency</TableHead>
                <TableHead className="text-right">Amount</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((tx) => (
                <TableRow key={tx.id}>
                  <TableCell className="text-muted-foreground">
                    {tx.date}
                  </TableCell>
                  <TableCell className="font-medium">{tx.type}</TableCell>
                  <TableCell>{tx.direction}</TableCell>
                  <TableCell>{tx.currency}</TableCell>
                  <TableCell className="text-right font-medium">
                    {tx.amount}
                  </TableCell>
                  <TableCell>{tx.status}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <div className="flex flex-wrap items-center justify-between gap-4">
        <p className="text-muted-foreground text-sm">
          Page 1 of 4 • 10 items per page
        </p>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            Previous
          </Button>
          <Button size="sm">1</Button>
          <Button variant="outline" size="sm">
            2
          </Button>
          <Button variant="outline" size="sm">
            3
          </Button>
          <Button variant="outline" size="sm">
            Next
          </Button>
        </div>
      </div>
    </div>
  )
}
