import type { CreateCollectionDTO, UpdateCollectionDTO } from "./dtos";
import type { ContentType } from "./types";
import type { PaginatedResult } from "$lib/db/simulator";
import { httpClient } from "$lib/api/client";

export class ContentTypesService {
    async findMany(
        params: {
            page?: number;
            pageSize?: number;
            search?: string;
        } = {},
    ): Promise<PaginatedResult<ContentType>> {
        return httpClient.get<PaginatedResult<ContentType>>("/content-types", {
            page: params.page,
            pageSize: params.pageSize ?? 20,
            search: params.search,
        });
    }

    async findOneBySlug(slug: string): Promise<ContentType | null> {
        return httpClient.get<ContentType | null>(`/content-types/${slug}`);
    }

    async create(data: CreateCollectionDTO): Promise<ContentType> {
        return httpClient.post<ContentType>("/content-types", data);
    }

    async update(id: string, data: UpdateCollectionDTO): Promise<ContentType> {
        return httpClient.patch<ContentType>(`/content-types/${id}`, data);
    }

    async delete(id: string): Promise<void> {
        return httpClient.delete<void>(`/content-types/${id}`);
    }
}

export const contentTypesService = new ContentTypesService();
