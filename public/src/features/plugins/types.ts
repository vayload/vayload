import type { Integration } from "$lib/types";

export type PluginIntegration = Integration;

export interface PluginsState {
    items: PluginIntegration[];
    loading: boolean;
}
