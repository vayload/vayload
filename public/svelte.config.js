import adapter from "@sveltejs/adapter-static";

/** @type {import('@sveltejs/kit').Config} */
const config = {
    kit: {
        // adapter: adapter({
        //     pages: "build",
        //     assets: "build",
        //     fallback: undefined,
        //     precompress: false,
        //     strict: true,
        // }),
        adapter: adapter({
            fallback: "index.html",
        }),
        alias: {
            $features: "src/features",
            "$features/*": "src/features/*",
        },
    },
};

export default config;
