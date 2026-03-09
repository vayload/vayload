<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
</script>

<div class="px-6 py-4">
    <SectionHeader
        title="Orders"
        subtitle="Monitor transactions and order status"
        breadcrumbs={["Workspace", "Orders"]}
    />
</div>

<div class="p-6 pt-0 flex-1 flex items-center justify-center">
    <div class="text-center space-y-2">
        <h3 class="text-lg font-semibold text-white">Order Tracking</h3>
        <p class="text-sm text-neutral-500 max-w-sm">
            View, process, and manage customer orders and fulfillment status in real-time.
        </p>
    </div>
</div>
