import type { FolderObject, FileObject, FileUpload, FolderCreate } from "$lib/types";
import { authStore } from "$features/auth";
import { mediaService, type IFilesService } from "./services";
import { appContext } from "../../shared/store.svelte";

interface CacheEntry {
    folders: FolderObject[];
    files: FileObject[];
    timestamp: number;
}

export class MediaStore {
    private service: IFilesService;
    private cache = new Map<string, CacheEntry>();
    private readonly CACHE_TTL = 60 * 1000;

    public currentFolders = $state<FolderObject[]>([]);
    public currentFiles = $state<FileObject[]>([]);
    public currentFolderId = $state<string | null>(null);
    public loading = $state(false);
    public error = $state<string | null>(null);

    private _breadcrumb = $state<Array<{ id: string | null; name: string }>>([{ id: null, name: "Root" }]);
    private _seen = new Set<string | null>();

    constructor(service: IFilesService) {
        this.service = service;
    }

    public get currentItems() {
        return {
            folders: this.currentFolders,
            files: this.currentFiles,
        };
    }

    public async uploadFile(file: FileUpload) {
        this.loading = true;
        try {
            const result = await this.service.uploadFile(file);
            this.currentFiles.push(result);
        } finally {
            this.loading = false;
        }
    }

    public async createFolder(folder: FolderCreate) {
        this.loading = true;
        try {
            const result = await this.service.createFolder(folder);
            this.currentFolders.push(result);
        } finally {
            this.loading = false;
        }
    }

    public async navigate(folderId: string | null = null, folderName: string | null = null, forceRefresh = false) {
        this.error = null;
        this.loading = true;

        try {
            const projectId = appContext.currentProjectId ?? "";
            console.log({ folderId, projectId });
            const result = await this.service.getFolderContents(folderId, projectId);
            this.currentFolderId = folderId;
            this.currentFolders = result.folders;
            this.currentFiles = result.files;
        } catch (err) {
            this.error = "Error loading directory";
        } finally {
            this.loading = false;
        }
    }

    public async filterByCategory(category: string) {
        this.loading = true;
        try {
            const result = await this.service.filterByCategory(category);
            this.currentFolders = result.folders;
            this.currentFiles = result.files;
            this.currentFolderId = "filter-results";
        } finally {
            this.loading = false;
        }
    }

    public async rename(id: string, newName: string, type: "file" | "folder") {
        if (type === "file") {
            const file = this.currentFiles.find((f) => f.id === id);
            if (file) file.name = newName;
        } else {
            const folder = this.currentFolders.find((f) => f.id === id);
            if (folder) folder.name = newName;
        }

        await this.service.rename(id, newName, type);
        this.invalidateCache(this.currentFolderId);
    }

    public async move(id: string, newParentId: string | null, type: "file" | "folder") {
        if (type === "file") {
            this.currentFiles = this.currentFiles.filter((f) => f.id !== id);
        } else {
            this.currentFolders = this.currentFolders.filter((f) => f.id !== id);
        }

        await this.service.move(id, newParentId, type);
        this.invalidateCache(this.currentFolderId);
        this.invalidateCache(newParentId);
    }

    public async delete(id: string, type: "file" | "folder") {
        if (type === "file") {
            this.currentFiles = this.currentFiles.filter((f) => f.id !== id);
        } else {
            this.currentFolders = this.currentFolders.filter((f) => f.id !== id);
        }

        await this.service.delete(id, type);
        this.invalidateCache(this.currentFolderId);
        if (type === "folder") this.invalidateCache(id);
    }

    public async find(query: string) {
        this.loading = true;
        try {
            const results = await this.service.find(query);
            this.currentFolders = results.folders;
            this.currentFiles = results.files;
            this.currentFolderId = "search-results";
        } finally {
            this.loading = false;
        }
    }

    public async grep(query: string) {
        this.loading = true;
        try {
            const results = await this.service.grep(query);
            this.currentFolders = [];
            this.currentFiles = results;
            this.currentFolderId = "grep-results";
        } finally {
            this.loading = false;
        }
    }

    public async getPreviewUrl(fileId: string): Promise<string | null> {
        const preview = await this.service.getPreview(fileId);
        return preview?.url || null;
    }

    private invalidateCache(folderId: string | null) {
        this.cache.delete(String(folderId));
    }

    public get breadcrumb() {
        return this._breadcrumb;
    }
}

export const filesStore = new MediaStore(mediaService);
