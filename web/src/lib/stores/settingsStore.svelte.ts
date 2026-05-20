// Svelte 5 store for application settings
class SettingsStore {
    settings = $state<Record<string, unknown>>({
        projectName: "E-commerce Main",
        projectId: "proj_829301",
        domain: "api.ecommerce.com",
        theme: "light",
        timezone: "UTC",
    });
    loading = $state(false);
    error = $state<string | null>(null);

    async loadSettings() {
        this.loading = true;
        this.error = null;
        try {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 400));
            // Settings already initialized
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load settings";
        } finally {
            this.loading = false;
        }
    }

    async updateSettings(data: Record<string, unknown>) {
        this.loading = true;
        this.error = null;
        try {
            await new Promise((resolve) => setTimeout(resolve, 500));
            this.settings = { ...this.settings, ...data };
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to update settings";
        } finally {
            this.loading = false;
        }
    }

    getSetting<T>(key: string, defaultValue: T): T {
        return (this.settings[key] as T) ?? defaultValue;
    }
}

export const settingsStore = new SettingsStore();
