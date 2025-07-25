"use client";

import { useEffect, useState } from "react";
import Script from "next/script";
import { createClient } from "@/utils/supabase/client";

type CredentialResponse = {
  credential: string;
  select_by: string;
  clientId: string;
};

export default function LoginPage() {
  const supabase = createClient();
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);

    // @ts-ignore
    window.handleSignInWithGoogle = async (response: CredentialResponse) => {
      const { data, error } = await supabase.auth.signInWithIdToken({
        provider: "google",
        token: response.credential,
      });

      if (error) {
        console.error("Supabase signIn error:", error);
      } else {
        console.log("Signed in!", data);
        // redirect or update UI here
      }
    };
  }, [supabase]);

  return (
    <>
      <Script
        src="https://accounts.google.com/gsi/client"
        strategy="beforeInteractive"
      />

      {isMounted && (
        <>
          <div
            id="g_id_onload"
            data-client_id={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID!}
            data-context="signin"
            data-ux_mode="popup"
            data-callback="handleSignInWithGoogle"
            data-auto_prompt="false"
          ></div>

          <div
            className="g_id_signin"
            data-type="icon"
            data-shape="square"
            data-theme="outline"
            data-text="signin_with"
            data-size="large"
          ></div>
        </>
      )}
    </>
  );
}
