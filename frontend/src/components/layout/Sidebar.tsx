"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import clsx from "clsx";

const navItems = [
  { href: "/", label: "Home", icon: "⌂" },
  { href: "/onboarding", label: "Profile", icon: "👤" },
  { href: "/diagnosis", label: "Diagnosis", icon: "◉" },
  { href: "/plan", label: "Plan", icon: "☰" },
  { href: "/practice", label: "Practice", icon: "⌨" },
  { href: "/mock-interview", label: "Mock", icon: "🎙" },
  { href: "/progress", label: "Progress", icon: "📈" },
  { href: "/settings", label: "Settings", icon: "⚙" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden w-56 shrink-0 border-r border-zinc-200 bg-zinc-50 lg:block dark:border-zinc-800 dark:bg-zinc-950">
      <div className="sticky top-0 flex h-screen flex-col px-3 py-5">
        <div className="mb-6 px-2">
          <p className="text-xs font-semibold uppercase tracking-wide text-teal-600">
            Interview Coach
          </p>
          <p className="mt-1 text-sm text-zinc-500">Recovery workspace</p>
        </div>
        <nav className="flex-1 space-y-1">
          {navItems.map((item) => {
            const active =
              item.href === "/"
                ? pathname === "/"
                : pathname.startsWith(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={clsx(
                  "flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  active
                    ? "bg-teal-600 text-white"
                    : "text-zinc-600 hover:bg-zinc-200 dark:text-zinc-300 dark:hover:bg-zinc-800",
                )}
              >
                <span aria-hidden>{item.icon}</span>
                {item.label}
              </Link>
            );
          })}
        </nav>
      </div>
    </aside>
  );
}

export function MobileNav() {
  const pathname = usePathname();
  const mobileItems = navItems.filter((item) =>
    ["/", "/practice", "/plan", "/progress", "/settings"].includes(item.href),
  );

  return (
    <nav className="fixed inset-x-0 bottom-0 z-40 border-t border-zinc-200 bg-white/95 backdrop-blur lg:hidden dark:border-zinc-800 dark:bg-zinc-950/95">
      <ul className="grid grid-cols-5">
        {mobileItems.map((item) => {
          const active =
            item.href === "/"
              ? pathname === "/"
              : pathname.startsWith(item.href);
          return (
            <li key={item.href}>
              <Link
                href={item.href}
                className={clsx(
                  "flex flex-col items-center gap-0.5 px-1 py-2 text-[10px] font-medium",
                  active ? "text-teal-600" : "text-zinc-500",
                )}
              >
                <span className="text-base" aria-hidden>
                  {item.icon}
                </span>
                {item.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
