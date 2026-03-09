<script lang="ts">
    import type { Collection } from "$lib/data";
    import { type Snippet } from "svelte";
    import { Settings } from "@lucide/svelte";

    interface Props {
        collection: Collection;
        icon: Snippet;
        fieldCount: number;
        entryCount: number;
        isSingle?: boolean;
    }

    let { collection, icon, fieldCount, entryCount, isSingle = false }: Props = $props();
</script>

<div
    class="group bg-card p-6 rounded-xl border hover:border-primary hover:shadow-md transition-all cursor-pointer relative overflow-hidden"
>
    <div class="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
        <a class="text-muted-foreground hover:text-primary" href={`/content-types/${collection?.slug}`}>
            <Settings size={18} />
        </a>
    </div>

    <div class="flex items-start gap-4 mb-4">
        <div
            class="p-3 bg-primary/10 text-primary rounded-lg group-hover:bg-primary group-hover:text-primary-foreground transition-colors"
        >
            {@render icon()}
        </div>
        <div>
            <h3 class="font-bold text-foreground">{collection?.name}</h3>
            <code class="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded mt-1 inline-block">
                {collection?.slug}
            </code>
        </div>
    </div>

    <div class="flex items-center justify-between text-sm text-muted-foreground mt-6 pt-4 border-t">
        <span>{fieldCount} Fields</span>
        <span class="flex items-center gap-1">
            {entryCount}
            {isSingle ? "Instance" : "Entries"}
        </span>
    </div>
</div>
