import type { FileObject, FolderObject, FileUpload, FolderCreate } from "$lib/types";
import { httpClient } from "$lib/api/client";

export interface ImagePreview {
    id: string;
    file_id: string;
    url: string;
}

export interface IFilesService {
    uploadFile(file: FileUpload): Promise<FileObject>;
    createFolder(folder: FolderCreate): Promise<FolderObject>;
    getFolderContents(
        folderId: string | null,
        projectId: string,
    ): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    filterByCategory(category: string): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    rename(id: string, newName: string, type: "file" | "folder"): Promise<void>;
    move(id: string, newParentId: string | null, type: "file" | "folder"): Promise<void>;
    delete(id: string, type: "file" | "folder"): Promise<void>;
    find(query: string): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    grep(query: string): Promise<FileObject[]>;
    getPreview(fileId: string): Promise<ImagePreview | null>;
}

export class MediaService implements IFilesService {
    async uploadFile(file: FileUpload): Promise<FileObject> {
        return httpClient.post<FileObject>("/media/upload", file);
    }

    async createFolder(folder: FolderCreate): Promise<FolderObject> {
        return httpClient.post<FolderObject>("/media/folders", folder);
    }

    async getFolderContents(folderId: string | null, projectId: string) {
        const [foldersResult, filesResult] = await Promise.all([
            httpClient.get<{ data: FolderObject[] }>("/media/folders", {
                parent_id: folderId as any,
                project_id: projectId,
            }),
            httpClient.get<{ data: FileObject[] }>("/media/files", {
                folder_id: folderId as any,
                project_id: projectId,
            }),
        ]);

        return {
            folders: foldersResult.data,
            files: filesResult.data,
        };
    }

    async filterByCategory(category: string) {
        const result = await httpClient.get<{ data: FileObject[] }>("/media/filter", {
            category,
        });
        return { folders: [], files: result.data };
    }

    async rename(id: string, newName: string, type: "file" | "folder") {
        await httpClient.patch("/media/rename", { id, name: newName, type });
    }

    async move(id: string, newParentId: string | null, type: "file" | "folder") {
        await httpClient.patch("/media/move", { id, parent_id: newParentId, type });
    }

    async delete(id: string, type: "file" | "folder") {
        await httpClient.delete("/media/delete", { id, type });
    }

    async find(query: string) {
        return httpClient.get<{ folders: FolderObject[]; files: FileObject[] }>("/media/find", {
            query,
        });
    }

    async grep(query: string) {
        return httpClient.get<FileObject[]>("/media/grep", {
            query,
        });
    }

    async getPreview(fileId: string) {
        return httpClient.get<ImagePreview | null>(`/media/preview/${fileId}`);
    }
}

export const mediaService = new MediaService();
