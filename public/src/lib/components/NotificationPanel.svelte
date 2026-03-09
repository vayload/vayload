<script lang="ts">
    import { notificationStore } from "$lib/stores/notificationStore.svelte";
    import { X, Info, CheckCircle2, AlertCircle } from "@lucide/svelte";

    interface Props {
        isOpen: boolean;
        onClose: () => void;
    }

    let { isOpen, onClose }: Props = $props();

    function formatTime(datetime: string): string {
        const date = new Date(datetime);
        const now = new Date();
        const diff = now.getTime() - date.getTime();
        const minutes = Math.floor(diff / 60000);
        const hours = Math.floor(diff / 3600000);
        const days = Math.floor(diff / 86400000);

        if (minutes < 60) return `${minutes}m ago`;
        if (hours < 24) return `${hours}h ago`;
        return `${days}d ago`;
    }

    async function handleMarkAllRead() {
        await notificationStore.markAllAsRead();
    }

    async function handleNotificationClick(id: string) {
        await notificationStore.markAsRead(id);
    }
</script>

{#if isOpen}
    <!-- Backdrop -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-30" onclick={onClose}></div>

    <!-- Panel -->
    <div
        class="absolute right-0 top-14 z-40 w-80 bg-neutral-900 border border-neutral-800 rounded-xl shadow-2xl overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200"
    >
        <div class="flex items-center justify-between p-4 border-b border-neutral-800 bg-neutral-900/50">
            <h3 class="text-sm font-semibold text-white">Notifications</h3>
            <button onclick={handleMarkAllRead} class="text-[10px] text-primary font-bold uppercase hover:underline">
                Mark all as read
            </button>
        </div>

        <div class="max-h-[400px] overflow-y-auto custom-scrollbar">
            {#if notificationStore.notifications.length === 0}
                <div class="p-8 text-center text-neutral-500">
                    <p class="text-sm">No notifications</p>
                </div>
            {:else}
                {#each notificationStore.notifications as notif}
                    <button
                        onclick={() => handleNotificationClick(notif.id)}
                        class="w-full p-4 border-b border-neutral-800 hover:bg-white/3 transition-colors cursor-pointer flex gap-3 text-left {notif.status ===
                        'unread'
                            ? 'bg-primary/5'
                            : ''}"
                    >
                        <div
                            class="mt-1 shrink-0 w-2 h-2 rounded-full {notif.status === 'unread'
                                ? 'bg-primary'
                                : 'bg-transparent'}"
                        ></div>
                        <div class="flex-1 min-w-0">
                            <p class="text-sm font-medium text-white line-clamp-1">{notif.title}</p>
                            <p class="text-xs text-neutral-500 mt-1 line-clamp-2 leading-relaxed">{notif.body}</p>
                            <p class="text-[10px] text-neutral-600 mt-2 font-medium uppercase tracking-wider">
                                {formatTime(notif.datetime)}
                            </p>
                        </div>
                    </button>
                {/each}
            {/if}
        </div>

        <div class="p-3 bg-neutral-950/50 border-t border-neutral-800 text-center">
            <button class="text-xs text-neutral-500 font-bold uppercase hover:text-white transition-colors">
                View All History
            </button>
        </div>
    </div>
{/if}
