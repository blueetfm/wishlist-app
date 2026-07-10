import Link from "next/link";

interface WishlistCardProps {
  title: string;
  description: string;
  href: string;
}

export function WishlistCard({ title, description, href }: WishlistCardProps) {
  return (
    <article className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm transition hover:shadow-md">
      <h3 className="text-xl font-semibold">{title}</h3>
      <p className="mt-2 text-slate-600">{description}</p>
      <Link href={href} className="mt-4 inline-flex rounded-xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white hover:bg-slate-800">
        View wishlist
      </Link>
    </article>
  );
}
