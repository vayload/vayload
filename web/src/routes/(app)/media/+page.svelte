<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import ScrollableTabs from "$lib/components/ScrollableTabs.svelte";
    import { filesStore } from "$lib/stores/files.store.svelte";
    import { Images, Activity as VideoIcon, FileText, Plus, Upload, FolderPen, MoreHorizontal } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import type { FileObject } from "$lib/data";

    import * as Breadcrumb from "$lib/components/ui/breadcrumb/index.js";

    let activeTab = $state("all");

    $effect(() => {
        // Solo navegar la primera vez si el store está vacío
        if (
            filesStore.currentFolderId === null &&
            filesStore.currentItems.files.length === 0 &&
            filesStore.currentItems.folders.length === 0
        ) {
            filesStore.navigate(null);
        }
    });

    function handleTabChange(tabId: string) {
        console.log(tabId);
        activeTab = tabId;
        if (tabId === "all") {
            filesStore.navigate(null);
        } else {
            filesStore.filterByCategory(tabId);
        }
    }

    const onHoverFile = async (file: FileObject) => {
        // filesStore.getPreviewUrl(file.id);
        if (file.category === "image") {
            const url = await filesStore.getPreviewUrl(file.id);
            console.log(url);
        }
    };

    const tabs = [
        { id: "all", label: "All Assets" },
        { id: "image", label: "Images" },
        { id: "video", label: "Videos" },
        { id: "document", label: "Documents" },
    ];
</script>

<div class="flex flex-col h-full">
    <SectionHeader
        title="Media Library"
        subtitle="Manage images, videos, and documents."
        breadcrumbs={["Workspace", "Assets"]}
    >
        {#snippet actions()}
            <Button variant="outline">
                <Plus size={16} class="mr-2" />
                New Folder
            </Button>
            <Button>
                <Upload size={16} class="mr-2" />
                Upload Asset
            </Button>
        {/snippet}
    </SectionHeader>

    <ScrollableTabs {tabs} {activeTab} onTabChange={handleTabChange} />

    <div class="py-6 flex-1 overflow-y-auto">
        <Breadcrumb.Root class="mb-6">
            <Breadcrumb.List>
                {#each filesStore.breadcrumb as crumb, index}
                    <Breadcrumb.Item>
                        {#if index === filesStore.breadcrumb.length - 1}
                            <Breadcrumb.Page>{crumb.name}</Breadcrumb.Page>
                        {:else}
                            <button
                                class="transition-colors hover:text-foreground"
                                onclick={() => filesStore.navigate(crumb.id, crumb.name)}
                            >
                                {crumb.name}
                            </button>
                        {/if}
                    </Breadcrumb.Item>
                    {#if index < filesStore.breadcrumb.length - 1}
                        <Breadcrumb.Separator />
                    {/if}
                {/each}
            </Breadcrumb.List>
        </Breadcrumb.Root>

        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {#if filesStore.loading}
                {#each Array(10) as _}
                    <div class="aspect-square bg-muted rounded-xl animate-pulse"></div>
                {/each}
            {:else}
                {#each filesStore.currentItems.files as file}
                    <div
                        class="group relative aspect-square bg-white border border-gray-200 rounded-xl overflow-hidden hover:shadow-md transition-all cursor-pointer"
                    >
                        <div
                            aria-label="File preview"
                            class="absolute inset-0 flex items-center justify-center bg-gray-100 text-gray-300"
                        >
                            {#if file.category === "image"}
                                {#if file.metadata.preview}
                                    <img
                                        src={file.metadata.preview}
                                        alt={file.name}
                                        class="w-full h-full object-cover"
                                    />
                                {:else}
                                    <Images size={32} />
                                {/if}
                            {:else if file.category === "video"}
                                <VideoIcon size={32} />
                            {:else}
                                <FileText size={32} />
                            {/if}
                        </div>
                        <div
                            class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors flex flex-col justify-end p-3 opacity-0 group-hover:opacity-100"
                        >
                            <div class="flex justify-between items-end text-white">
                                <div>
                                    <p class="text-xs font-medium truncate w-24">{file.name}</p>
                                    <p class="text-[10px] opacity-80">{(file.size / 1024 / 1024).toFixed(1)} MB</p>
                                </div>
                                <Button variant="ghost" size="icon" class="p-1 hover:bg-white/20 rounded h-auto">
                                    <MoreHorizontal size={16} />
                                </Button>
                            </div>
                        </div>
                        <div class="absolute top-2 left-2 opacity-0 group-hover:opacity-100 transition-opacity">
                            <Checkbox />
                        </div>
                    </div>
                {/each}
                {#each filesStore.currentItems.folders as folder}
                    <button
                        class="group aspect-square bg-primary/5 border border-primary/20 rounded-xl flex flex-col items-center justify-center cursor-pointer hover:bg-primary/10 hover:border-primary/30 transition-all"
                        onclick={() => filesStore.navigate(folder.id, folder.name)}
                    >
                        <FolderPen size={40} class="text-primary mb-2 fill-current opacity-80" />
                        <span class="text-sm font-medium text-foreground">{folder.name}</span>
                        <span class="text-xs text-muted-foreground">{folder.file_count} items</span>
                    </button>
                {/each}
            {/if}
        </div>
    </div>
</div>

<!-- <script lang="ts">
    import { onMount } from "svelte";
    import { Images, Activity, FileText, Plus, Upload, Folder, MoreHorizontal } from "@lucide/svelte";

    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { Breadcrumb, BreadcrumbItem } from "$lib/components/ui/breadcrumb";



    type AssetType = "image" | "video" | "doc";

    type Asset = {
        id: string;
        name: string;
        size: number;
        type: AssetType;
        folder: string; // ej: "Marketing/Product Shots"
    };


    let activeTab = $state<"all" | "images" | "videos" | "docs">("all");
    let search = $state("");
    let currentPath = $state<string[]>([]);
    let showUploadModal = $state(false);
    let showFolderModal = $state(false);
    let newFolderName = $state("");
    let uploadFiles = $state<FileList | null>(null);
    let loading = $state(false);

    let assets = $state<Asset[]>([]);

    onMount(() => {
        assets = [
            {
                id: "1",
                name: "hero.png",
                size: 1200000,
                type: "image",
                folder: "",
            },
            {
                id: "2",
                name: "ad-video.mp4",
                size: 5400000,
                type: "video",
                folder: "Marketing",
            },
            {
                id: "3",
                name: "brochure.pdf",
                size: 800000,
                type: "doc",
                folder: "Marketing",
            },
        ];
    });



    function getCurrentPathString() {
        return currentPath.join("/");
    }

    function enterFolder(folder: string) {
        currentPath = [...currentPath, folder];
    }

    function goToIndex(index: number) {
        currentPath = currentPath.slice(0, index + 1);
    }

    function goRoot() {
        currentPath = [];
    }

    function detectType(file: File): AssetType {
        if (file.type.startsWith("image")) return "image";
        if (file.type.startsWith("video")) return "video";
        return "doc";
    }

    function createFolder(name: string) {

        assets = [...assets];
    }

    async function uploadAssets(files: FileList) {
        loading = true;

        const newAssets: Asset[] = Array.from(files).map((file) => ({
            id: crypto.randomUUID(),
            name: file.name,
            size: file.size,
            type: detectType(file),
            folder: getCurrentPathString(),
        }));

        assets = [...assets, ...newAssets];

        loading = false;
    }



    const filteredAssets = $derived.by(() => {
        return assets.filter((asset) => {
            const pathMatch = asset.folder === getCurrentPathString();

            const searchMatch = asset.name.toLowerCase().includes(search.toLowerCase());

            const typeMatch =
                activeTab === "all" ||
                (activeTab === "images" && asset.type === "image") ||
                (activeTab === "videos" && asset.type === "video") ||
                (activeTab === "docs" && asset.type === "doc");

            return pathMatch && searchMatch && typeMatch;
        });
    });

    const folders = $derived.by(() => {
        const basePath = getCurrentPathString();
        const folderSet = new Set<string>();

        assets.forEach((asset) => {
            if (!asset.folder.startsWith(basePath)) return;

            const relative = asset.folder.slice(basePath.length).replace(/^\/?/, "");
            const segments = relative.split("/").filter(Boolean);

            if (segments.length > 0) {
                folderSet.add(segments[0]);
            }
        });

        return Array.from(folderSet);
    });
</script>

<div class="flex flex-col h-full">
    <div class="flex items-center justify-between pb-4">
        <div>
            <h1 class="text-2xl font-semibold">Media Library</h1>
            <p class="text-sm text-muted-foreground">Manage images, videos and documents</p>
        </div>

        <div class="flex gap-2">
            <Button variant="outline" onclick={() => (showFolderModal = true)}>
                <Plus size={16} class="mr-2" />
                New Folder
            </Button>

            <Button onclick={() => (showUploadModal = true)}>
                <Upload size={16} class="mr-2" />
                Upload
            </Button>
        </div>
    </div>

    <div class="flex gap-2 pb-4">
        <Button variant={activeTab === "all" ? "default" : "outline"} onclick={() => (activeTab = "all")}>All</Button>
        <Button variant={activeTab === "images" ? "default" : "outline"} onclick={() => (activeTab = "images")}
            >Images</Button
        >
        <Button variant={activeTab === "videos" ? "default" : "outline"} onclick={() => (activeTab = "videos")}
            >Videos</Button
        >
        <Button variant={activeTab === "docs" ? "default" : "outline"} onclick={() => (activeTab = "docs")}>Docs</Button
        >
    </div>

    <div class="flex justify-between items-center pb-4 gap-4">
        <Input placeholder="Search..." bind:value={search} class="max-w-sm" />

        {#if currentPath.length > 0}
            <Breadcrumb>
                <BreadcrumbItem>
                    <button onclick={goRoot}>Root</button>
                </BreadcrumbItem>

                {#each currentPath as segment, i}
                    <BreadcrumbItem>
                        <button onclick={() => goToIndex(i)}>
                            {segment}
                        </button>
                    </BreadcrumbItem>
                {/each}
            </Breadcrumb>
        {/if}
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        {#each folders as folder}
            <button
                onclick={() => enterFolder(folder)}
                class="aspect-square bg-primary/5 border border-primary/20 rounded-xl flex flex-col items-center justify-center hover:bg-primary/10 transition"
            >
                <Folder size={36} class="text-primary mb-2" />
                <span class="text-sm font-medium">{folder}</span>
            </button>
        {/each}

        {#each filteredAssets as asset}
            <div
                class="group relative aspect-square bg-white border rounded-xl overflow-hidden hover:shadow-md transition"
            >
                <div class="absolute inset-0 flex items-center justify-center bg-muted text-muted-foreground">
                    {#if asset.type === "image"}
                        <Images size={32} />
                    {:else if asset.type === "video"}
                        <Activity size={32} />
                    {:else}
                        <FileText size={32} />
                    {/if}
                </div>

                <div class="absolute top-2 left-2 opacity-0 group-hover:opacity-100">
                    <Checkbox />
                </div>

                <div
                    class="absolute bottom-0 w-full bg-black/40 text-white p-2 opacity-0 group-hover:opacity-100 transition"
                >
                    <p class="text-xs truncate">{asset.name}</p>
                    <p class="text-[10px] opacity-80">
                        {(asset.size / 1024 / 1024).toFixed(1)} MB
                    </p>
                </div>
            </div>
        {/each}
    </div>

    <Dialog.Root bind:open={showFolderModal}>
        <Dialog.Content class="sm:max-w-md">
            <Dialog.Header>
                <Dialog.Title>Create Folder</Dialog.Title>
            </Dialog.Header>

            <div class="space-y-4 py-4">
                <Label>Name</Label>
                <Input bind:value={newFolderName} />
            </div>

            <Dialog.Footer>
                <Button
                    onclick={() => {
                        if (!newFolderName) return;
                        createFolder(newFolderName);
                        newFolderName = "";
                        showFolderModal = false;
                    }}
                >
                    Create
                </Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>

    <Dialog.Root bind:open={showUploadModal}>
        <Dialog.Content class="sm:max-w-lg">
            <Dialog.Header>
                <Dialog.Title>Upload Files</Dialog.Title>
            </Dialog.Header>

            <div
                class="border-2 border-dashed rounded-xl p-10 text-center"
                ondrop={(e) => (uploadFiles = e.dataTransfer?.files ?? null)}
                ondragover={(e) => e.preventDefault()}
            >
                <p class="text-sm text-muted-foreground">Drag & drop files here</p>

                <input type="file" multiple bind:files={uploadFiles} class="hidden" id="fileInput" />

                <Button variant="secondary" class="mt-4" onclick={() => document.getElementById("fileInput")?.click()}>
                    Select Files
                </Button>
            </div>

            <Dialog.Footer>
                <Button
                    onclick={async () => {
                        if (!uploadFiles) return;
                        await uploadAssets(uploadFiles);
                        uploadFiles = null;
                        showUploadModal = false;
                    }}
                >
                    Upload
                </Button>
            </Dialog.Footer>
        </Dialog.Content>
    </Dialog.Root>
</div> -->
