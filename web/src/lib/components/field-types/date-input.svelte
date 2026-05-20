<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";
    import { Calendar } from "@lucide/svelte";

    let { value = $bindable() } = $props<{ value: string }>();
</script>

<div class="relative">
    <Calendar class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500" size={16} />
    <Input type="date" bind:value class="pl-10 bg-neutral-900/50 border-neutral-800 focus:ring-primary/20" />
</div>
