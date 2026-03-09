import { DateTime } from "$lib/shared/datetime";
import { usersTable } from "$lib/db/tables";
import type { UserDTO } from "./dtos";
import { type AuthUser } from "./types";
import { HttpError } from "$lib/shared/httpClient";

export class AuthService {
    /** Map raw DB/DTO → clean domain model */
    static toAuthUser(dto: UserDTO): AuthUser {
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
    }

    public async signUp(email: string, password: string, username: string): Promise<AuthUser> {
        const dto = await usersTable.findOneWhere({ email } as any);
        if (dto) throw new HttpError("Account already exists with that email.", {} as any);

        const newUser = await usersTable.create({
            email,
            password,
            first_name: "",
            last_name: "",
            username,
            avatar_url: "",
            is_super_admin: false,
        } as any);

        return AuthService.toAuthUser(newUser as unknown as UserDTO);
    }

    /**
     * Simulate password login.
     * In a real app this POSTs to /auth/login.
     */
    public async loginWithPassword(email: string, password: string): Promise<AuthUser> {
        const dto = await usersTable.findOneWhere({ email } as any);
        if (!dto) throw new HttpError("No account found with that email.", {} as any);
        if ((dto as any).password !== password) throw new HttpError("Incorrect password.", {} as any);

        await usersTable.update(dto.id, { last_sign_in_at: new Date().toISOString() } as any);

        return AuthService.toAuthUser(dto as unknown as UserDTO);
    }

    /**
     * Get the currently authenticated user (simulates /auth/me).
     * We persist user id in sessionStorage as a lightweight token.
     */
    public async getCurrentUser(): Promise<AuthUser | null> {
        const id = sessionStorage.getItem("cms_auth_token");
        if (!id) return null;
        const dto = await usersTable.findOne(id);
        if (!dto) return null;
        return AuthService.toAuthUser(dto as unknown as UserDTO);
    }

    public async logout(): Promise<void> {
        await new Promise((resolve) => setTimeout(resolve, 3000));

        sessionStorage.removeItem("cms_auth_token");
    }
}

export const authService = new AuthService();
