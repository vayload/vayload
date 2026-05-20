import type { DateTime } from "$lib/shared/datetime";

export interface User {
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

export interface UsersState {
    items: User[];
    selected: User | null;
    loading: boolean;
    total: number;
    page: number;
    pageSize: number;
    search: string;
}
