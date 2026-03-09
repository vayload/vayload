<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
</script>

<div class="px-6 py-4">
    <SectionHeader
        title="Triggers"
        subtitle="Define events that start your workflows"
        breadcrumbs={["Workspace", "Triggers"]}
    />
</div>

<div class="p-6 pt-0 flex-1 flex items-center justify-center">
    <div class="text-center space-y-2">
        <h3 class="text-lg font-semibold text-white">Event Triggers</h3>
        <p class="text-sm text-neutral-500 max-w-sm">
            Configure webhooks, scheduled tasks, and system events to trigger actions.
        </p>
    </div>
</div>
