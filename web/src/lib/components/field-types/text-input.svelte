<script lang="ts">
    import { Input } from "$lib/components/ui/input/index.js";

    let {
        value = $bindable(),
        placeholder = "",
        type = "text",
    } = $props<{
        value: any;
        placeholder?: string;
        type?: string;
    }>();
</script>

<Input
    type={type === "number" ? "number" : "text"}
    bind:value
    {placeholder}
    class="bg-neutral-900/50 border-neutral-800 focus:ring-primary/20"
/>
