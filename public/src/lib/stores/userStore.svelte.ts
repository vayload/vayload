import { fetchUsers, type User } from "$lib/data";

class UserStore {
    currentUser = $state<User | null>(null);
    users = $state<User[]>([]);
    loading = $state(false);
    error = $state<string | null>(null);

    get isAuthenticated() {
        return this.currentUser !== null;
    }

    get isSuperAdmin() {
        return this.currentUser?.is_super_admin ?? false;
    }

    async loadCurrentUser() {
        this.loading = true;
        this.error = null;
        try {
            const data = await fetchUsers();
            // Simulate logged in user (Ana Admin)
            this.currentUser = data[0];
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this.loading = false;
        }
    }

    async loadUsers() {
        this.loading = true;
        this.error = null;
        try {
            const data = await fetchUsers();
            this.users = data;
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load users";
        } finally {
            this.loading = false;
        }
    }

    async updateProfile(data: Partial<User>) {
        if (!this.currentUser) return;

        this.loading = true;
        try {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 500));
            this.currentUser = { ...this.currentUser, ...data };
        } finally {
            this.loading = false;
        }
    }

    /**
     * Fetches the current session.
     */
    public async fetchSession() {
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
            const data = await fetchUsers();
            this.currentUser = data[0];
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

    /**
     * Logs out the current user.
     */
    public logout() {
        this.currentUser = null;
    }
}

export const userStore = new UserStore();
