import { dashboardService } from "./services";
import type { DashboardState, DashboardActivity } from "./types";

class DashboardStore {
    #state = $state<DashboardState>({
        activities: [],
        loading: false,
    });

    get activities(): DashboardActivity[] {
        return this.#state.activities;
    }
    get loading(): boolean {
        return this.#state.loading;
    }

    async fetchActivities() {
        this.#state.loading = true;
        try {
            this.#state.activities = await dashboardService.fetchActivities();
        } finally {
            this.#state.loading = false;
        }
    }
}

export const dashboardStore = new DashboardStore();
