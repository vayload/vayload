export const ulid = () => {
    return crypto.getRandomValues(new Uint8Array(16)).toString();
};

export const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const withDelay = async <T>(fn: () => T, min = 300, max = 800): Promise<T> => {
    const ms = min + Math.random() * (max - min);
    await delay(ms);
    return fn();
};
