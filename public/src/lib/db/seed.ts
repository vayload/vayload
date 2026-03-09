/**
 * seed.ts
 *
 * Seeds the entire dummy database with realistic, relational data.
 * Call seedDatabase() once on app startup (idempotent by default).
 * Pass force=true to wipe and re-seed all tables.
 */

import { faker } from "@faker-js/faker";
import { ulid } from "./utils";
import {
    usersTable,
    projectsTable,
    rolesTable,
    permissionsTable,
    collectionsTable,
    entriesTable,
    entryFieldsTable,
    fileObjectsTable,
    foldersTable,
    auditLogsTable,
    notificationsTable,
    integrationsTable,
    activitiesTable,
} from "./tables";
import type {
    User,
    Project,
    Role,
    Permission,
    Collection,
    Entry,
    EntryField,
    FileObject,
    FolderObject,
    AuditLog,
    Notification,
    Integration,
    Activity,
    CollectionSchema,
    FieldSchema,
} from "../types";
import { FieldTypes } from "../types";

// ─── Deterministic seed so IDs are stable across sessions per browser ────────
faker.seed(42);

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

const pick = <T>(arr: T[]): T => faker.helpers.arrayElement(arr);
const pickMany = <T>(arr: T[], min = 1, max = 3): T[] => faker.helpers.arrayElements(arr, { min, max });
const isoDate = (d: Date) => d.toISOString();
const pastIso = () => isoDate(faker.date.past({ years: 2 }));
const recentIso = () => isoDate(faker.date.recent({ days: 30 }));

// ─────────────────────────────────────────────────────────────────────────────
// USERS (30 total)
// ─────────────────────────────────────────────────────────────────────────────

function buildUsers(): User[] {
    const superAdmin: User = {
        id: ulid(),
        first_name: "Ana",
        last_name: "Admin",
        username: "ana.admin",
        email: "admin@vayload.io",
        avatar_url: faker.image.avatar(),
        is_super_admin: true,
        last_sign_in_at: recentIso(),
        created_at: pastIso(),
        updated_at: recentIso(),
        password: "admin1234",
    };

    const editors: User[] = Array.from({ length: 5 }, (_, i) => ({
        id: ulid(),
        first_name: faker.person.firstName(),
        last_name: faker.person.lastName(),
        username: `editor.${i + 1}`,
        email: faker.internet.email(),
        avatar_url: faker.image.avatar(),
        is_super_admin: false,
        last_sign_in_at: recentIso(),
        created_at: pastIso(),
        updated_at: recentIso(),
        password: faker.internet.password(),
    }));

    const regular: User[] = Array.from({ length: 24 }, () => ({
        id: ulid(),
        first_name: faker.person.firstName(),
        last_name: faker.person.lastName(),
        username: faker.internet.username(),
        email: faker.internet.email(),
        avatar_url: Math.random() > 0.3 ? faker.image.avatar() : null,
        is_super_admin: false,
        last_sign_in_at: recentIso(),
        created_at: pastIso(),
        updated_at: recentIso(),
        password: faker.internet.password(),
    }));

    return [superAdmin, ...editors, ...regular];
}

// ─────────────────────────────────────────────────────────────────────────────
// PROJECTS (5)
// ─────────────────────────────────────────────────────────────────────────────

function buildProjects(users: User[]): Project[] {
    const presets = [
        { name: "E-Commerce Main", locale: "en-US" },
        { name: "Company Blog", locale: "en-US" },
        { name: "Marketing Site", locale: "es-ES" },
    ];

    const preset = presets.map((p) => ({
        id: ulid(),
        name: p.name,
        slug: faker.helpers.slugify(p.name.toLowerCase()),
        owner_id: pick(users).id,
        settings: { theme: pick(["light", "dark"]), versioning: true, timezone: "UTC" },
        locale: p.locale,
        created_at: pastIso(),
        updated_at: recentIso(),
    }));

    const extra: Project[] = Array.from({ length: 2 }, () => {
        const name = faker.company.name();
        return {
            id: ulid(),
            name,
            slug: faker.helpers.slugify(name.toLowerCase()),
            owner_id: pick(users).id,
            settings: { theme: pick(["light", "dark"]), versioning: false, timezone: "UTC" },
            locale: pick(["en-US", "es-ES", "fr-FR"]),
            created_at: pastIso(),
            updated_at: recentIso(),
        };
    });

    return [...preset, ...extra];
}

// ─────────────────────────────────────────────────────────────────────────────
// ROLES & PERMISSIONS
// ─────────────────────────────────────────────────────────────────────────────

function buildRoles(): Role[] {
    return [
        { id: ulid(), name: "Super Admin", description: "Full unrestricted access" },
        { id: ulid(), name: "Administrator", description: "Manages projects and users" },
        { id: ulid(), name: "Editor", description: "Can create and edit all content" },
        { id: ulid(), name: "Author", description: "Creates and edits own content only" },
        { id: ulid(), name: "Viewer", description: "Read-only access to content" },
        { id: ulid(), name: "API Consumer", description: "Public API read-only token" },
        { id: ulid(), name: "Manager", description: "Manages teams within a project" },
        { id: ulid(), name: "Auditor", description: "Read-only access to audit logs" },
    ];
}

function buildPermissions(): Permission[] {
    const resources = [
        "users",
        "roles",
        "projects",
        "collections",
        "entries",
        "media",
        "audit-logs",
        "integrations",
        "settings",
    ];
    const actions = ["create", "read", "update", "delete"];
    return resources.flatMap((resource) =>
        actions.map((action) => ({
            id: ulid(),
            action,
            resource,
        })),
    );
}

// ─────────────────────────────────────────────────────────────────────────────
// COLLECTIONS (content types, 2-3 per project)
// ─────────────────────────────────────────────────────────────────────────────

const textField = (label: string, required = false): FieldSchema => ({
    type: FieldTypes.TEXT,
    label,
    required,
    default_value: "",
});
const richTextField = (label: string): FieldSchema => ({
    type: FieldTypes.RICH_TEXT,
    label,
    required: false,
});
const numberField = (label: string): FieldSchema => ({
    type: FieldTypes.NUMBER,
    label,
    required: false,
});
const boolField = (label: string, def = false): FieldSchema => ({
    type: FieldTypes.BOOLEAN,
    label,
    required: false,
    default_value: def,
});
const mediaField = (label: string): FieldSchema => ({
    type: FieldTypes.MEDIA,
    label,
    required: false,
});
const relationField = (label: string, relation_to: string): FieldSchema => ({
    type: FieldTypes.RELATIONSHIP,
    label,
    required: false,
    relation_to,
});

const SCHEMAS: Record<string, CollectionSchema> = {
    products: {
        title: textField("Title", true),
        description: richTextField("Description"),
        price: numberField("Price"),
        sku: textField("SKU"),
        published: boolField("Published", false),
        cover: mediaField("Cover Image"),
        category: relationField("Category", "categories"),
    },
    categories: {
        title: textField("Title", true),
        description: richTextField("Description"),
        icon: mediaField("Icon"),
    },
    posts: {
        title: textField("Title", true),
        content: richTextField("Content"),
        excerpt: textField("Excerpt"),
        published: boolField("Published", false),
        cover: mediaField("Cover Image"),
        author: relationField("Author", "users"),
    },
    pages: {
        title: textField("Title", true),
        content: richTextField("Content"),
        slug: textField("Slug", true),
        published: boolField("Published", false),
        seo_title: textField("SEO Title"),
        seo_desc: textField("SEO Description"),
    },
    testimonials: {
        author_name: textField("Author Name", true),
        company: textField("Company"),
        quote: richTextField("Quote"),
        rating: numberField("Rating"),
        avatar: mediaField("Avatar"),
    },
    team: {
        name: textField("Full Name", true),
        role: textField("Role / Title"),
        bio: richTextField("Bio"),
        avatar: mediaField("Photo"),
    },
};

function buildCollections(projects: Project[]): Collection[] {
    const result: Collection[] = [];

    const presets: Array<{ name: string; slug: keyof typeof SCHEMAS; projectIdx: number; single: boolean }> = [
        { name: "Products", slug: "products", projectIdx: 0, single: false },
        { name: "Categories", slug: "categories", projectIdx: 0, single: false },
        { name: "Blog Posts", slug: "posts", projectIdx: 1, single: false },
        { name: "Pages", slug: "pages", projectIdx: 1, single: false },
        { name: "Testimonials", slug: "testimonials", projectIdx: 2, single: false },
        { name: "Team Members", slug: "team", projectIdx: 2, single: false },
    ];

    for (const p of presets) {
        const project = projects[p.projectIdx] ?? projects[0];
        result.push({
            id: ulid(),
            project_id: project.id,
            name: p.name,
            slug: p.slug as string,
            fields_schema: SCHEMAS[p.slug],
            settings: { versioning: true, draftable: true },
            entries: 0,
            single: p.single,
            created_at: pastIso(),
        });
    }

    // Extra dynamic collections for other projects
    const extraNames = ["FAQs", "Events", "Press Releases", "Gallery", "Downloads", "Docs"];
    for (let i = 0; i < 6; i++) {
        result.push({
            id: ulid(),
            project_id: pick(projects).id,
            name: extraNames[i],
            slug: faker.helpers.slugify(extraNames[i].toLowerCase()),
            fields_schema: SCHEMAS.posts,
            settings: { versioning: true, draftable: true },
            entries: 0,
            single: false,
            created_at: pastIso(),
        });
    }

    return result;
}

// ─────────────────────────────────────────────────────────────────────────────
// ENTRIES (200 total, linked to collections + users)
// ─────────────────────────────────────────────────────────────────────────────

function buildEntries(collections: Collection[], users: User[]): Entry[] {
    const entries: Entry[] = [];
    const targetPerCollection = Math.ceil(200 / collections.length);

    for (const col of collections) {
        for (let i = 0; i < targetPerCollection; i++) {
            const author = pick(users);
            const title = faker.commerce.productName();
            entries.push({
                id: ulid(),
                title,
                slug: faker.helpers.slugify(title.toLowerCase()) + "-" + i,
                author: `${author.first_name} ${author.last_name}`,
                content_type: col.name,
                collection_slug: col.slug,
                status: pick(["published", "draft", "archived", "scheduled"]),
                author_id: author.id,
                created_at: pastIso(),
                updated_at: recentIso(),
            });
        }
    }

    return entries;
}

// ─────────────────────────────────────────────────────────────────────────────
// ENTRY FIELDS (per entry, from schema)
// ─────────────────────────────────────────────────────────────────────────────

function buildEntryFields(entries: Entry[], collections: Collection[]): EntryField[] {
    const fields: EntryField[] = [];
    const colBySlug = Object.fromEntries(collections.map((c) => [c.slug, c]));

    for (const entry of entries) {
        const col = colBySlug[entry.collection_slug];
        if (!col?.fields_schema) continue;
        for (const [name, schema] of Object.entries(col.fields_schema)) {
            let value: any;
            switch (schema.type) {
                case FieldTypes.TEXT:
                    value = faker.lorem.words(3);
                    break;
                case FieldTypes.RICH_TEXT:
                    value = faker.lorem.paragraphs(2);
                    break;
                case FieldTypes.NUMBER:
                    value = faker.number.float({ min: 0, max: 1000, fractionDigits: 2 });
                    break;
                case FieldTypes.BOOLEAN:
                    value = faker.datatype.boolean();
                    break;
                case FieldTypes.MEDIA:
                    value = faker.image.url();
                    break;
                case FieldTypes.RELATIONSHIP:
                    value = null;
                    break; // resolved by service
                default:
                    value = null;
            }
            fields.push({
                id: ulid(),
                entry_id: entry.id,
                name,
                type: schema.type,
                required: schema.required ?? false,
                label: schema.label,
                default_value: schema.default_value ?? null,
                value,
            });
        }
    }

    return fields;
}

// ─────────────────────────────────────────────────────────────────────────────
// FOLDERS & FILES
// ─────────────────────────────────────────────────────────────────────────────

function buildFolders(projects: Project[], users: User[]): FolderObject[] {
    const folders: FolderObject[] = [];

    for (const project of projects) {
        const rootNames = ["Images", "Documents", "Videos", "Assets"];
        const roots: FolderObject[] = rootNames.map((name) => ({
            id: ulid(),
            owner_id: pick(users).id,
            project_id: project.id,
            parent_id: null,
            name,
            file_count: 0,
            subfolder_count: 0,
            path: `/${name}`,
            depth: 0,
            created_at: pastIso(),
            updated_at: recentIso(),
        }));
        folders.push(...roots);

        // Sub-folders
        for (const root of roots) {
            const subNames = [faker.word.adjective(), faker.word.adjective()];
            for (const sub of subNames) {
                folders.push({
                    id: ulid(),
                    owner_id: root.owner_id,
                    project_id: project.id,
                    parent_id: root.id,
                    name: sub,
                    file_count: 0,
                    subfolder_count: 0,
                    path: `${root.path}/${sub}`,
                    depth: 1,
                    created_at: pastIso(),
                    updated_at: recentIso(),
                });
            }
        }
    }

    return folders;
}

function buildFileObjects(projects: Project[], users: User[], folders: FolderObject[]): FileObject[] {
    const files: FileObject[] = [];
    const mimeMap: Record<string, "image" | "video" | "document"> = {
        "image/jpeg": "image",
        "image/png": "image",
        "image/webp": "image",
        "video/mp4": "video",
        "video/webm": "video",
        "application/pdf": "document",
        "text/plain": "document",
        "application/msword": "document",
    };
    const mimes = Object.keys(mimeMap) as (keyof typeof mimeMap)[];

    for (let i = 0; i < 100; i++) {
        const project = pick(projects);
        const mime = pick(mimes);
        const category = mimeMap[mime];
        const folder = Math.random() > 0.2 ? pick(folders.filter((f) => f.project_id === project.id)) : undefined;

        files.push({
            id: ulid(),
            owner_id: pick(users).id,
            project_id: project.id,
            folder_id: folder?.id ?? null,
            name: faker.system.fileName({ extensionCount: 1 }),
            mime_type: mime,
            category,
            size: faker.number.int({ min: 1024, max: 50_000_000 }),
            provider: "s3",
            provider_key: faker.string.uuid(),
            metadata: {
                width: category === "image" ? faker.number.int({ min: 100, max: 4000 }) : undefined,
                height: category === "image" ? faker.number.int({ min: 100, max: 4000 }) : undefined,
                duration: category === "video" ? faker.number.int({ min: 10, max: 7200 }) : undefined,
                textContent: category === "document" ? faker.lorem.paragraphs(2) : undefined,
            },
            created_at: pastIso(),
            updated_at: recentIso(),
        });
    }

    return files;
}

// ─────────────────────────────────────────────────────────────────────────────
// AUDIT LOGS (50)
// ─────────────────────────────────────────────────────────────────────────────

function buildAuditLogs(users: User[], entries: Entry[]): AuditLog[] {
    const actions = [
        "user.login",
        "user.logout",
        "user.create",
        "user.update",
        "user.delete",
        "entry.create",
        "entry.update",
        "entry.delete",
        "entry.publish",
        "collection.create",
        "collection.update",
        "media.upload",
        "media.delete",
        "integration.install",
        "integration.uninstall",
        "role.assign",
        "settings.update",
    ];

    return Array.from({ length: 50 }, () => {
        const user = pick(users);
        const action = pick(actions);
        let payload: Record<string, unknown> = {};

        if (action.startsWith("entry.")) {
            payload = { entry_id: pick(entries).id, entry_title: pick(entries).title };
        } else if (action.startsWith("user.")) {
            payload = { target_user_id: pick(users).id };
        } else {
            payload = { detail: faker.lorem.sentence() };
        }

        return {
            id: ulid(),
            actor_id: user.id,
            action,
            payload,
            ip_address: faker.internet.ip(),
            created_at: pastIso(),
        };
    });
}

// ─────────────────────────────────────────────────────────────────────────────
// NOTIFICATIONS (30)
// ─────────────────────────────────────────────────────────────────────────────

function buildNotifications(users: User[], projects: Project[]): Notification[] {
    const types = ["success", "info", "warning", "error"];
    const statuses: Notification["status"][] = ["sent", "read", "unread", "dismissed"];
    const templates = [
        { title: "Export completed", body: "Your content export is ready to download." },
        { title: "New team member", body: "A user requested access to your project." },
        { title: "Entry published", body: "An entry has been successfully published." },
        { title: "Update available", body: "A new platform version is available." },
        { title: "Media uploaded", body: "Your file has been uploaded successfully." },
        { title: "Access revoked", body: "A user's access has been revoked." },
        { title: "Backup completed", body: "Your project backup is ready." },
        { title: "Integration error", body: "An integration failed to sync." },
    ];

    return Array.from({ length: 30 }, () => {
        const tmpl = pick(templates);
        const dt = pastIso();
        return {
            id: ulid(),
            title: tmpl.title,
            body: tmpl.body,
            datetime: dt,
            status: pick(statuses),
            type: pick(types),
            user_id: pick(users).id,
            project_id: pick(projects).id,
            created_at: dt,
        };
    });
}

// ─────────────────────────────────────────────────────────────────────────────
// INTEGRATIONS (15)
// ─────────────────────────────────────────────────────────────────────────────

function buildIntegrations(): Integration[] {
    return [
        { id: ulid(), name: "Stripe", category: "Payments", description: "Accept payments globally.", installed: true },
        {
            id: ulid(),
            name: "PayPal",
            category: "Payments",
            description: "PayPal checkout integration.",
            installed: false,
        },
        {
            id: ulid(),
            name: "Google Analytics",
            category: "Analytics",
            description: "Track traffic and conversions.",
            installed: true,
        },
        {
            id: ulid(),
            name: "Mixpanel",
            category: "Analytics",
            description: "Product analytics platform.",
            installed: false,
        },
        {
            id: ulid(),
            name: "Mailchimp",
            category: "Email",
            description: "Email campaigns and automation.",
            installed: true,
        },
        {
            id: ulid(),
            name: "SendGrid",
            category: "Email",
            description: "Transactional email delivery.",
            installed: false,
        },
        {
            id: ulid(),
            name: "Slack",
            category: "Communication",
            description: "Team notifications via Slack.",
            installed: true,
        },
        {
            id: ulid(),
            name: "Zapier",
            category: "Automation",
            description: "Connect with 5000+ apps.",
            installed: true,
        },
        {
            id: ulid(),
            name: "Make",
            category: "Automation",
            description: "Visual integration builder.",
            installed: false,
        },
        {
            id: ulid(),
            name: "Cloudflare R2",
            category: "Storage",
            description: "S3-compatible object storage.",
            installed: true,
        },
        { id: ulid(), name: "AWS S3", category: "Storage", description: "Amazon cloud storage.", installed: false },
        {
            id: ulid(),
            name: "Algolia",
            category: "Search",
            description: "Powerful search-as-a-service.",
            installed: false,
        },
        {
            id: ulid(),
            name: "OpenAI",
            category: "AI",
            description: "AI content generation and analysis.",
            installed: false,
        },
        {
            id: ulid(),
            name: "Resend",
            category: "Email",
            description: "Modern email SDK for developers.",
            installed: true,
        },
        {
            id: ulid(),
            name: "Sentry",
            category: "Monitoring",
            description: "Error tracking and performance.",
            installed: true,
        },
    ];
}

// ─────────────────────────────────────────────────────────────────────────────
// ACTIVITIES (40)
// ─────────────────────────────────────────────────────────────────────────────

function buildActivities(users: User[], entries: Entry[]): Activity[] {
    const actionVerbs = ["published", "edited", "created", "deleted", "approved", "archived", "restored", "scheduled"];

    return Array.from({ length: 40 }, () => {
        const user = pick(users);
        const entry = pick(entries);
        const action = pick(actionVerbs);

        // Relative time strings
        const minutesAgo = faker.number.int({ min: 1, max: 10080 }); // 1 min to 1 week
        let timeStr: string;
        if (minutesAgo < 60) timeStr = `${minutesAgo}m ago`;
        else if (minutesAgo < 1440) timeStr = `${Math.floor(minutesAgo / 60)}h ago`;
        else timeStr = `${Math.floor(minutesAgo / 1440)}d ago`;

        return {
            id: ulid(),
            user: `${user.first_name} ${user.last_name}`,
            action,
            target: entry.title,
            time: timeStr,
        };
    });
}

// ─────────────────────────────────────────────────────────────────────────────
// Main seed function
// ─────────────────────────────────────────────────────────────────────────────

export function seedDatabase(force = false): void {
    // 1. Foundation data
    const users = buildUsers();
    usersTable.seed(users, force);

    const projects = buildProjects(users);
    projectsTable.seed(projects, force);

    const roles = buildRoles();
    rolesTable.seed(roles, force);

    const permissions = buildPermissions();
    permissionsTable.seed(permissions, force);

    // 2. Content types
    const collections = buildCollections(projects);
    collectionsTable.seed(collections, force);

    // 3. Content
    const entries = buildEntries(collections, users);
    entriesTable.seed(entries, force);

    const entryFields = buildEntryFields(entries, collections);
    entryFieldsTable.seed(entryFields, force);

    // Update entry counts on collections (in-memory only — seed re-reads anyway)
    for (const col of collections) {
        col.entries = entries.filter((e) => e.collection_slug === col.slug).length;
    }

    // 4. Media
    const folders = buildFolders(projects, users);
    foldersTable.seed(folders, force);

    const files = buildFileObjects(projects, users, folders);
    fileObjectsTable.seed(files, force);

    // 5. System
    const auditLogs = buildAuditLogs(users, entries);
    auditLogsTable.seed(auditLogs, force);

    const notifications = buildNotifications(users, projects);
    notificationsTable.seed(notifications, force);

    const integrations = buildIntegrations();
    integrationsTable.seed(integrations, force);

    const activities = buildActivities(users, entries);
    activitiesTable.seed(activities, force);

    if (import.meta.env.DEV) {
        console.groupCollapsed("[CMS] Database seeded");
        console.table({
            users: users.length,
            projects: projects.length,
            roles: roles.length,
            permissions: permissions.length,
            collections: collections.length,
            entries: entries.length,
            entryFields: entryFields.length,
            folders: folders.length,
            files: files.length,
            auditLogs: auditLogs.length,
            notifications: notifications.length,
            integrations: integrations.length,
            activities: activities.length,
        });
        console.groupEnd();
    }
}
