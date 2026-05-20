const ENCODING = "0123456789ABCDEFGHJKMNPQRSTVWXYZ";

function encodeTime(time: bigint, length: number) {
    let out = "";
    for (let i = length - 1; i >= 0; i--) {
        const mod = time % 32n;
        out = ENCODING[Number(mod)] + out;
        time = time / 32n;
    }
    return out;
}

function encodeRandom(length: number) {
    let out = "";
    const randomBytes = crypto.getRandomValues(new Uint8Array(length));

    for (let i = 0; i < length; i++) {
        out += ENCODING[randomBytes[i] % 32];
    }
    return out;
}

export function ulid() {
    const time = BigInt(Date.now());
    const timePart = encodeTime(time, 10); // 48 bits → 10 chars
    const randomPart = encodeRandom(16); // 80 bits → 16 chars

    return timePart + randomPart;
}

export const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const withDelay = async <T>(fn: () => T, min = 300, max = 800): Promise<T> => {
    const ms = min + Math.random() * (max - min);
    await delay(ms);
    return fn();
};
