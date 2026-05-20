export interface Role {
    id: string;
    name: string;
    description: string;
}

export interface Permission {
    id: string;
    action: string;
    resource: string;
}

export interface RolesState {
    roles: Role[];
    permissions: Permission[];
    loading: boolean;
    total: number;
}
