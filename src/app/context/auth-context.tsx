"use client";

import React, {
  createContext,
  useState,
  useContext,
  ReactNode,
  useMemo,
  useEffect,
} from "react";
import { supabase } from "@/app/lib/supabase";
import { buildUrl } from "@/lib/utils";
import type { Session, User } from "@supabase/supabase-js";

const BACKEND_BASE_URL = process.env.NEXT_PUBLIC_WISHLIST_API_BASE_URL ?? "http://localhost:4000";

interface AuthContextType {
  session: Session | null;
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;
  identityProvider: string | null;
  signInWithGoogle: () => Promise<void>;
  signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProviderWrapper = ({
  children,
}: {
  children: ReactNode;
}): React.JSX.Element => {
  const [user, setUser] = useState<User | null>(null);
  const [session, setSession] = useState<Session | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const client = supabase;

    const syncProfile = async (accessToken: string) => {
      try {
        const response = await fetch(buildUrl(BACKEND_BASE_URL, "/api/auth/profile"), {
          method: "POST",
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        });

        if (!response.ok) {
          console.error("Failed to sync profile", await response.text());
        }
      } catch (error) {
        console.error("Failed to sync profile", error);
      }
    };

    if (!client) {
      setLoading(false);
      return;
    }

    const setData = async () => {
      const {
        data: { session },
      } = await client.auth.getSession();

      setSession(session);
      setUser(session?.user ?? null);
      setLoading(false);

      if (session?.access_token) {
        void syncProfile(session.access_token);
      }
    };

    const {
      data: { subscription },
    } = client.auth.onAuthStateChange((_event, session) => {
      setSession(session);
      setUser(session?.user ?? null);
      setLoading(false);

      if (session?.access_token) {
        void syncProfile(session.access_token);
      }
    });

    void setData();

    return () => subscription.unsubscribe();
  }, []);

  const value = useMemo<AuthContextType>(() => {
    const provider =
      user?.app_metadata?.provider ??
      user?.identities?.[0]?.provider ??
      null;

    return {
      session,
      user,
      loading,
      isAuthenticated: Boolean(session && user),
      identityProvider: provider,
      signInWithGoogle: async () => {
        const client = supabase;

        if (!client) {
          throw new Error("Supabase is not configured for this environment");
        }

        const { error } = await client.auth.signInWithOAuth({
          provider: "google",
        });

        if (error) {
          throw error;
        }
      },
      signOut: async () => {
        const client = supabase;

        if (!client) {
          throw new Error("Supabase is not configured for this environment");
        }

        const { error } = await client.auth.signOut();

        if (error) {
          throw error;
        }
      },
    };
  }, [loading, session, user]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProviderWrapper");
  }

  return context;
};
