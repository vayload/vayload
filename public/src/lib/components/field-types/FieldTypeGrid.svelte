<script lang="ts">
    import { FieldTypes, FieldTypeMetadata } from "$lib/constants/field-types";
    import * as LucideIcons from "@lucide/svelte";
    import { Card, CardContent } from "$lib/components/ui/card";
    import type { Component } from "svelte";

    export let onSelect: (type: FieldTypes) => void = () => {};

    const fieldTypesList = Object.entries(FieldTypeMetadata).map(([type, meta]) => ({
        type: type as FieldTypes,
        ...meta,
        iconComponent: (LucideIcons as any)[meta.icon] as Component,
    }));
</script>

<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
    {#each fieldTypesList as fieldType}
        <button
            type="button"
            class="text-left transition-transform hover:scale-[1.02] active:scale-[0.98]"
            onclick={() => onSelect(fieldType.type)}
        >
            <Card class="h-full hover:border-primary hover:bg-primary/5 transition-colors cursor-pointer">
                <CardContent class="p-4 flex flex-col items-center text-center space-y-3">
                    <div class="p-3 rounded-full bg-muted text-primary">
                        <svelte:component this={fieldType.iconComponent} size={24} />
                    </div>
                    <div>
                        <h4 class="font-semibold text-sm">{fieldType.label}</h4>
                        <p class="text-xs text-muted-foreground mt-1 line-clamp-2">
                            {fieldType.description}
                        </p>
                    </div>
                </CardContent>
            </Card>
        </button>
    {/each}
</div>
