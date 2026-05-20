<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface BooleanFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: boolean;
        trueLabel?: string;
        falseLabel?: string;
        style: "toggle" | "checkbox" | "radio";
        validation?: {
            message?: string;
        };
    }

    let config: BooleanFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: false,
        trueLabel: "Yes",
        falseLabel: "No",
        style: "toggle",
        validation: {},
    };

    export let onUpdate: (config: BooleanFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">✓</span>
            Boolean Field
        </CardTitle>
        <CardDescription>Configure a boolean field with toggle, checkbox, or radio style.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="boolean-name">Field Name *</Label>
                <Input
                    id="boolean-name"
                    placeholder="e.g., is_active, published, featured"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="boolean-label">Display Label *</Label>
                <Input
                    id="boolean-label"
                    placeholder="e.g., Active, Published, Featured"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="boolean-style">Display Style</Label>
            <Input
                id="boolean-style"
                placeholder="toggle, checkbox, or radio"
                bind:value={config.style}
                oninput={handleChange}
            />
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="boolean-true-label">True Label</Label>
                <Input id="boolean-true-label" placeholder="Yes" bind:value={config.trueLabel} oninput={handleChange} />
            </div>
            <div class="space-y-2">
                <Label for="boolean-false-label">False Label</Label>
                <Input
                    id="boolean-false-label"
                    placeholder="No"
                    bind:value={config.falseLabel}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="boolean-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="boolean-required">Required Field</Label>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="boolean-default" bind:checked={config.defaultValue} onchange={handleChange} />
            <Label for="boolean-default">Default Value (True)</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.BOOLEAN,
                        required: config.required,
                        default: config.defaultValue,
                        trueLabel: config.trueLabel,
                        falseLabel: config.falseLabel,
                        style: config.style,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
