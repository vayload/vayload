import { ulid } from "./utils";

export const roles = [
    { id: ulid(), name: "owner", description: "Owner del proyecto" },
    { id: ulid(), name: "editor", description: "Editor" },
];

export const permissions = [
    { id: ulid(), action: "read", resource: "projects" },
    { id: ulid(), action: "update", resource: "projects" },
    { id: ulid(), action: "publish", resource: "entries" },
];

export const rolePermissions = [
    { role_id: ulid(), permission_id: ulid() },
    { role_id: ulid(), permission_id: ulid() },
    { role_id: ulid(), permission_id: ulid() },
    { role_id: ulid(), permission_id: ulid() },
    { role_id: ulid(), permission_id: ulid() },
];
