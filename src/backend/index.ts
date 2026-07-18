import "dotenv/config";
import fastify from "fastify";
import type { FastifyInstance } from "fastify";
import cors from "@fastify/cors";
import { serializerCompiler, validatorCompiler } from "fastify-type-provider-zod";
import authRoutes from "./routes/auth";
import wishlistRoutes from "./routes/wishlists";
import commentRoutes from "./routes/comments";
import scraperRoutes from "./routes/scraper";

const app = fastify({ logger: true });
app.setValidatorCompiler(validatorCompiler);
app.setSerializerCompiler(serializerCompiler);

const port = Number(process.env.PORT || 4000);

export default function registerRoutes(appInstance: FastifyInstance) {
  appInstance.register(cors, { origin: true });
  appInstance.register(authRoutes, { prefix: "/api/auth" });
  appInstance.register(wishlistRoutes, { prefix: "/api/wishlists" });
  appInstance.register(commentRoutes, { prefix: "/api" });
  appInstance.register(scraperRoutes, { prefix: "/api/scraper" });
}

async function main() {
  try {
    registerRoutes(app);
    await app.listen({ port, host: "0.0.0.0" });
  } catch (error) {
    app.log.error(error);
    process.exit(1);
  }
}

void main();
