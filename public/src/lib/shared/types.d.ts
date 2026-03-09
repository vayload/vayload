/**
 * Utils type
 */
declare type Flatten<T, D extends number = 1> = T extends any[]
    ? D extends 0
        ? T
        : T extends (infer U)[]
          ? Flatten<U, D extends 1 ? 0 : D>[]
          : T
    : T;

declare type Ref<F = any> = {
    current: F | null;
};

declare type ButtonEventMap = {
    onClick?: (event: MouseEvent) => any;
    onMouseEnter?: (event: MouseEvent) => void;
    onMouseLeave?: (event: MouseEvent) => void;
    onFocus?: (event: FocusEvent) => void;
    onBlur?: (event: FocusEvent) => void;
};

// Type definitions for form-data
declare type FormDataType =
    | Record<string | number, any>
    | Array<Record<string | number, any>>
    | Array<number | string>
    | FormData;

// DOM Type definitions
declare type SelectorTypes =
    | keyof HTMLElementTagNameMap
    | HTMLElement
    | Selector
    | Element
    | Array<HTMLElement>
    | NodeListOf<HTMLElement>
    | string;

declare class Selector<T extends HTMLElement = HTMLElement> {
    /**
     * Array of hosted elements selected by the selector.
     */
    protected hosted: Array<T> = null;

    /**
     * Constructor for Selector class.
     * @param selectors - Selector types.
     * @param single - Indicates if only a single element should be selected.
     */
    public constructor(selectors: SelectorTypes, single: boolean = false);

    /**
     * Internal method to select elements based on the provided selectors.
     * @param selectors - Selector types.
     * @param single - Indicates if only a single element should be selected.
     */
    protected select(selectors: SelectorTypes, single: boolean);

    /**
     * Selects the first element from the selected elements.
     * @returns The Selector instance for method chaining.
     */
    public first(): this;

    /**
     * Selects the last element from the selected elements.
     * @returns The Selector instance for method chaining.
     */
    public last(): this;

    /**
     * Selects the first child element of the selected elements.
     * @returns The Selector instance for method chaining.
     */
    public firstChild(): this;

    /**
     * Selects the parent element of the selected elements.
     * @returns A new Selector instance with the parent element selected.
     */
    public parent(): Selector<HTMLElement>;

    /**
     * Adds one or more CSS classes to the selected elements.
     * @param tokens - CSS class names.
     * @returns The Selector instance for method chaining.
     */
    public addClass(...tokens: string[]): this;

    /**
     * Removes one or more CSS classes from the selected elements.
     * @param tokens - CSS class names.
     * @returns The Selector instance for method chaining.
     */
    public removeClass(...tokens: string[]): this;

    /**
     * Checks if the first element in the selection has a specific CSS class.
     * @param classname - CSS class name to check.
     * @returns True if the class is present, false otherwise.
     */
    public hasClass(classname: string): boolean;

    /**
     * Toggles a CSS class on the selected elements.
     * @param classname - CSS class name to toggle.
     * @returns The Selector instance for method chaining.
     */
    public toggleClass(classname: string): this;

    /**
     * Adds or removes multiple CSS classes from the selected elements.
     * @param classnames - Quantum classes to add or remove.
     * @returns The Selector instance for method chaining.
     */
    public classes(...classnames: QuantumClassess[]): this;

    /**
     * Applies CSS styles to the selected elements.
     * @param styles - Selector styles to apply.
     * @returns The Selector instance for method chaining.
     */
    public styles(styles: SelectorStyles): this;

    /**
     * Alias for the `styles` method. Applies CSS styles to the selected elements.
     * @param styles - Selector styles to apply.
     * @returns The Selector instance for method chaining.
     */
    public css(styles: SelectorStyles): this;

    /**
     * Executes a provided function for each selected element.
     * @param func - Function to execute for each element.
     * @returns The Selector instance for method chaining.
     */
    public each(func: (this: Selector, el: HTMLElement, index: number) => void): this;

    /**
     * Internal method to iterate over the selected elements and execute a function.
     * @param func - Function to execute for each element.
     * @returns The Selector instance for method chaining.
     */
    protected walk(func: (element: HTMLElement, index: number) => void): this;

    /**
     * Internal method to iterate over a single element and execute a function.
     * @param fn - Function to execute for the element.
     * @returns The Selector instance for method chaining.
     */
    protected walkOne(fn: (element: T) => void): this;

    /**
     * Tries to retrieve a value from the first selected element based on the provided function.
     * @param fn - Function to retrieve a value from the element.
     * @param _default - Default value to return if the retrieval fails.
     * @returns The retrieved value or the default value if retrieval fails.
     */
    protected tryGet(fn: (element: T) => any, _default?: any): any;

    /**
     * Gets or sets attributes on the selected elements.
     * @param attrs - Selector attributes to get or set.
     * @returns The Selector instance for method chaining when used as a setter or the attribute value when used as a getter.
     */
    public attr(attrs: SelectorAttributes | string): this | string;

    /**
     * Gets or sets properties on the selected elements.
     * @param props - Properties to get or set.
     * @returns The property value when used as a getter or undefined when used as a setter.
     */
    public prop(props: Record<string, any> | string): any;

    /**
     * Adds an event listener to the selected elements.
     * @param type - Type of event to listen for.
     * @param listener - Event listener function.
     * @returns The Selector instance for method chaining.
     */
    public on(type: keyof HTMLElementEventMap, listener: EventHandler): this;

    /**
     * Removes an event listener from the selected elements.
     * @param type - Type of event to remove the listener from.
     * @returns The Selector instance for method chaining.
     */
    public off(type: keyof HTMLElementEventMap | string): this;

    /**
     * Triggers an event on the selected elements.
     * @param type - Type of event to trigger.
     * @returns The Selector instance for method chaining.
     */
    public trigger(type: keyof HTMLElementEventMap): this;

    /**
     * Sets focus on the first selected element.
     */
    public focus(): void;

    /**
     * Checks if the selection is empty (contains no elements).
     * @returns True if the selection is empty, false otherwise.
     */
    public isEmpty(): boolean;

    /**
     * Finds elements within the selected elements matching the provided selector.
     * @param selectors - Selector string to match elements.
     * @returns A new Selector instance with the matched elements.
     */
    public find<T extends HTMLElement = HTMLElement>(selectors: string): Selector<T>;

    /**
     * Finds the first element within the selected elements matching the provided selector.
     * @param selectors - Selector string to match elements.
     * @returns A new Selector instance with the matched element.
     */
    public findOne<T extends HTMLElement = HTMLElement>(selectors: string): Selector<T>;

    /**
     * Selects a single element matching the provided selector.
     * @param selectors - Selector string to match the element.
     * @returns A new Selector instance with the matched element.
     */
    public single<T extends HTMLElement = HTMLElement>(selectors: string): Selector<T>;

    /**
     * Gets the number of selected elements.
     */
    public get length(): number;

    /**
     * Gets the first selected element.
     */
    public get element(): T | null;

    /**
     * Gets all the selected elements as an array.
     */
    public get all(): T[];

    /**
     * Converts the selected elements to an array.
     * @returns Array of selected elements.
     */
    public toArray(): T[];

    /**
     * Checks if the selected element(s) are visible.
     * @param index - Index of the element to check visibility for, or 'all' to check all elements.
     * @returns True if the element(s) are visible, false otherwise.
     */
    public isVisible(index?: number | "all"): boolean;

    /**
     * Internal method to test if a specific element is visible.
     * @param element - Element to test visibility for.
     * @returns True if the element is visible, false otherwise.
     */
    private testIfIsVisible(element: T): boolean;

    /**
     * Checks if the Selector instance has selected any elements.
     * @returns True if elements are selected, false otherwise.
     */
    public hasElement(): boolean;

    /**
     * Gets the width of the first selected element.
     * @returns The width of the element.
     */
    public width(): number;

    /**
     * Gets the height of the first selected element.
     * @returns The height of the element.
     */
    public height(): number;

    /**
     * Gets the outer width of the first selected element.
     * @param margin - Indicates if the margin should be included in the calculation.
     * @returns The outer width of the element.
     */
    public outerWidth(margin?: boolean): number;

    /**
     * Gets the outer height of the first selected element.
     * @param margin - Indicates if the margin should be included in the calculation.
     * @returns The outer height of the element.
     */
    public outerHeight(margin?: boolean): number;

    /**
     * Positions the selected elements absolutely relative to a target element.
     * @param target - Target element to position relative to.
     * @param options - Positioning options.
     */
    public absolute(target: HTMLElement, options: { top?: number; left?: number }): void;

    /**
     * Positions the selected elements relatively relative to a target element.
     * @param target - Target element to position relative to.
     * @param options - Positioning options.
     */
    public relative(target: HTMLElement, options: { top?: number; left?: number }): void;

    /**
     * Gets the dimensions (width and height) of the first selected element.
     * @returns Object with width and height properties.
     */
    public dimensions(): { width: number; height: number };

    /**
     * Scrolls the selected elements.
     * @param duration - Scroll duration in milliseconds.
     * @param extra - Additional distance to scroll.
     * @param same - Indicates if the same distance should be scrolled for all elements.
     * @returns The Selector instance for method chaining.
     */
    public scroll(duration?: number, extra?: number, same?: boolean): this;

    /**
     * Scrolls the first selected element into view.
     * @param options - Options for scrolling into view.
     * @returns The Selector instance for method chaining.
     */
    public scrollIntoView(options?: boolean | ScrollIntoViewOptions): this;

    /**
     * Gets the host element(s) of the Selector instance.
     * @returns Array of host elements or null if no elements are selected.
     */
    public getHostElement(): T[] | null;
}

declare type SelectorStyles = Partial<Record<keyof CSSStyleDeclaration, string | number>> | Record<string, any>;

declare interface SelectorAttributes {
    accesskey?: string;
    class?: string;
    contenteditable?: boolean;
    dir?: "ltr" | "rtl" | "auto";
    hidden?: boolean;
    id?: string;
    lang?: string;
    spellcheck?: boolean;
    style?: string;
    tabindex?: number;
    title?: string;
    [x: string]: any;
}

declare type EventHandler = (this: HTMLElement, event: Event, element: HTMLElement) => void;

declare interface HTMLElementEventTarget extends HTMLElement {
    events: {
        [key: string]: EventHandler | ((event: Event) => void);
    };
}

declare type SubmitHandler = (
    values: Record<string, any>,
    event?: SubmitCustomEvent,
    form?: HTMLFormElement,
) => Promise<void> | void;

declare interface SubmitCustomEvent extends SubmitEvent {
    close: () => void;
}

declare interface HTMLInputElementValidate extends HTMLInputElement {
    validate: () => boolean;
    isValid: boolean;
    [x: string]: any;
}

/**
 * Cache storage interface
 */
declare interface CacheStorage {
    /**
     * Set given key with value into cache
     *
     * @param key - The human cache key
     * @param value - The value for stored
     * @param ttl - The time to live in seconds
     */
    set(key: string, value: any, ttl?: number | CacheOptions = {}): void;

    /**
     * Get given key from cache
     *
     * @param key - The human cache key
     */
    get(key: string): any;

    /**
     * Get given key from cache or fetch it using the provided fetcher function
     *
     * @param key - The human cache key
     * @param fetcher - The function to fetch data if not found in cache
     */
    getOrFetch(key: string, fetcher: CacheFetcher): Promise<any>;

    /**
     * Determine given key as found in cache
     *
     * @param key - The human cache key
     */
    has(key: string): boolean;

    /**
     * Remove given key from cache
     *
     * @param key - The human cache key
     */
    delete(key: string): void;

    /**
     * Clear all cache
     */
    clear(): void;
}

declare type CacheFetcher<F = any> = () => F | Promise<F>;

declare type CacheBucket = {
    exp: number;
    v: any;
};

declare type CacheOptions = {
    exp?: number;
    autoexpire?: boolean;
    encoded?: boolean;
};

declare type CacheEvents = Record<"set" | "delete" | "clear" | "expire" | "update" | "renovate", Function[]>;

declare interface CacheStorageOptions {
    name: string;
    storage?: Storage;
    expiry?: number;
    autoexpiration?: boolean;
}

// Events
declare type EventEmitterFunc<P> = (payload: P) => void;

declare type EventEmitterMap<T extends Record<any, any>> = Record<
    keyof T,
    Record<string, EventEmitterFunc<T[keyof T]>>
>;

declare type SingleOrArray<T> = T | T[];

// HTML
declare type ElementProps = {
    class?: string;
    children?: ElementChild | ElementChild[];
    html?: string;
    ref?: Ref<HTMLElement | SVGElement | null>;
    style?: Partial<CSSStyleDeclaration & { [key: string]: string | number }>;
    transition?: string;
} & { [key: string]: any };

declare type ElementChild = string | HTMLElement | SVGElement | Text;

declare type ElementTag<K extends keyof HTMLElementTagNameMap | `svg:${keyof SVGElementTagNameMap}`> =
    K extends keyof HTMLElementTagNameMap
        ? HTMLElementTagNameMap[K] & { transition?: TransitionElement }
        : K extends `svg:${infer Element}`
          ? Element extends keyof SVGElementTagNameMap
              ? SVGElementTagNameMap[Element]
              : never
          : never;

declare interface TransitionElement {
    enter(start?: () => void, end?: () => void): void;
    leave(start?: () => void, end?: () => void): void;
}

declare interface AnimateOptions {
    name?: string;
    mode?: "out-in" | "in-out";
}
