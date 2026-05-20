import { entriesService } from "./services";
import type { EntriesState, ContentEntry } from "./types";
import type { CreateEntryDTO, UpdateEntryDTO } from "./dtos";

class EntriesStore {
    #state = $state<EntriesState>({
        items: [],
        filtered: [],
        loading: false,
        error: null,
        statusFilter: "all",
        search: "",
        counts: {},
    });

    get items(): ContentEntry[] {
        return this.#state.items;
    }
    get filteredEntries(): ContentEntry[] {
        return this.#state.filtered;
    }
    get loading(): boolean {
        return this.#state.loading;
    }
    get error(): string | null {
        return this.#state.error;
    }
    get totalCount(): number {
        return this.#state.items.length;
    }
    get countByStatus(): Record<string, number> {
        return {
            all: this.#state.counts.all ?? this.#state.items.length,
            published: this.#state.counts.published ?? 0,
            draft: this.#state.counts.draft ?? 0,
            archived: this.#state.counts.archived ?? 0,
            scheduled: this.#state.counts.scheduled ?? 0,
        };
    }

    async loadEntries(filter: string = "all") {
        this.#state.loading = true;
        this.#state.error = null;
        this.#state.statusFilter = filter;
        try {
            const [result, counts] = await Promise.all([
                entriesService.findMany({ status: filter as any, search: this.#state.search, pageSize: 200 }),
                entriesService.countByStatus(),
            ]);

            this.#state.items = result.data;
            this.#state.counts = counts;
            this.applyFilters();
        } catch (err) {
            this.#state.error = err instanceof Error ? err.message : "Failed to load entries";
        } finally {
            this.#state.loading = false;
        }
    }

    searchEntries(query: string) {
        this.#state.search = query;
        this.applyFilters();
    }

    private applyFilters() {
        let result = [...this.#state.items];

        if (this.#state.search) {
            const q = this.#state.search.toLowerCase();
            result = result.filter(
                (e) =>
                    e.title.toLowerCase().includes(q) ||
                    e.content_type.toLowerCase().includes(q) ||
                    e.author.toLowerCase().includes(q),
            );
        }

        this.#state.filtered = result;
    }

    async createEntry(data: CreateEntryDTO) {
        this.#state.loading = true;
        try {
            const entry = await entriesService.create(data);
            this.#state.items = [entry, ...this.#state.items];
            this.applyFilters();
            this.#state.counts = await entriesService.countByStatus();
            return entry;
        } finally {
            this.#state.loading = false;
        }
    }

    async updateEntry(id: string, data: UpdateEntryDTO) {
        this.#state.loading = true;
        try {
            const entry = await entriesService.update(id, data);
            this.#state.items = this.#state.items.map((e) => (e.id === id ? entry : e));
            this.applyFilters();
            return entry;
        } finally {
            this.#state.loading = false;
        }
    }

    async deleteEntry(id: string) {
        this.#state.loading = true;
        try {
            await entriesService.delete(id);
            this.#state.items = this.#state.items.filter((e) => e.id !== id);
            this.applyFilters();
            this.#state.counts = await entriesService.countByStatus();
        } finally {
            this.#state.loading = false;
        }
    }
}

export const entriesStore = new EntriesStore();
