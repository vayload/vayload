<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface TonesFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: string;
        palette: string[];
        customColors: boolean;
        validation?: {
            message?: string;
        };
    }

    let config: TonesFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: "",
        palette: ["#000000", "#FF0000", "#00FF00", "#0000FF", "#FFFF00"],
        customColors: false,
        validation: {},
    };

    export let onUpdate: (config: TonesFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">🎨</span>
            Tones Field
        </CardTitle>
        <CardDescription>Configure a color picker field with predefined palette options.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="tones-name">Field Name *</Label>
                <Input
                    id="tones-name"
                    placeholder="e.g., color, theme, accent"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="tones-label">Display Label *</Label>
                <Input
                    id="tones-label"
                    placeholder="e.g., Color, Theme, Accent"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="tones-palette">Color Palette</Label>
            <Input
                id="tones-palette"
                placeholder="#000000,#FF0000,#00FF00,#0000FF"
                bind:value={config.palette}
                oninput={handleChange}
            />
        </div>

        <div class="space-y-2">
            <Label for="tones-default">Default Color</Label>
            <Input id="tones-default" placeholder="#000000" bind:value={config.defaultValue} oninput={handleChange} />
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="tones-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="tones-required">Required Field</Label>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="tones-custom" bind:checked={config.customColors} onchange={handleChange} />
            <Label for="tones-custom">Allow Custom Colors</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.TONES,
                        required: config.required,
                        default: config.defaultValue,
                        palette: config.palette,
                        customColors: config.customColors,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
