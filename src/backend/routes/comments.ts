import type { FastifyInstance } from "fastify";

export default async function commentRoutes(fastify: FastifyInstance) {
  fastify.post("/comments", async (request, reply) => {
    const body = request.body as {
      wishlistId?: string;
      itemId?: string;
      parentId?: string;
      content?: string;
      authorName?: string;
    };

    if (!body.content) {
      return reply.code(400).send({ error: "Content is required" });
    }

    return reply.code(201).send({
      comment: {
        id: "new-comment-id",
        wishlistId: body.wishlistId || null,
        itemId: body.itemId || null,
        parentId: body.parentId || null,
        authorName: body.authorName || null,
        content: body.content,
      },
    });
  });

  fastify.get("/wishlists/:wishlistId/comments", async (request) => {
    const { wishlistId } = request.params as { wishlistId: string };
    return {
      wishlistId,
      comments: [],
    };
  });

  fastify.get("/items/:itemId/comments", async (request) => {
    const { itemId } = request.params as { itemId: string };
    return {
      itemId,
      comments: [],
    };
  });
}
