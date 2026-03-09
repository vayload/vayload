<script lang="ts">
    import Sidebar from "$lib/components/Sidebar.svelte";
    import Header from "$lib/components/Header.svelte";
    import CreateProjectModal from "$lib/components/CreateProjectModal.svelte";
    import { appContext } from "$lib/stores/app-context.svelte";
    import { userStore } from "$lib/stores/userStore.svelte";
    import { notificationStore } from "$lib/stores/notificationStore.svelte";
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import type { Snippet } from "svelte";
    import { fly } from "svelte/transition";

    interface Props {
        children: Snippet;
    }

    let { children }: Props = $props();

    let isMobileMenuOpen = $state(false);
    let showCreateProject = $state(false);

    onMount(async () => {
        await Promise.all([
            appContext.fetchProjects(),
            userStore.loadCurrentUser(),
            notificationStore.loadNotifications(),
        ]);
    });

    $effect(() => {
        if (appContext.needsOnboarding) {
            showCreateProject = true;
        }
    });

    const getPageTitle = (pathname: string): string => {
        const map: Record<string, string> = {
            "/dashboard": "Dashboard",
            "/dashboard/content-types": "Content Type Builder",
            "/dashboard/entries": "Content Entries",
            "/dashboard/media": "Media Library",
            "/dashboard/users": "User Management",
            "/dashboard/roles": "Roles & Permissions",
            "/dashboard/integrations": "Integrations",
            "/dashboard/settings": "Project Settings",
            "/dashboard/audit": "Audit Logs",
        };
        return map[pathname] || "Dashboard";
    };

    let pageTitle = $derived(getPageTitle(page.url.pathname));

    let ready = $derived(appContext.initialized && !appContext.needsOnboarding && appContext.currentProject !== null);
</script>

<CreateProjectModal bind:open={showCreateProject} closable={!appContext.needsOnboarding} />

<div class="flex h-screen font-sans overflow-hidden text-neutral-300 bg-background">
    <Sidebar bind:isMobileOpen={isMobileMenuOpen} />
    <div class="flex-1 flex flex-col min-w-0 md:pl-[72px] lg:pl-64 transition-all duration-300">
        <Header title={pageTitle} bind:isMobileOpen={isMobileMenuOpen} />

        {#if !appContext.initialized || appContext.loadings.projects}
            <main class="flex-1 flex items-center justify-center">
                <div class="flex flex-col items-center gap-4">
                    <span class="size-8 animate-spin rounded-full border-3 border-primary border-t-transparent"></span>
                    <p class="text-sm text-muted-foreground">Loading your workspace…</p>
                </div>
            </main>
        {:else if ready}
            {#key page.url.pathname}
                <main
                    class="flex-1 overflow-y-auto px-9 py-8"
                    in:fly={{ duration: 350, opacity: 0, y: 8, delay: 300 }}
                    out:fly={{ duration: 250, opacity: 0 }}
                >
                    {@render children()}
                </main>
            {/key}
        {:else}
            <main class="flex-1 flex items-center justify-center">
                <div class="flex flex-col items-center gap-3 text-center px-6">
                    <p class="text-muted-foreground text-sm">Please create a project to get started.</p>
                </div>
            </main>
        {/if}
    </div>
</div>
