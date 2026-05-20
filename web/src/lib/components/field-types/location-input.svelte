<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
    import { MapPin } from "@lucide/svelte";

    let { value = $bindable() } = $props<{ value: { lat: number; lng: number } }>();
</script>

<div class="grid grid-cols-2 gap-3">
    <div class="space-y-1.5">
        <Label class="text-[10px] uppercase text-neutral-500 font-bold tracking-tight">Latitude</Label>
        <div class="relative">
            <MapPin class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-600" size={14} />
            <Input
                type="number"
                step="0.000001"
                bind:value={value.lat}
                placeholder="0.000000"
                class="pl-9 bg-neutral-900/50 border-neutral-800"
            />
        </div>
    </div>
    <div class="space-y-1.5">
        <Label class="text-[10px] uppercase text-neutral-500 font-bold tracking-tight">Longitude</Label>
        <div class="relative">
            <MapPin class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-600" size={14} />
            <Input
                type="number"
                step="0.000001"
                bind:value={value.lng}
                placeholder="0.000000"
                class="pl-9 bg-neutral-900/50 border-neutral-800"
            />
        </div>
    </div>
</div>
