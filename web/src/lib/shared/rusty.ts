/**
 * 🦀 Rusty - Rust-style error handling for TypeScript
 *
 * Brings Result<T, E> and Option<T> with full type safety and ergonomic APIs
 */

/**
 * Represents either success (Ok) or failure (Err)
 * Similar to Rust's Result<T, E>
 */
export type Result<T, E = Error> = Ok<T, E> | Err<T, E>;

class Ok<T, E> {
    readonly _tag = "Ok" as const;

    constructor(readonly value: T) {}

    public isOk(): this is Ok<T, E> {
        return true;
    }

    public isErr(): this is Err<T, E> {
        return false;
    }

    /**
     * Unwraps the value, throws if Err (like Rust's unwrap)
     * ⚠️ Use with caution in production
     */
    public unwrap(): T {
        return this.value;
    }

    /**
     * Returns the value or a default
     */
    public unwrapOr(_defaultValue: T): T {
        return this.value;
    }

    /**
     * Returns the value or computes it from a function
     */
    public unwrapOrElse(_fn: (err: E) => T): T {
        return this.value;
    }

    /**
     * Transforms the Ok value, leaves Err untouched
     */
    public map<U>(fn: (value: T) => U): Result<U, E> {
        return ok(fn(this.value));
    }

    /**
     * Transforms the Err value, leaves Ok untouched
     */
    public mapErr<F>(_fn: (err: E) => F): Result<T, F> {
        return ok(this.value);
    }

    /**
     * Chains operations that return Results (flatMap)
     */
    public andThen<U>(fn: (value: T) => Result<U, E>): Result<U, E> {
        return fn(this.value);
    }

    /**
     * Returns this Result if Ok, otherwise returns the other
     */
    public or(_other: Result<T, E>): Result<T, E> {
        return this;
    }

    /**
     * Pattern matching like Rust's match
     */
    public match<U>(patterns: { ok: (value: T) => U; err: (error: E) => U }): U {
        // Force ok, because this object represent success operation
        return patterns.ok(this.value);
    }

    /**
     * Execute side effects without changing the Result
     */
    public inspect(fn: (value: T) => void): Result<T, E> {
        fn(this.value);
        return this;
    }

    public inspectErr(_fn: (err: E) => void): Result<T, E> {
        return this;
    }

    /**
     * Converts to Option, discarding error
     */
    public ok(): Option<T> {
        return some(this.value);
    }

    /**
     * Converts to Option of the error
     */
    public err(): Option<E> {
        return none();
    }
}

class Err<T, E> {
    readonly _tag = "Err" as const;

    constructor(readonly error: E) {}

    isOk(): this is Ok<T, E> {
        return false;
    }

    isErr(): this is Err<T, E> {
        return true;
    }

    unwrap(): never {
        throw new Error(`Called unwrap() on Err: ${JSON.stringify(this.error)}`);
    }

    unwrapOr(defaultValue: T): T {
        return defaultValue;
    }

    unwrapOrElse(fn: (err: E) => T): T {
        return fn(this.error);
    }

    map<U>(_fn: (value: T) => U): Result<U, E> {
        return err(this.error);
    }

    mapErr<F>(fn: (err: E) => F): Result<T, F> {
        return err(fn(this.error));
    }

    andThen<U>(_fn: (value: T) => Result<U, E>): Result<U, E> {
        return err(this.error);
    }

    or(other: Result<T, E>): Result<T, E> {
        return other;
    }

    match<U>(patterns: { ok: (value: T) => U; err: (error: E) => U }): U {
        return patterns.err(this.error);
    }

    inspect(_fn: (value: T) => void): Result<T, E> {
        return this;
    }

    inspectErr(fn: (err: E) => void): Result<T, E> {
        fn(this.error);
        return this;
    }

    ok(): Option<T> {
        return none();
    }

    err(): Option<E> {
        return some(this.error);
    }
}

/**
 * Creates an Ok Result
 */
export function ok<T, E = never>(value: T): Result<T, E> {
    return new Ok(value);
}

/**
 * Creates an Err Result
 */
export function err<T = never, E = Error>(error: E): Result<T, E> {
    return new Err(error);
}

// ============================================================================
// OPTION TYPE
// ============================================================================

/**
 * Represents optional values: Some(value) or None
 * Similar to Rust's Option<T>
 */
export type Option<T> = Some<T> | None;

class Some<T> {
    readonly _tag = "Some" as const;

    constructor(readonly value: T) {}

    isSome(): this is Some<T> {
        return true;
    }

    isNone(): this is None {
        return false;
    }

    unwrap(): T {
        return this.value;
    }

    unwrapOr(_defaultValue: T): T {
        return this.value;
    }

    unwrapOrElse(_fn: () => T): T {
        return this.value;
    }

    map<U>(fn: (value: T) => U): Option<U> {
        return some(fn(this.value));
    }

    andThen<U>(fn: (value: T) => Option<U>): Option<U> {
        return fn(this.value);
    }

    or(_other: Option<T>): Option<T> {
        return this;
    }

    match<U>(patterns: { some: (value: T) => U; none: () => U }): U {
        return patterns.some(this.value);
    }

    inspect(fn: (value: T) => void): Option<T> {
        fn(this.value);
        return this;
    }

    /**
     * Converts Option to Result
     */
    okOr<E>(error: E): Result<T, E> {
        return ok(this.value);
    }

    okOrElse<E>(_fn: () => E): Result<T, E> {
        return ok(this.value);
    }

    /**
     * Filters the Option based on a predicate
     */
    filter(predicate: (value: T) => boolean): Option<T> {
        return predicate(this.value) ? this : none();
    }
}

class None {
    readonly _tag = "None" as const;

    isSome(): this is Some<never> {
        return false;
    }

    isNone(): this is None {
        return true;
    }

    unwrap(): never {
        throw new Error("Called unwrap() on None");
    }

    unwrapOr<T>(defaultValue: T): T {
        return defaultValue;
    }

    unwrapOrElse<T>(fn: () => T): T {
        return fn();
    }

    map<U>(_fn: (value: never) => U): Option<U> {
        return none();
    }

    andThen<U>(_fn: (value: never) => Option<U>): Option<U> {
        return none();
    }

    or<T>(other: Option<T>): Option<T> {
        return other;
    }

    match<U>(patterns: { some: (value: never) => U; none: () => U }): U {
        return patterns.none();
    }

    inspect(_fn: (value: never) => void): Option<never> {
        return this;
    }

    okOr<T, E>(error: E): Result<T, E> {
        return err(error);
    }

    okOrElse<T, E>(fn: () => E): Result<T, E> {
        return err(fn());
    }

    filter(_predicate: (value: never) => boolean): Option<never> {
        return this;
    }
}

const CACHE_NONE = new None();

/**
 * Creates a Some Option
 */
export function some<T>(value: T): Option<T> {
    return new Some(value);
}

/**
 * Creates a None Option
 */
export function none<T = never>(): Option<T> {
    return CACHE_NONE;
}

/**
 * Converts nullable values to Option
 */
export function fromNullable<T>(value: T | null | undefined): Option<T> {
    return value != null ? some(value) : none();
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

/**
 * Wraps a synchronous function in Result
 */
export function tryCatch<T, E = Error>(fn: () => T, mapError?: (error: unknown) => E): Result<T, E> {
    try {
        return ok(fn());
    } catch (error) {
        return err(mapError ? mapError(error) : (error as E));
    }
}

/**
 * Wraps an async function in Result
 */
export async function tryCatchAsync<T, E = Error>(
    fn: () => Promise<T>,
    mapError?: (error: unknown) => E,
): Promise<Result<T, E>> {
    try {
        const value = await fn();
        return ok(value);
    } catch (error) {
        return err(mapError ? mapError(error) : (error as E));
    }
}

/**
 * Converts a Promise to Result
 */
export async function fromPromise<T, E = Error>(
    promise: Promise<T>,
    mapError?: (error: unknown) => E,
): Promise<Result<T, E>> {
    try {
        const value = await promise;
        return ok(value);
    } catch (error) {
        return err(mapError ? mapError(error) : (error as E));
    }
}

/**
 * Combines multiple Results into one
 * Returns Ok only if all are Ok, otherwise returns first Err
 */
export function combine<T extends readonly Result<any, any>[]>(
    results: T,
): Result<
    { [K in keyof T]: T[K] extends Result<infer U, any> ? U : never },
    T[number] extends Result<any, infer E> ? E : never
> {
    const values: any[] = [];

    for (const result of results) {
        if (result.isErr()) {
            return err(result.error);
        }
        values.push(result.value);
    }

    return ok(values as any);
}

/**
 * Collects all Ok values, discarding Errs
 */
export function collectOk<T, E>(results: Result<T, E>[]): T[] {
    return results.filter((r) => r.isOk()).map((r) => r.unwrap());
}

/**
 * Collects all Err values, discarding Oks
 */
export function collectErr<T, E>(results: Result<T, E>[]): E[] {
    return results.filter((r) => r.isErr()).map((r) => r.error);
}

/**
 * Partitions Results into [oks, errs]
 */
export function partition<T, E>(results: Result<T, E>[]): [T[], E[]] {
    const oks: T[] = [];
    const errs: E[] = [];

    for (const result of results) {
        if (result.isOk()) {
            oks.push(result.value);
        } else {
            errs.push(result.error);
        }
    }

    return [oks, errs];
}

// ============================================================================
// TYPE GUARDS & UTILITIES
// ============================================================================

/**
 * Type guard for Result
 */
export function isResult<T, E>(value: unknown): value is Result<T, E> {
    return value instanceof Ok || value instanceof Err;
}

/**
 * Type guard for Option
 */
export function isOption<T>(value: unknown): value is Option<T> {
    return value instanceof Some || value instanceof None;
}

export default {
    ok,
    err,
    some,
    none,
    fromNullable,

    tryCatch,
    tryCatchAsync,
    fromPromise,
    combine,
    collectOk,
    collectErr,
    partition,

    isResult,
    isOption,
};
