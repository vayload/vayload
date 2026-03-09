import { authService } from "./services";
import type { AuthUser } from "./types";

const USER_KEY = "cms_auth_token";

class AuthStore {
    private currentUser = $state<AuthUser | null>(null);
    private _loading = $state(false);
    private _error = $state<string | null>(null);
    private permissions = $state<Record<string, boolean>>({});

    public get isAuthenticated() {
        return this.currentUser !== null;
    }

    public get isSuperAdmin() {
        return this.currentUser?.isSuperAdmin ?? false;
    }

    public async signUp(email: string, password: string, username: string) {
        this._loading = true;
        this._error = null;
        try {
            const data = await authService.signUp(email, password, username);
            if (data) {
                this.currentUser = data;
                sessionStorage.setItem(USER_KEY, data.id);
            } else {
                this._error = "Failed to sign up";
            }
        } catch (err) {
            this._error = err instanceof Error ? err.message : "Failed to sign up";
        } finally {
            this._loading = false;
        }
    }

    /**
     * Fetches the current session.
     */
    public async fetchSession() {
        this._loading = true;
        this._error = null;
        try {
            if (this.currentUser) {
                return;
            }

            const raw = await authService.getCurrentUser();
            if (raw) {
                this.currentUser = raw;
            } else {
                this.currentUser = null;
            }
        } catch (err) {
            this._error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this._loading = false;
        }
    }

    /**
     * Logs in a user with their email and password.
     *
     * @param email
     * @param password
     */
    public async loginWithPassword(email: string, password: string) {
        this._loading = true;
        this._error = null;
        try {
            const data = await authService.loginWithPassword(email, password);
            if (data) {
                this.currentUser = data;
                sessionStorage.setItem("cms_auth_token", data.id);
            } else {
                this._error = "Failed to login";
            }
        } catch (err) {
            this._error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this._loading = false;
        }
    }

    /**
     * Logs in a user with their OAuth provider.
     *
     * @param provider
     */
    public async loginWithOAuth(provider: string) {
        this._loading = true;
        this._error = null;
        try {
            const data = await authService.getCurrentUser();
            this.currentUser = data;
        } catch (err) {
            this._error = err instanceof Error ? err.message : "Failed to load user";
        } finally {
            this._loading = false;
        }
    }

    public hasPermission(permission: string) {
        return this.permissions[permission] ?? false;
    }

    /**
     * Logs out the current user.
     */
    public async logout() {
        await authService.logout();

        this.currentUser = null;
    }

    public get user() {
        return this.currentUser;
    }
}

export const authStore = new AuthStore();
