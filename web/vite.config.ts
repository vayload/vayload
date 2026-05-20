import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
    plugins: [sveltekit(), tailwindcss()],
    // User for reverse proxy with nginx|apache
    server: {
        port: 8090,
        strictPort: true,
        allowedHosts: ["localhost", "app.vayload.dev", "app.vayload.me", "dash.vayload.dev", "dash.vayload.me"],
        hmr: {
            protocol: "ws",
            host: "localhost",
            port: 8090,
            path: "vite",
        },
    },
    define: {
        __DEV__: JSON.stringify(true),
    },
});
