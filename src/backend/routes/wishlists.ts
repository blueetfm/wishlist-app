import type { FastifyInstance } from "fastify";
import { db } from "../db/client";
import { wishlists } from "../db/schema";
import { lookupRequestUser, upsertProfileForRequest } from "./auth";

export default async function wishlistRoutes(fastify: FastifyInstance) {
  fastify.get("/", async () => {
    return { wishlists: [] };
  });

  fastify.get("/me", async (request, reply) => {
    const user = await lookupRequestUser(request);
    console.log("User from /me route:", user);

    if (!user) {
      return reply.code(401).send({ error: "Missing or invalid bearer token" });
    }

    return {
      userId: user.id,
      wishlists: [],
    };
  });

  fastify.post("/new", async (request, reply) => {
    const body = request.body as {
      title?: string;
      description?: string;
    };

    const user = await lookupRequestUser(request);
    const normalizedTitle = body.title?.trim();

    if (!user) {
      return reply.code(401).send({ error: "Missing or invalid bearer token" });
    }

    await upsertProfileForRequest(request);

    if (!normalizedTitle) {
      return reply.code(400).send({ error: "Title is required" });
    }

    try {
      const [createdWishlist] = await db.insert(wishlists).values({
        title: normalizedTitle,
        description: body.description ?? null,
        user_id: user.id,
        share_token: crypto.randomUUID(),
        created_at: new Date(),
      }).returning({
        id: wishlists.id,
        userId: wishlists.user_id,
        title: wishlists.title,
        description: wishlists.description,
        createdAt: wishlists.created_at,
      });

      if (!createdWishlist) {
        return reply.code(500).send({ error: "Failed to create wishlist" });
      }

      return reply.code(201).send({
        id: createdWishlist.id,
        userId: createdWishlist.userId,
        title: createdWishlist.title,
        description: createdWishlist.description,
        createdAt: createdWishlist.createdAt,
      });
    } catch (error) {
      fastify.log.error(error);

      return reply.code(500).send({
        error: "Failed to create wishlist",
        details: error instanceof Error ? error.message : "Unknown database error",
      });
    }
  });

  fastify.get("/shared/:shareToken/items", async (request) => {
    const { shareToken } = request.params as { shareToken: string };
    return {
      shareToken,
      items: [],
    };
  });

  fastify.get("/:wishlistId/items", async (request) => {
    const { wishlistId } = request.params as { wishlistId: string };
    return {
      wishlistId,
      items: [],
    };
  });

  fastify.post("/:wishlistId/items", async (request, reply) => {
    const { wishlistId } = request.params as { wishlistId: string };
    const body = request.body as {
      name?: string;
      productUrl?: string;
      imageUrl?: string;
      price?: number;
    };

    if (!body.name) {
      return reply.code(400).send({ error: "Item name is required" });
    }

    return reply.code(201).send({
      item: {
        id: "new-item-id",
        wishlistId,
        ...body,
      },
    });
  });

  fastify.post("/shared/:shareToken/items/:itemId/claim", async (request, reply) => {
    const { shareToken, itemId } = request.params as {
      shareToken: string;
      itemId: string;
    };
    const body = request.body as { guestName?: string; message?: string };

    if (!body.guestName && !body.message) {
      return reply.code(400).send({ error: "guestName or message is required" });
    }

    return reply.code(201).send({
      claim: {
        id: "new-claim-id",
        itemId,
        guestName: body.guestName || null,
        message: body.message || null,
      },
    });
  });
}
