import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  webpack(config) {
    // Add rule to handle .svg files with @svgr/webpack
    config.module.rules.push({
      test: /\.svg$/i,
      issuer: /\.[jt]sx?$/,
      use: ["@svgr/webpack"],
    });

    return config;
  },

  reactStrictMode: true,
};

export default nextConfig;
