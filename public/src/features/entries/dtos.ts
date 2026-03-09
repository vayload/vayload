export interface EntryDTO {
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

export interface EntryFieldDTO {
    id: string;
    entry_id: string;
    name: string;
    type: string;
    required: boolean;
    label: string;
    default_value: any;
    value: any;
}

export interface CreateEntryDTO {
    title: string;
    slug?: string;
    collection_slug: string;
    content_type: string;
    status?: "published" | "draft" | "archived" | "scheduled";
    author_id: string;
    author: string;
}
