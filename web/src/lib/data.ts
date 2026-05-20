import { faker } from "@faker-js/faker";

// Crockford's Base32 (sin I, L, O, U)
const ENCODING = "0123456789ABCDEFGHJKMNPQRSTVWXYZ";

function encodeTime(time: bigint, length: number) {
    let out = "";
    for (let i = length - 1; i >= 0; i--) {
        const mod = time % 32n;
        out = ENCODING[Number(mod)] + out;
        time = time / 32n;
    }
    return out;
}

function encodeRandom(length: number) {
    let out = "";
    const randomBytes = crypto.getRandomValues(new Uint8Array(length));

    for (let i = 0; i < length; i++) {
        out += ENCODING[randomBytes[i] % 32];
    }
    return out;
}

export function ulid() {
    const time = BigInt(Date.now());
    const timePart = encodeTime(time, 10); // 48 bits → 10 chars
    const randomPart = encodeRandom(16); // 80 bits → 16 chars

    return timePart + randomPart;
}

export function delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

export enum FieldTypes {
    TEXT = "text",
    RICH_TEXT = "rich_text",
    NUMBER = "number",
    DATE = "date",
    BOOLEAN = "boolean",
    RELATIONSHIP = "relationship",
    MEDIA = "media",
    TONES = "tones",
    LOCATION = "location",
    JSON = "json",
}

export interface Role {
    id: string;
    name: string;
    description: string;
}

export interface Permission {
    id: string;
    action: string;
    resource: string;
}

export interface Notification {
    id: string;
    title: string;
    body: string;
    datetime: string;
    status: "sent" | "read" | "unread" | "dismissed" | "failed";
    type: string;
    user_id: string;
    project_id: string;
    created_at: string;
}

export enum FileCategory {
    IMAGE = "image",
    VIDEO = "video",
    DOCUMENT = "document",
}

export interface FolderObject {
    id: string;
    owner_id: string;
    project_id: string;
    parent_id: string | null;
    name: string;
    file_count: number;
    subfolder_count: number;
    path: string;
    depth: number;
    created_at: string | Date;
    updated_at: string | Date;
}

export interface FileObject {
    id: string;
    owner_id: string;
    project_id: string;
    folder_id: string | null;
    name: string;
    mime_type: string;
    category: "image" | "video" | "document";
    size: number;
    provider: "local" | "s3" | "r2" | "gcs";
    provider_key: string;
    metadata: Record<string, any>;
    created_at: string | Date;
    updated_at: string | Date;
}

export interface FileUpload {
    file: File;
    owner_id: string;
    project_id: string;
    folder_id: string | null;
}

export interface FolderCreate {
    name: string;
    owner_id: string;
    project_id: string;
    parent_id: string | null;
}

export interface Integration {
    id: string;
    name: string;
    category: string;
    description: string;
    installed: boolean;
    icon?: string;
}

export interface AuditLog {
    id: string;
    actor_id: string;
    action: string;
    payload: Record<string, unknown>;
    ip_address: string;
    created_at: string;
}

export interface Activity {
    id: string;
    user: string;
    action: string;
    target: string;
    time: string;
}

export interface FieldSchema {
    type: FieldTypes;
    required?: boolean;
    label: string;
    name?: string;
    default_value?: any;
    relation_to?: string;
    config?: Record<string, unknown>;
}

export type CollectionSchema = Record<string, FieldSchema>;

export interface User {
    id: string;
    first_name: string;
    last_name: string;
    username: string;
    email: string;
    avatar_url: string | null;
    is_super_admin: boolean;
    last_sign_in_at: string;
    created_at: string;
    updated_at: string;
    password?: string; // For dummy purposes only; in real systems, use hashed passwords
}

export interface Project {
    id: string;
    name: string;
    slug: string;
    owner_id: string;
    settings: Record<string, unknown>;
    locale: string;
    created_at: string;
    updated_at: string;
}

export interface ProjectInput {
    name: string;
    settings: Record<string, unknown>;
    locale: string;
}

export interface Collection {
    id: string;
    project_id: string;
    name: string;
    slug: string;
    fields_schema: CollectionSchema;
    settings: Record<string, unknown>;
    entries: number;
    single: boolean;
    created_at: string;
}

export interface Entry {
    id: string;
    title: string;
    slug: string;
    author: string;
    content_type: string;
    collection_slug: string;
    status: "published" | "draft" | "archived" | "scheduled";
    author_id: string;
    created_at: string;
    updated_at: string;
}

export interface EntryField {
    id: string;
    entry_id: string;
    name: string;
    type: FieldTypes;
    required: boolean;
    label: string;
    default_value: any;
    value: any;
}

export interface EntryWithFields extends Entry {
    fields: EntryField[];
}

const USER_ADMIN_ID = ulid();
const USER_EDITOR_ID = ulid();
const USER_AUTHOR_ID = ulid();

const PROJECT_ECOM_ID = ulid();
const PROJECT_BLOG_ID = ulid();

export const ROLES: Role[] = [
    { id: ulid(), name: "Administrator", description: "Full access to all system features." },
    { id: ulid(), name: "Editor", description: "Can manage content but not system settings." },
    { id: ulid(), name: "Author", description: "Can create and edit their own content." },
    { id: ulid(), name: "Public API", description: "Read-only access for public endpoints." },
    { id: ulid(), name: "Viewer", description: "Read-only access to content." },
    { id: ulid(), name: "Manager", description: "Manages teams and projects." },
];

export const NOTIFICATIONS: Notification[] = [
    {
        id: ulid(),
        title: "Export completed",
        body: "Your content export is ready to download.",
        datetime: "2024-01-31T10:56:00Z",
        status: "unread",
        type: "success",
        user_id: USER_ADMIN_ID,
        project_id: PROJECT_ECOM_ID,
        created_at: "2024-01-31T10:56:00Z",
    },
    {
        id: ulid(),
        title: "New team member",
        body: 'Carlos requested access to "E-commerce Main".',
        datetime: "2024-01-31T09:00:00Z",
        status: "unread",
        type: "info",
        user_id: USER_ADMIN_ID,
        project_id: PROJECT_ECOM_ID,
        created_at: "2024-01-31T09:00:00Z",
    },
    {
        id: ulid(),
        title: "Update available",
        body: "A new version of the platform is available.",
        datetime: "2024-02-15T14:30:00Z",
        status: "read",
        type: "info",
        user_id: USER_EDITOR_ID,
        project_id: PROJECT_BLOG_ID,
        created_at: "2024-02-15T14:30:00Z",
    },
    {
        id: ulid(),
        title: "Entry published",
        body: 'Your entry "New Product Launch" has been published.',
        datetime: "2024-03-01T11:45:00Z",
        status: "sent",
        type: "success",
        user_id: USER_AUTHOR_ID,
        project_id: PROJECT_ECOM_ID,
        created_at: "2024-03-01T11:45:00Z",
    },
    {
        id: ulid(),
        title: "Access revoked",
        body: "Access to project Blog has been revoked for a user.",
        datetime: "2024-03-10T16:20:00Z",
        status: "unread",
        type: "warning",
        user_id: USER_ADMIN_ID,
        project_id: PROJECT_BLOG_ID,
        created_at: "2024-03-10T16:20:00Z",
    },
];

// export const FILE_OBJECTS: FileObject[] = [
//     {
//         id: ulid(),
//         owner_id: USER_ADMIN_ID,
//         project_id: PROJECT_ECOM_ID,
//         name: "IMG_001.jpg",
//         mime_type: "image/jpeg",
//         category: "image",
//         size: 1258291,
//         folder: "Marketing",
//         created_at: "2024-01-20T10:00:00Z",
//         updated_at: "2024-01-20T10:00:00Z",
//     },
//     {
//         id: ulid(),
//         owner_id: USER_EDITOR_ID,
//         project_id: PROJECT_BLOG_ID,
//         name: "document.pdf",
//         mime_type: "application/pdf",
//         category: "document",
//         size: 2048000,
//         folder: "Reports",
//         created_at: "2024-02-05T12:30:00Z",
//         updated_at: "2024-02-05T12:30:00Z",
//     },
//     {
//         id: ulid(),
//         owner_id: USER_AUTHOR_ID,
//         project_id: PROJECT_ECOM_ID,
//         name: "video_demo.mp4",
//         mime_type: "video/mp4",
//         category: "video",
//         size: 52428800,
//         folder: "Products",
//         created_at: "2024-03-15T09:45:00Z",
//         updated_at: "2024-03-15T09:45:00Z",
//     },
//     {
//         id: ulid(),
//         owner_id: USER_ADMIN_ID,
//         project_id: PROJECT_BLOG_ID,
//         name: "banner.png",
//         mime_type: "image/png",
//         category: "image",
//         size: 1048576,
//         folder: "Assets",
//         created_at: "2024-04-01T14:00:00Z",
//         updated_at: "2024-04-01T14:00:00Z",
//     },
// ];

export const INTEGRATIONS: Integration[] = [
    {
        id: ulid(),
        name: "Stripe",
        category: "Payment",
        description: "Accept payments globally with the world's best processing platform.",
        installed: true,
    },
    {
        id: ulid(),
        name: "Google Analytics",
        category: "Analytics",
        description: "Get insights into your users and traffic.",
        installed: false,
    },
    {
        id: ulid(),
        name: "Mailchimp",
        category: "Email Marketing",
        description: "Send newsletters and manage email campaigns.",
        installed: true,
    },
    {
        id: ulid(),
        name: "Slack",
        category: "Communication",
        description: "Integrate notifications with your team chat.",
        installed: false,
    },
    {
        id: ulid(),
        name: "Zapier",
        category: "Automation",
        description: "Connect apps and automate workflows.",
        installed: true,
    },
];

export const ACTIVITIES: Activity[] = [
    {
        id: ulid(),
        user: "Ana Admin",
        action: "published",
        target: "Summer Sale",
        time: "10m ago",
    },
    {
        id: ulid(),
        user: "Editor Eve",
        action: "edited",
        target: "Blog Post: Tech Trends",
        time: "2h ago",
    },
    {
        id: ulid(),
        user: "Author Alex",
        action: "created",
        target: "New Product Entry",
        time: "1d ago",
    },
    {
        id: ulid(),
        user: "Manager Mike",
        action: "approved",
        target: "User Access Request",
        time: "3d ago",
    },
    {
        id: ulid(),
        user: "Admin Ana",
        action: "deleted",
        target: "Old Archive",
        time: "1w ago",
    },
];

export const AUDIT_LOGS: AuditLog[] = [
    {
        id: ulid(),
        actor_id: USER_ADMIN_ID,
        action: "user.login",
        payload: { method: "password" },
        ip_address: "192.168.1.100",
        created_at: "2024-01-31T10:58:00Z",
    },
    {
        id: ulid(),
        actor_id: USER_EDITOR_ID,
        action: "entry.update",
        payload: { entry_id: "some-entry-id", changes: "title" },
        ip_address: "192.168.1.101",
        created_at: "2024-02-15T14:32:00Z",
    },
    {
        id: ulid(),
        actor_id: USER_AUTHOR_ID,
        action: "entry.create",
        payload: { collection: "products" },
        ip_address: "192.168.1.102",
        created_at: "2024-03-01T11:47:00Z",
    },
    {
        id: ulid(),
        actor_id: USER_ADMIN_ID,
        action: "user.access.revoked",
        payload: { user_id: "revoked-user-id" },
        ip_address: "192.168.1.100",
        created_at: "2024-03-10T16:22:00Z",
    },
    {
        id: ulid(),
        actor_id: USER_EDITOR_ID,
        action: "integration.install",
        payload: { integration: "Mailchimp" },
        ip_address: "192.168.1.101",
        created_at: "2024-04-05T09:15:00Z",
    },
];

export async function fetchNotifications(): Promise<Notification[]> {
    await delay(300);
    return [...NOTIFICATIONS];
}

// export async function fetchFileObjects(): Promise<FileObject[]> {
//     await delay(500);
//     return [...FILE_OBJECTS];
// }

export async function fetchIntegrations(): Promise<Integration[]> {
    await delay(400);
    return [...INTEGRATIONS];
}

export async function fetchActivities(): Promise<Activity[]> {
    await delay(300);
    return [...ACTIVITIES];
}

export async function fetchAuditLogs(): Promise<AuditLog[]> {
    await delay(400);
    return [...AUDIT_LOGS];
}

export async function fetchRoles(): Promise<Role[]> {
    await delay(200);
    return [...ROLES];
}

const adminUser: User = {
    id: USER_ADMIN_ID,
    first_name: "Admin",
    last_name: "User",
    username: "admin",
    email: "admin@example.com",
    avatar_url: null,
    is_super_admin: true,
    last_sign_in_at: faker.date.recent().toISOString(),
    created_at: faker.date.past().toISOString(),
    updated_at: faker.date.recent().toISOString(),
    password: "adminpass",
};

const editorUser: User = {
    id: USER_EDITOR_ID,
    first_name: "Editor",
    last_name: "Eve",
    username: "editor",
    email: "editor@example.com",
    avatar_url: faker.image.avatar(),
    is_super_admin: false,
    last_sign_in_at: faker.date.recent().toISOString(),
    created_at: faker.date.past().toISOString(),
    updated_at: faker.date.recent().toISOString(),
    password: "editorpass",
};

const authorUser: User = {
    id: USER_AUTHOR_ID,
    first_name: "Author",
    last_name: "Alex",
    username: "author",
    email: "author@example.com",
    avatar_url: faker.image.avatar(),
    is_super_admin: false,
    last_sign_in_at: faker.date.recent().toISOString(),
    created_at: faker.date.past().toISOString(),
    updated_at: faker.date.recent().toISOString(),
    password: "authorpass",
};

export const USERS: User[] = [
    adminUser,
    editorUser,
    authorUser,
    ...Array.from({ length: 7 }).map(() => ({
        id: ulid(),
        first_name: faker.person.firstName(),
        last_name: faker.person.lastName(),
        username: faker.internet.username(),
        email: faker.internet.email(),
        avatar_url: faker.image.avatar(),
        is_super_admin: faker.datatype.boolean({ probability: 0.1 }),
        last_sign_in_at: faker.date.recent().toISOString(),
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
        password: faker.internet.password(),
    })),
];

export const PROJECTS: Project[] = [
    {
        id: PROJECT_ECOM_ID,
        name: "E-commerce Main",
        slug: "e-commerce-main",
        owner_id: USER_ADMIN_ID,
        settings: { theme: "light", timezone: "UTC" },
        locale: "en-US",
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: PROJECT_BLOG_ID,
        name: "Company Blog",
        slug: "company-blog",
        owner_id: USER_EDITOR_ID,
        settings: { theme: "dark", timezone: "PST" },
        locale: "en-US",
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    ...Array.from({ length: 3 }).map(() => {
        const name = faker.company.name();
        return {
            id: ulid(),
            name,
            slug: faker.helpers.slugify(name.toLowerCase()),
            owner_id: faker.helpers.arrayElement(USERS).id,
            settings: {
                theme: faker.helpers.arrayElement(["light", "dark"]),
                timezone: "UTC",
            },
            locale: faker.helpers.arrayElement(["en-US", "es-ES", "fr-FR"]),
            created_at: faker.date.past().toISOString(),
            updated_at: faker.date.recent().toISOString(),
        };
    }),
];

const categorySchema: CollectionSchema = {
    title: {
        type: FieldTypes.TEXT,
        required: true,
        label: "Title",
        default_value: "",
    },
    description: {
        type: FieldTypes.RICH_TEXT,
        required: false,
        label: "Description",
    },
};

const productSchema: CollectionSchema = {
    title: {
        type: FieldTypes.TEXT,
        required: true,
        label: "Title",
        default_value: "",
    },
    description: {
        type: FieldTypes.RICH_TEXT,
        required: false,
        label: "Description",
    },
    price: {
        type: FieldTypes.NUMBER,
        required: false,
        label: "Price",
    },
    published: {
        type: FieldTypes.BOOLEAN,
        required: false,
        label: "Published",
        default_value: false,
    },
    cover: {
        type: FieldTypes.MEDIA,
        required: false,
        label: "Cover",
    },
    category: {
        type: FieldTypes.RELATIONSHIP,
        required: false,
        label: "Category",
        relation_to: "categories",
    },
};

const blogPostSchema: CollectionSchema = {
    title: {
        type: FieldTypes.TEXT,
        required: true,
        label: "Title",
        default_value: "",
    },
    content: {
        type: FieldTypes.RICH_TEXT,
        required: true,
        label: "Content",
    },
    published: {
        type: FieldTypes.BOOLEAN,
        required: false,
        label: "Published",
        default_value: false,
    },
    author: {
        type: FieldTypes.RELATIONSHIP,
        required: false,
        label: "Author",
        relation_to: "users", // Assuming users can be related, but for demo; adjust as needed
    },
};

export const COLLECTIONS: Collection[] = [
    {
        id: ulid(),
        project_id: PROJECT_ECOM_ID,
        name: "Categories",
        slug: "categories",
        fields_schema: categorySchema,
        settings: { versioning: true },
        entries: 0, // Will be updated later
        single: false,
        created_at: faker.date.past().toISOString(),
    },
    {
        id: ulid(),
        project_id: PROJECT_ECOM_ID,
        name: "Products",
        slug: "products",
        fields_schema: productSchema,
        settings: { versioning: true },
        entries: 0,
        single: false,
        created_at: faker.date.past().toISOString(),
    },
    {
        id: ulid(),
        project_id: PROJECT_BLOG_ID,
        name: "Posts",
        slug: "posts",
        fields_schema: blogPostSchema,
        settings: { versioning: false },
        entries: 0,
        single: false,
        created_at: faker.date.past().toISOString(),
    },
    // Add more generic collections for other projects
    ...Array.from({ length: 5 }).map((_, i) => ({
        id: ulid(),
        project_id: faker.helpers.arrayElement(PROJECTS).id,
        name: faker.commerce.department(),
        slug: `collection-${i + 1}`,
        fields_schema: productSchema, // Reuse for variety
        settings: { versioning: faker.datatype.boolean() },
        entries: faker.number.int({ min: 5, max: 20 }),
        single: false,
        created_at: faker.date.past().toISOString(),
    })),
];

const categoryEntries: Entry[] = [
    {
        id: ulid(),
        title: "Electronics",
        slug: "electronics",
        author: adminUser.first_name,
        content_type: "Categories",
        collection_slug: "categories",
        status: "published",
        author_id: adminUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Clothing",
        slug: "clothing",
        author: editorUser.first_name,
        content_type: "Categories",
        collection_slug: "categories",
        status: "published",
        author_id: editorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Books",
        slug: "books",
        author: authorUser.first_name,
        content_type: "Categories",
        collection_slug: "categories",
        status: "published",
        author_id: authorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Home Appliances",
        slug: "home-appliances",
        author: adminUser.first_name,
        content_type: "Categories",
        collection_slug: "categories",
        status: "draft",
        author_id: adminUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
];

const productEntries: Entry[] = [
    {
        id: ulid(),
        title: "Smartphone XYZ",
        slug: "smartphone-xyz",
        author: authorUser.first_name,
        content_type: "Products",
        collection_slug: "products",
        status: "published",
        author_id: authorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "T-Shirt Basic",
        slug: "t-shirt-basic",
        author: editorUser.first_name,
        content_type: "Products",
        collection_slug: "products",
        status: "published",
        author_id: editorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Novel Adventure",
        slug: "novel-adventure",
        author: adminUser.first_name,
        content_type: "Products",
        collection_slug: "products",
        status: "scheduled",
        author_id: adminUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Blender Pro",
        slug: "blender-pro",
        author: authorUser.first_name,
        content_type: "Products",
        collection_slug: "products",
        status: "archived",
        author_id: authorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Laptop Ultra",
        slug: "laptop-ultra",
        author: editorUser.first_name,
        content_type: "Products",
        collection_slug: "products",
        status: "published",
        author_id: editorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
];

const blogPostEntries: Entry[] = [
    {
        id: ulid(),
        title: "Tech Trends 2024",
        slug: "tech-trends-2024",
        author: authorUser.first_name,
        content_type: "Posts",
        collection_slug: "posts",
        status: "published",
        author_id: authorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Fashion Tips",
        slug: "fashion-tips",
        author: editorUser.first_name,
        content_type: "Posts",
        collection_slug: "posts",
        status: "draft",
        author_id: editorUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
    {
        id: ulid(),
        title: "Book Reviews",
        slug: "book-reviews",
        author: adminUser.first_name,
        content_type: "Posts",
        collection_slug: "posts",
        status: "published",
        author_id: adminUser.id,
        created_at: faker.date.past().toISOString(),
        updated_at: faker.date.recent().toISOString(),
    },
];

export const ENTRIES: Entry[] = [
    ...categoryEntries,
    ...productEntries,
    ...blogPostEntries,
    // Add more for other collections if needed
];

// Update entries count in collections
COLLECTIONS.forEach((collection) => {
    collection.entries = ENTRIES.filter((e) => e.collection_slug === collection.slug).length;
});

function generateRandomValue(type: FieldTypes, schema: FieldSchema) {
    switch (type) {
        case FieldTypes.TEXT:
            return faker.lorem.words(3);
        case FieldTypes.RICH_TEXT:
            return faker.lorem.paragraphs(2);
        case FieldTypes.NUMBER:
            return faker.number.int({ min: 0, max: 1000 });
        case FieldTypes.BOOLEAN:
            return faker.datatype.boolean();
        case FieldTypes.MEDIA:
            return faker.image.url();
        case FieldTypes.RELATIONSHIP:
            if (schema.relation_to) {
                const relatedEntries = ENTRIES.filter((e) => e.collection_slug === schema.relation_to);
                if (relatedEntries.length > 0) {
                    return faker.helpers.arrayElement(relatedEntries).id;
                }
            }
            return null;
        default:
            return null;
    }
}

export async function loginWithPassword(email: string, password: string): Promise<User | null> {
    await delay(500);
    const user = USERS.find((u) => u.email === email && u.password === password);
    if (user) {
        user.last_sign_in_at = new Date().toISOString();
        // Log the login in audit logs (simulate)
        AUDIT_LOGS.push({
            id: ulid(),
            actor_id: user.id,
            action: "user.login",
            payload: { method: "password", email },
            ip_address: faker.internet.ip(),
            created_at: new Date().toISOString(),
        });
        return { ...user, password: undefined }; // Don't return password
    }
    return null;
}

export async function fetchUsers(): Promise<User[]> {
    await delay(300);
    return USERS.map((u) => ({ ...u, password: undefined })); // Strip passwords
}

export async function fetchProjects(): Promise<Project[]> {
    await delay(300);
    return [...PROJECTS];
}

export async function fetchCollections(): Promise<Collection[]> {
    await delay(300);
    return [...COLLECTIONS];
}

export async function fetchCollectionBySlug(slug: string): Promise<Collection | undefined> {
    await delay(400);
    return COLLECTIONS.find((c) => c.slug === slug) || COLLECTIONS[0];
}

export async function fetchEntries(filter?: string): Promise<Entry[]> {
    await delay(300);
    if (filter) {
        return ENTRIES.filter((e) => e.collection_slug === filter);
    }

    return [...ENTRIES];
}

export async function fetchEntriesByStatus(status: string): Promise<Entry[]> {
    await delay(300);
    if (status === "all") {
        return [...ENTRIES];
    }

    return ENTRIES.filter((e) => e.status === status);
}

export async function fetchEntryWithFields(slug: string): Promise<EntryWithFields | undefined> {
    await delay(300);

    const entry = ENTRIES.find((e) => e.slug === slug);
    if (!entry) return undefined;

    const collection = COLLECTIONS.find((c) => c.slug === entry.collection_slug);
    if (!collection) return undefined;

    const fields: EntryField[] = Object.entries(collection.fields_schema).map(([name, schema]) => ({
        id: ulid(),
        entry_id: entry.id,
        name,
        type: schema.type,
        required: schema.required ?? false,
        label: schema.label,
        default_value: schema.default_value ?? null,
        value: schema.default_value !== undefined ? schema.default_value : generateRandomValue(schema.type, schema),
    }));

    return {
        ...entry,
        fields,
    };
}

export const getUser = () => {
    return JSON.parse(localStorage.getItem("auth") || "{}") as User;
};
