/**
 * tables.ts
 *
 * Exports one StoreSimulator instance per entity type.
 * All features import from this file — never instantiate simulators directly.
 *
 * This is the single source of truth for the dummy database tables.
 */

import { StoreSimulator } from "./simulator";
import type {
    User,
    Project,
    Role,
    Permission,
    Collection,
    Entry,
    EntryField,
    FileObject,
    FolderObject,
    AuditLog,
    Notification,
    Integration,
    Activity,
} from "../types";

export const usersTable = new StoreSimulator<User>("users");
export const projectsTable = new StoreSimulator<Project>("projects");
export const rolesTable = new StoreSimulator<Role>("roles");
export const permissionsTable = new StoreSimulator<Permission>("permissions");
export const collectionsTable = new StoreSimulator<Collection>("collections");
export const entriesTable = new StoreSimulator<Entry>("entries");
export const entryFieldsTable = new StoreSimulator<EntryField>("entry_fields");
export const fileObjectsTable = new StoreSimulator<FileObject>("file_objects");
export const foldersTable = new StoreSimulator<FolderObject>("folders");
export const auditLogsTable = new StoreSimulator<AuditLog>("audit_logs");
export const notificationsTable = new StoreSimulator<Notification>("notifications");
export const integrationsTable = new StoreSimulator<Integration>("integrations");
export const activitiesTable = new StoreSimulator<Activity>("activities");
