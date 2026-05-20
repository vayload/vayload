import type { Activity } from "$lib/types";

export type DashboardActivity = Activity;

export interface DashboardState {
    activities: DashboardActivity[];
    loading: boolean;
}
