"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"

import { cn } from "@/lib/utils"

const navItems = [
  { href: "/", label: "Dashboard" },
  { href: "/transfer", label: "Transfer" },
  { href: "/exchange", label: "Exchange" },
  { href: "/transactions", label: "Transactions" },
]

export function Navigation() {
  const pathname = usePathname()

  return (
    <header className="border-border/70 bg-background/80 sticky top-0 z-40 border-b backdrop-blur">
      <div className="mx-auto flex w-full max-w-6xl items-center justify-between gap-6 px-6 py-4">
        <Link href="/" className="text-lg font-semibold tracking-tight">
          Inquire Bank
        </Link>
        <nav className="flex flex-wrap items-center gap-3 text-sm">
          {navItems.map((item) => {
            const isActive =
              item.href === "/" ? pathname === "/" : pathname.startsWith(item.href)

            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "text-muted-foreground hover:text-foreground rounded-full px-3 py-1 transition",
                  isActive && "bg-muted text-foreground"
                )}
              >
                {item.label}
              </Link>
            )
          })}
        </nav>
      </div>
    </header>
  )
}
