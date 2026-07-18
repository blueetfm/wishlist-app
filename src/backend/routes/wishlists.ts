import type { FastifyInstance } from "fastify";
import { db } from "../db/client";
import { and, eq } from "drizzle-orm";
import { wishlists, wishlistItems, claims, comments } from "../db/schema";
import { lookupRequestUser, upsertProfileForRequest } from "./auth";
import { z } from "zod";

/* ==================== REQUEST SCHEMAS ==================== */

const CreateWishlistRequestSchema = z.object({
  title: z.string().min(1, "Title is required").max(255, "Title must be 255 characters or less"),
  description: z.string().optional(),
});

const CreateWishlistItemRequestSchema = z.object({
  name: z.string().min(1, "Item name is required").max(255, "Item name must be 255 characters or less"),
  product_url: z.string().url("Invalid product URL").max(2048, "Product URL must be 2048 characters or less").optional().or(z.literal("")),
  image_url: z.string().url("Invalid image URL").max(2048, "Image URL must be 2048 characters or less").optional().or(z.literal("")),
  price: z.string().refine(
    (val) => !val || !isNaN(parseFloat(val)),
    "Price must be a valid number"
  ).optional(),
});

const UpdateWishlistItemRequestSchema = z.object({
  name: z.string().min(1, "Item name is required").max(255, "Item name must be 255 characters or less"),
  product_url: z.string().url("Invalid product URL").max(2048, "Product URL must be 2048 characters or less").optional().or(z.literal("")),
  image_url: z.string().url("Invalid image URL").max(2048, "Image URL must be 2048 characters or less").optional().or(z.literal("")),
  price: z.string().refine(
    (val) => !val || !isNaN(parseFloat(val)),
    "Price must be a valid number"
  ).optional(),
});

const ParamsWithWishlistIdSchema = z.object({
  wishlistId: z.string()
});

const ParamsWithWishlistAndItemIdSchema = z.object({
  wishlistId: z.string(),
  itemId: z.string(),
});

/* ==================== RESPONSE SCHEMAS ==================== */

const WishlistResponseSchema = z.object({
  id: z.string(),
  userId: z.string(),
  title: z.string(),
  description: z.string().nullable(),
  createdAt: z.date(),
});

const WishlistItemResponseSchema = z.object({
  id: z.string(),
  wishlist_id: z.string(),
  name: z.string(),
  product_url: z.string().nullable(),
  image_url: z.string().nullable(),
  price: z.string().nullable(),
  created_at: z.date(),
});

const UserWishlistsResponseSchema = z.object({
  userId: z.string(),
  wishlists: z.array(WishlistResponseSchema),
});

const ErrorResponseSchema = z.object({
  error: z.string(),
  details: z.string().optional(),
});

/* ==================== TYPE EXPORTS ==================== */

export type CreateWishlistRequest = z.infer<typeof CreateWishlistRequestSchema>;
export type CreateWishlistItemRequest = z.infer<typeof CreateWishlistItemRequestSchema>;
export type UpdateWishlistItemRequest = z.infer<typeof UpdateWishlistItemRequestSchema>;
export type WishlistResponse = z.infer<typeof WishlistResponseSchema>;
export type WishlistItemResponse = z.infer<typeof WishlistItemResponseSchema>;
export type UserWishlistsResponse = z.infer<typeof UserWishlistsResponseSchema>;
export type ErrorResponse = z.infer<typeof ErrorResponseSchema>;

/* ==================== ROUTES ==================== */
export default async function wishlistRoutes(fastify: FastifyInstance) {
  // Fetch user's wishlists
  fastify.get(
    "/me",
    {
      schema: {
        response: {
          200: UserWishlistsResponseSchema,
          401: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const user = await lookupRequestUser(request);
      console.log("User from /me route:", user);

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      return {
        userId: user.id,
        wishlists: [],
      };
    }
  );

  // Create a new wishlist
  fastify.post(
    "/new",
    {
      schema: {
        body: CreateWishlistRequestSchema,
        response: {
          201: WishlistResponseSchema,
          400: ErrorResponseSchema,
          401: ErrorResponseSchema,
          500: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const body = request.body as CreateWishlistRequest;

      const user = await lookupRequestUser(request);
      const normalizedTitle = body.title?.trim();

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      if (!normalizedTitle) {
        return reply.code(400).send({ error: "Title is required" });
      }

      try {
        // insert wishlist into DB - drizzle will automatcialy geneate UUID
        const [createdWishlist] = await db
        .insert(wishlists).values({
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
    }
  );

  // List wishlist items (owner view)
  fastify.get(
    "/:wishlistId/items",
    {
      schema: {
        params: ParamsWithWishlistIdSchema,
        response: {
          200: z.array(WishlistItemResponseSchema),
          401: ErrorResponseSchema,
          404: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const { wishlistId } = request.params as { wishlistId: string };
      const user = await lookupRequestUser(request);

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      // fetch the wishlist if exists and belongs to the user
      const wishlist = await db
        .select()
        .from(wishlists)
        .where(
          and(
            eq(wishlists.id, wishlistId),
            eq(wishlists.user_id, user.id)
          )
        );

      if (!wishlist) {
        return reply.code(404).send({ error: "Wishlist not found" });
      }

      // fetch and return the items for the wishlist
      const items = await db
        .select()
        .from(wishlistItems)
        .where(eq(wishlistItems.wishlist_id, wishlistId));

      return reply.code(200).send(items);
    }
  );

  // Add wishlist item as owner
  fastify.post(
    "/:wishlistId/items",
    {
      schema: {
        params: ParamsWithWishlistIdSchema,
        body: z.array(CreateWishlistItemRequestSchema),
        response: {
          201: WishlistItemResponseSchema,
          400: ErrorResponseSchema,
          401: ErrorResponseSchema,
          404: ErrorResponseSchema,
          500: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const { wishlistId } = request.params as { wishlistId: string };

      const items = request.body as CreateWishlistItemRequest[];
      const user = await lookupRequestUser(request);

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      // fetch the wishlist if exists and belongs to the user
      const wishlist = await db
        .select()
        .from(wishlists)
        .where(
          and(
            eq(wishlists.id, wishlistId),
            eq(wishlists.user_id, user.id)
          )
        );

      if (!wishlist) {
        return reply.code(404).send({ error: "Wishlist not found" });
      }

      // Insert items into the wishlist - drizzle will automatcialy geneate UUID
      try {
        const insertedItems = await db.transaction(async (tx) => {
          return await tx
            .insert(wishlistItems)
            .values(
              items.map((item) => ({
                wishlist_id: wishlistId,
                name: item.name, 
                product_url: item.product_url ?? null,
                image_url: item.image_url ?? null,
                price: item.price ?? null,
                created_at: new Date(),
              }))
            )
            .returning();
        });

        return reply.code(201).send(insertedItems);
      } catch (error) {
        fastify.log.error(error);
        return reply.code(500).send({
          error: "Failed to add items to wishlist",
          details: error instanceof Error ? error.message : "Unknown database error",
        });
      }
          }
  );

  // Edit a particular item in a wishlist as owner
  fastify.put(
    "/:wishlistId/items/:itemId",
    {
      schema: {
        params: ParamsWithWishlistAndItemIdSchema,
        body: UpdateWishlistItemRequestSchema,
        response: {
          200: WishlistItemResponseSchema,
          401: ErrorResponseSchema,
          404: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const { wishlistId, itemId } = request.params as { wishlistId: string; itemId: string };
      const item = request.body as UpdateWishlistItemRequest;
      const user = await lookupRequestUser(request);

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      // fetch the wishlist if exists and belongs to the user
      const wishlist = await db
        .select()
        .from(wishlists)
        .where(
          and(
            eq(wishlists.id, wishlistId),
            eq(wishlists.user_id, user.id)
          )
        );

      if (!wishlist) {
        return reply.code(404).send({ error: "Wishlist not found" });
      }

      // fetch the item if exists and belongs to the wishlist
      const dbItem = await db
        .select()
        .from(wishlistItems)
        .where(
          and(
            eq(wishlistItems.id, itemId),
            eq(wishlistItems.wishlist_id, wishlistId)
          )
        );
      
        if (!dbItem) {
          return reply.code(404).send({ error: "Item not found" });
        }
      
      // Update the item
      const updatedItem = await db
        .update(wishlistItems)
        .set({
          name: item.name,
          product_url: item.product_url ?? null,
          image_url: item.image_url ?? null,
          price: item.price ?? null,
        })
        .where(eq(wishlistItems.id, itemId))
        .returning();

      return reply.code(200).send(updatedItem);
    }
  );

  // Delete a particular item in a wishlist as owner
  fastify.delete(
    "/:wishlistId/items/:itemId",
    {
      schema: {
        params: ParamsWithWishlistAndItemIdSchema,
        response: {
          204: z.null(),
          401: ErrorResponseSchema,
          404: ErrorResponseSchema,
        },
      },
    },
    async (request, reply) => {
      const { wishlistId, itemId } = request.params as { wishlistId: string; itemId: string };
      const user = await lookupRequestUser(request);

      if (!user) {
        return reply.code(401).send({ error: "Missing or invalid bearer token" });
      }

      // fetch the wishlist if exists and belongs to the user
      const wishlist = await db
        .select()
        .from(wishlists)
        .where(
          and(
            eq(wishlists.id, wishlistId),
            eq(wishlists.user_id, user.id)
          )
        );

      if (!wishlist) {
        return reply.code(404).send({ error: "Wishlist not found" });
      }

      // fetch the item if exists and belongs to the wishlist
      const dbItem = await db
        .select()
        .from(wishlistItems)
        .where(
          and(
            eq(wishlistItems.id, itemId),
            eq(wishlistItems.wishlist_id, wishlistId)
          )
        );

      if (!dbItem) {
        return reply.code(404).send({ error: "Item not found" });
      }

      // Delete the item
      await db
        .delete(wishlistItems)
        .where(eq(wishlistItems.id, itemId));

      return reply.code(204).send();

    }
  );
}
