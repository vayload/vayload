export interface AuditLogDTO {
    id: string;
    actor_id: string;
    action: string;
    payload: Record<string, unknown>;
    ip_address: string;
    created_at: string;
}
