<script lang="ts">
    import { Textarea } from "$lib/components/ui/textarea/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Type, Link as LinkIcon, Image as ImageIcon, Code } from "@lucide/svelte";

    let {
        value = $bindable(),
        placeholder = "",
        onOpenMedia = () => {},
    } = $props<{
        value: string;
        placeholder?: string;
        onOpenMedia?: () => void;
    }>();
</script>

<div class="border border-neutral-800 rounded-xl overflow-hidden bg-neutral-900/30">
    <div class="flex items-center gap-1 p-1.5 border-b border-neutral-800 bg-neutral-900/80">
        <Button variant="ghost" size="icon" class="h-8 w-8 text-neutral-400 hover:text-white"><Type size={14} /></Button
        >
        <div class="w-px h-4 bg-neutral-800 mx-1"></div>
        <Button variant="ghost" size="icon" class="h-8 w-8 font-bold text-neutral-400 hover:text-white">B</Button>
        <Button variant="ghost" size="icon" class="h-8 w-8 italic text-neutral-400 hover:text-white">I</Button>
        <Button variant="ghost" size="icon" class="h-8 w-8 underline text-neutral-400 hover:text-white">U</Button>
        <div class="w-px h-4 bg-neutral-800 mx-1"></div>
        <Button variant="ghost" size="icon" class="h-8 w-8 text-neutral-400 hover:text-white"
            ><LinkIcon size={14} /></Button
        >
        <Button variant="ghost" size="icon" class="h-8 w-8 text-neutral-400 hover:text-white" onclick={onOpenMedia}>
            <ImageIcon size={14} />
        </Button>
        <Button variant="ghost" size="icon" class="h-8 w-8 text-neutral-400 hover:text-white ml-auto"
            ><Code size={14} /></Button
        >
    </div>
    <Textarea
        bind:value
        class="border-none focus-visible:ring-0 min-h-[200px] resize-y bg-transparent p-4"
        {placeholder}
    />
</div>
