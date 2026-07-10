import type { FastifyInstance } from "fastify";

export default async function scraperRoutes(fastify: FastifyInstance) {
  fastify.post("/extract", async (request, reply) => {
    const body = request.body as { url?: string };

    if (!body.url) {
      return reply.code(400).send({ error: "URL is required" });
    }

    return {
      url: body.url,
      name: "Extracted product name",
      price: 29.99,
      imageUrl: null,
    };
  });
}
