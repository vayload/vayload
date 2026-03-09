<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { Select, SelectContent, SelectItem, SelectTrigger } from "$lib/components/ui/select";
    import { Button } from "$lib/components/ui/button";
    import { Plus, Trash2, Settings2, ArrowLeft, Files } from "@lucide/svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import FieldTypeGrid from "./FieldTypeGrid.svelte";
    import { FieldTypes, FieldTypeMetadata } from "$lib/constants/field-types";
    import { goto } from "$app/navigation";

    interface Field {
        name: string;
        type: string;
        required: boolean;
        label?: string;
        config?: any;
    }

    interface ContentType {
        id: string;
        name: string;
        slug: string;
        description?: string;
        fields?: Field[];
        fields_schema?: Record<string, any>;
        settings: {
            published: boolean;
            revision: boolean;
            preview: boolean;
        };
    }

    interface Props {
        collection: ContentType;
        onUpdate: (collection: ContentType) => void;
    }

    let { collection: propCollection, onUpdate }: Props = $props();

    let isAddingField = $state(false);
    let selectedFieldType: FieldTypes | null = $state(null);
    let collection = $state<ContentType>({} as any);

    let newField = $state<Field>({
        name: "",
        type: "",
        required: false,
        label: "",
    });

    $effect(() => {
        collection = propCollection;
        if (collection.fields_schema && (!collection.fields || collection.fields.length === 0)) {
            collection.fields = Object.entries(collection.fields_schema).map(([name, data]) => ({
                name,
                ...(data as any),
            }));
        }
    });

    function openAddField(type: FieldTypes) {
        selectedFieldType = type;
        newField = {
            name: "",
            type: type,
            required: false,
            label: "",
        };

        isAddingField = true;
    }

    function handleAddField() {
        if (!newField.name || !collection) return;

        if (!collection.fields) {
            collection.fields = [];
        }

        collection.fields = [...collection.fields, { ...newField }];

        if (!collection.fields_schema) {
            collection.fields_schema = {};
        }
        collection.fields_schema[newField.name] = {
            type: newField.type,
            required: newField.required,
            label: newField.label || newField.name,
        };

        onUpdate(collection);
        isAddingField = false;
        selectedFieldType = null;
    }

    function removeField(fieldName: string) {
        if (!collection || !collection.fields) return;

        collection.fields = collection.fields.filter((f) => f.name !== fieldName);
        if (collection.fields_schema) {
            delete collection.fields_schema[fieldName];
        }

        onUpdate(collection);
    }
</script>

<div class="space-y-6">
    <div class="flex items-center justify-between">
        <div>
            <h2 class="text-2xl font-bold">Content Types</h2>
            <p class="text-muted-foreground">Manage your data structures and fields.</p>
        </div>
        <Button variant="outline" onclick={() => goto("/content-types")}>
            <ArrowLeft size={16} class="mr-2" />
            Back to List
        </Button>
    </div>

    <!-- {#if !collection && !isCreating}
        <Card>
            <CardHeader>
                <CardTitle>Available Collections</CardTitle>
                <CardDescription>Select a content type to manage its schema.</CardDescription>
            </CardHeader>
            <CardContent>
                {#if collections.length === 0}
                    <div class="text-center py-12 border-2 border-dashed rounded-xl">
                        <p class="text-muted-foreground">No content types found.</p>
                        <Button variant="link" onclick={() => (isCreating = true)}>Create your first one</Button>
                    </div>
                {:else}
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {#each collections as collection}
                            <div
                                role="button"
                                tabindex="0"
                                aria-label="Select collection"
                                class="p-5 border rounded-xl hover:border-primary hover:bg-primary/5 transition-all group cursor-pointer"
                                onclick={() => handleSelect(collection)}
                            >
                                <div class="flex items-start justify-between">
                                    <div class="space-y-1">
                                        <h3 class="font-bold text-lg group-hover:text-primary transition-colors">
                                            {collection.name}
                                        </h3>
                                        <p
                                            class="text-xs font-mono text-muted-foreground bg-muted w-fit px-2 py-0.5 rounded italic"
                                        >
                                            {collection.slug}
                                        </p>
                                        {#if collection.description}
                                            <p class="text-sm text-muted-foreground line-clamp-2 mt-2">
                                                {collection.description}
                                            </p>
                                        {/if}
                                    </div>
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        class="text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            handleDelete(collection.id);
                                        }}
                                    >
                                        <Trash2 size={18} />
                                    </Button>
                                </div>
                                <div class="mt-4 flex items-center gap-3 text-xs font-medium text-muted-foreground">
                                    <div class="flex items-center gap-1.5">
                                        <Settings2 size={14} />
                                        <span>{Object.keys(collection.fields_schema || {}).length} Fields</span>
                                    </div>
                                    {#if collection.settings.published}
                                        <span class="px-2 py-0.5 bg-green-100 text-green-700 rounded-full"
                                            >Published</span
                                        >
                                    {/if}
                                </div>
                            </div>
                        {/each}

                        <button
                            class="p-5 border-2 border-dashed rounded-xl hover:border-primary hover:text-primary hover:bg-primary/5 transition-all text-muted-foreground flex flex-col items-center justify-center gap-2 group"
                            onclick={() => (isCreating = true)}
                        >
                            <div class="p-3 rounded-full bg-muted group-hover:bg-primary/10 transition-colors">
                                <Plus size={24} />
                            </div>
                            <span class="font-semibold">Create New Content Type</span>
                        </button>
                    </div>
                {/if}
            </CardContent>
        </Card>
    {/if} -->

    <!-- {#if isCreating}
        <Card>
            <CardHeader>
                <CardTitle>Create New Content Type</CardTitle>
                <CardDescription>Define a new collection to store your content.</CardDescription>
            </CardHeader>
            <CardContent class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-2">
                        <Label for="collection-name">Display Name *</Label>
                        <Input
                            id="collection-name"
                            placeholder="e.g. Blog Post"
                            bind:value={newCollection.name}
                            oninput={handleNameChange}
                        />
                    </div>
                    <div class="space-y-2">
                        <Label for="collection-slug">API ID (Slug) *</Label>
                        <Input id="collection-slug" placeholder="e.g. blog_post" bind:value={newCollection.slug} />
                    </div>
                </div>

                <div class="space-y-2">
                    <Label for="collection-description">Description</Label>
                    <Input
                        id="collection-description"
                        placeholder="What is this collection for?"
                        bind:value={newCollection.description}
                    />
                </div>

                <div class="space-y-4">
                    <Label>Permissions & Options</Label>
                    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <div class="flex items-center justify-between p-3 border rounded-lg">
                            <Label for="published" class="cursor-pointer">Published</Label>
                            <Switch id="published" bind:checked={newCollection.settings.published} />
                        </div>
                        <div class="flex items-center justify-between p-3 border rounded-lg">
                            <Label for="revision" class="cursor-pointer">Revisions</Label>
                            <Switch id="revision" bind:checked={newCollection.settings.revision} />
                        </div>
                        <div class="flex items-center justify-between p-3 border rounded-lg">
                            <Label for="preview" class="cursor-pointer">Preview</Label>
                            <Switch id="preview" bind:checked={newCollection.settings.preview} />
                        </div>
                    </div>
                </div>

                <div class="flex items-center gap-3 pt-4 border-t">
                    <Button onclick={handleCreate} disabled={!newCollection.name || !newCollection.slug}>
                        Create & Continue
                    </Button>
                    <Button variant="ghost" onclick={() => (isCreating = false)}>Cancel</Button>
                </div>
            </CardContent>
        </Card>
    {/if} -->

    <!-- {#if collection && !isCreating} -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
            <Card>
                <CardHeader class="flex flex-row items-center justify-between space-y-0">
                    <div>
                        <CardTitle>Fields</CardTitle>
                        <CardDescription>Manage the schema for {collection.name}.</CardDescription>
                    </div>

                    <Dialog.Root>
                        <Dialog.Trigger>
                            <Button size="sm">
                                <Plus size={16} class="mr-2" />
                                Add Field
                            </Button>
                        </Dialog.Trigger>
                        <Dialog.Content class="max-w-4xl max-h-[90vh] overflow-y-auto">
                            <Dialog.Header>
                                <Dialog.Title>Select Field Type</Dialog.Title>
                                <Dialog.Description
                                    >Choose the type of field you want to add to this collection.</Dialog.Description
                                >
                            </Dialog.Header>
                            <div class="py-6">
                                <FieldTypeGrid onSelect={openAddField} />
                            </div>
                        </Dialog.Content>
                    </Dialog.Root>
                </CardHeader>
                <CardContent>
                    <div class="space-y-3">
                        {#if !collection.fields || collection.fields.length === 0}
                            <div class="text-center py-12 border-2 border-dashed rounded-xl">
                                <p class="text-muted-foreground">No fields defined yet.</p>
                            </div>
                        {:else}
                            <div class="border rounded-xl divide-y">
                                {#each collection.fields as field}
                                    <div
                                        class="p-4 flex items-center justify-between hover:bg-muted/30 transition-colors"
                                    >
                                        <div class="flex items-center gap-4">
                                            <div class="p-2 rounded-lg bg-muted text-primary">
                                                {#if FieldTypeMetadata[field.type as FieldTypes]}
                                                    <!-- {@render lucideIcons[
                                                                FieldTypeMetadata[field.type as FieldTypes].icon
                                                            ]} -->
                                                {:else}
                                                    <Settings2 size={18} />
                                                {/if}
                                            </div>
                                            <div>
                                                <div class="flex items-center gap-2">
                                                    <span class="font-bold">{field.label || field.name}</span>
                                                    {#if field.required}
                                                        <span
                                                            class="text-[10px] font-bold uppercase tracking-wider bg-red-100 text-red-600 px-1.5 py-0.5 rounded"
                                                            >Required</span
                                                        >
                                                    {/if}
                                                </div>
                                                <p class="text-xs font-mono text-muted-foreground mt-0.5">
                                                    {field.name} • {field.type}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex items-center gap-1">
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                class="text-muted-foreground hover:text-primary"
                                            >
                                                <Settings2 size={16} />
                                            </Button>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                class="text-destructive hover:bg-destructive/10"
                                                onclick={() => removeField(field.name)}
                                            >
                                                <Trash2 size={16} />
                                            </Button>
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Schema Preview</CardTitle>
                    <CardDescription>The underlying JSON schema for this collection.</CardDescription>
                </CardHeader>
                <CardContent>
                    <div class="p-4 bg-slate-950 rounded-xl overflow-x-auto">
                        <pre class="text-sm text-green-400 font-mono"><code
                                >{JSON.stringify(collection.fields_schema, null, 2)}</code
                            ></pre>
                    </div>
                </CardContent>
            </Card>
        </div>

        <div class="space-y-6">
            <Card>
                <CardHeader>
                    <CardTitle>Collection Settings</CardTitle>
                </CardHeader>
                <CardContent class="space-y-4">
                    <div class="space-y-2">
                        <Label>Display Name</Label>
                        <Input bind:value={collection.name} oninput={() => onUpdate(collection as any)} />
                    </div>
                    <div class="space-y-2">
                        <Label>API ID</Label>
                        <Input value={collection.slug} disabled class="font-mono bg-muted" />
                    </div>
                    <div class="space-y-2">
                        <Label>Description</Label>
                        <Input bind:value={collection.description} oninput={() => onUpdate(collection as any)} />
                    </div>

                    <!-- <div class="pt-4 space-y-3">
                            <div class="flex items-center justify-between">
                                <Label>Published</Label>
                                <Switch
                                    bind:checked={collection.settings.published}
                                    onchange={() => onUpdate(collection as any)}
                                />
                            </div>
                            <div class="flex items-center justify-between">
                                <Label>Revisions</Label>
                                <Switch
                                    bind:checked={collection.settings.revision}
                                    onchange={() => onUpdate(collection as any)}
                                />
                            </div>
                            <div class="flex items-center justify-between">
                                <Label>Preview</Label>
                                <Switch
                                    bind:checked={collection.settings.preview}
                                    onchange={() => onUpdate(collection as any)}
                                />
                            </div>
                        </div> -->
                </CardContent>
            </Card>
        </div>
    </div>
    <!-- {/if} -->
</div>

<Dialog.Root bind:open={isAddingField}>
    <Dialog.Content class="z-60">
        <Dialog.Header>
            <Dialog.Title>
                Configure {selectedFieldType ? FieldTypeMetadata[selectedFieldType].label : ""} Field
            </Dialog.Title>
            <Dialog.Description>Set the name and basic settings for your new field.</Dialog.Description>
        </Dialog.Header>
        <div class="space-y-4 py-4 z-50">
            <div class="space-y-2">
                <Label for="field-label">Display Label</Label>
                <Input
                    id="field-label"
                    placeholder="e.g. Meta Description"
                    bind:value={newField.label}
                    oninput={() => (newField.name = (newField.label || "").toLowerCase().replace(/\s/g, "_"))}
                />
            </div>
            <div class="space-y-2">
                <Label for="field-name">API ID (Unique Name)</Label>
                <Input id="field-name" placeholder="meta_description" bind:value={newField.name} class="font-mono" />
            </div>
            <div class="flex items-center space-x-2 pt-2">
                <Switch id="field-required" bind:checked={newField.required} />
                <Label for="field-required" class="text-sm font-normal">Required field</Label>
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => ((isAddingField = false), (selectedFieldType = null))}
                >Cancel</Button
            >
            <Button onclick={handleAddField} disabled={!newField.name}>Add Field</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
