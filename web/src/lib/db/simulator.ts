import { withDelay } from "./utils";

// ─── Types ──────────────────────────────────────────────────────────────────

export type WhereCondition<T> = {
    [K in keyof T]?: T[K] | ((val: T[K]) => boolean);
};

export type SortOrder = "asc" | "desc";

export interface QueryParams<T> {
    where?: WhereCondition<T>;
    sort?: { field: keyof T; order: SortOrder };
    page?: number;
    pageSize?: number;
}

export interface PaginatedResult<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
    hasNextPage: boolean;
    hasPrevPage: boolean;
}

// ─── Simulator ───────────────────────────────────────────────────────────────

export class StoreSimulator<T extends { id: string }> {
    private readonly key: string;

    constructor(tableName: string) {
        this.key = `cms::${tableName}`;
    }

    // ── Internal ──────────────────────────────────────────────────────────

    private read(): T[] {
        if (typeof window === "undefined") return [];
        try {
            return JSON.parse(localStorage.getItem(this.key) ?? "[]") as T[];
        } catch {
            return [];
        }
    }

    private write(items: T[]): void {
        if (typeof window === "undefined") return;
        localStorage.setItem(this.key, JSON.stringify(items));
    }

    /** Used by seed.ts to populate the table once. */
    public seed(items: T[], force = false): void {
        if (force || this.read().length === 0) {
            this.write(items);
        }
    }

    public clear(): void {
        if (typeof window !== "undefined") localStorage.removeItem(this.key);
    }

    // ── API (all async to simulate network) ───────────────────────────────

    async findMany(params: QueryParams<T> = {}): Promise<PaginatedResult<T>> {
        return withDelay(
            () => {
                let items = this.read();

                // Filter
                if (params.where) {
                    items = items.filter((item) =>
                        Object.entries(params.where!).every(([k, condition]) => {
                            const val = (item as any)[k];
                            return typeof condition === "function" ? (condition as Function)(val) : val === condition;
                        }),
                    );
                }

                // Sort
                if (params.sort) {
                    const { field, order } = params.sort;
                    items = [...items].sort((a, b) => {
                        const av = a[field];
                        const bv = b[field];
                        const cmp = av < bv ? -1 : av > bv ? 1 : 0;
                        return order === "asc" ? cmp : -cmp;
                    });
                }

                const total = items.length;
                const page = Math.max(1, params.page ?? 1);
                const pageSize = Math.max(1, params.pageSize ?? 20);
                const totalPages = Math.ceil(total / pageSize);

                // Paginate
                const start = (page - 1) * pageSize;
                const data = items.slice(start, start + pageSize);

                return {
                    data,
                    total,
                    page,
                    pageSize,
                    totalPages,
                    hasNextPage: page < totalPages,
                    hasPrevPage: page > 1,
                };
            },
            200,
            500,
        );
    }

    async findOne(id: string): Promise<T | null> {
        return withDelay(
            () => {
                return this.read().find((i) => i.id === id) ?? null;
            },
            100,
            300,
        );
    }

    async findOneWhere(where: WhereCondition<T>): Promise<T | null> {
        return withDelay(
            () => {
                return (
                    this.read().find((item) =>
                        Object.entries(where).every(([k, condition]) => {
                            const val = (item as any)[k];
                            return typeof condition === "function" ? (condition as Function)(val) : val === condition;
                        }),
                    ) ?? null
                );
            },
            100,
            300,
        );
    }

    async create(data: Omit<T, "id"> & { id?: string }): Promise<T> {
        return withDelay(
            () => {
                const items = this.read();
                const newItem = { ...data, id: data.id ?? crypto.randomUUID() } as T;
                this.write([...items, newItem]);
                return newItem;
            },
            200,
            400,
        );
    }

    async update(id: string, data: Partial<Omit<T, "id">>): Promise<T> {
        return withDelay(
            () => {
                const items = this.read();
                const idx = items.findIndex((i) => i.id === id);
                if (idx === -1) throw new Error(`[StoreSimulator] Record not found: ${id}`);
                const updated = { ...items[idx], ...data } as T;
                items[idx] = updated;
                this.write(items);
                return updated;
            },
            200,
            400,
        );
    }

    async delete(id: string): Promise<void> {
        return withDelay(
            () => {
                const items = this.read().filter((i) => i.id !== id);
                this.write(items);
            },
            150,
            350,
        );
    }

    async count(where?: WhereCondition<T>): Promise<number> {
        return withDelay(
            () => {
                let items = this.read();
                if (where) {
                    items = items.filter((item) =>
                        Object.entries(where).every(([k, condition]) => {
                            const val = (item as any)[k];
                            return typeof condition === "function" ? (condition as Function)(val) : val === condition;
                        }),
                    );
                }
                return items.length;
            },
            50,
            150,
        );
    }
}
