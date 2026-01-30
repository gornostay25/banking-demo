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

export default function TransferPage() {
  return (
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
            Choose a recipient, currency, and amount to send.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="recipient">Recipient</Label>
            <Select>
              <SelectTrigger id="recipient">
                <SelectValue placeholder="Select recipient" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="alex">Alex Morgan</SelectItem>
                <SelectItem value="jamie">Jamie Rivera</SelectItem>
                <SelectItem value="taylor">Taylor Price</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-muted-foreground text-xs">
              You can also enter an account ID in the next step.
            </p>
          </div>

          <div className="grid gap-4 sm:grid-cols-[1fr_160px]">
            <div className="space-y-2">
              <Label htmlFor="transfer-amount">Amount</Label>
              <Input
                id="transfer-amount"
                placeholder="0.00"
                aria-invalid="true"
              />
              <p className="text-destructive text-xs">
                Please enter a valid amount greater than 0.
              </p>
            </div>
            <div className="space-y-2">
              <Label>Currency</Label>
              <Select defaultValue="usd">
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="usd">USD</SelectItem>
                  <SelectItem value="eur">EUR</SelectItem>
                </SelectContent>
              </Select>
              <p className="text-muted-foreground text-xs">
                Available: 12,480.50 USD
              </p>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="reference">Reference</Label>
            <Input id="reference" placeholder="Add a note (optional)" />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-muted-foreground text-xs">
            Transfers are processed instantly.
          </p>
          <Button>Submit transfer</Button>
        </CardFooter>
      </Card>
    </div>
  )
}
