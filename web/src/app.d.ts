// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        // interface Locals {}
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }

    declare const __DEV__: boolean;
    declare class DateTime extends Date implements Date {
        public get isToday(): boolean {}
        public format(formatString: string, locale?: string): string {}
    }

    interface Window {
        setTimeout: typeof setTimeout;
        clearTimeout: typeof clearTimeout;
        msMaxTouchPoints: number;
        DateTime: typeof DateTime;
        [x: string]: any;
    }
}

export {};
