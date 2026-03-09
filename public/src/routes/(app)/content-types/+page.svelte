<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import ContentTypeCard from "$lib/components/ContentTypeCard.svelte";
    import { fetchCollections, type Collection } from "$lib/data";
    import { FileText, Tags, Users, Filter, Globe, Search, Plus } from "@lucide/svelte";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import type { Component } from "svelte";

    import { FieldTypes } from "$lib/constants/field-types";

    let collections = $state<(Collection & { icon: Component })[]>([]);
    let loading = $state(true);
    let isDialogOpen = $state(false);

    let contentType = $state({
        name: "",
        description: "",
        single: false,
    });

    const defaultFields = [
        { name: "index", type: FieldTypes.TEXT, required: true },
        { name: "title", type: FieldTypes.TEXT, required: true },
        { name: "description", type: FieldTypes.TEXT, required: false },
        { name: "meta", type: FieldTypes.JSON, label: "SEO Meta", required: false },
        { name: "json_schema", type: FieldTypes.JSON, label: "Org JSON Schema", required: false },
    ];

    $effect(() => {
        fetchCollections().then((res) => {
            collections = res.map((collection) => ({
                ...collection,
                icon: matchIconByType(collection.slug),
            }));

            loading = false;
        });
    });

    function handleCreate() {
        if (!contentType.name) return;

        const newCollection: any = {
            id: crypto.randomUUID(),
            name: contentType.name,
            slug: contentType.name.toLowerCase().replace(/\s/g, "_"),
            description: contentType.description,
            single: contentType.single,
            fields: defaultFields.length,
            entries: 0,
            fields_schema: defaultFields.reduce((acc: any, field) => {
                acc[field.name] = { type: field.type, required: field.required, label: field.label || field.name };
                return acc;
            }, {}),
            created_at: new Date().toISOString(),
            icon: matchIconByType(contentType.name.toLowerCase()),
        };

        collections = [...collections, newCollection];
        isDialogOpen = false;

        // Reset form
        contentType = {
            name: "",
            description: "",
            single: false,
        };
    }

    const matchIconByType = (type: string) => {
        if (type.includes("post")) return FileText;
        if (type.includes("product")) return Tags;
        if (type.includes("author")) return Users;
        if (type.includes("category")) return Filter;
        if (type.includes("page")) return Globe;
        if (type.includes("seo")) return Search;
        return FileText;
    };
</script>

<div class="pb-8">
    <SectionHeader
        title="Content Types"
        subtitle="Define the schema and structure of your content."
        breadcrumbs={["Workspace", "Content Types"]}
    >
        {#snippet actions()}
            <Button onclick={() => (isDialogOpen = true)}>
                <Plus size={16} class="mr-2" />
                Create Type
            </Button>
        {/snippet}
    </SectionHeader>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#if loading}
            {#each Array(6) as _}
                <div class="bg-card p-6 rounded-xl border animate-pulse">
                    <div class="h-12 bg-muted rounded mb-4"></div>
                    <div class="h-4 bg-muted rounded w-3/4"></div>
                </div>
            {/each}
        {:else}
            {#each collections as collection}
                <ContentTypeCard
                    {collection}
                    fieldCount={Object.keys(collection.fields_schema).length}
                    entryCount={collection.entries}
                >
                    {#snippet icon()}
                        <collection.icon size={24} />
                    {/snippet}
                </ContentTypeCard>
            {/each}
        {/if}
    </div>
</div>

<Dialog.Root bind:open={isDialogOpen}>
    <Dialog.Content>
        <Dialog.Header>
            <Dialog.Title>Create Content Type</Dialog.Title>
        </Dialog.Header>
        <div class="space-y-4 py-4">
            <div class="space-y-2">
                <Label for="name">Display Name</Label>
                <Input id="name" placeholder="e.g. Blog Post" bind:value={contentType.name} />
            </div>
            <div class="space-y-2">
                <Label for="api-id">API ID</Label>
                <Input
                    id="api-id"
                    placeholder="blog_post"
                    disabled
                    value={contentType.name.toLowerCase().replace(/\s/g, "_")}
                    class="font-mono text-sm"
                />
                <p class="text-xs text-gray-500">Generated automatically from display name.</p>
            </div>
            <div class="flex items-center gap-2">
                <Checkbox id="single" bind:checked={contentType.single} />
                <Label for="single" class="text-sm font-normal">This is a Single Type (e.g. Homepage)</Label>
            </div>
        </div>
        <Dialog.Footer>
            <Button variant="outline" onclick={() => (isDialogOpen = false)}>Cancel</Button>
            <Button onclick={handleCreate}>Create Type</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>
