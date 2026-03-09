<script lang="ts">
    import { page } from "$app/state";
    import { fetchCollectionBySlug, fetchEntryWithFields, type EntryWithFields } from "$lib/data";
    import FieldInput from "$lib/components/field-types/FieldInput.svelte";
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Save, Eye, History, Globe, Trash2, ArrowLeft, X } from "@lucide/svelte";

    let collection = $state<EntryWithFields>({} as EntryWithFields);
    let entryData = $state<Record<string, any>>({});
    let loading = $state(true);

    $effect(() => {
        if (page.params.slug) {
            fetchEntryWithFields(page.params.slug).then((res) => {
                if (!res) return;
                collection = res;
                loading = false;
            });
        }
    });
</script>

<div>
    {#if loading}
        <div class="flex-1 flex items-center justify-center">
            <div class="flex flex-col items-center gap-4">
                <div class="w-10 h-10 border-2 border-primary/20 border-t-primary rounded-full animate-spin"></div>
                <p class="text-neutral-500 font-medium animate-pulse">Loading collection schema...</p>
            </div>
        </div>
    {:else if collection.id}
        <div class="px-6 py-4">
            <SectionHeader
                title={collection.title}
                subtitle="Create or edit content entry"
                breadcrumbs={["Workspace", "Entries", collection.title]}
            >
                {#snippet actions()}
                    <div class="flex items-center gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            class="bg-neutral-900 border-neutral-800 hover:bg-neutral-800"
                        >
                            <Eye size={16} class="mr-2" />
                            Preview
                        </Button>
                        <Button variant="default" size="sm" class="shadow-lg shadow-primary/20">
                            <Save size={16} class="mr-2" />
                            Save Entry
                        </Button>
                    </div>
                {/snippet}
            </SectionHeader>
        </div>

        <div class="flex-1 overflow-auto p-6 md:p-8 md:pt-0">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <div class="lg:col-span-2 space-y-8 pb-20">
                    <!-- <Card.Root class="bg-neutral-950 border-neutral-800 shadow-xl shadow-black/50 overflow-hidden">
                        <Card.Header>
                            <Card.Title class="text-2xl font-bold text-white tracking-tight">Main Content</Card.Title>
                            <Card.Description class="text-neutral-500">
                                Fill in the details for this {collection.title} entry.
                            </Card.Description>
                        </Card.Header>
                        <Card.Content class="space-y-5">
                            {#each collection.fields as field, index (field.id)}
                                <FieldInput
                                    field={field.label}
                                    type={field.type}
                                    options={{
                                        required: field.required,
                                        label: field.label,
                                        placeholder: field.default_value,
                                    }}
                                    bind:value={field.value}
                                />
                            {/each}
                        </Card.Content>
                    </Card.Root> -->

                    <!-- <div class="p-8 pb-4">
                        <h2 class="text-2xl font-bold text-white tracking-tight">Main Content</h2>
                        <p class="text-neutral-500">
                            Fill in the details for this {collection.title} entry.
                        </p>
                    </div> -->
                    <div class="space-y-5">
                        {#each collection.fields as field, index (field.id)}
                            <FieldInput
                                field={field.label}
                                type={field.type}
                                options={{
                                    required: field.required,
                                    label: field.label,
                                    placeholder: field.default_value,
                                }}
                                bind:value={field.value}
                            />
                        {/each}
                    </div>
                </div>

                <div class="space-y-6">
                    <Card.Root class="bg-neutral-950 border-neutral-800 overflow-hidden shadow-lg shadow-black/20">
                        <Card.Content class="space-y-5">
                            <div class="flex items-center justify-between">
                                <span class="text-xs text-neutral-500 flex items-center gap-2">
                                    <Globe size={14} /> Status
                                </span>
                                <Badge
                                    variant="secondary"
                                    class="bg-emerald-500/10 text-emerald-500 border-emerald-500/20 text-[10px] font-bold uppercase"
                                >
                                    Draft
                                </Badge>
                            </div>
                            <div class="flex items-center justify-between text-xs text-neutral-500">
                                <span class="flex items-center gap-2"><History size={14} /> Last saved</span>
                                <span>Just now</span>
                            </div>
                            <div class="pt-2 flex flex-col gap-2">
                                <Button
                                    variant="secondary"
                                    class="w-full bg-neutral-900 border border-neutral-800 hover:bg-neutral-800 text-xs h-9"
                                >
                                    Save as Draft
                                </Button>
                                <Button class="w-full text-xs h-9">Publish Now</Button>
                            </div>
                        </Card.Content>
                    </Card.Root>

                    <Card.Root class="bg-neutral-950 border-neutral-800 overflow-hidden shadow-lg shadow-black/20">
                        <Card.Header>
                            <Card.Title class="text-sm">SEO Preview</Card.Title>
                        </Card.Header>
                        <Card.Content class="pb-5 space-y-4">
                            <div
                                class="p-4 rounded-xl bg-neutral-900/30 border border-neutral-800 space-y-1.5 shadow-inner shadow-black/20"
                            >
                                <div class="text-blue-500 text-sm font-medium hover:underline cursor-pointer truncate">
                                    {entryData.title || entryData.name || "Untitled Entry"} | Vayload
                                </div>
                                <div class="text-emerald-600 text-[10px] truncate opacity-80">
                                    https://vayload.io/{collection.slug}/...
                                </div>
                                <div class="text-neutral-500 text-[11px] line-clamp-2 leading-relaxed">
                                    {entryData.description ||
                                        entryData.summary ||
                                        entryData.body ||
                                        "No description provided for SEO. Search engines will automatically generate a snippet..."}
                                </div>
                            </div>
                        </Card.Content>
                    </Card.Root>

                    <div class="flex gap-3">
                        <Button
                            variant="outline"
                            class="flex-1 bg-neutral-950 border-neutral-800 hover:bg-neutral-900 text-xs"
                        >
                            Duplicate
                        </Button>
                        <Button
                            variant="outline"
                            class="flex-1 bg-neutral-950 border-neutral-800 hover:bg-destructive/10 hover:border-destructive/20 text-neutral-500 hover:text-destructive transition-all text-xs"
                        >
                            <Trash2 size={14} class="mr-2" />
                            Delete
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    {:else}
        <div class="flex-1 flex flex-col items-center justify-center p-8 text-center space-y-4">
            <div class="p-6 rounded-full bg-neutral-900/50 border border-neutral-800">
                <X size={40} class="text-neutral-700" />
            </div>
            <div class="space-y-1">
                <h2 class="text-xl font-bold text-white tracking-tight">Collection Not Found</h2>
                <p class="text-neutral-500 max-w-sm text-sm">
                    The collection you are trying to access does not exist or has been deleted.
                </p>
            </div>
            <Button href="/entries" variant="outline" class="gap-2 mt-2 bg-neutral-900 border-neutral-800">
                <ArrowLeft size={16} />
                Back to Entries
            </Button>
        </div>
    {/if}
</div>
