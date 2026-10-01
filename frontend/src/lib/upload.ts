import { supabase } from "@/lib/supabase";

const BUCKET = "item-images";
const MAX_FILE_SIZE_BYTES = 5 * 1024 * 1024;
const ALLOWED_MIME_TYPES = ["image/jpeg", "image/png", "image/webp"];

export class UploadError extends Error {}

/**
 * Uploads a user-selected product photo to the `item-images` Supabase Storage
 * bucket under a per-user path (bucket RLS should restrict writes to
 * `auth.uid() = <userId>`) and returns its public URL.
 */
export async function uploadItemImage(file: File, userId: string): Promise<string> {
  if (!ALLOWED_MIME_TYPES.includes(file.type)) {
    throw new UploadError("Unsupported file type. Please upload a JPEG, PNG, or WEBP image.");
  }
  if (file.size > MAX_FILE_SIZE_BYTES) {
    throw new UploadError("Image is too large. Maximum size is 5MB.");
  }

  const extension = file.name.split(".").pop()?.toLowerCase() || "jpg";
  const path = `${userId}/${crypto.randomUUID()}.${extension}`;

  const { error } = await supabase.storage.from(BUCKET).upload(path, file, {
    cacheControl: "3600",
    upsert: false,
    contentType: file.type,
  });
  if (error) {
    throw new UploadError(error.message);
  }

  return supabase.storage.from(BUCKET).getPublicUrl(path).data.publicUrl;
}
