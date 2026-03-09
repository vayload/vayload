<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
</script>

<div class="px-6 py-4">
    <SectionHeader
        title="Live Feed"
        subtitle="Real-time activity and platform events"
        breadcrumbs={["Workspace", "Live Feed"]}
    />
</div>

<div class="p-6 pt-0 flex-1 flex items-center justify-center">
    <div class="text-center space-y-2">
        <h3 class="text-lg font-semibold text-white">Activity Monitoring</h3>
        <p class="text-sm text-neutral-500 max-w-sm">
            This page will display real-time events, logs, and platform activity as they happen.
        </p>
    </div>
</div>
