import { rolesService } from "./services";
import type { RolesState, Role, Permission } from "./types";

class RolesStore {
    #state = $state<RolesState>({
        roles: [],
        permissions: [],
        loading: false,
        total: 0,
    });

    get roles(): Role[] {
        return this.#state.roles;
    }
    get permissions(): Permission[] {
        return this.#state.permissions;
    }
    get loading(): boolean {
        return this.#state.loading;
    }
    get total(): number {
        return this.#state.total;
    }

    async fetchRoles() {
        this.#state.loading = true;
        try {
            this.#state.roles = await rolesService.findAllRoles();
            this.#state.total = this.#state.roles.length;
        } finally {
            this.#state.loading = false;
        }
    }

    async fetchPermissions() {
        this.#state.permissions = await rolesService.findAllPermissions();
    }

    async fetchAll() {
        await Promise.all([this.fetchRoles(), this.fetchPermissions()]);
    }

    async createRole(data: { name: string; description: string }) {
        const role = await rolesService.createRole(data);
        this.#state.roles = [...this.#state.roles, role];
        this.#state.total += 1;
        return role;
    }

    async updateRole(id: string, data: Partial<{ name: string; description: string }>) {
        const role = await rolesService.updateRole(id, data);
        this.#state.roles = this.#state.roles.map((r) => (r.id === id ? role : r));
        return role;
    }

    async deleteRole(id: string) {
        await rolesService.deleteRole(id);
        this.#state.roles = this.#state.roles.filter((r) => r.id !== id);
        this.#state.total -= 1;
    }

    get permissionsMatrix(): Record<string, Permission[]> {
        return this.#state.permissions.reduce<Record<string, Permission[]>>((acc, p) => {
            (acc[p.resource] ??= []).push(p);
            return acc;
        }, {});
    }
}

export const rolesStore = new RolesStore();
