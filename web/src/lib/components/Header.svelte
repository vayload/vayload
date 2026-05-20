<script lang="ts">
    import { notificationStore } from "$lib/stores/notificationStore.svelte";
    import { userStore } from "$lib/stores/userStore.svelte";
    import { Search, Bell, Menu, ChevronDown } from "@lucide/svelte";
    import NotificationPanel from "./NotificationPanel.svelte";
    import * as InputGroup from "$lib/components/ui/input-group/index.js";
    import SearchIcon from "@lucide/svelte/icons/search";

    interface Props {
        title: string;
        isMobileOpen?: boolean;
    }

    let { title, isMobileOpen = $bindable(false) }: Props = $props();

    let isNotifOpen = $state(false);
    let searchQuery = $state("");

    function handleNotificationClick() {
        isNotifOpen = !isNotifOpen;
    }

    function closeNotifications() {
        isNotifOpen = false;
    }
</script>

<header class="sticky top-0 z-20 backdrop-blur-md border-b h-16 px-6 flex items-center justify-between">
    <div class="flex items-center gap-4">
        <button onclick={() => (isMobileOpen = !isMobileOpen)} class="md:hidden p-2 rounded-lg hover:bg-gray-100">
            <Menu size={20} />
        </button>
        <h1 class="text-lg font-semibold hidden md:block">{title}</h1>
    </div>

    <div class="flex items-center gap-4">
        <InputGroup.Root>
            <InputGroup.Input placeholder="Jump to..." bind:value={searchQuery} />
            <InputGroup.Addon>
                <SearchIcon />
            </InputGroup.Addon>
        </InputGroup.Root>

        <div class="h-6 w-px hidden md:block"></div>

        <div class="relative">
            <button
                onclick={handleNotificationClick}
                class="relative p-2 rounded-lg transition-colors {notificationStore.hasUnread
                    ? 'text-indigo-600 bg-indigo-50'
                    : 'text-gray-500 hover:bg-gray-100'}"
            >
                <Bell size={20} />
                {#if notificationStore.hasUnread}
                    <span class="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full ring-2 ring-white"></span>
                {/if}
            </button>

            <NotificationPanel isOpen={isNotifOpen} onClose={closeNotifications} />
        </div>

        <button
            class="flex items-center gap-3 pl-2 pr-1 py-1 rounded-full hover:bg-gray-50 border border-transparent hover:border-gray-200 transition-all ml-2"
        >
            <div class="text-right hidden sm:block leading-tight">
                <div class="text-sm font-medium text-gray-900">
                    {userStore.currentUser?.first_name ?? "User"}
                    {userStore.currentUser?.last_name ?? ""}
                </div>
                <div class="text-xs text-gray-500">
                    {userStore.currentUser?.is_super_admin ? "Super Admin" : "User"}
                </div>
            </div>
            <div
                class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shadow-sm ring-2 ring-white"
            >
                {userStore.currentUser?.first_name?.[0] ?? "U"}{userStore.currentUser?.last_name?.[0] ?? ""}
            </div>
            <ChevronDown size={14} class="text-gray-400 hidden sm:block mr-1" />
        </button>
    </div>
</header>
