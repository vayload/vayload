<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface NumberFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: number;
        min?: number;
        max?: number;
        step?: number;
        precision?: number;
        validation?: {
            message?: string;
        };
    }

    let config: NumberFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: undefined,
        min: undefined,
        max: undefined,
        step: 1,
        precision: 0,
        validation: {},
    };

    export let onUpdate: (config: NumberFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">🔢</span>
            Number Field
        </CardTitle>
        <CardDescription>Configure a numeric input field with validation options.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="number-name">Field Name *</Label>
                <Input
                    id="number-name"
                    placeholder="e.g., price, quantity, rating"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="number-label">Display Label *</Label>
                <Input
                    id="number-label"
                    placeholder="e.g., Price, Quantity, Rating"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="number-default">Default Value</Label>
                <Input
                    id="number-default"
                    type="number"
                    placeholder="0"
                    bind:value={config.defaultValue}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="number-step">Step</Label>
                <Input id="number-step" type="number" placeholder="1" bind:value={config.step} oninput={handleChange} />
            </div>
        </div>

        <div class="grid grid-cols-3 gap-4">
            <div class="space-y-2">
                <Label for="number-min">Minimum Value</Label>
                <Input id="number-min" type="number" placeholder="0" bind:value={config.min} oninput={handleChange} />
            </div>
            <div class="space-y-2">
                <Label for="number-max">Maximum Value</Label>
                <Input id="number-max" type="number" placeholder="100" bind:value={config.max} oninput={handleChange} />
            </div>
            <div class="space-y-2">
                <Label for="number-precision">Decimal Places</Label>
                <Input
                    id="number-precision"
                    type="number"
                    placeholder="0"
                    bind:value={config.precision}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="number-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="number-required">Required Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.NUMBER,
                        required: config.required,
                        default: config.defaultValue,
                        min: config.min,
                        max: config.max,
                        step: config.step,
                        precision: config.precision,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
