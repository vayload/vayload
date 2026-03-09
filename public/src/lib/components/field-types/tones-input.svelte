<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";
    import { Palette } from "@lucide/svelte";

    let { value = $bindable() } = $props<{ value: string }>();
</script>

<div class="flex flex-wrap items-center gap-3 p-3 rounded-xl border border-neutral-800 bg-neutral-900/30">
    {#each ["#ef4444", "#f59e0b", "#10b981", "#3b82f6", "#8b5cf6", "#ec4899"] as color}
        <button
            onclick={() => (value = color)}
            class="w-7 h-7 rounded-full border-2 transition-all hover:scale-110 shadow-sm {value === color
                ? 'border-white ring-2 ring-primary/20'
                : 'border-transparent'}"
            style="background-color: {color}"
            title={color}
        ></button>
    {/each}
    <div class="w-px h-6 bg-neutral-800 mx-1"></div>
    <div class="relative">
        <Palette class="absolute left-2.5 top-1/2 -translate-y-1/2 text-neutral-500" size={14} />
        <Input
            type="text"
            bind:value
            class="pl-8 w-28 h-8 text-xs bg-neutral-900 border-neutral-800 uppercase"
            placeholder="#000000"
        />
    </div>
</div>
