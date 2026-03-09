<script lang="ts">
    import { Label } from "$lib/components/ui/label/index.js";
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { FieldTypes } from "$lib/constants/field-types";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Image as ImageIcon } from "@lucide/svelte";

    // Import modular inputs
    import BooleanInput from "./boolean-input.svelte";
    import TextInput from "./text-input.svelte";
    import RichTextInput from "./rich-text-input.svelte";
    import MediaInput from "./media-input.svelte";
    import RelationshipInput from "./relationship-input.svelte";
    import DateInput from "./date-input.svelte";
    import JsonInput from "./json-input.svelte";
    import LocationInput from "./location-input.svelte";
    import TonesInput from "./tones-input.svelte";

    let {
        field,
        value = $bindable(),
        options = {},
        type = "text",
    } = $props<{
        field: string;
        value: any;
        options?: any;
        type?: string;
    }>();

    // Media gallery state (shared by RichText and Media inputs)
    let showMediaGallery = $state(false);
    const mockImages = [
        "https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=500&auto=format&fit=crop&q=60",
        "https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=500&auto=format&fit=crop&q=60",
        "https://images.unsplash.com/photo-1488590528505-98d2b5aba04b?w=500&auto=format&fit=crop&q=60",
        "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=500&auto=format&fit=crop&q=60",
    ];

    function selectImage(url: string) {
        if (type === FieldTypes.MEDIA) {
            value = url;
        } else {
            // For rich text, we'd normally insert into the editor
            // For now just append to value or log
            value += `\n![Image](${url})`;
        }
        showMediaGallery = false;
    }

    const label = options.label || field;
    const placeholder = options.placeholder || `Enter ${label.toLowerCase()}...`;
</script>

<div class="space-y-2">
    <div class="flex items-center justify-between">
        <Label class="text-sm font-medium text-neutral-300">{label}</Label>
        {#if options.required}
            <Badge variant="outline" class="text-[10px] h-4 px-1 text-neutral-500 border-neutral-800">Required</Badge>
        {/if}
    </div>

    {#if type === FieldTypes.BOOLEAN}
        <BooleanInput bind:value />
    {:else if type === FieldTypes.TEXT || type === "string" || type === "text"}
        <TextInput bind:value {placeholder} />
    {:else if type === FieldTypes.NUMBER}
        <TextInput bind:value {placeholder} type="number" />
    {:else if type === FieldTypes.DATE}
        <DateInput bind:value />
    {:else if type === FieldTypes.RICH_TEXT || type === "richtext" || type === "rich_text"}
        <RichTextInput bind:value {placeholder} onOpenMedia={() => (showMediaGallery = true)} />
    {:else if type === FieldTypes.MEDIA}
        <MediaInput bind:value onOpenMedia={() => (showMediaGallery = true)} />
    {:else if type === FieldTypes.RELATIONSHIP || type === "relation"}
        <RelationshipInput bind:value />
    {:else if type === FieldTypes.JSON}
        <JsonInput bind:value {placeholder} />
    {:else if type === FieldTypes.LOCATION}
        <LocationInput bind:value />
    {:else if type === FieldTypes.TONES}
        <TonesInput bind:value />
    {:else}
        <TextInput bind:value {placeholder} />
    {/if}

    {#if options.description}
        <p class="text-[11px] text-neutral-500">{options.description}</p>
    {/if}
</div>

<!-- Global Media Gallery for this field's context -->
<Dialog.Root bind:open={showMediaGallery}>
    <Dialog.Content class="sm:max-w-[700px] bg-neutral-950 border-neutral-800 p-0 overflow-hidden">
        <div class="p-6 border-b border-neutral-800 bg-neutral-900/20">
            <h3 class="text-lg font-semibold flex items-center gap-2 text-white">
                <ImageIcon class="text-primary" size={20} />
                Asset Library
            </h3>
            <p class="text-sm text-neutral-500">Select an image to insert into your content</p>
        </div>

        <div class="p-6 bg-neutral-950">
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
                {#each mockImages as img}
                    <button
                        onclick={() => selectImage(img)}
                        class="relative aspect-square rounded-lg overflow-hidden border-2 border-transparent hover:border-primary transition-all group active:scale-95"
                    >
                        <img
                            src={img}
                            alt="Gallery item"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                        />
                        <div
                            class="absolute inset-0 bg-primary/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                        >
                            <span class="bg-primary text-white text-[10px] font-bold px-2 py-1 rounded">SELECT</span>
                        </div>
                    </button>
                {/each}
                <button
                    class="aspect-square rounded-lg border-2 border-dashed border-neutral-800 flex flex-col items-center justify-center gap-2 text-neutral-500 hover:border-neutral-600 hover:text-neutral-400 hover:bg-neutral-900/50 transition-all group"
                >
                    <div class="p-2 rounded-full bg-neutral-900 group-hover:bg-neutral-800 transition-colors">
                        <ImageIcon size={20} />
                    </div>
                    <span class="text-[10px] uppercase font-bold tracking-wider">Upload</span>
                </button>
            </div>
        </div>

        <div class="p-4 bg-neutral-900/50 flex justify-end gap-3 border-t border-neutral-800">
            <button
                onclick={() => (showMediaGallery = false)}
                class="px-4 py-2 text-sm font-medium text-neutral-400 hover:text-white transition-colors"
            >
                Cancel
            </button>
        </div>
    </Dialog.Content>
</Dialog.Root>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: #262626;
        border-radius: 10px;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb:hover {
        background: #404040;
    }
</style>
