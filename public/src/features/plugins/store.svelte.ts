import { pluginsService } from "./services";
import type { PluginsState, PluginIntegration } from "./types";

class PluginsStore {
    #state = $state<PluginsState>({
        items: [],
        loading: false,
    });

    get items(): PluginIntegration[] {
        return this.#state.items;
    }
    get loading(): boolean {
        return this.#state.loading;
    }

    async fetchAll() {
        this.#state.loading = true;
        try {
            this.#state.items = await pluginsService.findAll();
        } finally {
            this.#state.loading = false;
        }
    }

    async toggleInstall(id: string, installed: boolean) {
        const updated = await pluginsService.update(id, { installed });
        this.#state.items = this.#state.items.map((item) => (item.id === id ? updated : item));
        return updated;
    }
}

export const pluginsStore = new PluginsStore();
