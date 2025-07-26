import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    webpack: (config) => {
        // Add rule for SVG files
        config.module.rules.push({
          test: /\.svg$/,
          issuer: /\.[jt]sx?$/,   // Only apply to JS/TS files importing SVGs
          use: ["@svgr/webpack"],
        });
        return config;
    },
    reactStrictMode: true,
};

export default nextConfig;