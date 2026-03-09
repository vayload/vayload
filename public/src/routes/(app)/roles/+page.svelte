<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import { ROLES, type Role } from "$lib/data";
    import { Shield } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Card from "$lib/components/ui/card";
    import { Badge } from "$lib/components/ui/badge";
    import * as Avatar from "$lib/components/ui/avatar";

    const roleData = [
        { role: ROLES[0], users: 3, type: "System" },
        { role: ROLES[1], users: 12, type: "Custom" },
        { role: ROLES[2], users: 5, type: "Custom" },
        { role: ROLES[3], users: 0, type: "System" },
    ];
</script>

<div class="flex flex-col h-full">
    <SectionHeader
        title="Roles & Permissions"
        subtitle="Define what users can do in the system."
        breadcrumbs={["System", "Roles & ACL"]}
    >
        {#snippet actions()}
            <Button>Create Role</Button>
        {/snippet}
    </SectionHeader>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each roleData as { role, users, type }}
            <Card.Root class="hover:shadow-md transition-all">
                <Card.Header>
                    <div class="flex justify-between items-start mb-4">
                        <div class="p-2 bg-muted rounded-lg text-muted-foreground">
                            <Shield size={20} />
                        </div>
                        <Badge variant={type === "System" ? "default" : "secondary"}>{type}</Badge>
                    </div>
                    <Card.Title>{role.name}</Card.Title>
                    <Card.Description class="h-10">{role.description}</Card.Description>
                </Card.Header>
                <Card.Content>
                    <div class="flex items-center justify-between pt-4 border-t border-gray-50">
                        <div class="flex -space-x-2">
                            {#each Array(Math.min(users, 3)) as _, u}
                                <Avatar.Root class="w-6 h-6 border-2 border-card">
                                    <Avatar.Fallback class="bg-muted"></Avatar.Fallback>
                                </Avatar.Root>
                            {/each}
                            {#if users > 3}
                                <Avatar.Root class="w-6 h-6 border-2 border-card">
                                    <Avatar.Fallback class="bg-muted text-muted-foreground text-[8px]">
                                        +{users - 3}
                                    </Avatar.Fallback>
                                </Avatar.Root>
                            {/if}
                        </div>
                        <Button variant="link" size="sm" class="text-primary hover:text-primary/80">
                            Edit Permissions
                        </Button>
                    </div>
                </Card.Content>
            </Card.Root>
        {/each}
    </div>
</div>
