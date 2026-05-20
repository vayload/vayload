import { auditLogsService } from "./services";
import type { AuditLogsState, AuditLogEntry } from "./types";

class AuditLogsStore {
    #state = $state<AuditLogsState>({
        items: [],
        loading: false,
        total: 0,
    });

    get items(): AuditLogEntry[] {
        return this.#state.items;
    }
    get loading(): boolean {
        return this.#state.loading;
    }
    get total(): number {
        return this.#state.total;
    }

    async fetch(search = "") {
        this.#state.loading = true;
        try {
            const result = await auditLogsService.findMany({ search, pageSize: 100 });
            this.#state.items = result.data;
            this.#state.total = result.total;
        } finally {
            this.#state.loading = false;
        }
    }
}

export const auditLogsStore = new AuditLogsStore();
