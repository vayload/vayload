import { faker } from "@faker-js/faker";
import { ulid } from "$lib/data"; // Asegúrate de ajustar tus imports
import type { FolderObject, FileObject, FileUpload, FolderCreate } from "./../data";

export interface ImagePreview {
    id: string;
    file_id: string;
    url: string; // URL de Unsplash/Pexels
}

export interface IFilesService {
    uploadFile(file: FileUpload): Promise<FileObject>;
    createFolder(folder: FolderCreate): Promise<FolderObject>;
    getFolderContents(
        folderId: string | null,
        projectId: string,
        ownerId: string,
    ): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    filterByCategory(category: string): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    rename(id: string, newName: string, type: "file" | "folder"): Promise<void>;
    move(id: string, newParentId: string | null, type: "file" | "folder"): Promise<void>;
    delete(id: string, type: "file" | "folder"): Promise<void>;
    find(query: string): Promise<{ folders: FolderObject[]; files: FileObject[] }>;
    grep(query: string): Promise<FileObject[]>;
    getPreview(fileId: string): Promise<ImagePreview | null>;
}

const DELAY = 300; // Simular latencia de red

export class DummyFilesService implements IFilesService {
    private folders: FolderObject[] = [];
    private files: FileObject[] = [];
    private imagePreviews: ImagePreview[] = [];

    public async uploadFile(file: FileUpload): Promise<FileObject> {
        await this.simulateDelay();
        const fileObject = this.createMockFile(file.owner_id, file.project_id, file.folder_id);
        this.files.push(fileObject);
        return fileObject;
    }

    public async createFolder(folder: FolderCreate): Promise<FolderObject> {
        await this.simulateDelay();
        const folderObject = this.createMockFolder(folder.owner_id, folder.project_id, folder.parent_id, folder.name);
        this.folders.push(folderObject);
        return folderObject;
    }

    private generateSeedData(ownerId: string, projectId: string) {
        // Generar 20 Carpetas (4 en la raíz, 4 hijas cada una)
        let rootFolders: FolderObject[] = [];
        for (let i = 0; i < 4; i++) {
            rootFolders.push(this.createMockFolder(ownerId, projectId, null, faker.lorem.word()));
        }

        for (let i = 0; i < 4; i++) {
            this.files.push(this.createMockFile(ownerId, projectId, null));
        }

        let subFolders: FolderObject[] = [];
        for (const root of rootFolders) {
            for (let j = 0; j < 4; j++) {
                subFolders.push(
                    this.createMockFolder(ownerId, projectId, root.id, faker.lorem.word(), root.path, root.depth),
                );
            }
        }

        this.folders = [...rootFolders, ...subFolders];

        // Generar 100 Archivos (5 en cada carpeta)
        for (const folder of this.folders) {
            for (let k = 0; k < 5; k++) {
                const file = this.createMockFile(ownerId, projectId, folder.id);
                this.files.push(file);

                // Si es imagen, crear su registro en la tabla de previews
                if (file.category === "image") {
                    this.imagePreviews.push({
                        id: ulid(),
                        file_id: file.id,
                        url: faker.image.url(), // URL aleatoria tipo Unsplash/Flickr
                    });
                }
            }
        }
    }

    private createMockFolder(
        ownerId: string,
        projectId: string,
        parentId: string | null,
        name: string,
        parentPath = "",
        parentDepth = -1,
    ): FolderObject {
        const path = parentPath ? `${parentPath}/${name}` : `/${name}`;
        return {
            id: ulid(),
            owner_id: ownerId,
            project_id: projectId,
            parent_id: parentId,
            name,
            file_count: 5,
            subfolder_count: parentId === null ? 4 : 0,
            path,
            depth: parentDepth + 1,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    private createMockFile(ownerId: string, projectId: string, folderId: string | null): FileObject {
        const mimeType = faker.helpers.arrayElement([
            "image/jpeg",
            "image/png",
            "video/mp4",
            "application/pdf",
            "text/plain",
        ]);
        const category = mimeType.startsWith("image") ? "image" : mimeType.startsWith("video") ? "video" : "document";

        // Agregar contenido de texto falso en metadatos para simular el "grep"
        const mockContent = category === "document" ? faker.lorem.paragraphs(2) : "";

        return {
            id: ulid(),
            owner_id: ownerId,
            project_id: projectId,
            folder_id: folderId,
            name: faker.system.fileName(),
            mime_type: mimeType,
            category,
            size: faker.number.int({ min: 1024, max: 50000000 }),
            provider: "s3",
            provider_key: faker.string.uuid(),
            metadata: { textContent: mockContent, preview: faker.image.url() },
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        };
    }

    private async simulateDelay() {
        return new Promise((resolve) => setTimeout(resolve, DELAY));
    }

    async getFolderContents(folderId: string | null, projectId: string, ownerId: string) {
        await this.simulateDelay();
        if (this.folders.length === 0) {
            this.generateSeedData(ownerId, projectId);
        }

        return {
            folders: this.folders.filter((f) => f.parent_id === folderId),
            files: this.files.filter((f) => f.folder_id === folderId),
        };
    }

    async filterByCategory(category: string): Promise<{ folders: FolderObject[]; files: FileObject[] }> {
        await this.simulateDelay();

        return { folders: [], files: this.files.filter((f) => f.category === category) };
    }

    async rename(id: string, newName: string, type: "file" | "folder") {
        await this.simulateDelay();
        const item = type === "file" ? this.files.find((f) => f.id === id) : this.folders.find((f) => f.id === id);
        if (item) {
            item.name = newName;
            item.updated_at = new Date().toISOString();
        }
    }

    async move(id: string, newParentId: string | null, type: "file" | "folder") {
        await this.simulateDelay();
        if (type === "file") {
            const file = this.files.find((f) => f.id === id);
            if (file) file.folder_id = newParentId;
        } else {
            const folder = this.folders.find((f) => f.id === id);
            if (folder) folder.parent_id = newParentId;
            // Nota: En un backend real aquí recalcularías paths recursivamente.
        }
    }

    async delete(id: string, type: "file" | "folder") {
        await this.simulateDelay();
        if (type === "file") {
            this.files = this.files.filter((f) => f.id !== id);
        } else {
            // Eliminación en cascada simple para el mock
            this.folders = this.folders.filter((f) => f.id !== id && f.parent_id !== id);
            this.files = this.files.filter((f) => f.folder_id !== id);
        }
    }

    async find(query: string) {
        await this.simulateDelay();
        const q = query.toLowerCase();
        return {
            folders: this.folders.filter((f) => f.name.toLowerCase().includes(q)),
            files: this.files.filter((f) => f.name.toLowerCase().includes(q)),
        };
    }

    async grep(query: string) {
        await this.simulateDelay();
        const q = query.toLowerCase();
        // Busca dentro del mockContent de los metadatos
        return this.files.filter((f) => f.metadata?.textContent?.toLowerCase().includes(q));
    }

    async getPreview(fileId: string) {
        await this.simulateDelay();
        return this.imagePreviews.find((p) => p.file_id === fileId) || null;
    }
}

export const filesService = new DummyFilesService();
