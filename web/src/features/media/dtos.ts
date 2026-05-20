export interface FileObjectDTO {
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
    created_at: string;
    updated_at: string;
}

export interface FolderDTO {
    id: string;
    owner_id: string;
    project_id: string;
    parent_id: string | null;
    name: string;
    file_count: number;
    subfolder_count: number;
    path: string;
    depth: number;
    created_at: string;
    updated_at: string;
}
