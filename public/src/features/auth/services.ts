import { DateTime } from "$lib/shared/datetime";
import type { UserDTO } from "./dtos";
import { type AuthUser } from "./types";
import { HttpError } from "$lib/shared/httpClient";
import { httpClient } from "$lib/api/client";
import { getCookie, deleteCookie } from "$lib/shared/cookies";

const toAuthUser = (dto: UserDTO): AuthUser => {
    return {
        id: dto.id,
        firstName: dto.first_name,
        lastName: dto.last_name,
        fullName: `${dto.first_name} ${dto.last_name}`,
        username: dto.username,
        email: dto.email,
        avatarUrl: dto.avatar_url,
        isSuperAdmin: dto.is_super_admin,
        lastSignInAt: new DateTime(dto.last_sign_in_at),
        createdAt: new DateTime(dto.created_at),
        updatedAt: new DateTime(dto.updated_at),
    };
};

export class AuthService {
    public async signUp(email: string, password: string, username: string): Promise<AuthUser> {
        const newUser = await httpClient.post<UserDTO>("/auth/signup", {
            email,
            password,
            username,
        });

        return toAuthUser(newUser as unknown as UserDTO);
    }

    /**
     * Simulate password login.
     * In a real app this POSTs to /auth/login.
     */
    public async loginWithPassword(email: string, password: string): Promise<AuthUser> {
        const dto = await httpClient.post<UserDTO | null>("/auth/login", { email, password });
        if (!dto) throw new HttpError("No account found with that email.", {} as any);
        if ((dto as any).password !== password) throw new HttpError("Incorrect password.", {} as any);

        return toAuthUser(dto as unknown as UserDTO);
    }

    /**
     * Get the currently authenticated user (simulates /auth/me).
     * We persist user id in cookies as a lightweight token.
     */
    public async getCurrentUser(): Promise<AuthUser | null> {
        const id = getCookie("cms_auth_token");
        if (!id) return null;
        const dto = await httpClient.get<UserDTO | null>("/auth/me", { id });
        if (!dto) return null;
        return toAuthUser(dto as unknown as UserDTO);
    }

    public async logout(): Promise<void> {
        await new Promise((resolve) => setTimeout(resolve, 3000));

        deleteCookie("cms_auth_token");
    }
}

export const authService = new AuthService();
