<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface LocationFieldConfig {
        name: string;
        label: string;
        required: boolean;
        defaultValue?: {
            latitude: number;
            longitude: number;
            address?: string;
        };
        mapProvider: "openstreetmap" | "google" | "mapbox";
        showCoordinates: boolean;
        showAddress: boolean;
        validation?: {
            message?: string;
        };
    }

    let config: LocationFieldConfig = {
        name: "",
        label: "",
        required: false,
        defaultValue: {
            latitude: 0,
            longitude: 0,
            address: "",
        },
        mapProvider: "openstreetmap",
        showCoordinates: true,
        showAddress: true,
        validation: {},
    };

    export let onUpdate: (config: LocationFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">📍</span>
            Location Field
        </CardTitle>
        <CardDescription>Configure a location field with map integration.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="location-name">Field Name *</Label>
                <Input
                    id="location-name"
                    placeholder="e.g., location, address, coordinates"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="location-label">Display Label *</Label>
                <Input
                    id="location-label"
                    placeholder="e.g., Location, Address, Coordinates"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="location-provider">Map Provider</Label>
            <Input
                id="location-provider"
                placeholder="openstreetmap, google, mapbox"
                bind:value={config.mapProvider}
                oninput={handleChange}
            />
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="location-lat">Default Latitude</Label>
                <!-- <Input
                    id="location-lat"
                    type="number"
                    step="any"
                    placeholder="0.0000"
                    bind:value={config.defaultValue?.latitude}
                    oninput={handleChange}
                /> -->
                <Input id="location-lat" type="number" step="any" placeholder="0.0000" oninput={handleChange} />
            </div>
            <div class="space-y-2">
                <Label for="location-lng">Default Longitude</Label>
                <!-- <Input
                    id="location-lng"
                    type="number"
                    step="any"
                    placeholder="0.0000"
                    bind:value={config.defaultValue?.longitude}
                    oninput={handleChange}
                /> -->
                <Input id="location-lng" type="number" step="any" placeholder="0.0000" oninput={handleChange} />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="location-address">Default Address</Label>
            <!-- <Input
                id="location-address"
                placeholder="123 Main St, City, Country"
                bind:value={config.defaultValue?.address}
                on:input={handleChange}
            /> -->
            <Input id="location-address" placeholder="123 Main St, City, Country" oninput={handleChange} />
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="location-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="location-required">Required Field</Label>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="location-coords" bind:checked={config.showCoordinates} onchange={handleChange} />
            <Label for="location-coords">Show Coordinates</Label>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="location-addr" bind:checked={config.showAddress} onchange={handleChange} />
            <Label for="location-addr">Show Address Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.LOCATION,
                        required: config.required,
                        default: config.defaultValue,
                        mapProvider: config.mapProvider,
                        showCoordinates: config.showCoordinates,
                        showAddress: config.showAddress,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
