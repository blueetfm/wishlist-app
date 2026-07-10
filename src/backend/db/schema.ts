import {
  AnyPgColumn,
  pgTable,
  uuid,
  varchar,
  text,
  timestamp,
  numeric,
} from "drizzle-orm/pg-core";

export const profiles = pgTable("profiles", {
  id: uuid("id").primaryKey(),
  username: text("username"),
  avatar_url: text("avatar_url"),
  created_at: timestamp("created_at", { withTimezone: true }).defaultNow(),
});

export const wishlists = pgTable("wishlists", {
  id: uuid("id").defaultRandom().primaryKey(),
  title: varchar("title", { length: 255 }),
  description: text("description"),
  user_id: uuid("user_id").references(() => profiles.id),
  share_token: uuid("share_token").defaultRandom(),
  created_at: timestamp("created_at").defaultNow(),
});

export const wishlistItems = pgTable("wishlist_items", {
  id: uuid("id").defaultRandom().primaryKey(),
  wishlist_id: uuid("wishlist_id").references(() => wishlists.id),
  name: varchar("name", { length: 255 }),
  product_url: varchar("product_url", { length: 2048 }),
  image_url: varchar("image_url", { length: 2048 }),
  price: numeric("price"),
  created_at: timestamp("created_at").defaultNow(),
});

export const claims = pgTable("claims", {
  id: uuid("id").defaultRandom().primaryKey(),
  item_id: uuid("item_id").references(() => wishlistItems.id),
  guest_name: varchar("guest_name", { length: 255 }),
  message: text("message"),
  claimed_at: timestamp("claimed_at").defaultNow(),
});

export const comments = pgTable("comments", {
  id: uuid("id").defaultRandom().primaryKey(),
  wishlist_id: uuid("wishlist_id").references(() => wishlists.id),
  item_id: uuid("item_id").references(() => wishlistItems.id),
  parent_id: uuid("parent_id").references((): AnyPgColumn => comments.id),
  author_name: varchar("author_name", { length: 255 }),
  author_user_id: uuid("author_user_id").references(() => profiles.id),
  content: text("content"),
  created_at: timestamp("created_at").defaultNow(),
});
