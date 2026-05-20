<script lang="ts">
    import { projectStore } from "$lib/stores/projectStore.svelte";
    import {
        ChevronDown,
        Check,
        Plus,
        Home,
        LayoutDashboard,
        FolderPen,
        Images,
        Users,
        Shield,
        Plug,
        FileText,
        Settings,
    } from "@lucide/svelte";
    import { page } from "$app/state";

    interface Props {
        isMobileOpen?: boolean;
    }

    let { isMobileOpen = $bindable(false) }: Props = $props();

    const menuGroups = [
        {
            id: "platform",
            label: "Platform",
            items: [
                { id: "dashboard", label: "Dashboard", icon: Home, href: "/" },
                {
                    id: "content-types",
                    label: "Content Types",
                    icon: LayoutDashboard,
                    href: "/content-types",
                },
                { id: "entries", label: "Content Entries", icon: FolderPen, href: "/entries" },
                { id: "assets", label: "Media Library", icon: Images, href: "/media" },
            ],
        },
        {
            id: "system",
            label: "System",
            items: [
                { id: "users", label: "User Management", icon: Users, href: "/users" },
                { id: "roles", label: "Roles & ACL", icon: Shield, href: "/roles" },
                { id: "integrations", label: "Integrations", icon: Plug, href: "/integrations" },
                { id: "audit", label: "Audit Logs", icon: FileText, href: "/audit" },
                { id: "settings", label: "Settings", icon: Settings, href: "/settings" },
            ],
        },
    ];

    let isDropdownOpen = $state(false);

    function handleProjectSwitch(projectId: string) {
        projectStore.setCurrentProject(projectId);
        isDropdownOpen = false;
    }

    function isActive(href: string) {
        return page.url.pathname === href;
    }
</script>

{#if isMobileOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 bg-black/50 z-30 md:hidden" onclick={() => (isMobileOpen = false)}></div>
{/if}

<aside
    class="fixed bg-neutral-900 border-r border-solid border-neutral-900 inset-y-0 left-0 z-40 flex flex-col transition-all duration-300 {isMobileOpen
        ? 'translate-x-0 w-64'
        : '-translate-x-full md:translate-x-0 md:w-[72px] lg:w-64'}"
>
    <div class="relative h-16 border-b border-neutral-900 shrink-0">
        <button
            onclick={() => (isDropdownOpen = !isDropdownOpen)}
            class="w-full h-full flex items-center justify-center lg:justify-start px-0 lg:px-4 hover:bg-white/5 transition-colors"
        >
            <div class="flex items-center gap-3 w-full justify-between lg:justify-start">
                <div class="px-2 flex items-center gap-3">
                    <svg width="40" height="32" viewBox="0 0 40 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M0 24L20 0V14L40 0L20 24V14V10L0 24Z" fill="#FF6347" />
                    </svg>
                    <span>Vayload</span>
                </div>
                <div class="hidden lg:block text-left min-w-0 flex-1">
                    <!-- <div class="text-sm font-medium text-neutral-200 truncate">
                        {projectStore.currentProject?.name ?? "Select Project"}
                    </div> -->
                    <div class="text-[10px] uppercase tracking-wider text-neutral-500">
                        {projectStore.currentProject?.locale ?? ""}
                    </div>
                </div>
                <ChevronDown
                    size={14}
                    class="hidden lg:block transition-transform {isDropdownOpen ? 'rotate-180' : ''}"
                />
            </div>
        </button>

        <!-- Dropdown -->
        {#if isDropdownOpen}
            <div
                class="absolute top-14 left-2 right-2 bg-neutral-800 border border-neutral-700 rounded-xl shadow-2xl z-50 overflow-hidden"
            >
                <div class="py-2">
                    {#each projectStore.projects as proj}
                        <button
                            onclick={() => handleProjectSwitch(proj.id)}
                            class="w-full text-left px-4 py-2.5 text-sm hover:bg-neutral-700 flex items-center justify-between group transition-colors"
                        >
                            <div>
                                <div
                                    class="font-medium {proj.id === projectStore.currentProject?.id
                                        ? 'text-white'
                                        : 'text-neutral-400 group-hover:text-neutral-200'}"
                                >
                                    {proj.name}
                                </div>
                                <div class="text-[10px] text-neutral-500">{proj.locale}</div>
                            </div>
                            {#if proj.id === projectStore.currentProject?.id}
                                <Check size={14} class="text-orange-400" />
                            {/if}
                        </button>
                    {/each}
                </div>
                <div class="border-t border-neutral-700 p-2">
                    <button
                        class="w-full flex items-center gap-2 px-3 py-2 text-xs font-medium text-orange-400 hover:bg-orange-500/10 rounded-lg transition-colors"
                    >
                        <Plus size={14} /> Create New Project
                    </button>
                </div>
            </div>
        {/if}
    </div>

    <div class="flex-1 overflow-y-auto py-6 space-y-8">
        {#each menuGroups as group}
            <div class="px-3">
                <div class="hidden lg:block px-3 mb-2 text-xs font-semibold uppercase tracking-wider text-neutral-600">
                    {group.label}
                </div>
                <div class="space-y-1">
                    {#each group.items as item}
                        {@const active = isActive(item.href)}
                        <a
                            href={item.href}
                            onclick={() => (isMobileOpen = false)}
                            class="group relative flex items-center lg:justify-start justify-center w-full p-2.5 rounded-lg transition-all duration-200 {active
                                ? 'bg-orange-500/10 text-orange-400'
                                : 'hover:bg-white/5 hover:text-neutral-200'}"
                            title={item.label}
                        >
                            <item.icon size={20} strokeWidth={active ? 2.5 : 2} />
                            <span class="hidden lg:block ml-3 text-sm font-medium">
                                {item.label}
                            </span>
                            {#if active}
                                <div
                                    class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-orange-500 rounded-r-full lg:hidden"
                                ></div>
                            {/if}
                        </a>
                    {/each}
                </div>
            </div>
        {/each}
    </div>

    <div class="p-4 border-t border-neutral-800 lg:hidden">
        <button
            class="flex items-center gap-3 w-full p-2 rounded-lg hover:bg-white/5 transition-colors text-neutral-400"
        >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                />
            </svg>
            <span class="font-medium text-sm">Logout</span>
        </button>
    </div>
</aside>
