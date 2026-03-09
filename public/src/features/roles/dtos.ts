export interface RoleDTO {
    id: string;
    name: string;
    description: string;
}

export interface PermissionDTO {
    id: string;
    action: string;
    resource: string;
}
