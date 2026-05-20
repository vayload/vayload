<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import { fetchUsers, type User } from "$lib/data";
    import { MoreHorizontal } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";
    import * as Avatar from "$lib/components/ui/avatar";
    import { onMount } from "svelte";

    let users = $state<User[]>([]);
    let loading = $state(true);

    onMount(async () => {
        users = await fetchUsers();
        loading = false;
    });

    const userRoles = {
        1: "Owner",
        2: "Developer",
        3: "Editor",
    };

    const userStatuses = {
        1: "Active",
        2: "Active",
        3: "Invited",
    };

    const lastLogin = {
        1: "2 mins ago",
        2: "1 day ago",
        3: "-",
    };

    function getInitials(firstName: string, lastName: string): string {
        return `${firstName[0]}${lastName[0]}`;
    }
</script>

<div class="pb-8">
    <SectionHeader title="User Management" subtitle="Control access to your project." breadcrumbs={["System", "Users"]}>
        {#snippet actions()}
            <Button>Invite User</Button>
        {/snippet}
    </SectionHeader>

    <div class="bg-card border rounded-xl shadow-sm">
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head>User</Table.Head>
                    <Table.Head>Role</Table.Head>
                    <Table.Head>Status</Table.Head>
                    <Table.Head>Last Login</Table.Head>
                    <Table.Head></Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#if loading}
                    {#each Array(3) as _}
                        <Table.Row>
                            <Table.Cell colspan={5}>
                                <div class="h-12 bg-muted animate-pulse rounded"></div>
                            </Table.Cell>
                        </Table.Row>
                    {/each}
                {:else}
                    {#each users as user}
                        <Table.Row class="hover:bg-muted/50">
                            <Table.Cell>
                                <div class="flex items-center">
                                    <Avatar.Root class="w-8 h-8 mr-3">
                                        <Avatar.Fallback class="bg-primary/10 text-primary text-xs font-bold">
                                            {getInitials(user.first_name, user.last_name)}
                                        </Avatar.Fallback>
                                    </Avatar.Root>
                                    <div>
                                        <div class="text-sm font-medium text-foreground">
                                            {user.first_name}
                                            {user.last_name}
                                        </div>
                                        <div class="text-xs text-muted-foreground">{user.email}</div>
                                    </div>
                                </div>
                            </Table.Cell>
                            <Table.Cell class="text-sm text-muted-foreground">
                                {userRoles[user.id as keyof typeof userRoles]}
                            </Table.Cell>
                            <Table.Cell>
                                <Badge
                                    variant={userStatuses[user.id as keyof typeof userStatuses] === "Active"
                                        ? "default"
                                        : "secondary"}
                                >
                                    {userStatuses[user.id as keyof typeof userStatuses]}
                                </Badge>
                            </Table.Cell>
                            <Table.Cell class="text-sm text-muted-foreground">
                                {lastLogin[user.id as keyof typeof lastLogin]}
                            </Table.Cell>
                            <Table.Cell class="text-right">
                                <Button variant="ghost" size="icon">
                                    <MoreHorizontal size={18} />
                                </Button>
                            </Table.Cell>
                        </Table.Row>
                    {/each}
                {/if}
            </Table.Body>
        </Table.Root>
    </div>
</div>
