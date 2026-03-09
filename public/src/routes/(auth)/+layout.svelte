<script lang="ts">
    import { page } from "$app/state";
    import { fly } from "svelte/transition";

    let { children } = $props();
</script>

<div
    class="relative flex min-h-svh items-center justify-center p-6 md:p-10
         bg-linear-gradient-to-br from-neutral-500/20 via-neutral-500/20 to-neutral-500/20"
>
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
        <div class="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-neutral-500/30 blur-3xl"></div>
        <div class="absolute top-1/3 -right-32 h-96 w-96 rounded-full bg-neutral-500/30 blur-3xl"></div>
        <div class="absolute bottom-0 left-1/3 h-96 w-96 rounded-full bg-neutral-500/20 blur-3xl"></div>
    </div>

    {#key page.url.pathname}
        <main class="relative w-full" in:fly={{ duration: 350, opacity: 0, y: 10 }}>
            {@render children()}
        </main>
    {/key}
</div>
