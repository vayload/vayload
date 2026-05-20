import type { Entry, EntryWithFields } from "$lib/types";
import type { CreateEntryDTO, UpdateEntryDTO } from "./dtos";
import type { PaginatedResult } from "$lib/db/simulator";
import { httpClient } from "$lib/api/client";

export class EntriesService {
    async findMany(
        params: {
            status?: Entry["status"] | "all";
            search?: string;
            page?: number;
            pageSize?: number;
        } = {},
    ): Promise<PaginatedResult<Entry>> {
        return httpClient.get<PaginatedResult<Entry>>("/entries", {
            status: params.status ?? "all",
            search: params.search,
            page: params.page,
            pageSize: params.pageSize ?? 50,
        });
    }

    async countByStatus(): Promise<Record<string, number>> {
        return httpClient.get<Record<string, number>>("/entries/counts");
    }

    async findOneBySlug(slug: string): Promise<Entry | null> {
        return httpClient.get<Entry | null>(`/entries/${slug}`);
    }

    async findWithFields(slug: string): Promise<EntryWithFields | null> {
        return httpClient.get<EntryWithFields | null>(`/entries/${slug}/fields`);
    }

    async create(data: CreateEntryDTO): Promise<Entry> {
        return httpClient.post<Entry>("/entries", data);
    }

    async update(id: string, data: UpdateEntryDTO): Promise<Entry> {
        return httpClient.patch<Entry>(`/entries/${id}`, data);
    }

    async delete(id: string): Promise<void> {
        return httpClient.delete<void>(`/entries/${id}`);
    }
}

export const entriesService = new EntriesService();
