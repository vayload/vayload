<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import StatCard from "$lib/components/StatCard.svelte";
    import ActivityFeed from "$lib/components/ActivityFeed.svelte";
    import { fetchActivities, type Activity } from "$lib/data";
    import { Activity as ActivityIcon, Users, FileText, Clock } from "@lucide/svelte";
    import { onMount } from "svelte";
    import Chart from "$lib/components/chart.svelte";

    let activities = $state<Activity[]>([]);
    let loading = $state(true);

    onMount(async () => {
        activities = await fetchActivities();
        loading = false;
    });

    const chartData = [40, 70, 45, 90, 60, 80, 50, 40, 70, 45, 90, 60, 75, 50, 85];
</script>

<div class="pb-8">
    <SectionHeader
        title="Overview"
        subtitle="Welcome back, Ana. Here's what's happening today."
        breadcrumbs={["Dashboard", "Overview"]}
    />

    <div class="space-y-8">
        <!-- Stats Row -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard
                label="Total Revenue"
                value="$45,231"
                change="+20.1%"
                trend="up"
                color="text-primary bg-primary/10"
            >
                {#snippet icon()}
                    <ActivityIcon size={18} />
                {/snippet}
            </StatCard>
            <StatCard label="Active Users" value="+2,350" change="+18.1%" trend="up" color="text-primary bg-primary/10">
                {#snippet icon()}
                    <Users size={18} />
                {/snippet}
            </StatCard>
            <StatCard label="Entries" value="12,234" change="+19%" trend="up" color="text-primary bg-primary/10">
                {#snippet icon()}
                    <FileText size={18} />
                {/snippet}
            </StatCard>
            <StatCard label="Avg. Response" value="24ms" change="-4%" trend="down" color="text-primary bg-primary/10">
                {#snippet icon()}
                    <Clock size={18} />
                {/snippet}
            </StatCard>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div class="lg:col-span-2">
                <Chart />
            </div>

            {#if loading}
                <div class="bg-card rounded-xl border shadow-sm p-6">
                    <div class="animate-pulse space-y-4">
                        <div class="h-4 bg-muted rounded w-1/2"></div>
                        <div class="h-4 bg-muted rounded"></div>
                        <div class="h-4 bg-muted rounded"></div>
                    </div>
                </div>
            {:else}
                <ActivityFeed {activities} />
            {/if}
        </div>
    </div>
</div>
