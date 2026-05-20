import { settingsService } from "./services";
import type { SettingsState, ProjectSettings, AppNotification } from "./types";
import type { UpdateProjectDTO } from "./dtos";

class SettingsStore {
    #state = $state<SettingsState>({
        project: null,
        notifications: [],
        unreadCount: 0,
        loading: false,
    });

    get project(): ProjectSettings | null {
        return this.#state.project;
    }
    get notifications(): AppNotification[] {
        return this.#state.notifications;
    }
    get unreadCount(): number {
        return this.#state.unreadCount;
    }
    get loading(): boolean {
        return this.#state.loading;
    }

    async loadProject(id: string) {
        this.#state.loading = true;
        try {
            this.#state.project = await settingsService.getProject(id);
        } finally {
            this.#state.loading = false;
        }
    }

    async loadNotifications(userId: string) {
        this.#state.notifications = await settingsService.getNotifications(userId);
        this.#state.unreadCount = this.#state.notifications.filter((n) => n.status === "unread").length;
    }

    async updateProject(data: UpdateProjectDTO) {
        if (!this.#state.project) return;
        this.#state.project = await settingsService.updateProject(this.#state.project.id, data);
    }

    async markRead(id: string) {
        await settingsService.markNotificationRead(id);
        this.#state.notifications = this.#state.notifications.map((n) =>
            n.id === id ? { ...n, status: "read" as const } : n,
        );
        this.#state.unreadCount = this.#state.notifications.filter((n) => n.status === "unread").length;
    }

    async dismiss(id: string) {
        await settingsService.dismissNotification(id);
        this.#state.notifications = this.#state.notifications.filter((n) => n.id !== id);
        this.#state.unreadCount = this.#state.notifications.filter((n) => n.status === "unread").length;
    }
}

export const settingsStore = new SettingsStore();
