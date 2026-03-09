<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import ScrollableTabs from "$lib/components/ScrollableTabs.svelte";
    import { entriesStore } from "$lib/stores/entriesStores.svelte";
    import { Check, FileText, DatabaseBackup, Clock, Search, Filter, Plus, MoreHorizontal } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import { Checkbox } from "$lib/components/ui/checkbox";
    import { goto } from "$app/navigation";

    let activeTab = $state("all");
    let searchQuery = $state("");

    $effect(() => {
        entriesStore.loadEntries(activeTab);
    });

    function handleTabChange(tabId: string) {
        activeTab = tabId;
    }

    function handleSearch(e: Event) {
        const target = e.target as HTMLInputElement;
        entriesStore.searchEntries(target.value);
    }

    const tabs = [
        { id: "all", label: "All Entries", count: entriesStore.countByStatus.all },
        {
            id: "published",
            label: "Published",
            count: entriesStore.countByStatus.published,
        },
        {
            id: "draft",
            label: "Drafts",
            count: entriesStore.countByStatus.draft,
        },
        {
            id: "archived",
            label: "Archived",
            count: entriesStore.countByStatus.archived,
        },
        {
            id: "scheduled",
            label: "Scheduled",
            count: entriesStore.countByStatus.scheduled,
        },
    ];

    function getBadgeColor(status: string) {
        switch (status) {
            case "published":
                return "default";
            case "draft":
                return "secondary";
            case "archived":
                return "outline";
            case "scheduled":
                return "default";
            default:
                return "secondary";
        }
    }
</script>

<div class="flex flex-col h-full">
    <SectionHeader
        title="Entries"
        subtitle="Manage your content across all collections."
        breadcrumbs={["Workspace", "Entries"]}
    >
        {#snippet actions()}
            <Button>
                <Plus size={16} class="mr-2" />
                Create New
            </Button>
        {/snippet}
    </SectionHeader>

    <ScrollableTabs {tabs} {activeTab} onTabChange={handleTabChange} />

    <div class="flex-1 overflow-y-auto">
        <div class="rounded-xl shadow-sm overflow-hidden">
            <div class="p-4 flex flex-wrap gap-4 justify-between items-center">
                <div class="flex items-center gap-2">
                    <div class="relative">
                        <Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-400" />
                        <Input
                            type="text"
                            placeholder="Search entries..."
                            class="pl-9 pr-4 w-64"
                            oninput={handleSearch}
                        />
                    </div>
                    <Button variant="outline" size="sm">
                        <Filter size={16} class="mr-2" />
                        Filters
                    </Button>
                </div>
                <div class="text-sm text-neutral-500">
                    Showing 1-{entriesStore.filteredEntries.length} of {entriesStore.totalCount}
                </div>
            </div>

            <Table.Root>
                <Table.Header>
                    <Table.Row>
                        <Table.Head class="w-12">
                            <Checkbox />
                        </Table.Head>
                        <Table.Head>Title</Table.Head>
                        <Table.Head>Content Type</Table.Head>
                        <Table.Head>Status</Table.Head>
                        <Table.Head>Last Updated</Table.Head>
                        <Table.Head class="text-right"></Table.Head>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {#if entriesStore.loading}
                        {#each Array(5) as _}
                            <Table.Row>
                                <Table.Cell colspan={6}>
                                    <div class="h-12animate-pulse rounded"></div>
                                </Table.Cell>
                            </Table.Row>
                        {/each}
                    {:else}
                        {#each entriesStore.filteredEntries as entry}
                            <Table.Row class="group cursor-pointer">
                                <Table.Cell>
                                    <Checkbox />
                                </Table.Cell>
                                <Table.Cell>
                                    <div class="text-sm font-medium text-neutral-300">{entry.title}</div>
                                    <div class="text-xs text-neutral-500">by {entry.author}</div>
                                </Table.Cell>
                                <Table.Cell class="text-sm text-neutral-500">{entry.content_type}</Table.Cell>
                                <Table.Cell>
                                    <Badge variant={getBadgeColor(entry.status)}>
                                        {entry.status}
                                    </Badge>
                                </Table.Cell>
                                <Table.Cell class="text-sm text-neutral-500">{entry.updated_at}</Table.Cell>
                                <Table.Cell class="text-right">
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        class="opacity-0 group-hover:opacity-100 transition-all"
                                        onclick={() => {
                                            console.log(entry);
                                            goto(`/entries/${entry.slug}`);
                                        }}
                                    >
                                        <MoreHorizontal size={18} />
                                    </Button>
                                </Table.Cell>
                            </Table.Row>
                        {/each}
                    {/if}
                </Table.Body>
            </Table.Root>

            <div class="p-4 flex justify-center">
                <Button variant="ghost" size="sm">Load More</Button>
            </div>
        </div>
    </div>
</div>
