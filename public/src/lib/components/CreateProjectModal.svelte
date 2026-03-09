<script lang="ts">
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import * as Field from "$lib/components/ui/field/index.js";
    import * as Select from "$lib/components/ui/select/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Button } from "$lib/components/ui/button/index.js";
    import { appContext } from "$lib/stores/app-context.svelte";
    import { fly } from "svelte/transition";
    import FolderPlusIcon from "@lucide/svelte/icons/folder-plus";
    import GlobeIcon from "@lucide/svelte/icons/globe";
    import ClockIcon from "@lucide/svelte/icons/clock";
    import type { ProjectInput } from "$lib/data";

    interface Props {
        open: boolean;
        closable?: boolean;
    }

    let { open = $bindable(true), closable = false }: Props = $props();

    let name = $state("");
    let locale = $state("en-US");
    let timezone = $state("UTC");
    let submitting = $state(false);
    let error: string | null = $state(null);

    const locales = [
        { value: "en-US", label: "English (US)" },
        { value: "es-ES", label: "Español (ES)" },
        { value: "fr-FR", label: "Français (FR)" },
        { value: "pt-BR", label: "Português (BR)" },
        { value: "de-DE", label: "Deutsch (DE)" },
        { value: "ja-JP", label: "日本語 (JP)" },
    ];

    const timezones = [
        { value: "UTC", label: "UTC" },
        { value: "America/New_York", label: "Eastern (EST)" },
        { value: "America/Chicago", label: "Central (CST)" },
        { value: "America/Denver", label: "Mountain (MST)" },
        { value: "America/Los_Angeles", label: "Pacific (PST)" },
        { value: "Europe/London", label: "London (GMT)" },
        { value: "Europe/Berlin", label: "Berlin (CET)" },
        { value: "Asia/Tokyo", label: "Tokyo (JST)" },
    ];

    const handleSubmit = async (event: Event) => {
        event.preventDefault();
        error = null;

        if (!name.trim()) {
            error = "Project name is required.";
            return;
        }

        if (name.trim().length < 3) {
            error = "Project name must be at least 3 characters.";
            return;
        }

        submitting = true;

        try {
            const input: ProjectInput = {
                name: name.trim(),
                locale,
                settings: { timezone },
            };

            const project = await appContext.createProject(input);

            if (project) {
                open = false;
            } else {
                error = "Failed to create project. Please try again.";
            }
        } catch (err) {
            error = err instanceof Error ? err.message : "An unexpected error occurred.";
        } finally {
            submitting = false;
        }
    };
</script>

<Dialog.Root
    bind:open
    onOpenChange={(v) => {
        if (!closable) open = true;
    }}
>
    <Dialog.Content
        showCloseButton={closable}
        onInteractOutside={(e) => {
            if (!closable) e.preventDefault();
        }}
        onEscapeKeydown={(e) => {
            if (!closable) e.preventDefault();
        }}
        class="sm:max-w-lg py-7"
    >
        <Dialog.Header>
            <div class="flex items-center gap-3">
                <div>
                    <Dialog.Title>Create your first project</Dialog.Title>
                    <Dialog.Description class="pt-2">
                        Every workspace starts with a project. Set it up to begin managing your content.
                    </Dialog.Description>
                </div>
            </div>
        </Dialog.Header>

        <form onsubmit={handleSubmit} class="space-y-5 pt-2">
            <Field.Group>
                <Field.Field>
                    <Field.Label for="project-name">Project Name</Field.Label>
                    <Input
                        id="project-name"
                        type="text"
                        placeholder="My Awesome Project"
                        bind:value={name}
                        autofocus
                        disabled={submitting}
                    />
                    <Field.Description>This will be used as the display name for your project.</Field.Description>
                </Field.Field>

                <div class="grid grid-cols-2 gap-4">
                    <Field.Field>
                        <Field.Label>
                            <span class="flex items-center gap-1.5">
                                <GlobeIcon class="size-3.5 text-muted-foreground" />
                                Locale
                            </span>
                        </Field.Label>
                        <Select.Root type="single" bind:value={locale}>
                            <Select.Trigger disabled={submitting}>
                                {locales.find((l) => l.value === locale)?.label ?? "Select locale"}
                            </Select.Trigger>
                            <Select.Content>
                                {#each locales as loc}
                                    <Select.Item value={loc.value}>{loc.label}</Select.Item>
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    </Field.Field>

                    <Field.Field>
                        <Field.Label>
                            <span class="flex items-center gap-1.5">
                                <ClockIcon class="size-3.5 text-muted-foreground" />
                                Timezone
                            </span>
                        </Field.Label>
                        <Select.Root type="single" bind:value={timezone}>
                            <Select.Trigger disabled={submitting}>
                                {timezones.find((t) => t.value === timezone)?.label ?? "Select timezone"}
                            </Select.Trigger>
                            <Select.Content>
                                {#each timezones as tz}
                                    <Select.Item value={tz.value}>{tz.label}</Select.Item>
                                {/each}
                            </Select.Content>
                        </Select.Root>
                    </Field.Field>
                </div>
            </Field.Group>

            {#if error}
                <p class="text-sm text-destructive text-center" in:fly={{ y: 4, duration: 200 }}>
                    {error}
                </p>
            {/if}

            <Dialog.Footer>
                <Button type="submit" class="w-full" disabled={submitting}>
                    {#if submitting}
                        <span class="flex items-center gap-2">
                            <span class="size-4 animate-spin rounded-full border-2 border-current border-t-transparent"
                            ></span>
                            Creating project…
                        </span>
                    {:else}
                        Create Project
                    {/if}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
