<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
</script>

<div class="px-6 py-4">
    <SectionHeader
        title="Inventory"
        subtitle="Track stock and product availability"
        breadcrumbs={["Workspace", "Inventory"]}
    />
</div>

<div class="p-6 pt-0 flex-1 flex items-center justify-center">
    <div class="text-center space-y-2">
        <h3 class="text-lg font-semibold text-white">Stock Management</h3>
        <p class="text-sm text-neutral-500 max-w-sm">
            Monitor inventory levels across your collections and set up low-stock alerts.
        </p>
    </div>
</div>
