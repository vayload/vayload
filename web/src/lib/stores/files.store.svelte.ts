import { type FolderObject, type FileObject, getUser, type FileUpload, type FolderCreate } from "$lib/data";
import { filesService, type IFilesService } from "./../dumy/files";
import { appContext } from "./app-context.svelte";

interface CacheEntry {
    folders: FolderObject[];
    files: FileObject[];
    timestamp: number;
}

export class FilesStore {
    private service: IFilesService;
    private cache = new Map<string, CacheEntry>();
    private readonly CACHE_TTL = 60 * 1000; // 1 Minuto de TTL

    // Svelte 5 States
    public currentFolders = $state<FolderObject[]>([]);
    public currentFiles = $state<FileObject[]>([]);
    public currentFolderId = $state<string | null>(null);
    public loading = $state(false);
    public error = $state<string | null>(null);

    private _breadcrumb = $state<Array<{ id: string | null; name: string }>>([{ id: null, name: "Root" }]);
    private _seen = new Set();

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

    /**
     * Navegación Cacheada
     */
    public async navigate(folderId: string | null = null, folderName: string | null = null, forceRefresh = false) {
        this.error = null;

        if (folderName) {
            if (!this._seen.has(folderId)) {
                this._breadcrumb.push({ id: folderId, name: folderName });
                this._seen.add(folderId);
            }
        } else if (folderId === null) {
            this._breadcrumb = [{ id: null, name: "Root" }];
        }

        const cacheKey = String(folderId);

        if (!forceRefresh && this.cache.has(cacheKey)) {
            const entry = this.cache.get(cacheKey)!;
            if (Date.now() - entry.timestamp < this.CACHE_TTL) {
                this.currentFolderId = folderId;
                this.currentFolders = entry.folders;
                this.currentFiles = entry.files;
                return;
            }
        }

        this.loading = true;
        try {
            const result = await this.service.getFolderContents(folderId, getUser().id, appContext.currentProjectId!);

            this.cache.set(cacheKey, {
                folders: result.folders,
                files: result.files,
                timestamp: Date.now(),
            });

            this.currentFolderId = folderId;
            this.currentFolders = result.folders;
            this.currentFiles = result.files;
        } catch (err) {
            this.error = "Error al cargar el directorio";
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

    /**
     * Operaciones VFS (Con invalidación de caché inteligente)
     */
    public async rename(id: string, newName: string, type: "file" | "folder") {
        // Actualización optimista (UI primero)
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
        // Optimista: quitar del current view
        if (type === "file") {
            this.currentFiles = this.currentFiles.filter((f) => f.id !== id);
        } else {
            this.currentFolders = this.currentFolders.filter((f) => f.id !== id);
        }

        await this.service.move(id, newParentId, type);
        this.invalidateCache(this.currentFolderId);
        this.invalidateCache(newParentId); // Invalidamos el destino para que lo pida fresco si navegamos a él
    }

    public async delete(id: string, type: "file" | "folder") {
        // Optimista
        if (type === "file") {
            this.currentFiles = this.currentFiles.filter((f) => f.id !== id);
        } else {
            this.currentFolders = this.currentFolders.filter((f) => f.id !== id);
        }

        await this.service.delete(id, type);
        this.invalidateCache(this.currentFolderId);
        if (type === "folder") this.invalidateCache(id); // Limpiar posibles cachés hijos
    }

    /**
     * Búsqueda por nombre de archivo/carpeta
     */
    public async find(query: string) {
        this.loading = true;
        try {
            const results = await this.service.find(query);
            this.currentFolders = results.folders;
            this.currentFiles = results.files;
            this.currentFolderId = "search-results"; // Estado virtual
        } finally {
            this.loading = false;
        }
    }

    /**
     * Búsqueda profunda (Grep) dentro del contenido o metadatos
     */
    public async grep(query: string) {
        this.loading = true;
        try {
            const results = await this.service.grep(query);
            this.currentFolders = []; // Grep solo retorna archivos
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
        const lastFolder = this._breadcrumb.findIndex((folder) => folder.id === this.currentFolderId);
        return this._breadcrumb.slice(0, lastFolder + 1);
    }
}

export const filesStore = new FilesStore(filesService);
