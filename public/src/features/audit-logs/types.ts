import type { AuditLog } from "$lib/types";

export type AuditLogEntry = AuditLog;

export interface AuditLogsState {
    items: AuditLogEntry[];
    loading: boolean;
    total: number;
}
