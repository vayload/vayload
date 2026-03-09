<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import ScrollableTabs from "$lib/components/ScrollableTabs.svelte";
    import { fetchIntegrations, type Integration } from "$lib/data";
    import { CreditCard, Activity as AnalyticsIcon, Mail, DatabaseBackup } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import { onMount } from "svelte";

    let activeTab = $state("all");
    let integrations = $state<Integration[]>([]);
    let loading = $state(true);

    onMount(async () => {
        integrations = await fetchIntegrations();
        loading = false;
    });

    function handleTabChange(tabId: string) {
        activeTab = tabId;
    }

    const tabs = [
        { id: "all", label: "All Apps" },
        { id: "payment", label: "Payment" },
        { id: "analytics", label: "Analytics" },
        { id: "comm", label: "Communications" },
        { id: "storage", label: "Storage" },
    ];
</script>

<div class="flex flex-col h-full">
    <SectionHeader
        title="Integrations Marketplace"
        subtitle="Connect your project with third-party services."
        breadcrumbs={["System", "Integrations"]}
    />

    <ScrollableTabs {tabs} {activeTab} onTabChange={handleTabChange} />

    <div class="flex-1 py-6 overflow-y-auto">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#if loading}
                {#each Array(6) as _}
                    <Card.Root>
                        <Card.Content class="p-6">
                            <div class="h-32 bg-muted animate-pulse rounded"></div>
                        </Card.Content>
                    </Card.Root>
                {/each}
            {:else}
                {#each integrations as app}
                    <Card.Root class="flex flex-col justify-between hover:border-indigo-400 transition-colors">
                        <Card.Content class="p-6">
                            <div class="flex justify-between items-start mb-4">
                                <div
                                    class="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center text-gray-500 font-bold text-lg"
                                >
                                    {app.name[0]}
                                </div>
                                {#if app.installed}
                                    <Badge variant="default">Installed</Badge>
                                {/if}
                            </div>
                            <h4 class="font-bold text-gray-900">{app.name}</h4>
                            <p class="text-sm text-gray-500 mt-2">{app.description}</p>
                        </Card.Content>
                        <Card.Footer>
                            <Button variant={app.installed ? "outline" : "default"} class="w-full">
                                {app.installed ? "Configure" : "Install"}
                            </Button>
                        </Card.Footer>
                    </Card.Root>
                {/each}
            {/if}
        </div>
    </div>
</div>
