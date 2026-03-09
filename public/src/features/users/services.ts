import { DateTime } from "$lib/shared/datetime";
import { usersTable } from "$lib/db/tables";
import type { UserDTO, CreateUserDTO, UpdateUserDTO } from "./dtos";
import type { User } from "./types";
import type { PaginatedResult } from "$lib/db/simulator";

export class UsersService {
    static toUser(dto: UserDTO): User {
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

    async findMany(
        params: {
            page?: number;
            pageSize?: number;
            search?: string;
        } = {},
    ): Promise<PaginatedResult<User>> {
        const where: any = {};
        if (params.search) {
            const q = params.search.toLowerCase();
            // Simulate a full-text search by checking email/username
            where.email = (v: string) => v.toLowerCase().includes(q) || false;
        }

        const result = await usersTable.findMany({
            where: params.search
                ? {
                      email: (v: string) => v.toLowerCase().includes(params.search!.toLowerCase()),
                  }
                : undefined,
            sort: { field: "created_at", order: "desc" },
            page: params.page,
            pageSize: params.pageSize ?? 20,
        });

        return {
            ...result,
            data: result.data.map((d) => UsersService.toUser(d as unknown as UserDTO)),
        };
    }

    async findOne(id: string): Promise<User | null> {
        const dto = await usersTable.findOne(id);
        if (!dto) return null;
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async create(data: CreateUserDTO): Promise<User> {
        const now = new Date().toISOString();
        const dto = await usersTable.create({
            ...data,
            avatar_url: data.avatar_url ?? null,
            is_super_admin: data.is_super_admin ?? false,
            last_sign_in_at: now,
            created_at: now,
            updated_at: now,
        } as any);
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async update(id: string, data: UpdateUserDTO): Promise<User> {
        const dto = await usersTable.update(id, {
            ...data,
            updated_at: new Date().toISOString(),
        } as any);
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async delete(id: string): Promise<void> {
        return usersTable.delete(id);
    }
}

export const usersService = new UsersService();
