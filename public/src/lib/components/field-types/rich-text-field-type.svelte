<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Select, SelectContent, SelectItem, SelectTrigger } from "$lib/components/ui/select";
    import { FieldTypes } from "$lib/constants/field-types";

    interface RichTextFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: string;
        placeholder?: string;
        maxLength?: number;
        editor: "basic" | "standard" | "full";
        toolbar: string[];
        allowImages: boolean;
        allowLinks: boolean;
        allowTables: boolean;
        allowCode: boolean;
        validation?: {
            message?: string;
        };
    }

    let config: RichTextFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: "",
        placeholder: "",
        maxLength: undefined,
        editor: "standard",
        toolbar: ["bold", "italic", "underline", "list", "link"],
        allowImages: true,
        allowLinks: true,
        allowTables: false,
        allowCode: false,
        validation: {
            message: "",
        },
    };

    const props: { config: RichTextFieldConfig; onUpdate: (config: RichTextFieldConfig) => void } = $props();

    const editorOptions = [
        { value: "basic", label: "Basic (Simple formatting)" },
        { value: "standard", label: "Standard (Common features)" },
        { value: "full", label: "Full (All features)" },
    ];

    const toolbarOptions = [
        { value: "bold", label: "Bold" },
        { value: "italic", label: "Italic" },
        { value: "underline", label: "Underline" },
        { value: "strikethrough", label: "Strikethrough" },
        { value: "heading", label: "Headings" },
        { value: "list", label: "Lists" },
        { value: "quote", label: "Quote" },
        { value: "link", label: "Link" },
        { value: "image", label: "Image" },
        { value: "table", label: "Table" },
        { value: "code", label: "Code" },
        { value: "color", label: "Text Color" },
    ];

    function handleChange() {
        if (props.onUpdate) {
            props.onUpdate(config);
        }
    }

    function toggleToolbarOption(option: string) {
        const index = config.toolbar.indexOf(option);
        if (index > -1) {
            config.toolbar.splice(index, 1);
        } else {
            config.toolbar.push(option);
        }
        handleChange();
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">📄</span>
            Rich Text Field
        </CardTitle>
        <CardDescription>Configure a rich text editor with formatting options.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="richtext-name">Field Name *</Label>
                <Input
                    id="richtext-name"
                    placeholder="e.g., body, content, description"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="richtext-label">Display Label *</Label>
                <Input
                    id="richtext-label"
                    placeholder="e.g., Body, Content, Description"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="richtext-editor">Editor Type</Label>
            <!-- <Select bind:value={config.editor as any} onValueChange={handleChange} items={editorOptions}>
                <SelectTrigger>
                    <SelectValue placeholder="Select editor type" />
                </SelectTrigger>
                <SelectContent>
                    {#each editorOptions as option}
                        <SelectItem value={option.value}>{option.label}</SelectItem>
                    {/each}
                </SelectContent>
            </Select> -->
            <!-- <Select
                items={editorOptions}
                bind:value={config.editor as any}
                onValueChange={handleChange}
            /> -->
        </div>

        <div class="space-y-2">
            <Label>Toolbar Options</Label>
            <div class="grid grid-cols-3 gap-2">
                {#each toolbarOptions as option}
                    <Button
                        variant={config.toolbar.includes(option.value) ? "default" : "outline"}
                        size="sm"
                        class="h-8 text-xs"
                        onclick={() => toggleToolbarOption(option.value)}
                    >
                        {option.label}
                    </Button>
                {/each}
            </div>
        </div>

        <div class="space-y-2">
            <Label for="richtext-placeholder">Placeholder</Label>
            <Input
                id="richtext-placeholder"
                placeholder="Enter rich text content..."
                bind:value={config.placeholder}
                oninput={handleChange}
            />
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="richtext-default">Default Value</Label>
                <Input
                    id="richtext-default"
                    placeholder="Default rich text"
                    bind:value={config.defaultValue}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="richtext-max-length">Max Length</Label>
                <Input
                    id="richtext-max-length"
                    type="number"
                    placeholder="5000"
                    bind:value={config.maxLength}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="richtext-validation-message">Validation Error Message</Label>
            <Input id="richtext-validation-message" placeholder="This field is required" oninput={handleChange} />
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="flex items-center space-x-2">
                <Switch id="richtext-required" bind:checked={config.required} onchange={handleChange} />
                <Label for="richtext-required">Required Field</Label>
            </div>
            <div class="flex items-center space-x-2">
                <Switch id="richtext-images" bind:checked={config.allowImages} onchange={handleChange} />
                <Label for="richtext-images">Allow Images</Label>
            </div>
            <div class="flex items-center space-x-2">
                <Switch id="richtext-links" bind:checked={config.allowLinks} onchange={handleChange} />
                <Label for="richtext-links">Allow Links</Label>
            </div>
            <div class="flex items-center space-x-2">
                <Switch id="richtext-tables" bind:checked={config.allowTables} onchange={handleChange} />
                <Label for="richtext-tables">Allow Tables</Label>
            </div>
            <div class="flex items-center space-x-2">
                <Switch id="richtext-code" bind:checked={config.allowCode} onchange={handleChange} />
                <Label for="richtext-code">Allow Code Blocks</Label>
            </div>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.RICH_TEXT,
                        required: config.required,
                        default: config.defaultValue,
                        maxLength: config.maxLength,
                        editor: config.editor,
                        toolbar: config.toolbar,
                        allowImages: config.allowImages,
                        allowLinks: config.allowLinks,
                        allowTables: config.allowTables,
                        allowCode: config.allowCode,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
