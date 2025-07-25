"use client";

import { supabase } from "@/lib/supabase";
import { useState, useEffect } from "react";
import type { Session } from "@supabase/supabase-js";

export default function Home() {
  // https://www.youtube.com/watch?v=2SEz6SK_ekE
  // https://supabase.com/docs/guides/auth/social-login/auth-google?queryGroups=platform&platform=web#google-pre-built
  const [session, setSession] = useState<Session | null>(null);

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

  const signIn = async () => {
    supabase.auth.signInWithOAuth({
      provider: "google",
    });
  };
  const signOut = async () => {
    const { error } = await supabase.auth.signOut();
  };

  if (!session) {
    return (
      <button onClick={signIn}>Sign in with Google</button>
      // <Auth supabaseClient={supabase} appearance={{ theme: ThemeSupa }} />;
    );
  } else {
    return (
      <div>
        Logged in!
        <button onClick={signOut}>Sign out</button>
      </div>
    );
  }
}
