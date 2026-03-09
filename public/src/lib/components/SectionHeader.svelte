<script lang="ts">
    import { Home, ChevronRight } from "@lucide/svelte";
    import type { Snippet } from "svelte";

    interface Props {
        title: string;
        subtitle?: string;
        breadcrumbs?: string[];
        actions?: Snippet;
    }

    let { title, subtitle, breadcrumbs = [], actions }: Props = $props();

    function slugify(text: string): string {
        return text
            .toLowerCase()
            .trim()
            .replace(/[\s_]+/g, "-")
            .replace(/[^\w-]+/g, "")
            .replace(/--+/g, "-");
    }

    function buildPath(index: number): string {
        const segments = breadcrumbs.slice(1, index + 1).map((c) => slugify(c));
        return "/" + segments.join("/");
    }
</script>

<div class="flex flex-col gap-4 mb-6">
    {#if breadcrumbs.length}
        <div class="flex items-center gap-2 text-xs text-muted-foreground">
            <a href="/" class="flex items-center">
                <Home size={12} />
            </a>

            {#each breadcrumbs as crumb, i}
                <ChevronRight size={12} class="text-muted-foreground/50" />

                {#if i === breadcrumbs.length - 1}
                    <span class="text-foreground font-medium">
                        {crumb}
                    </span>
                {:else}
                    <a href={buildPath(i)} class="hover:text-foreground transition-colors">
                        {crumb}
                    </a>
                {/if}
            {/each}
        </div>
    {/if}

    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
            <h2 class="text-2xl font-bold text-foreground tracking-tight">
                {title}
            </h2>

            {#if subtitle}
                <p class="text-muted-foreground text-sm mt-1">
                    {subtitle}
                </p>
            {/if}
        </div>

        {#if actions}
            <div class="flex gap-3">
                {@render actions()}
            </div>
        {/if}
    </div>
</div>
