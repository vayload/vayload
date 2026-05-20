import type { AuditLog } from "$lib/types";
import type { PaginatedResult } from "$lib/db/simulator";
import { httpClient } from "$lib/api/client";

export class AuditLogsService {
    async findMany(
        params: {
            page?: number;
            pageSize?: number;
            search?: string;
        } = {},
    ): Promise<PaginatedResult<AuditLog>> {
        return httpClient.get<PaginatedResult<AuditLog>>("/audit-logs", {
            page: params.page,
            pageSize: params.pageSize ?? 50,
            search: params.search,
        });
    }
}

export const auditLogsService = new AuditLogsService();
