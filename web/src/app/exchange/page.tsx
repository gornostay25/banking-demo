import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export default function ExchangePage() {
  return (
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
        <CardContent className="space-y-6">
          <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
            <div className="space-y-2">
              <Label htmlFor="exchange-amount">Amount</Label>
              <Input id="exchange-amount" placeholder="0.00" />
              <p className="text-muted-foreground text-xs">
                Minimum exchange amount: 10.00
              </p>
            </div>
            <div className="space-y-2">
              <Label>From</Label>
              <Select defaultValue="usd">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="usd">USD</SelectItem>
                  <SelectItem value="eur">EUR</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
            <div className="space-y-2">
              <Label>Converted amount</Label>
              <Input value="0.00" readOnly />
              <p className="text-muted-foreground text-xs">
                You receive 0.00 in the target wallet.
              </p>
            </div>
            <div className="space-y-2">
              <Label>To</Label>
              <Select defaultValue="eur">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="eur">EUR</SelectItem>
                  <SelectItem value="usd">USD</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-muted-foreground text-xs">
            Exchange rate: 1 USD = 0.92 EUR
          </p>
          <Button>Exchange funds</Button>
        </CardFooter>
      </Card>
    </div>
  )
}
