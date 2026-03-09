import fs from "fs";
import path from "path";

const FEATURES_BASE = "d:/worklance/vayload.io/vayload/public/src/features";

const features = [
    "auth",
    "settings",
    "users",
    "roles",
    "dashboard",
    "content-types",
    "entries",
    "media",
    "audit-logs",
    "plugins",
];

const filesToCreate = ["dtos.ts", "types.ts", "services.ts", "store.svelte.ts", "index.ts"];

function scaffold() {
    features.forEach((feature) => {
        const dir = path.join(FEATURES_BASE, feature);
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, { recursive: true });
            console.log(`Created directory: ${dir}`);
        }

        filesToCreate.forEach((file) => {
            const filePath = path.join(dir, file);
            if (!fs.existsSync(filePath)) {
                fs.writeFileSync(filePath, "// Skeleton file for " + feature + " " + file + "\n");
                console.log(`Created file: ${filePath}`);
            } else {
                console.log(`Skipped existing file: ${filePath}`);
            }
        });
    });
}

scaffold();
