import { usersService } from "./services";
import type { UsersState, User } from "./types";
import type { CreateUserDTO, UpdateUserDTO } from "./dtos";
import type { PaginatedResult } from "$lib/db/simulator";

class UsersStore {
    #state = $state<UsersState>({
        items: [],
        selected: null,
        loading: false,
        total: 0,
        page: 1,
        pageSize: 20,
        search: "",
    });

    get items(): User[] {
        return this.#state.items;
    }
    get selected(): User | null {
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
            const result = await usersService.findMany({ page, pageSize: this.#state.pageSize, search });
            this.#state.items = result.data;
            this.#state.total = result.total;
            this.#state.page = page;
            this.#state.search = search;
        } finally {
            this.#state.loading = false;
        }
    }

    async select(id: string) {
        this.#state.loading = true;
        try {
            this.#state.selected = await usersService.findOne(id);
        } finally {
            this.#state.loading = false;
        }
    }

    async create(data: CreateUserDTO) {
        const user = await usersService.create(data);
        this.#state.items = [user, ...this.#state.items];
        this.#state.total += 1;
        return user;
    }

    async update(id: string, data: UpdateUserDTO) {
        const user = await usersService.update(id, data);
        this.#state.items = this.#state.items.map((u) => (u.id === id ? user : u));
        if (this.#state.selected?.id === id) this.#state.selected = user;
        return user;
    }

    async delete(id: string) {
        await usersService.delete(id);
        this.#state.items = this.#state.items.filter((u) => u.id !== id);
        this.#state.total -= 1;
    }
}

export const usersStore = new UsersStore();
