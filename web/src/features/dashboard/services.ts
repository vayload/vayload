import type { Activity } from "$lib/types";
import { httpClient } from "$lib/api/client";

export class DashboardService {
    async fetchActivities(): Promise<Activity[]> {
        const response = await httpClient.get<{ data: Activity[] }>("/dashboard/activities");
        return response.data;
    }
}

export const dashboardService = new DashboardService();
