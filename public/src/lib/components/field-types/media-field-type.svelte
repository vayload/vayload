<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface MediaFieldConfig {
        name: string;
        label: string;
        required: boolean;
        multiple: boolean;
        allowedTypes: string[];
        maxSize: number;
        maxFiles: number;
        uploadFolder: string;
        validation?: {
            message?: string;
        };
    }

    let config: MediaFieldConfig = {
        name: "",
        label: "",
        required: false,
        multiple: false,
        allowedTypes: ["image/*", "video/*"],
        maxSize: 10485760, // 10MB
        maxFiles: 1,
        uploadFolder: "uploads",
        validation: {},
    };

    export let onUpdate: (config: MediaFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">🎬</span>
            Media Field
        </CardTitle>
        <CardDescription>Configure a media upload field for images, videos, and files.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="media-name">Field Name *</Label>
                <Input
                    id="media-name"
                    placeholder="e.g., image, video, attachment"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="media-label">Display Label *</Label>
                <Input
                    id="media-label"
                    placeholder="e.g., Image, Video, Attachment"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="space-y-2">
            <Label for="media-types">Allowed File Types</Label>
            <Input
                id="media-types"
                placeholder="image/*,video/*,application/pdf"
                bind:value={config.allowedTypes}
                oninput={handleChange}
            />
        </div>

        <div class="grid grid-cols-3 gap-4">
            <div class="space-y-2">
                <Label for="media-size">Max Size (bytes)</Label>
                <Input
                    id="media-size"
                    type="number"
                    placeholder="10485760"
                    bind:value={config.maxSize}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="media-files">Max Files</Label>
                <Input
                    id="media-files"
                    type="number"
                    placeholder="1"
                    bind:value={config.maxFiles}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="media-folder">Upload Folder</Label>
                <Input
                    id="media-folder"
                    placeholder="uploads"
                    bind:value={config.uploadFolder}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="media-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="media-required">Required Field</Label>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="media-multiple" bind:checked={config.multiple} onchange={handleChange} />
            <Label for="media-multiple">Allow Multiple Files</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.MEDIA,
                        required: config.required,
                        multiple: config.multiple,
                        allowedTypes: config.allowedTypes,
                        maxSize: config.maxSize,
                        maxFiles: config.maxFiles,
                        uploadFolder: config.uploadFolder,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
