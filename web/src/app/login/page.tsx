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

export default function LoginPage() {
  return (
    <div className="flex min-h-[70vh] w-full items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Welcome back</CardTitle>
          <CardDescription>
            Sign in to manage your wallets and transactions.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input id="email" type="email" placeholder="you@inquire.bank" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" placeholder="••••••••" />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col gap-3">
          <Button className="w-full">Sign in</Button>
          <p className="bg-muted/60 text-muted-foreground text-xs leading-relaxed rounded-md border px-3 py-2 w-full">
            <span className="font-medium text-foreground">Emails:</span>{" "}
            <span className="select-all">user1@test.com</span>&nbsp;&bull;&nbsp;
            <span className="select-all">user2@test.com</span>&nbsp;&bull;&nbsp;
            <span className="select-all">user3@test.com</span>&nbsp;
            <br />
            <span className="font-medium text-foreground">Password:</span>{" "}
            <span className="select-all">password</span>
          </p>
        </CardFooter>
      </Card>
    </div>
  )
}
