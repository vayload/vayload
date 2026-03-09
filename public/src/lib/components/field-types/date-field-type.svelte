<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { Select, SelectContent, SelectItem, SelectTrigger } from "$lib/components/ui/select";
    import { FieldTypes } from "$lib/constants/field-types";

    interface DateFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: string;
        min?: string;
        max?: string;
        format: "date" | "datetime" | "time";
        validation?: {
            message?: string;
        };
    }

    let config: DateFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: "",
        min: "",
        max: "",
        format: "date",
        validation: {},
    };

    export let onUpdate: (config: DateFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">📅</span>
            Date Field
        </CardTitle>
        <CardDescription>Configure a date, datetime, or time input field.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="date-name">Field Name *</Label>
                <Input
                    id="date-name"
                    placeholder="e.g., published_at, created_date, birthday"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="date-label">Display Label *</Label>
                <Input
                    id="date-label"
                    placeholder="e.g., Published Date, Created At, Birthday"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="date-format">Date Format</Label>
            <!-- <Select bind:value={config.format as any}>
                <SelectTrigger>
                    <span>{config.format || "date"}</span>
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="date">Date Only</SelectItem>
                    <SelectItem value="datetime">Date & Time</SelectItem>
                    <SelectItem value="time">Time Only</SelectItem>
                </SelectContent>
            </Select> -->
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="date-default">Default Value</Label>
                <Input
                    id="date-default"
                    type={config.format === "time" ? "time" : "date"}
                    bind:value={config.defaultValue}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="date-min">Minimum Date</Label>
                <Input
                    id="date-min"
                    type={config.format === "time" ? "time" : "date"}
                    bind:value={config.min}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="date-max">Maximum Date</Label>
            <Input
                id="date-max"
                type={config.format === "time" ? "time" : "date"}
                bind:value={config.max}
                oninput={handleChange}
            />
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="date-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="date-required">Required Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.DATE,
                        required: config.required,
                        default: config.defaultValue,
                        min: config.min,
                        max: config.max,
                        format: config.format,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
