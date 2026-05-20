import { contentTypesService } from "./services";
import type { ContentTypesState, ContentType } from "./types";
import type { CreateCollectionDTO, UpdateCollectionDTO } from "./dtos";

class ContentTypesStore {
    #state = $state<ContentTypesState>({
        items: [],
        selected: null,
        loading: false,
        total: 0,
        page: 1,
        pageSize: 20,
        search: "",
    });

    get items(): ContentType[] {
        return this.#state.items;
    }
    get selected(): ContentType | null {
        return this.#state.selected;
    }
    get loading(): boolean {
        return this.#state.loading;
    }
    get total(): number {
        return this.#state.total;
    }
    get page(): number {
        return this.#state.page;
    }
    get pageSize(): number {
        return this.#state.pageSize;
    }
    get search(): string {
        return this.#state.search;
    }
    get totalPages(): number {
        return Math.ceil(this.#state.total / this.#state.pageSize);
    }

    async fetch(page = this.#state.page, search = this.#state.search) {
        this.#state.loading = true;
        try {
            const result = await contentTypesService.findMany({ page, pageSize: this.#state.pageSize, search });
            this.#state.items = result.data;
            this.#state.total = result.total;
            this.#state.page = page;
            this.#state.search = search;
        } finally {
            this.#state.loading = false;
        }
    }

    async selectBySlug(slug: string) {
        this.#state.loading = true;
        try {
            this.#state.selected = await contentTypesService.findOneBySlug(slug);
        } finally {
            this.#state.loading = false;
        }
    }

    async create(data: CreateCollectionDTO) {
        const collection = await contentTypesService.create(data);
        this.#state.items = [collection, ...this.#state.items];
        this.#state.total += 1;
        return collection;
    }

    async update(id: string, data: UpdateCollectionDTO) {
        const collection = await contentTypesService.update(id, data);
        this.#state.items = this.#state.items.map((c) => (c.id === id ? collection : c));
        if (this.#state.selected?.id === id) this.#state.selected = collection;
        return collection;
    }

    async delete(id: string) {
        await contentTypesService.delete(id);
        this.#state.items = this.#state.items.filter((c) => c.id !== id);
        this.#state.total -= 1;
        if (this.#state.selected?.id === id) this.#state.selected = null;
    }
}

export const contentTypesStore = new ContentTypesStore();
