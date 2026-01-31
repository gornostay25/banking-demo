"use client";

import { LogOut, Menu } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/", label: "Dashboard" },
  { href: "/transfer", label: "Transfer" },
  { href: "/exchange", label: "Exchange" },
  { href: "/transactions", label: "Transactions" },
];

export function Navigation() {
  const pathname = usePathname();
  const { isAuthenticated, logout, isLoggingOut } = useAuth();
  const [isOpen, setIsOpen] = React.useState(false);

  const closeMenu = () => setIsOpen(false);
  const handleLogout = async () => {
    await logout();
    closeMenu();
  };

  return (
    <header className="border-border/70 bg-background/80 sticky top-0 z-40 border-b backdrop-blur">
      <div className="mx-auto flex w-full max-w-6xl items-center justify-between gap-6 px-6 py-4">
        <Link href="/" className="text-lg font-semibold tracking-tight">
          Banking Demo
        </Link>
        {isAuthenticated ? (
          <>
            <nav className="hidden flex-wrap items-center gap-3 text-sm md:flex">
              {navItems.map((item) => {
                const isActive =
                  item.href === "/"
                    ? pathname === "/"
                    : pathname.startsWith(item.href);

                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={cn(
                      "text-muted-foreground hover:text-foreground rounded-full px-3 py-1 transition",
                      isActive && "bg-muted text-foreground",
                    )}
                  >
                    {item.label}
                  </Link>
                );
              })}
            </nav>
            <div className="hidden items-center gap-2 md:flex">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleLogout}
                disabled={isLoggingOut}
              >
                <LogOut className="size-4" />
                {isLoggingOut ? "Logging out..." : "Logout"}
              </Button>
            </div>
            <Drawer open={isOpen} onOpenChange={setIsOpen} direction="right">
              <DrawerTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  className="md:hidden"
                  aria-label="Open navigation menu"
                >
                  <Menu className="size-4" />
                </Button>
              </DrawerTrigger>
              <DrawerContent className="flex flex-col gap-6">
                <DrawerHeader>
                  <DrawerTitle>Navigation</DrawerTitle>
                </DrawerHeader>
                <nav className="flex flex-col gap-2 text-sm">
                  {navItems.map((item) => {
                    const isActive =
                      item.href === "/"
                        ? pathname === "/"
                        : pathname.startsWith(item.href);

                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={closeMenu}
                        className={cn(
                          "text-muted-foreground hover:text-foreground rounded-md px-3 py-2 transition",
                          isActive && "bg-muted text-foreground",
                        )}
                      >
                        {item.label}
                      </Link>
                    );
                  })}
                </nav>
                <Button
                  variant="ghost"
                  size="sm"
                  className="justify-start"
                  onClick={handleLogout}
                  disabled={isLoggingOut}
                >
                  <LogOut className="size-4" />
                  {isLoggingOut ? "Logging out..." : "Logout"}
                </Button>
              </DrawerContent>
            </Drawer>
          </>
        ) : null}
      </div>
    </header>
  );
}
