<script lang="ts">
    import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card";
    import { Label } from "$lib/components/ui/label";
    import { Switch } from "$lib/components/ui/switch";
    import { Input } from "$lib/components/ui/input";
    import { FieldTypes } from "$lib/constants/field-types";

    interface RelationshipFieldConfig {
        name: string;
        label: string;
        required: boolean;
        targetCollection: string;
        relationshipType: "one-to-one" | "one-to-many" | "many-to-one" | "many-to-many";
        displayField: string;
        filterBy?: string;
        sortBy?: string;
        limit?: number;
        validation?: {
            message?: string;
        };
    }

    let config: RelationshipFieldConfig = {
        name: "",
        label: "",
        required: false,
        targetCollection: "",
        relationshipType: "one-to-many",
        displayField: "title",
        filterBy: "",
        sortBy: "",
        limit: undefined,
        validation: {},
    };

    export let onUpdate: (config: RelationshipFieldConfig) => void = () => {};

    function handleChange() {
        if (onUpdate) {
            onUpdate(config);
        }
    }
</script>

<Card>
    <CardHeader>
        <CardTitle class="flex items-center gap-2">
            <span class="text-lg">🔗</span>
            Relationship Field
        </CardTitle>
        <CardDescription>Configure a relationship field to link with other collections.</CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="rel-name">Field Name *</Label>
                <Input
                    id="rel-name"
                    placeholder="e.g., author, category, tags"
                    bind:value={config.name}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="rel-label">Display Label *</Label>
                <Input
                    id="rel-label"
                    placeholder="e.g., Author, Category, Tags"
                    bind:value={config.label}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="rel-target">Target Collection *</Label>
                <Input
                    id="rel-target"
                    placeholder="e.g., users, posts, categories"
                    bind:value={config.targetCollection}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="rel-type">Relationship Type</Label>
                <Input
                    id="rel-type"
                    placeholder="one-to-one, one-to-many, etc."
                    bind:value={config.relationshipType}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="rel-display">Display Field</Label>
                <Input
                    id="rel-display"
                    placeholder="e.g., title, name, email"
                    bind:value={config.displayField}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="rel-limit">Limit (Optional)</Label>
                <Input id="rel-limit" type="number" placeholder="10" bind:value={config.limit} oninput={handleChange} />
            </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
                <Label for="rel-filter">Filter By</Label>
                <Input
                    id="rel-filter"
                    placeholder="e.g., status=published"
                    bind:value={config.filterBy}
                    oninput={handleChange}
                />
            </div>
            <div class="space-y-2">
                <Label for="rel-sort">Sort By</Label>
                <Input
                    id="rel-sort"
                    placeholder="e.g., created_at desc"
                    bind:value={config.sortBy}
                    oninput={handleChange}
                />
            </div>
        </div>

        <div class="flex items-center space-x-2">
            <Switch id="rel-required" bind:checked={config.required} onchange={handleChange} />
            <Label for="rel-required">Required Field</Label>
        </div>

        <div class="text-sm text-muted-foreground mt-4 p-3 bg-muted rounded">
            <strong>Database Schema:</strong>
            <pre class="text-xs mt-2">{JSON.stringify(
                    {
                        name: config.name || "field_name",
                        type: FieldTypes.RELATIONSHIP,
                        required: config.required,
                        targetCollection: config.targetCollection,
                        relationshipType: config.relationshipType,
                        displayField: config.displayField,
                        filterBy: config.filterBy,
                        sortBy: config.sortBy,
                        limit: config.limit,
                    },
                    null,
                    2,
                )}</pre>
        </div>
    </CardContent>
</Card>
