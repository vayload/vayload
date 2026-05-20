import type { Integration } from "$lib/types";
import { httpClient } from "$lib/api/client";

export class PluginsService {
    async findAll(): Promise<Integration[]> {
        const result = await httpClient.get<{ data: Integration[] }>("/plugins");
        return result.data;
    }

    async update(id: string, data: Partial<Integration>): Promise<Integration> {
        return httpClient.patch<Integration>(`/plugins/${id}`, data);
    }
}

export const pluginsService = new PluginsService();
