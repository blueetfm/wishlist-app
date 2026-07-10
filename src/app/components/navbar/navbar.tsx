"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/app/context/auth-context";
import { Alert } from "../alert";
import { Button } from "@/app/components/ui/button";
import { NavigationSheet } from "./navigation-sheet";

// https://www.shadcnui-blocks.com/blocks/navbar-01
const Navbar = () => {
  // https://www.youtube.com/watch?v=2SEz6SK_ekE
  // https://supabase.com/docs/guides/auth/social-login/auth-google?queryGroups=platform&platform=web#google-pre-built
  const {
    user,
    loading,
    isAuthenticated,
    identityProvider,
    signInWithGoogle,
    signOut,
  } =
    useAuth();
  const [signingIn, setSigningIn] = useState(false);
  const [showSignOutAlert, setShowSignOutAlert] = useState(false);

  const handleSignInWithGoogle = async () => {
    try {
      setSigningIn(true);
      await signInWithGoogle();
    } finally {
      setSigningIn(false);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
      setShowSignOutAlert(true);
      setTimeout(() => setShowSignOutAlert(false), 3000);
    } catch {
      setShowSignOutAlert(false);
    }
  };

  return (
    <nav className="border-b bg-background">
      <div className="mx-auto flex h-16 max-w-screen-xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <img
          src="/wishlist.svg"
          alt="Wishlist logo"
          className="h-auto w-auto"
        />

        {/* <NavMenu className="hidden md:block" /> */}

        {showSignOutAlert && <Alert alertText="You are now signed out!" />}
        <div className="flex items-center gap-3">
          <Button asChild variant="outline" className="hidden sm:inline-flex">
            <Link href="/wishlists/new">Create a Wishlist</Link>
          </Button>

          {isAuthenticated ? (
            <Button onClick={handleSignOut}>Sign out</Button>
          ) : (
            <Button onClick={handleSignInWithGoogle} disabled={loading || signingIn}>
              {loading ? "Checking identity..." : signingIn ? "Signing you in..." : "Sign in with Google"}
            </Button>
          )}

          {isAuthenticated ? (
            <p className="hidden text-xs text-slate-600 sm:block">
              {user?.email ?? "Signed in"}
              {/* {identityProvider ? ` (${identityProvider})` : ""} */}
            </p>
          ) : null}

          {/* Mobile Menu */}
          <div className="md:hidden">
            <NavigationSheet />
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
