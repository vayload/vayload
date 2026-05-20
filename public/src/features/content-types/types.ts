import type { Collection } from "$lib/types";

export type ContentType = Collection;

export interface ContentTypesState {
    items: ContentType[];
    selected: ContentType | null;
    loading: boolean;
    total: number;
    page: number;
    pageSize: number;
    search: string;
}
