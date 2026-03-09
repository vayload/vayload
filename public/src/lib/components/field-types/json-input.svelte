<script lang="ts">
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { Code } from "@lucide/svelte";

    let { value = $bindable(), placeholder = "" } = $props<{
        value: string;
        placeholder?: string;
    }>();
</script>

<div class="relative group">
    <Code class="absolute right-3 top-3 text-neutral-600 group-hover:text-primary transition-colors" size={16} />
    <Textarea
        bind:value
        class="font-mono text-xs bg-neutral-900/50 border-neutral-800 min-h-[120px]"
        placeholder={placeholder || '{ "key": "value" }'}
    />
</div>
