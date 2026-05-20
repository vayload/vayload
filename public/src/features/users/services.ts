import { DateTime } from "$lib/shared/datetime";
import type { UserDTO, CreateUserDTO, UpdateUserDTO } from "./dtos";
import type { User } from "./types";
import type { PaginatedResult } from "$lib/db/simulator";
import { httpClient } from "$lib/api/client";

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
        const result = await httpClient.get<PaginatedResult<UserDTO>>("/users", {
            page: params.page,
            pageSize: params.pageSize ?? 20,
            search: params.search,
        });

        return {
            ...result,
            data: result.data.map((d) => UsersService.toUser(d as unknown as UserDTO)),
        };
    }

    async findOne(id: string): Promise<User | null> {
        const dto = await httpClient.get<UserDTO | null>(`/users/${id}`);
        if (!dto) return null;
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async create(data: CreateUserDTO): Promise<User> {
        const dto = await httpClient.post<UserDTO>("/users", data);
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async update(id: string, data: UpdateUserDTO): Promise<User> {
        const dto = await httpClient.patch<UserDTO>(`/users/${id}`, data);
        return UsersService.toUser(dto as unknown as UserDTO);
    }

    async delete(id: string): Promise<void> {
        return httpClient.delete<void>(`/users/${id}`);
    }
}

export const usersService = new UsersService();
