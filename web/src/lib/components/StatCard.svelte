<script lang="ts">
    import { type Snippet } from "svelte";

    interface Props {
        label: string;
        value: string;
        change: string;
        trend: "up" | "down";
        icon: Snippet;
        color: string;
    }

    let { label, value, change, trend, icon, color }: Props = $props();
</script>

<div class="bg-card p-6 rounded-xl border shadow-sm hover:shadow-md transition-all">
    <div class="flex items-center justify-between mb-4">
        <span class="text-sm font-medium text-muted-foreground">{label}</span>
        <div class="p-2 rounded-lg {color}">
            {@render icon()}
        </div>
    </div>
    <div class="flex items-baseline gap-2">
        <span class="text-2xl font-bold text-foreground">{value}</span>
        <span class="text-xs font-medium {trend === 'up' ? 'text-green-600' : 'text-muted-foreground'}">
            {change}
        </span>
    </div>
</div>
