import { fetchEntries, fetchEntriesByStatus, ulid, type Entry } from "$lib/data";

// Svelte 5 store for content entries
class EntriesStore {
    entries = $state<Entry[]>([]);
    filteredEntries = $state<Entry[]>([]);
    loading = $state(false);
    error = $state<string | null>(null);
    currentFilter = $state<string>("all");
    searchQuery = $state<string>("");

    private counts = $state<Record<string, number>>({});

    public get countByStatus() {
        return {
            all: this.entries.length,
            published: this.counts.published || 0,
            draft: this.counts.draft || 0,
            archived: this.counts.archived || 0,
            scheduled: this.counts.scheduled || 0,
        };
    }

    get totalCount() {
        return this.entries.length;
    }

    get publishedCount() {
        return this.entries.filter((e) => e.status === "published").length;
    }

    get draftCount() {
        return this.entries.filter((e) => e.status === "draft").length;
    }

    get archivedCount() {
        return this.entries.filter((e) => e.status === "archived").length;
    }

    get scheduledCount() {
        return this.entries.filter((e) => e.status === "scheduled").length;
    }

    async loadEntries(filter: string = "all") {
        this.loading = true;
        this.error = null;
        this.currentFilter = filter;
        try {
            // const data = await fetchEntriesByStatus(filter);
            const [data, entries] = await Promise.all([fetchEntriesByStatus(filter), fetchEntries()]);

            this.entries = data;
            this.applyFilters();

            const counts = entries.reduce(
                (acc, entry) => {
                    acc[entry.status] = (acc[entry.status] || 0) + 1;
                    return acc;
                },
                {} as Record<string, number>,
            );

            this.counts = counts;
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load entries";
        } finally {
            this.loading = false;
        }
    }

    searchEntries(query: string) {
        this.searchQuery = query;
        this.applyFilters();
    }

    private applyFilters() {
        let result = [...this.entries];

        // Apply search filter
        if (this.searchQuery) {
            const query = this.searchQuery.toLowerCase();
            result = result.filter(
                (e) =>
                    e.title.toLowerCase().includes(query) ||
                    e.content_type.toLowerCase().includes(query) ||
                    e.author.toLowerCase().includes(query),
            );
        }

        this.filteredEntries = result;
    }

    async createEntry(data: Omit<Entry, "id" | "created_at">) {
        this.loading = true;
        try {
            await new Promise((resolve) => setTimeout(resolve, 500));
            const newEntry: Entry = {
                ...data,
                id: ulid(),
                slug: data.title
                    .toLowerCase()
                    .replace(/\s+/g, "-")
                    .replace(/[^a-z0-9-]/g, ""),
                created_at: new Date().toISOString(),
            };
            this.entries = [...this.entries, newEntry];
            this.applyFilters();
            return newEntry;
        } finally {
            this.loading = false;
        }
    }

    async updateEntry(id: string, data: Partial<Entry>) {
        this.loading = true;
        try {
            await new Promise((resolve) => setTimeout(resolve, 500));
            const index = this.entries.findIndex((e) => e.id === id);
            if (index !== -1) {
                this.entries[index] = { ...this.entries[index], ...data };
                this.entries = [...this.entries];
                this.applyFilters();
            }
        } finally {
            this.loading = false;
        }
    }
}

export const entriesStore = new EntriesStore();
