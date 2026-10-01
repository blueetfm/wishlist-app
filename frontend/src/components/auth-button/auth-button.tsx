"use client";

import { useEffect, useState } from "react";
import type { Session } from "@supabase/supabase-js";
import { supabase } from "@/lib/supabase";
import { Alert } from "@/components/alert";
import { StylizedButton } from "@/components/ui/stylized-button";

interface AuthButtonProps {
  className?: string;
}

/** Self-contained Google sign-in/sign-out button, tracking its own Supabase session. */
export function AuthButton({ className = "" }: AuthButtonProps) {
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
      setTimeout(() => setShowSignOutAlert(false), 3000);
    }
  };

  return (
    <>
      {showSignOutAlert && <Alert alertText="You are now signed out!" />}
      <StylizedButton size="sm" onClick={session ? signOut : signInWithGoogle} className={className}>
        {session ? "Sign Out" : signingIn ? "Signing you in..." : "Sign In"}
      </StylizedButton>
    </>
  );
}

export default AuthButton;
