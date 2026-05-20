import { Arr, hasOwnProperty, isEmpty, isNil } from "./utils";

// ============================================================================================================================
// ================================================== JSON UTILITIES ==========================================================
// ============================================================================================================================

const INVALID_JSON_STRING = /^\s*["[{]|^\s*-?\d[\d.]{0,14}\s*$/;
const SUSPECT_PROTO =
    /"(?:_|\\u0{2}5[Ff]){2}(?:p|\\u0{2}70)(?:r|\\u0{2}72)(?:o|\\u0{2}6[Ff])(?:t|\\u0{2}74)(?:o|\\u0{2}6[Ff])(?:_|\\u0{2}5[Ff]){2}"\s*:/;
const SUSPECT_CONSTRUCTOR =
    /"(?:c|\\u0063)(?:o|\\u006[Ff])(?:n|\\u006[Ee])(?:s|\\u0073)(?:t|\\u0074)(?:r|\\u0072)(?:u|\\u0075)(?:c|\\u0063)(?:t|\\u0074)(?:o|\\u006[Ff])(?:r|\\u0072)"\s*:/;

/**
 * The package provide stringify and parse Json options with safe verification
 *
 * @namespace Json
 */
export const Json = {
    /**
     * Transform javascrit object to string JSON with circular reference prevention
     *
     * @param data - The data for stringify
     * @param _default - The default value if stringify have a error
     */
    stringify: <T = any>(data: T, _default: string | null = null): string | null => {
        const seen = new WeakSet();

        const replacer = (_: any, value: any) => {
            if (typeof value === "object" && value !== null) {
                if (seen.has(value)) {
                    return "[Circular]";
                }
                seen.add(value);
            } else if (typeof value === "function") {
                return value.toString();
            } else if (typeof value === "symbol") {
                return value.toString();
            }

            return value;
        };

        try {
            return JSON.stringify(data, replacer);
        } catch (error: any) {
            console.error(`Error al serializar el objeto: ${error.message}`);
            return _default;
        }
    },

    /**
     * Parse json with sanitity content
     *
     * @param value - string represent json or other type
     * @param strict - indicates validation is strict
     */
    parse: <F = any>(value: string, strict: boolean = false): F => {
        if (typeof value !== "string") return value as F;

        const val = value.toLowerCase().trim();
        if (val === "true") return true as F;
        if (val === "false") return false as F;
        if (val === "null") return null as F;
        if (val === "nan") return Number.NaN as F;
        if (val === "infinity") return Number.POSITIVE_INFINITY as F;
        if (val === "undefined") return undefined as F;

        if (!INVALID_JSON_STRING.test(value)) {
            if (strict) {
                throw new SyntaxError("Invalid JSON");
            }
            return value as F;
        }

        try {
            if (SUSPECT_PROTO.test(value) || SUSPECT_CONSTRUCTOR.test(value)) {
                return JSON.parse(value, (key: any, value: any) => {
                    if (key === "__proto__") {
                        return;
                    }
                    if (key === "constructor" && value && typeof value === "object" && "prototype" in value) {
                        return;
                    }
                    return value;
                });
            }
            return JSON.parse(value) as F;
        } catch (error: any) {
            if (strict) {
                throw error;
            }
            return value as F;
        }
    },
};

// ============================================================================================================================
// =============================================== FORM-DATA UTILITIES ========================================================
// ============================================================================================================================

declare type FormDataType =
    | Record<string | number, any>
    | Array<Record<string | number, any>>
    | Array<number | string>
    | FormData;

const recursiveMultipart = (data: any, key?: string | number, fd = new FormData(), seen = new WeakSet<any>()) => {
    if (seen.has(data)) return fd;
    const type = typeof data;

    if (type === "object" && data !== null) {
        seen.add(data);
        if (Array.isArray(data)) {
            data.forEach((value, index) => {
                recursiveMultipart(value, `${key}[${index}]`, fd, seen);
            });
        } else {
            for (const prop in data) {
                if (hasOwnProperty(data, prop)) {
                    const value = data[prop];
                    recursiveMultipart(value, key ? `${key}[${prop}]` : prop, fd, seen);
                }
            }
        }
    } else {
        if (type === "boolean") {
            fd.append(String(key || 0), data ? "1" : "0");
        } else {
            fd.append(String(key || 0), data.toString());
        }
    }

    return fd;
};

/**
 * The FormData serializer and unserializer from complex JavaScript datatypes
 *
 * @namespace formData
 */
export const formData = {
    /**
     * Serialize given data into a FormData object.
     *
     * @param data The literal JavaScript object to serialize.
     */
    serialize: (data: FormDataType): FormData => {
        if (data instanceof FormData) {
            return data;
        }

        if (isNil(data) || isEmpty(data)) {
            return new FormData();
        }

        return recursiveMultipart(data);
    },

    /**
     * Deserialize a FormData object into a JavaScript object.
     *
     * @param data The FormData object to deserialize.
     */
    unserialize: (data: FormData): Record<string | number, any> => {
        let result = {};
        for (const [key, value] of data) {
            const keys = key.split(".");
            let currentObj: Record<string, any> = result;

            for (let i = 0; i < keys.length - 1; i++) {
                const currentKey = keys[i];
                currentObj[currentKey] = currentObj[currentKey] || {};
                currentObj = currentObj[currentKey];
            }

            const lastKey = keys[keys.length - 1];
            if (lastKey.endsWith("[]")) {
                const arrayKey = lastKey.slice(0, -2);
                currentObj[arrayKey] = currentObj[arrayKey] || [];
                currentObj[arrayKey].push(value);
            } else {
                currentObj[lastKey] = value;
            }
        }

        return result;
    },
};

/**
 * The package for Base64 encoding and decoding.
 */
const BASE64_CHARS = [
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
    'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
    'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a',
    'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
    'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's',
    't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1',
    '2', '3', '4', '5', '6', '7', '8', '9', '+',
    '/', '='
]; // prettier-ignore
const BASE64_CODES = Arr.flip(BASE64_CHARS);
const BASE64_REGEX = /^(?:[A-Za-z\d+\/]{4})*?(?:[A-Za-z\d+\/]{2}(?:==)?|[A-Za-z\d+\/]{3}=?)?$/;

const sanitizeBase64 = (s: string): string => s.replace(/[^A-Za-z0-9\+\/]/g, "");
const makeBase64UriSafe = (s: string): string => s.replace(/=/g, "").replace(/[+\/]/g, (m0) => (m0 == "+" ? "-" : "_"));
const unwrapBase64Uri = (a: string): string => sanitizeBase64(a.replace(/[-_]/g, (m0) => (m0 == "-" ? "+" : "/")));

/**
 * Base64 polyfill for the `btoa` function.
 *
 * @param {string} bin - The binary input string to encode.
 */
const btoaPolyfill = (bin: string): string => {
    let u32: number,
        c0: number,
        c1: number,
        c2: number,
        asc: string = "";
    const pad = bin.length % 3;
    for (let i = 0; i < bin.length; ) {
        if ((c0 = bin.charCodeAt(i++)) > 255 || (c1 = bin.charCodeAt(i++)) > 255 || (c2 = bin.charCodeAt(i++)) > 255)
            throw new TypeError("invalid character found");
        u32 = (c0 << 16) | (c1 << 8) | c2;
        asc +=
            BASE64_CHARS[(u32 >> 18) & 63] +
            BASE64_CHARS[(u32 >> 12) & 63] +
            BASE64_CHARS[(u32 >> 6) & 63] +
            BASE64_CHARS[u32 & 63];
    }

    return pad ? asc.slice(0, pad - 3) + "===".substring(pad) : asc;
};

/**
 * Base64 polyfill for the `atob` function.
 *
 * @param {string} asc - The ASCII input string to decode.
 */
const atobPolyfill = (asc: string): string => {
    asc = asc.replace(/\s+/g, "");
    if (!BASE64_REGEX.test(asc)) throw new TypeError("malformed base64.");
    asc += "==".slice(2 - (asc.length & 3));
    let u24: number,
        bin = "",
        r1: number,
        r2: number;
    for (let i = 0; i < asc.length; ) {
        u24 =
            (BASE64_CODES[asc.charAt(i++)] << 18) |
            (BASE64_CODES[asc.charAt(i++)] << 12) |
            ((r1 = BASE64_CODES[asc.charAt(i++)]) << 6) |
            (r2 = BASE64_CODES[asc.charAt(i++)]);
        bin +=
            r1 === 64
                ? String.fromCharCode((u24 >> 16) & 255)
                : r2 === 64
                  ? String.fromCharCode((u24 >> 16) & 255, (u24 >> 8) & 255)
                  : String.fromCharCode((u24 >> 16) & 255, (u24 >> 8) & 255, u24 & 255);
    }

    return bin;
};

const _atob = typeof atob === "function" ? (asc: string) => atob(sanitizeBase64(asc)) : atobPolyfill;
const _btoa = typeof btoa === "function" ? (bin: string) => btoa(bin) : btoaPolyfill;

/**
 * Base64 utility functions for encoding and decoding.
 *
 * @namespace Base64
 */
export const Base64 = {
    /**
     * Encodes a string to Base64.
     *
     * @param {string} src - The input string to encode.
     * @param {boolean} [urlsafe=false] - Whether to make the output URL-safe.
     */
    encode: (src: string, urlsafe: boolean = false): string => {
        const source = _btoa(src);

        return urlsafe ? makeBase64UriSafe(source) : source;
    },

    /**
     * Decodes a Base64-encoded string.
     *
     * @param {string} src - The Base64-encoded input string.
     */
    decode: (src: string): string => {
        return _atob(unwrapBase64Uri(src));
    },
};
