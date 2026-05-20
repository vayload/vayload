<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
</script>

<div class="px-6 py-4">
    <SectionHeader
        title="Inbox"
        subtitle="Manage your messages and notifications"
        breadcrumbs={["Workspace", "Inbox"]}
    />
</div>

<div class="p-6 pt-0 flex-1 flex items-center justify-center">
    <div class="text-center space-y-2">
        <h3 class="text-lg font-semibold text-white">No Messages</h3>
        <p class="text-sm text-neutral-500 max-w-sm">
            Your inbox is clear. New notifications and system messages will appear here.
        </p>
    </div>
</div>
