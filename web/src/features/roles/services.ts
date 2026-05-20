import type { RoleDTO, PermissionDTO } from "./dtos";
import type { Role, Permission } from "./types";
import { httpClient } from "$lib/api/client";

export class RolesService {
    static toRole(dto: RoleDTO): Role {
        return { id: dto.id, name: dto.name, description: dto.description };
    }

    static toPermission(dto: PermissionDTO): Permission {
        return { id: dto.id, action: dto.action, resource: dto.resource };
    }

    async findAllRoles(): Promise<Role[]> {
        const result = await httpClient.get<{ data: RoleDTO[] }>("/roles");
        return result.data.map(RolesService.toRole);
    }

    async findRole(id: string): Promise<Role | null> {
        const dto = await httpClient.get<RoleDTO | null>(`/roles/${id}`);
        return dto ? RolesService.toRole(dto) : null;
    }

    async createRole(data: { name: string; description: string }): Promise<Role> {
        const dto = await httpClient.post<RoleDTO>("/roles", data);
        return RolesService.toRole(dto);
    }

    async updateRole(id: string, data: Partial<{ name: string; description: string }>): Promise<Role> {
        const dto = await httpClient.patch<RoleDTO>(`/roles/${id}`, data);
        return RolesService.toRole(dto);
    }

    async deleteRole(id: string): Promise<void> {
        return httpClient.delete<void>(`/roles/${id}`);
    }

    async findAllPermissions(): Promise<Permission[]> {
        const result = await httpClient.get<{ data: PermissionDTO[] }>("/permissions");
        return result.data.map(RolesService.toPermission);
    }

    /** Permissions grouped by resource for easy UI rendering */
    async getPermissionsMatrix(): Promise<Record<string, Permission[]>> {
        const permissions = await this.findAllPermissions();
        return permissions.reduce<Record<string, Permission[]>>((acc, p) => {
            (acc[p.resource] ??= []).push(p);
            return acc;
        }, {});
    }
}

export const rolesService = new RolesService();
