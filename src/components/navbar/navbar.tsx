"use client";

import { supabase } from "@/lib/supabase";
import { useState, useEffect } from "react";
import type { Session } from "@supabase/supabase-js";
import Link from "next/link";
import { Alert } from "../alert";
import { Button } from "@/components/ui/button";
import { Logo } from "./logo";
// import { NavMenu } from "./nav-menu";
import { NavigationSheet } from "./navigation-sheet";

// https://www.shadcnui-blocks.com/blocks/navbar-01
const Navbar = () => {
  // https://www.youtube.com/watch?v=2SEz6SK_ekE
  // https://supabase.com/docs/guides/auth/social-login/auth-google?queryGroups=platform&platform=web#google-pre-built
  const [session, setSession] = useState<Session | null>(null);
  const [signingIn, setSigningIn] = useState(false);
  const [showSignOutAlert, setShowSignOutAlert] = useState(false);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
    });

    return () => subscription.unsubscribe();
  }, []);

  const signInWithGoogle = async () => {
    setSigningIn(true);
    supabase.auth.signInWithOAuth({
      provider: "google",
    });

  };
  const signOut = async () => {
    const { error } = await supabase.auth.signOut();

    if (!error) {
      setShowSignOutAlert(true);
      setTimeout(() => 
        setShowSignOutAlert(false), 3000)
    }
  };

  return (
    <div className="min-h-screen bg-muted">
      <nav className="h-16 bg-background border-b">
        <div className="h-full flex items-center justify-between max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
          <Logo />
          {/* <WishlistSvg /> */}

          {/* <NavMenu className="hidden md:block" /> */}

          {showSignOutAlert && <Alert alertText="You are now signed out!"/>}
          <div className="flex items-center gap-3">
            <Button variant="outline" className="hidden sm:inline-flex">
              Create a Wishlist
            </Button>

            {session ? 
            <Button onClick={signOut}>Sign out</Button> :
            <Button onClick={signInWithGoogle}>{signingIn ? "Signing you in...": "Sign in with Google"}</Button>
            }
            

            {/* Mobile Menu */}
            <div className="md:hidden">
              <NavigationSheet />
            </div>
          </div>
        </div>
      </nav>
    </div>
  );
};

export default Navbar;
