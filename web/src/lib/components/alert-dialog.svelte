<script lang="ts">
    import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";

    type Props = {
        title: string;
        description: string;
        cancelText?: string;
        actionText?: string;
        onCancel?: () => void;
        onAction?: () => void;
    };

    let {
        title,
        description,
        cancelText = "Cancel",
        actionText = "Continue",
        onCancel = () => {},
        onAction = () => {},
    }: Props = $props();
</script>

<AlertDialog.Root>
    <AlertDialog.Trigger class={buttonVariants({ variant: "outline" })}>{title}</AlertDialog.Trigger>
    <AlertDialog.Content>
        <AlertDialog.Header>
            <AlertDialog.Title>{title}</AlertDialog.Title>
            <AlertDialog.Description>
                {description}
            </AlertDialog.Description>
        </AlertDialog.Header>
        <AlertDialog.Footer>
            <AlertDialog.Cancel onclick={onCancel}>{cancelText}</AlertDialog.Cancel>
            <AlertDialog.Action onclick={onAction}>{actionText}</AlertDialog.Action>
        </AlertDialog.Footer>
    </AlertDialog.Content>
</AlertDialog.Root>
