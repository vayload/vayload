<script lang="ts">
    import { Switch } from "$lib/components/ui/switch/index.js";
    import { Label } from "$lib/components/ui/label/index.js";

    let { value = $bindable() } = $props<{ value: boolean }>();
</script>

<div class="flex items-center space-x-3 p-3 rounded-xl border border-neutral-800 bg-neutral-900/30">
    <Switch bind:checked={value} />
    <span class="text-sm text-neutral-400">{value ? "Enabled" : "Disabled"}</span>
</div>
