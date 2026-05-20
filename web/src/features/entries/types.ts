import type { Entry, EntryField, EntryWithFields } from "$lib/types";

export type ContentEntry = Entry;
export type ContentEntryField = EntryField;
export type ContentEntryWithFields = EntryWithFields;

export interface EntriesState {
    items: ContentEntry[];
    filtered: ContentEntry[];
    loading: boolean;
    error: string | null;
    statusFilter: string;
    search: string;
    counts: Record<string, number>;
}
