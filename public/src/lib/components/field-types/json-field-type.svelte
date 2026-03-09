<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface JsonFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: object;
        schema?: object;
        validation?: {
            message?: string;
        };
    }

    let config: JsonFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: {},
        schema: {},
        validation: {},
    };

    export let onUpdate: (config: JsonFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">🔧</span>
            JSON Field
        </CardTitle>
        <CardDescription>Configure a JSON field for structured data storage.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="json-name">Field Name *</Label>
                <Input
                    id="json-name"
                    placeholder="e.g., metadata, settings, config"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="json-label">Display Label *</Label>
                <Input
                    id="json-label"
                    placeholder="e.g., Metadata, Settings, Config"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="json-default">Default Value (JSON)</Label>
            <Input
                id="json-default"
                placeholder={JSON.stringify({ key: "value" })}
                bind:value={config.defaultValue}
                oninput={handleChange}
            />
        </div>

        <div class="space-y-2">
            <Label for="json-schema">JSON Schema</Label>
            <Input
                id="json-schema"
                placeholder={JSON.stringify({ type: "object", properties: {} })}
                bind:value={config.schema}
                oninput={handleChange}
            />
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="json-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="json-required">Required Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.JSON,
                        required: config.required,
                        default: config.defaultValue,
                        schema: config.schema,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
