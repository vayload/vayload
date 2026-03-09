/**
 * Auth domain models — used inside the Svelte app.
 * camelCase + proper Date objects.
 */
import type { DateTime } from "$lib/shared/datetime";

export interface AuthUser {
    id: string;
    firstName: string;
    lastName: string;
    fullName: string;
    username: string;
    email: string;
    avatarUrl: string | null;
    isSuperAdmin: boolean;
    lastSignInAt: DateTime;
    createdAt: DateTime;
    updatedAt: DateTime;
}

export interface AuthState {
    user: AuthUser | null;
    isAuthenticated: boolean;
    token: string | null;
}

export class AuthError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "AuthError";
    }
}
