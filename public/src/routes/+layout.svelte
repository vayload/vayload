<script lang="ts">
    import "@fontsource-variable/geist";
    import "../app.css";
    import type { Snippet } from "svelte";

    import { start, inc, done, subscribe } from "$lib/packages/progress";
    import { beforeNavigate, afterNavigate } from "$app/navigation";

    let active = $state(false);
    let progress = $state(0);

    beforeNavigate(() => {
        start();
    });

    afterNavigate(() => {
        done();
    });

    interface Props {
        children: Snippet;
    }

    let { children }: Props = $props();

    $effect(() => {
        const unsubscribe = subscribe((value, act) => {
            active = act;
            progress = value;
        });

        return unsubscribe;
    });
</script>

{@render children()}

{#if active}
    <div class="bar" style="transform: scaleX({progress})"></div>
{/if}

<style>
    .bar {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: linear-gradient(90deg, #ffae00, #ff5e00);
        transform-origin: left;
        transform: scaleX(0);
        pointer-events: none;
        z-index: 9999;
        transition: transform 0.2s ease-out;
    }
</style>
