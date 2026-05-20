import { rolesTable, permissionsTable } from "$lib/db/tables";
import type { RoleDTO, PermissionDTO } from "./dtos";
import type { Role, Permission } from "./types";

export class RolesService {
    static toRole(dto: RoleDTO): Role {
        return { id: dto.id, name: dto.name, description: dto.description };
    }

    static toPermission(dto: PermissionDTO): Permission {
        return { id: dto.id, action: dto.action, resource: dto.resource };
    }

    async findAllRoles(): Promise<Role[]> {
        const result = await rolesTable.findMany({ sort: { field: "name", order: "asc" }, pageSize: 100 });
        return result.data.map(RolesService.toRole);
    }

    async findRole(id: string): Promise<Role | null> {
        const dto = await rolesTable.findOne(id);
        return dto ? RolesService.toRole(dto) : null;
    }

    async createRole(data: { name: string; description: string }): Promise<Role> {
        const dto = await rolesTable.create(data);
        return RolesService.toRole(dto);
    }

    async updateRole(id: string, data: Partial<{ name: string; description: string }>): Promise<Role> {
        const dto = await rolesTable.update(id, data);
        return RolesService.toRole(dto);
    }

    async deleteRole(id: string): Promise<void> {
        return rolesTable.delete(id);
    }

    async findAllPermissions(): Promise<Permission[]> {
        const result = await permissionsTable.findMany({ sort: { field: "resource", order: "asc" }, pageSize: 200 });
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
