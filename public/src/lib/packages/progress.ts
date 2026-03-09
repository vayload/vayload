type Listener = (value: number, active: boolean) => void;

const MIN = 0.08;
const MAX = 0.994;
const TRICKLE_SPEED = 200;

let value: number | null = null;
let active = false;
let timer: number | null = null;
let listeners = new Set<Listener>();

/* utils */
const clamp = (n: number, min: number, max: number) => Math.min(Math.max(n, min), max);

function calculateIncrement(n: number) {
    if (n < 0.2) return 0.1;
    if (n < 0.5) return 0.04;
    if (n < 0.8) return 0.02;
    if (n < 0.99) return 0.005;
    return 0;
}

/* subscription */
export function subscribe(fn: Listener) {
    listeners.add(fn);
    fn(value ?? 0, active);
    return () => listeners.delete(fn);
}

function notify() {
    listeners.forEach((fn) => fn(value ?? 0, active));
}

/* API */

export function start() {
    if (value !== null) return;

    value = MIN;
    active = true;
    notify();

    timer = window.setInterval(() => inc(), TRICKLE_SPEED);
}

export function inc(amount?: number) {
    if (value === null) {
        start();
        return;
    }

    const delta = typeof amount === "number" ? amount : calculateIncrement(value);

    value = clamp(value + delta, MIN, MAX);
    notify();
}

export function done() {
    if (value === null) return;

    if (timer) {
        clearInterval(timer);
        timer = null;
    }

    // placebo jump
    value = clamp(value + 0.3 + Math.random() * 0.5, 0, 1);
    notify();

    setTimeout(() => {
        value = null;
        active = false;
        notify();
    }, 200);
}
