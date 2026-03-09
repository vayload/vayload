export interface ProjectDTO {
    id: string;
    name: string;
    slug: string;
    owner_id: string;
    settings: Record<string, unknown>;
    locale: string;
    created_at: string;
    updated_at: string;
}

export interface UpdateProjectDTO {
    name?: string;
    slug?: string;
    settings?: Record<string, unknown>;
    locale?: string;
}

export interface NotificationDTO {
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
