import { createClient } from "@supabase/supabase-js";
import { sql } from "drizzle-orm";
import type { FastifyInstance, FastifyRequest } from "fastify";
import type { User as SupabaseUser } from "@supabase/supabase-js";
import { db } from "../db/client";

// minimal user for application use, derived from Supabase user
type ClientAuthUser = {
  id: string;
  email: string | null;
  role: string | null;
};

// full Supabase user type, for internal use only
type InternalAuthUser = {
  user: SupabaseUser;
  clientUser: ClientAuthUser;
};

let authClient: ReturnType<typeof createClient> | null = null;

// Create a Supabase client for server-side use, using the secret key
function getAuthClient(): ReturnType<typeof createClient> {
  if (authClient) {
    return authClient;
  }

  const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL;
  const supabaseApiKey =
    process.env.SUPABASE_SECRET_KEY ?? process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY;

  if (!supabaseUrl || !supabaseApiKey) {
    throw new Error(
      "Missing Supabase auth env vars. Set NEXT_PUBLIC_SUPABASE_URL and SUPABASE_SECRET_KEY",
    );
  }

  authClient = createClient(supabaseUrl, supabaseApiKey, {
    auth: {
      autoRefreshToken: false,
      persistSession: false,
    },
  });

  return authClient;
}

// Read the bearer token from the Authorization header of the request
function readBearerToken(request: FastifyRequest): string | null {
  const authHeader = request.headers.authorization;

  if (!authHeader) {
    return null;
  }

  const [scheme, token] = authHeader.split(" ");

  if (!scheme || !token || scheme.toLowerCase() !== "bearer") {
    return null;
  }

  return token;
}

// Request to return the full internal auth user for the request, or null if not authenticated
async function getInternalAuthUser(request: FastifyRequest): Promise<InternalAuthUser | null> {
  const token = readBearerToken(request);

  if (!token) {
    return null;
  }

  const supabase = getAuthClient();
  const { data, error } = await supabase.auth.getUser(token);

  if (error || !data.user) {
    return null;
  }

  return {
    user: data.user,
    clientUser: {
      id: data.user.id,
      email: data.user.email ?? null,
      role: (data.user.app_metadata?.role as string | undefined) ?? null,
    },
  };
}

// Resolve the username and avatar URL from the Supabase user metadata
function resolveProfileFields(user: SupabaseUser) {
  const metadata = user.user_metadata ?? {};

  const username =
    (typeof metadata.user_name === "string" && metadata.user_name) ||
    (typeof metadata.full_name === "string" && metadata.full_name) ||
    (typeof metadata.name === "string" && metadata.name) ||
    (typeof metadata.preferred_username === "string" && metadata.preferred_username) ||
    user.email?.split("@")[0] ||
    user.id;

  const avatarUrl =
    (typeof metadata.avatar_url === "string" && metadata.avatar_url) ||
    (typeof metadata.picture === "string" && metadata.picture) ||
    null;

  return { username, avatarUrl };
}

// Request to return the minimal authenticated user for the request, or null if not authenticated
export async function lookupRequestUser(request: FastifyRequest): Promise<ClientAuthUser | null> {
  const authContext = await getInternalAuthUser(request);
  return authContext?.clientUser ?? null;
}

// Upsert the user's profile in the database based on the request
export async function upsertProfileForRequest(request: FastifyRequest) {
  const authContext = await getInternalAuthUser(request);

  if (!authContext) {
    return null;
  }

  const { user } = authContext;

  const { username, avatarUrl } = resolveProfileFields(user);

  await db.execute(sql`
    insert into profiles (id, username, avatar_url)
    values (${user.id}, ${username}, ${avatarUrl})
    on conflict (id) do update set
      username = excluded.username,
      avatar_url = excluded.avatar_url
  `);

  return {
    id: user.id,
    username,
    avatarUrl,
  };
}

// Fastify plugin to register auth routes
export default async function authRoutes(fastify: FastifyInstance) {
  fastify.post("/profile", async (request, reply) => {
    const profile = await upsertProfileForRequest(request);

    if (!profile) {
      return reply.code(401).send({ error: "Missing or invalid bearer token" });
    }

    return reply.code(200).send({ profile });
  });
}
