<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Image as ImageIcon } from "@lucide/svelte";

    let { value = $bindable(), onOpenMedia = () => {} } = $props<{
        value: string | null;
        onOpenMedia?: () => void;
    }>();
</script>

<div class="flex flex-col gap-3">
    {#if value}
        <div class="relative group w-full aspect-video rounded-xl overflow-hidden border border-neutral-800 bg-black">
            <img src={value} alt="Media preview" class="w-full h-full object-contain" />
            <div
                class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2"
            >
                <Button variant="secondary" size="sm" onclick={onOpenMedia}>Replace</Button>
                <Button variant="destructive" size="sm" onclick={() => (value = null)}>Remove</Button>
            </div>
        </div>
    {:else}
        <button
            onclick={onOpenMedia}
            class="w-full aspect-video rounded-xl border border-dashed border-neutral-800 hover:border-neutral-600 hover:bg-neutral-900/50 transition-all flex flex-col items-center justify-center gap-2 group"
        >
            <div class="p-3 rounded-full bg-neutral-900 group-hover:bg-neutral-800 transition-colors">
                <ImageIcon size={24} class="text-neutral-500" />
            </div>
            <span class="text-sm text-neutral-400">Select or upload media</span>
        </button>
    {/if}
</div>
