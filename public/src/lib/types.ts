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

export function delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
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
    password?: string; // For dummy purposes only;
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
