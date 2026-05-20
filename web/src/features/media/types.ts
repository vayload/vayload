import type { FileObject, FolderObject } from "$lib/types";

export interface MediaState {
    folders: FolderObject[];
    files: FileObject[];
    loading: boolean;
    error: string | null;
    currentFolderId: string | null;
}
