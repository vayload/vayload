import { fetchNotifications, type Notification } from "$lib/data";

// Svelte 5 store for notifications
class NotificationStore {
    notifications = $state<Notification[]>([]);
    loading = $state(false);
    error = $state<string | null>(null);

    get unreadCount() {
        return this.notifications.filter((n) => n.status === "unread").length;
    }

    get hasUnread() {
        return this.unreadCount > 0;
    }

    async loadNotifications() {
        this.loading = true;
        this.error = null;
        try {
            const data = await fetchNotifications();
            this.notifications = data;
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load notifications";
        } finally {
            this.loading = false;
        }
    }

    async markAsRead(id: string) {
        const notification = this.notifications.find((n) => n.id === id);
        if (notification) {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 200));
            notification.status = "read";
            this.notifications = [...this.notifications];
        }
    }

    async markAllAsRead() {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 300));
        this.notifications = this.notifications.map((n) => ({
            ...n,
            status: n.status === "unread" ? "read" : n.status,
        }));
    }

    async dismiss(id: string) {
        const notification = this.notifications.find((n) => n.id === id);
        if (notification) {
            await new Promise((resolve) => setTimeout(resolve, 200));
            notification.status = "dismissed";
            this.notifications = [...this.notifications];
        }
    }
}

export const notificationStore = new NotificationStore();
