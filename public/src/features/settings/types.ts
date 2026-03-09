import type { DateTime } from "$lib/shared/datetime";

export interface ProjectSettings {
    id: string;
    name: string;
    slug: string;
    ownerId: string;
    settings: Record<string, unknown>;
    locale: string;
    createdAt: DateTime;
    updatedAt: DateTime;
}

export interface AppNotification {
    id: string;
    title: string;
    body: string;
    datetime: DateTime;
    status: "sent" | "read" | "unread" | "dismissed" | "failed";
    type: string;
    userId: string;
    projectId: string;
    createdAt: DateTime;
}

export interface SettingsState {
    project: ProjectSettings | null;
    notifications: AppNotification[];
    unreadCount: number;
    loading: boolean;
}
