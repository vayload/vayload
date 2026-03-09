import type { CollectionSchema } from "$lib/types";

export interface CollectionDTO {
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

export interface CreateCollectionDTO {
    project_id: string;
    name: string;
    slug: string;
    fields_schema?: CollectionSchema;
    settings?: Record<string, unknown>;
    single?: boolean;
}
