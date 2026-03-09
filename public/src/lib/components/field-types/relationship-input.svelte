<script lang="ts">
    import { Badge } from "$lib/components/ui/badge/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { X } from "@lucide/svelte";

    let { value = $bindable() } = $props<{ value: any[] }>();
</script>

<div class="flex flex-wrap gap-2 p-3 border border-neutral-800 rounded-xl bg-neutral-900/30 min-h-[46px]">
    {#if value && Array.isArray(value)}
        {#each value as item}
            <Badge variant="secondary" class="gap-1 bg-neutral-800 text-neutral-300 hover:bg-neutral-700">
                {item}
                <button onclick={() => (value = value.filter((v: any) => v !== item))} class="hover:text-white">
                    <X size={10} />
                </button>
            </Badge>
        {/each}
    {/if}
    <Button
        variant="ghost"
        size="sm"
        class="h-6 px-2 text-[10px] text-neutral-500 hover:text-neutral-300 uppercase font-bold tracking-wider"
    >
        + Link Entry
    </Button>
</div>
