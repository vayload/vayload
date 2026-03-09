<!-- <script lang="ts">
    import { page } from "$app/state";
    import ContentTypeSelector from "$lib/components/field-types/content-type-selector.svelte";
    import { fetchCollectionBySlug } from "$lib/data";

    let collection = $state<any>(null);
    let loading = $state(true);

    $effect(() => {
        if (page.params.slug) {
            fetchCollectionBySlug(page.params.slug).then((res) => {
                collection = res;
                loading = false;
            });
        }
    });

    const onUpdate = () => {
        console.log("update");
    };
</script>

<div>
    {#if loading}
        <p>Loading...</p>
    {:else}
        <ContentTypeSelector {collection} {onUpdate} />
    {/if}
</div> -->

<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { Select, SelectContent, SelectItem, SelectTrigger } from "$lib/components/ui/select";
    import { Button } from "$lib/components/ui/button";
    import {
        Plus,
        Trash2,
        Settings2,
        ArrowLeft,
        Files,
        Text,
        Check,
        Calendar,
        Link,
        Palette,
        MapPin,
        Code,
        FileText,
        Image,
    } from "@lucide/svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { FieldTypes, FieldTypeMetadata } from "$lib/constants/field-types";
    import { goto } from "$app/navigation";

    import { page } from "$app/state";
    import { fetchCollectionBySlug, type Collection, type FieldSchema } from "$lib/data";
    import FieldTypeGrid from "$lib/components/field-types/FieldTypeGrid.svelte";
    import Number from "@tabler/icons-svelte/icons/number";

    let collection = $state<Collection>({} as any);
    let loading = $state(true);
    let status = $state({
        error: false,
        message: "",
    });

    $effect(() => {
        if (page.params.slug) {
            fetchCollectionBySlug(page.params.slug).then((res) => {
                if (!res) {
                    status.error = true;
                    status.message = "Collection not found";

                    return;
                }

                collection = res;

                loading = false;
            });
        }
    });

    const onUpdate = (data: any) => {
        console.log("update", data);
    };

    let isAddingField = $state(false);
    let selectedFieldType: FieldTypes | null = $state(null);

    let newField = $state<FieldSchema>({
        type: FieldTypes.TEXT,
        required: false,
        label: "",
        default_value: null,
        relation_to: "",
        config: {},
    });

    function openAddField(type: FieldTypes) {
        selectedFieldType = type;
        newField = {
            type: type as any,
            required: false,
            label: "",
            default_value: null,
            relation_to: "",
            config: {},
        };

        isAddingField = true;
    }

    function handleAddField() {
        if (!newField.label || !collection) return;

        if (!collection.fields_schema) {
            collection.fields_schema = {};
        }

        const fieldName = newField.label.toLowerCase().replace(/[^a-z0-9_]/g, "_");
        if (collection.fields_schema[fieldName]) {
            status.error = true;
            status.message = "Field already exists";

            return;
        }

        collection.fields_schema[fieldName] = {
            type: newField.type as any,
            required: newField.required,
            label: newField.label,
            default_value: newField.default_value,
            relation_to: newField.relation_to,
            config: newField.config,
        };

        onUpdate(collection);
        isAddingField = false;
        selectedFieldType = null;
    }

    function removeField(fieldName: string) {
        if (!collection || !collection.fields_schema) return;

        delete collection.fields_schema[fieldName];

        onUpdate(collection);
    }

    const iconForFieldType: Record<FieldTypes, any> = {
        [FieldTypes.TEXT]: Text,
        [FieldTypes.NUMBER]: Number,
        [FieldTypes.BOOLEAN]: Check,
        [FieldTypes.DATE]: Calendar,
        [FieldTypes.RELATIONSHIP]: Link,
        [FieldTypes.MEDIA]: Image,
        [FieldTypes.TONES]: Palette,
        [FieldTypes.LOCATION]: MapPin,
        [FieldTypes.JSON]: Code,
        [FieldTypes.RICH_TEXT]: FileText,
    };
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
                        {#if Object.keys(collection.fields_schema || {}).length === 0}
                            <div class="text-center py-12 border-2 border-dashed rounded-xl">
                                <p class="text-muted-foreground">No fields defined yet.</p>
                            </div>
                        {:else}
                            <div class="border rounded-xl divide-y">
                                {#each Object.entries(collection.fields_schema || {}) as [name, field]}
                                    <div
                                        class="p-4 flex items-center justify-between hover:bg-muted/30 transition-colors"
                                    >
                                        <div class="flex items-center gap-4">
                                            <div class="p-2 rounded-lg bg-muted text-primary">
                                                {#if iconForFieldType[field.type]}
                                                    {@render iconForFieldType[field.type]({ size: 18 })}
                                                {:else}
                                                    <Settings2 size={18} />
                                                {/if}
                                            </div>
                                            <div>
                                                <div class="flex items-center gap-2">
                                                    <span class="font-bold">{field.label || name}</span>
                                                    {#if field.required}
                                                        <span
                                                            class="text-[10px] font-bold uppercase tracking-wider bg-red-100 text-red-600 px-1.5 py-0.5 rounded"
                                                            >Required</span
                                                        >
                                                    {/if}
                                                </div>
                                                <p class="text-xs font-mono text-muted-foreground mt-0.5">
                                                    {field.label} • {field.type}
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
                                                onclick={() => removeField(name)}
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
                    <!-- <div class="space-y-2">
                        <Label>Description</Label>
                        <Input bind:value={collection.description} oninput={() => onUpdate(collection as any)} />
                    </div> -->

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
