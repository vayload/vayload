import { fetchUsers, loginWithPassword, type User } from "$lib/data";

class AuthStore {
    private currentUser = $state<User | null>(null);
    private loading = $state(false);
    private error = $state<string | null>(null);
    private permissions = $state<Record<string, boolean>>({});

    public get isAuthenticated() {
        return this.currentUser !== null;
    }

    public get isSuperAdmin() {
        return this.currentUser?.is_super_admin ?? false;
    }

    /**
     * Fetches the current session.
     */
    public async fetchSession() {
        this.loading = true;
        this.error = null;
        try {
            if (this.currentUser) {
                return;
            }

            await new Promise((resolve) => setTimeout(resolve, 3000));
            const raw = localStorage.getItem("auth");
            if (raw) {
                this.currentUser = JSON.parse(raw);
            } else {
                this.currentUser = null;
            }
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this.loading = false;
        }
    }

    /**
     * Logs in a user with their email and password.
     *
     * @param email
     * @param password
     */
    public async loginWithPassword(email: string, password: string) {
        this.loading = true;
        this.error = null;
        try {
            const data = await loginWithPassword(email, password);
            if (data) {
                this.currentUser = data;
                localStorage.setItem("auth", JSON.stringify(this.currentUser));
            } else {
                this.error = "Failed to login";
            }
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this.loading = false;
        }
    }

    /**
     * Logs in a user with their OAuth provider.
     *
     * @param provider
     */
    public async loginWithOAuth(provider: string) {
        this.loading = true;
        this.error = null;
        try {
            const data = await fetchUsers();
            this.currentUser = data[0];
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this.loading = false;
        }
    }

    public hasPermission(permission: string) {
        return this.permissions[permission] ?? false;
    }

    /**
     * Logs out the current user.
     */
    public logout() {
        this.currentUser = null;
    }

    public get user() {
        return this.currentUser;
    }
}

export const authStore = new AuthStore();
