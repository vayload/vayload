<script lang="ts">
    import { type Snippet } from "svelte";
    import * as Tabs from "$lib/components/ui/tabs";

    interface Tab {
        id: string;
        label: string;
        count?: number;
        icon?: Snippet;
    }

    interface Props {
        tabs: Tab[];
        activeTab: string;
        onTabChange: (id: string) => void;
    }

    let { tabs, activeTab, onTabChange }: Props = $props();
</script>

<Tabs.Root value={activeTab} onValueChange={onTabChange} class="w-full">
    <Tabs.List>
        {#each tabs as tab}
            {@const isActive = activeTab === tab.id}
            <Tabs.Trigger value={tab.id}>
                <span class="flex items-center gap-2">
                    {#if tab.icon}
                        {@render tab.icon()}
                    {/if}
                    {tab.label}
                    {#if tab.count !== undefined}
                        <span class="text-xs px-2 py-0.5 rounded-full">
                            {tab.count}
                        </span>
                    {/if}
                </span>
            </Tabs.Trigger>
        {/each}
    </Tabs.List>
</Tabs.Root>
