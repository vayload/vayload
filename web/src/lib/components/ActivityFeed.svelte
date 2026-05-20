<script lang="ts">
    import type { Activity } from "$lib/data";

    interface Props {
        activities: Activity[];
    }

    let { activities }: Props = $props();
</script>

<div class="bg-card rounded-xl border shadow-sm p-6">
    <h3 class="text-lg font-semibold text-foreground mb-6">Activity Feed</h3>
    <div
        class="space-y-6 relative before:absolute before:left-2 before:top-2 before:bottom-2 before:w-0.5 before:bg-border"
    >
        {#each activities as item}
            <div class="relative pl-8">
                <div
                    class="absolute left-0 top-1.5 w-4 h-4 rounded-full border-2 border-card bg-primary shadow-sm z-10"
                ></div>
                <p class="text-sm text-foreground">
                    <span class="font-medium">{item.user}</span>
                    {item.action}
                    <span class="font-medium">{item.target}</span>
                </p>
                <p class="text-xs text-muted-foreground mt-0.5">{item.time}</p>
            </div>
        {/each}
    </div>
</div>
