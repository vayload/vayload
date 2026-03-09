<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import ScrollableTabs from "$lib/components/ScrollableTabs.svelte";
    import { Settings as SettingsIcon, Lock, CreditCard, Users, Code, Mail, Save } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Label } from "$lib/components/ui/label";
    import * as Alert from "$lib/components/ui/alert";
    import { appContext } from "$lib/stores/app-context.svelte";
    import EmptyProject from "$lib/components/empty-project.svelte";
    import AlertDialog from "$lib/components/alert-dialog.svelte";

    let activeTab = $state("general");
    const project = $derived(appContext.currentProject);

    function handleTabChange(tabId: string) {
        activeTab = tabId;
    }

    const tabs = [
        { id: "general", label: "General Details" },
        { id: "security", label: "Security & Auth" },
        { id: "billing", label: "Plan & Billing" },
        { id: "team", label: "Team Members" },
        { id: "api", label: "API Keys" },
        { id: "email", label: "Email Templates" },
    ];
</script>

<div class="flex flex-col h-full">
    <SectionHeader
        title="Project Settings"
        subtitle="Manage configuration, security, and billing."
        breadcrumbs={["System", "Settings"]}
    >
        {#snippet actions()}
            <Button>
                <Save size={16} class="mr-2" />
                Save Changes
            </Button>
        {/snippet}
    </SectionHeader>

    <ScrollableTabs {tabs} {activeTab} onTabChange={handleTabChange} />
    {#if !project}
        <EmptyProject />
    {:else}
        <div class="flex-1 overflow-y-auto py-6">
            <h2 class="text-xl font-bold text-foreground mb-6 capitalize">{activeTab} Configuration</h2>

            <div class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-2">
                        <Label for="project-name">Project Name</Label>
                        <Input id="project-name" value={project.name} />
                    </div>
                    <div class="space-y-2">
                        <Label for="project-id">Project ID</Label>
                        <Input id="project-id" value={project.id} disabled class="bg-muted text-muted-foreground" />
                    </div>
                </div>

                <div class="space-y-2">
                    <Label for="domain">Primary Domain</Label>
                    <div class="flex">
                        <span
                            class="inline-flex items-center px-3 rounded-l-lg border border-r-0 bg-muted text-muted-foreground text-sm"
                        >
                            https://
                        </span>
                        <Input id="domain" value="api.ecommerce.com" class="rounded-l-none" />
                    </div>
                    <p class="text-xs text-muted-foreground mt-2">
                        This domain will be used for your public API endpoints.
                    </p>
                </div>

                <div class="pt-6 border-t border-gray-100">
                    <h3 class="text-sm font-medium text-foreground mb-4">Danger Zone</h3>
                    <Alert.Root variant="destructive">
                        <Alert.Title>Delete Project</Alert.Title>
                        <Alert.Description class="flex items-center justify-between">
                            <span class="text-sm">This action cannot be undone.</span>
                            <AlertDialog
                                title="Delete Project"
                                description="This action cannot be undone. This will permanently delete your account and remove your data from our servers."
                                cancelText="Cancel"
                                actionText="Continue"
                                onCancel={() => {}}
                                onAction={() => {
                                    appContext.deleteProject(project.id);
                                }}
                            />
                        </Alert.Description>
                    </Alert.Root>
                </div>
            </div>
        </div>
    {/if}
</div>
