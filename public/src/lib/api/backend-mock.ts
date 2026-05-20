import {
    activitiesTable,
    auditLogsTable,
    collectionsTable,
    entriesTable,
    entryFieldsTable,
    fileObjectsTable,
    foldersTable,
    integrationsTable,
    notificationsTable,
    permissionsTable,
    projectsTable,
    rolesTable,
    usersTable,
} from "$lib/db/tables";
import type {
    Activity,
    AuditLog,
    Collection,
    Entry,
    EntryField,
    FileObject,
    FolderObject,
    Integration,
    Notification,
    Permission,
    Project,
    User,
} from "$lib/types";
import { HttpError } from "$lib/shared/httpClient";
import { getCookie, setCookie, deleteCookie } from "$lib/shared/cookies";

type Method = "GET" | "POST" | "PATCH" | "DELETE" | "PUT";

type Context = {
    method: Method;
    path: string;
    params: Record<string, string>;
    query: Record<string, any>;
    body?: any;
    userId?: string | null;
};

type Handler = (ctx: Context) => Promise<any>;

type Route = {
    method: Method;
    pattern: string;
    handler: Handler;
};

const routes: Route[] = [];

function route(method: Method, pattern: string, handler: Handler) {
    routes.push({ method, pattern, handler });
}

function match(pattern: string, path: string) {
    const p1 = pattern.split("/").filter(Boolean);
    const p2 = path.split("/").filter(Boolean);

    if (p1.length !== p2.length) return null;

    const params: Record<string, string> = {};

    for (let i = 0; i < p1.length; i++) {
        const a = p1[i];
        const b = p2[i];

        if (a.startsWith(":")) {
            params[a.slice(1)] = b;
        } else if (a !== b) {
            return null;
        }
    }

    return params;
}

function getUserFromCookies() {
    const token = getCookie("access_token");
    if (!token) return null;

    try {
        const payload = JSON.parse(atob(token.split(".")[1] || token));
        return payload.user_id;
    } catch {
        // Fallback for non-jwt format if user manually set it
        return token;
    }
}

export interface ApiRequest {
    method: Method;
    path: string;
    query?: Record<string, string | number | boolean | undefined>;
    body?: unknown;
}

export async function handleRequest<T>(req: ApiRequest): Promise<T> {
    const url = req.path.split("?")[0];

    for (const r of routes) {
        if (r.method !== req.method) continue;

        const params = match(r.pattern, url);

        if (!params) continue;

        const ctx: Context = {
            method: req.method,
            path: url,
            params,
            query: req.query ?? {},
            body: req.body,
            userId: getUserFromCookies(),
        };

        return r.handler(ctx);
    }

    throw new Error(`[API Simulator] No route for ${req.method} ${req.path}`);
}

const slugify = (str: string): string => {
    return str
        .replace(/^\s+|\s+$/g, "")
        .toLowerCase()
        .replace(/[^a-z0-9 -]/g, "")
        .replace(/\s+|-+/g, "-");
};

// --- AUTH ---

route("POST", "/auth/signup", async (ctx) => {
    return (await usersTable.create(ctx.body)) as any;
});

route("POST", "/auth/login", async (ctx) => {
    const { email } = ctx.body;
    const user = await usersTable.findOneWhere({ email } as any);
    if (!user) throw new HttpError("Invalid credentials", {} as any);

    const accessToken = btoa(JSON.stringify({ user_id: user.id, type: "access" }));
    const refreshToken = btoa(JSON.stringify({ user_id: user.id, type: "refresh" }));

    setCookie("access_token", accessToken, 7);
    setCookie("refresh_token", refreshToken, 7);
    setCookie("cms_auth_token", user.id, 7); // compatibility

    return user;
});

// --- USERS ---

route("GET", "/users", async (ctx) => {
    const { page, pageSize, search } = ctx.query;
    return await usersTable.findMany({
        where: search ? { email: (v: string) => v.toLowerCase().includes(String(search).toLowerCase()) } : undefined,
        sort: { field: "created_at", order: "desc" },
        page: page ? Number(page) : undefined,
        pageSize: pageSize ? Number(pageSize) : 20,
    });
});

route("GET", "/users/:id", async (ctx) => {
    return await usersTable.findOne(ctx.params.id);
});

route("POST", "/users", async (ctx) => {
    if (!ctx.body) throw new HttpError("Missing body", {} as any);
    const data = ctx.body as User;
    if (await usersTable.findOneWhere({ email: data.email } as any))
        throw new HttpError("User already exists", {} as any);
    if (await usersTable.findOneWhere({ username: data.username } as any))
        throw new HttpError("User already exists", {} as any);

    const now = new Date().toISOString();
    return await usersTable.create({
        ...data,
        avatar_url: data.avatar_url ?? null,
        is_super_admin: data.is_super_admin ?? false,
        last_sign_in_at: now,
        created_at: now,
        updated_at: now,
    } as any);
});

route("PATCH", "/users/:id", async (ctx) => {
    return await usersTable.update(ctx.params.id, {
        ...ctx.body,
        updated_at: new Date().toISOString(),
    });
});

route("DELETE", "/users/:id", async (ctx) => {
    return await usersTable.delete(ctx.params.id);
});

// --- ROLES & PERMISSIONS ---

route("GET", "/roles", async () => {
    return await rolesTable.findMany({
        sort: { field: "name", order: "asc" },
        pageSize: 100,
    });
});

route("GET", "/roles/:id", async (ctx) => {
    return await rolesTable.findOne(ctx.params.id);
});

route("POST", "/roles", async (ctx) => {
    return await rolesTable.create(ctx.body);
});

route("PATCH", "/roles/:id", async (ctx) => {
    return await rolesTable.update(ctx.params.id, ctx.body);
});

route("DELETE", "/roles/:id", async (ctx) => {
    return await rolesTable.delete(ctx.params.id);
});

route("GET", "/permissions", async () => {
    return await permissionsTable.findMany({
        sort: { field: "resource", order: "asc" },
        pageSize: 200,
    });
});

// --- PROJECTS ---

route("GET", "/projects", async (ctx) => {
    const { page, pageSize, search } = ctx.query;
    return await projectsTable.findMany({
        where: search ? { name: (v: string) => v.toLowerCase().includes(String(search).toLowerCase()) } : undefined,
        sort: { field: "created_at", order: "desc" },
        page: page ? Number(page) : undefined,
        pageSize: pageSize ? Number(pageSize) : 50,
    });
});

route("GET", "/projects/:id", async (ctx) => {
    return await projectsTable.findOne(ctx.params.id);
});

route("POST", "/projects", async (ctx) => {
    const data = ctx.body as Partial<Project>;
    const now = new Date().toISOString();
    return await projectsTable.create({
        ...data,
        slug: data.slug ?? slugify(data.name ?? ""),
        created_at: now,
        updated_at: now,
    } as any);
});

route("PATCH", "/projects/:id", async (ctx) => {
    return await projectsTable.update(ctx.params.id, {
        ...ctx.body,
        updated_at: new Date().toISOString(),
    });
});

route("DELETE", "/projects/:id", async (ctx) => {
    return await projectsTable.delete(ctx.params.id);
});

// --- NOTIFICATIONS & PLUGINS ---

route("GET", "/notifications", async (ctx) => {
    const userId = String(ctx.query.userId ?? "");
    return await notificationsTable.findMany({
        where: { user_id: userId } as any,
        sort: { field: "created_at", order: "desc" },
        pageSize: 50,
    });
});

route("PATCH", "/notifications/:id", async (ctx) => {
    return await notificationsTable.update(ctx.params.id, ctx.body);
});

route("GET", "/plugins", async () => {
    return await integrationsTable.findMany({
        sort: { field: "name", order: "asc" },
        pageSize: 200,
    });
});

route("PATCH", "/plugins/:id", async (ctx) => {
    return await integrationsTable.update(ctx.params.id, ctx.body);
});

// --- DASHBOARD & AUDIT LOGS ---

route("GET", "/dashboard/activities", async () => {
    return await activitiesTable.findMany({ sort: { field: "id", order: "desc" }, pageSize: 50 });
});

// --- CONTENT TYPES & ENTRIES ---

route("GET", "/content-types", async (ctx) => {
    const { page, pageSize, search } = ctx.query;
    return await collectionsTable.findMany({
        where: search ? { name: (v: string) => v.toLowerCase().includes(String(search).toLowerCase()) } : undefined,
        sort: { field: "created_at", order: "desc" },
        page: page ? Number(page) : undefined,
        pageSize: pageSize ? Number(pageSize) : 20,
    });
});

route("GET", "/content-types/:slug", async (ctx) => {
    return await collectionsTable.findOneWhere({ slug: ctx.params.slug } as any);
});

route("POST", "/content-types", async (ctx) => {
    const data = ctx.body as Partial<Collection>;
    const now = new Date().toISOString();
    return await collectionsTable.create({
        project_id: data.project_id!,
        name: data.name!,
        slug: data.slug!,
        fields_schema: data.fields_schema ?? {},
        settings: data.settings ?? {},
        entries: data.entries ?? 0,
        single: data.single ?? false,
        created_at: now,
    });
});

route("PATCH", "/content-types/:id", async (ctx) => {
    return await collectionsTable.update(ctx.params.id, ctx.body);
});

route("DELETE", "/content-types/:id", async (ctx) => {
    return await collectionsTable.delete(ctx.params.id);
});

route("GET", "/entries", async (ctx) => {
    const { status, search, page, pageSize } = ctx.query;
    const where: any = {};
    if (status && status !== "all") where.status = status;
    const result = (await entriesTable.findMany({
        where: Object.keys(where).length ? where : undefined,
        sort: { field: "updated_at", order: "desc" },
        page: page ? Number(page) : undefined,
        pageSize: pageSize ? Number(pageSize) : 50,
    })) as any;

    let data = result.data as Entry[];
    if (search) {
        const q = String(search).toLowerCase();
        data = data.filter(
            (entry) =>
                entry.title.toLowerCase().includes(q) ||
                entry.author.toLowerCase().includes(q) ||
                entry.content_type.toLowerCase().includes(q),
        );
    }

    return {
        ...result,
        data,
        total: search ? data.length : result.total,
    };
});

route("GET", "/entries/counts", async () => {
    const [all, published, draft, archived, scheduled] = await Promise.all([
        entriesTable.count(),
        entriesTable.count({ status: "published" } as any),
        entriesTable.count({ status: "draft" } as any),
        entriesTable.count({ status: "archived" } as any),
        entriesTable.count({ status: "scheduled" } as any),
    ]);
    return { all, published, draft, archived, scheduled };
});

route("GET", "/entries/:slug/fields", async (ctx) => {
    const entry = (await entriesTable.findOneWhere({ slug: ctx.params.slug } as any)) as Entry | null;
    if (!entry) return null;

    const fieldsResult = (await entryFieldsTable.findMany({
        where: { entry_id: entry.id } as any,
        pageSize: 500,
    })) as any;

    let fields = fieldsResult.data as EntryField[];
    if (fields.length === 0) {
        const collection = (await collectionsTable.findOneWhere({
            slug: entry.collection_slug,
        } as any)) as Collection | null;
        if (collection?.fields_schema) {
            const created: EntryField[] = [];
            for (const [name, schema] of Object.entries(collection.fields_schema)) {
                const field = (await entryFieldsTable.create({
                    entry_id: entry.id,
                    name,
                    type: schema.type,
                    required: schema.required ?? false,
                    label: schema.label,
                    default_value: schema.default_value ?? null,
                    value: schema.default_value ?? null,
                } as any)) as EntryField;
                created.push(field);
            }
            fields = created;
        }
    }
    return { ...entry, fields };
});

route("GET", "/entries/:slug", async (ctx) => {
    return await entriesTable.findOneWhere({ slug: ctx.params.slug } as any);
});

route("POST", "/entries", async (ctx) => {
    const data = ctx.body as Partial<Entry>;
    const now = new Date().toISOString();
    const slug = data.slug ?? slugify(data.title ?? "");
    return await entriesTable.create({
        title: data.title!,
        slug: slug!,
        author: data.author!,
        content_type: data.content_type!,
        collection_slug: data.collection_slug!,
        status: data.status ?? "draft",
        author_id: data.author_id!,
        created_at: now,
        updated_at: now,
    });
});

route("PATCH", "/entries/:id", async (ctx) => {
    return await entriesTable.update(ctx.params.id, {
        ...ctx.body,
        updated_at: new Date().toISOString(),
    });
});

route("DELETE", "/entries/:id", async (ctx) => {
    return await entriesTable.delete(ctx.params.id);
});

// --- MEDIA ---

route("GET", "/media/folders", async (ctx) => {
    const parentId = (ctx.query.parent_id as string) ?? null;
    const projectId = String(ctx.query.project_id ?? "");
    return await foldersTable.findMany({
        where: { parent_id: parentId, project_id: projectId, owner_id: ctx.userId } as any,
        sort: { field: "name", order: "asc" },
        pageSize: 500,
    });
});

route("GET", "/media/files", async (ctx) => {
    const folderId = (ctx.query.folder_id as string) ?? null;
    const projectId = String(ctx.query.project_id ?? "");
    return await fileObjectsTable.findMany({
        where: { folder_id: folderId, project_id: projectId, owner_id: ctx.userId } as any,
        sort: { field: "created_at", order: "desc" },
        pageSize: 500,
    });
});

route("GET", "/media/filter", async (ctx) => {
    const category = String(ctx.query.category ?? "");
    return await fileObjectsTable.findMany({
        where: { category } as any,
        sort: { field: "created_at", order: "desc" },
        pageSize: 500,
    });
});

route("POST", "/media/upload", async (ctx) => {
    const payload = ctx.body as { file: File } & Partial<FileObject>;
    const file = payload.file;
    const category = file.type.startsWith("image") ? "image" : file.type.startsWith("video") ? "video" : "document";
    const now = new Date().toISOString();
    return await fileObjectsTable.create({
        owner_id: payload.owner_id,
        project_id: payload.project_id,
        folder_id: payload.folder_id ?? null,
        name: file.name,
        mime_type: file.type || "application/octet-stream",
        category,
        size: file.size,
        provider: "local",
        provider_key: file.name,
        metadata: {
            preview: category === "image" ? URL.createObjectURL(file) : undefined,
        },
        created_at: now,
        updated_at: now,
    } as any);
});

route("POST", "/media/folders", async (ctx) => {
    const data = ctx.body as Partial<FolderObject>;
    const parent = data.parent_id ? await foldersTable.findOne(data.parent_id) : null;
    const depth = parent ? parent.depth + 1 : 0;
    const path = parent ? `${parent.path}/${data.name}` : `/${data.name}`;
    const created = (await foldersTable.create({
        ...data,
        depth,
        path,
        file_count: data.file_count ?? 0,
        subfolder_count: data.subfolder_count ?? 0,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
    } as any)) as FolderObject;
    if (parent) {
        await foldersTable.update(parent.id, { subfolder_count: parent.subfolder_count + 1 } as any);
    }
    return created;
});

route("PATCH", "/media/rename", async (ctx) => {
    const data = ctx.body as { id: string; name: string; type: "file" | "folder" };
    if (data.type === "file") {
        return await fileObjectsTable.update(data.id, { name: data.name } as any);
    }
    return await foldersTable.update(data.id, { name: data.name } as any);
});

route("PATCH", "/media/move", async (ctx) => {
    const data = ctx.body as { id: string; parent_id: string | null; type: "file" | "folder" };
    if (data.type === "file") {
        return await fileObjectsTable.update(data.id, { folder_id: data.parent_id } as any);
    }
    const parent = data.parent_id ? await foldersTable.findOne(data.parent_id) : null;
    const depth = parent ? parent.depth + 1 : 0;
    const path = parent ? `${parent.path}/` : "/";
    return await foldersTable.update(data.id, {
        parent_id: data.parent_id,
        depth,
        path,
    } as any);
});

route("DELETE", "/media/delete", async (ctx) => {
    const id = String(ctx.query.id ?? "");
    const type = String(ctx.query.type ?? "file");
    if (type === "file") return await fileObjectsTable.delete(id);

    const folders = (await foldersTable.findMany({ pageSize: 2000 })) as any;
    const files = (await fileObjectsTable.findMany({ pageSize: 5000 })) as any;

    const folderIdsToDelete = new Set<string>();
    folderIdsToDelete.add(id);
    let changed = true;
    while (changed) {
        changed = false;
        for (const folder of folders.data as FolderObject[]) {
            if (folder.parent_id && folderIdsToDelete.has(folder.parent_id) && !folderIdsToDelete.has(folder.id)) {
                folderIdsToDelete.add(folder.id);
                changed = true;
            }
        }
    }
    for (const folderId of folderIdsToDelete) {
        await foldersTable.delete(folderId);
    }
    for (const file of files.data as FileObject[]) {
        if (file.folder_id && folderIdsToDelete.has(file.folder_id)) {
            await fileObjectsTable.delete(file.id);
        }
    }
    return undefined;
});

route("GET", "/media/find", async (ctx) => {
    const q = String(ctx.query.query ?? "").toLowerCase();
    const folders = (await foldersTable.findMany({ pageSize: 2000 })) as any;
    const files = (await fileObjectsTable.findMany({ pageSize: 5000 })) as any;
    return {
        folders: (folders.data as FolderObject[]).filter((f) => f.name.toLowerCase().includes(q)),
        files: (files.data as FileObject[]).filter((f) => f.name.toLowerCase().includes(q)),
    };
});

route("GET", "/media/grep", async (ctx) => {
    const q = String(ctx.query.query ?? "").toLowerCase();
    const files = (await fileObjectsTable.findMany({ pageSize: 5000 })) as any;
    return (files.data as FileObject[]).filter((f) =>
        (f.metadata?.textContent as string | undefined)?.toLowerCase().includes(q),
    );
});

route("GET", "/media/preview/:id", async (ctx) => {
    const id = ctx.params.id;
    const file = (await fileObjectsTable.findOne(id)) as FileObject | null;
    const url = (file?.metadata as any)?.preview;
    return url ? { id, file_id: id, url } : null;
});
