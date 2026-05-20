<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { FieldTypes } from "$lib/constants/field-types";

    interface TextFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: string;
        placeholder?: string;
        maxLength?: number;
        minLength?: number;
        pattern?: string;
        validation?: {
            message?: string;
        };
    }

    let config: TextFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: "",
        placeholder: "",
        maxLength: undefined,
        minLength: undefined,
        pattern: "",
        validation: {},
    };

    export let onUpdate: (config: TextFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">📝</span>
            Text Field
        </CardTitle>
        <CardDescription>Configure a text input field for single-line text content.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="text-name">Field Name *</Label>
                <Input
                    id="text-name"
                    placeholder="e.g., title, name, email"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="text-label">Display Label *</Label>
                <Input
                    id="text-label"
                    placeholder="e.g., Title, Name, Email"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="text-placeholder">Placeholder</Label>
            <Input
                id="text-placeholder"
                placeholder="Enter placeholder text"
                bind:value={config.placeholder}
                oninput={handleChange}
            />
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="text-default">Default Value</Label>
                <Input
                    id="text-default"
                    placeholder="Default text"
                    bind:value={config.defaultValue}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="text-pattern">Validation Pattern (Regex)</Label>
                <Input
                    id="text-pattern"
                    placeholder="e.g., ^[a-zA-Z0-9]+$"
                    bind:value={config.pattern}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="text-min-length">Min Length</Label>
                <Input
                    id="text-min-length"
                    type="number"
                    placeholder="0"
                    bind:value={config.minLength}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="text-max-length">Max Length</Label>
                <Input
                    id="text-max-length"
                    type="number"
                    placeholder="255"
                    bind:value={config.maxLength}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="text-validation-message">Validation Error Message</Label>
            <Input id="text-validation-message" placeholder="This field is required" oninput={handleChange} />
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="text-required" value={config.required} onchange={handleChange} />
            <Label for="text-required">Required Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.TEXT,
                        required: config.required,
                        default: config.defaultValue,
                        maxLength: config.maxLength,
                        minLength: config.minLength,
                        pattern: config.pattern,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
