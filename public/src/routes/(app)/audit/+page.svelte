<script lang="ts">
    import SectionHeader from "$lib/components/SectionHeader.svelte";
    import { AUDIT_LOGS, type AuditLog } from "$lib/data";
    import { Search, Filter } from "@lucide/svelte";
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import * as Table from "$lib/components/ui/table";
    import { Badge } from "$lib/components/ui/badge";

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleString("en-US", {
            month: "short",
            day: "numeric",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function getActionBadge(action: string) {
        if (action.includes("login")) return "default";
        if (action.includes("publish")) return "default";
        if (action.includes("delete")) return "destructive";
        return "secondary";
    }
</script>

<div class="pb-8">
    <SectionHeader title="Audit Logs" subtitle="Track all system activity." breadcrumbs={["System", "Audit Logs"]} />

    <div class="bg-card border rounded-xl shadow-sm overflow-hidden">
        <div class="p-4 border-b flex flex-wrap gap-4 justify-between items-center">
            <div class="flex items-center gap-2">
                <div class="relative">
                    <Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                    <Input type="text" placeholder="Search logs..." class="pl-9 pr-4 w-64" />
                </div>
                <Button variant="outline" size="sm">
                    <Filter size={16} class="mr-2" />
                    Filters
                </Button>
            </div>
            <div class="text-sm text-muted-foreground">Showing {AUDIT_LOGS.length} entries</div>
        </div>

        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head>Timestamp</Table.Head>
                    <Table.Head>Action</Table.Head>
                    <Table.Head>Actor</Table.Head>
                    <Table.Head>IP Address</Table.Head>
                    <Table.Head>Details</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#each AUDIT_LOGS as log}
                    <Table.Row class="hover:bg-muted/50">
                        <Table.Cell class="text-sm text-muted-foreground">
                            {formatDate(log.created_at)}
                        </Table.Cell>
                        <Table.Cell>
                            <Badge variant={getActionBadge(log.action)}>
                                {log.action}
                            </Badge>
                        </Table.Cell>
                        <Table.Cell class="text-sm text-foreground">User #{log.actor_id}</Table.Cell>
                        <Table.Cell class="text-sm text-muted-foreground font-mono">{log.ip_address}</Table.Cell>
                        <Table.Cell class="text-sm text-muted-foreground">
                            {JSON.stringify(log.payload)}
                        </Table.Cell>
                    </Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>

        <div class="p-4 border-t flex justify-center">
            <Button variant="ghost" size="sm">Load More</Button>
        </div>
    </div>
</div>
