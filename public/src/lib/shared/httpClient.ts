import { tryCatchAsync } from "./rusty";

export type FileLike = {
    data: Buffer | Uint8Array | ArrayBuffer | Blob;
    originalname: string;
    mimetype: string;
    size: number;
};

export type BrowserFileLike = {
    data: Blob;
    originalname: string;
    mimetype: string;
    size: number;
};

export type ResponseType =
    | "json" // default
    | "text"
    | "blob"
    | "arraybuffer"
    | "stream" // ReadableStream<Uint8Array>
    | "document"; // parse HTML/XML as Document

export interface HttpRequestConfig extends RequestInit {
    url?: string;
    baseURL?: string;
    params?: Record<string, any>;
    data?: any;
    timeout?: number;
    headers?: Record<string, string>;
    responseType?: ResponseType;
    withCredentials?: boolean;
}

export interface HttpResponse<T = any> {
    data: T;
    status: number;
    statusText: string;
    headers: Headers;
    config: HttpRequestConfig;
}

export enum HttpMethod {
    GET = "get",
    POST = "post",
    DELETE = "delete",
    PATCH = "patch",
}

export class HttpError<T = any> extends Error {
    isHttpError = true;
    config: HttpRequestConfig;
    code?: string;
    request?: any;
    response?: HttpResponse<T>;

    constructor(message: string, config: HttpRequestConfig, response?: HttpResponse<T>, code?: string) {
        super(message);
        this.config = config;
        this.response = response;
        this.code = code;
    }
}

type Fulfilled<V> = (value: V) => V | Promise<V>;
type Rejected<E> = (error: E) => any;

interface InterceptorHandler<V, E> {
    fulfilled: Fulfilled<V>;
    rejected?: Rejected<E>;
}

class InterceptorManager<V, E> {
    private handlers: (InterceptorHandler<V, E> | null)[] = [];

    use(fulfilled: Fulfilled<V>, rejected?: Rejected<E>): number {
        this.handlers.push({ fulfilled, rejected });
        return this.handlers.length - 1;
    }

    private setHandlers(handlers: (InterceptorHandler<V, E> | null)[]) {
        this.handlers = handlers;
    }

    eject(id: number): void {
        if (this.handlers[id]) {
            this.handlers[id] = null;
        }
    }

    async run<E = any>(value: V | E, isError = false): Promise<V> {
        let current: any = value;

        for (const handler of this.handlers) {
            if (!handler) continue;

            try {
                if (isError) {
                    if (handler.rejected) {
                        current = await handler.rejected(current);
                    } else {
                        throw current;
                    }
                } else {
                    current = await handler.fulfilled(current);
                }
            } catch (err: any) {
                if (handler.rejected) {
                    current = await handler.rejected(err);
                } else {
                    throw err;
                }
            }
        }

        return current;
    }

    clear(): void {
        this.handlers = [];
    }

    public clone(): InterceptorManager<V, E> {
        const clone = new InterceptorManager<V, E>();
        clone.setHandlers([...this.handlers]);
        return clone;
    }
}

export class HttpClient {
    private defaults: HttpRequestConfig = {
        timeout: 15000,
    };

    public readonly interceptors = {
        request: new InterceptorManager<HttpRequestConfig, HttpError>(),
        response: new InterceptorManager<HttpResponse, HttpError>(),
    };

    constructor(config?: HttpRequestConfig) {
        if (config) {
            this.defaults = { ...this.defaults, ...config };
        }
    }

    public static create(config?: HttpRequestConfig): HttpClient {
        return new HttpClient(config);
    }

    public clone(): HttpClient {
        const clone = new HttpClient(this.defaults);
        clone.interceptors.request = this.interceptors.request.clone();
        clone.interceptors.response = this.interceptors.response.clone();
        return clone;
    }

    private buildURL(config: HttpRequestConfig): string {
        let fullURL = (config.baseURL || "") + (config.url || "");

        if (!config.params) return fullURL;

        const params = new URLSearchParams();
        Object.entries(config.params).forEach(([key, value]) => {
            if (value !== null && value !== undefined) {
                params.append(key, String(value));
            }
        });

        const query = params.toString();
        return query ? `${fullURL}?${query}` : fullURL;
    }

    private async _request<T = any>(config: HttpRequestConfig): Promise<HttpResponse<T>> {
        let mergedConfig: HttpRequestConfig = {
            ...this.defaults,
            ...config,
            headers: {
                ...(this.defaults.headers || {}),
                ...(config.headers || {}),
            },
        };

        mergedConfig = await this.interceptors.request.run(mergedConfig);

        const controller = new AbortController();
        const timeoutMs = mergedConfig.timeout ?? 15000;
        const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

        let body: BodyInit | null | undefined = undefined;
        const userSetContentType = "Content-Type" in (mergedConfig.headers || {});

        if (mergedConfig.data != null && !["GET", "HEAD"].includes((mergedConfig.method || "GET").toUpperCase())) {
            const data = mergedConfig.data;

            if (data instanceof FormData) {
                body = data;
                delete mergedConfig.headers?.["Content-Type"];
            } else if (data instanceof URLSearchParams) {
                body = data;
                if (!userSetContentType) {
                    mergedConfig.headers = {
                        ...mergedConfig.headers,
                        "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
                    };
                }
            } else if (
                typeof data === "object" &&
                data !== null &&
                !(
                    data instanceof Blob ||
                    data instanceof ArrayBuffer ||
                    ArrayBuffer.isView(data) ||
                    data instanceof ReadableStream ||
                    data instanceof FormData ||
                    data instanceof URLSearchParams
                )
            ) {
                body = JSON.stringify(data);
                if (!userSetContentType) {
                    mergedConfig.headers = {
                        ...mergedConfig.headers,
                        "Content-Type": "application/json;charset=UTF-8",
                    };
                }
            } else {
                body = data as BodyInit;
                if (!userSetContentType) {
                    mergedConfig.headers = {
                        ...mergedConfig.headers,
                        "Content-Type": "application/octet-stream",
                    };
                }
            }
        }

        let fetchResponse: Response;

        try {
            fetchResponse = await fetch(this.buildURL(mergedConfig), {
                ...mergedConfig,
                body,
                headers: mergedConfig.headers ?? {},
                signal: controller.signal,
                credentials: mergedConfig.withCredentials ? "include" : "omit",
            });
        } catch (err: any) {
            clearTimeout(timeoutId);
            const code = err.name === "AbortError" ? "ECONNABORTED" : undefined;
            throw new HttpError("Network Error", mergedConfig, undefined, code);
        }

        clearTimeout(timeoutId);

        const contentType = fetchResponse.headers.get("content-type") || "";
        const responseType = mergedConfig.responseType ?? "json";

        let data: any;

        try {
            switch (responseType) {
                case "json":
                    if (contentType.includes("application/json") || contentType.includes("+json")) {
                        data = await fetchResponse.json();
                    } else {
                        data = await fetchResponse.text();
                    }
                    break;

                case "text":
                    data = await fetchResponse.text();
                    break;

                case "blob":
                    data = await fetchResponse.blob();
                    break;

                case "arraybuffer":
                    data = await fetchResponse.arrayBuffer();
                    break;

                case "stream":
                    data = fetchResponse.body;
                    break;

                case "document":
                    const text = await fetchResponse.text();
                    const parser = new DOMParser();
                    const mime = contentType.includes("xml") ? "application/xml" : "text/html";
                    data = parser.parseFromString(text, mime);
                    break;

                default:
                    data = await fetchResponse.text();
            }
        } catch (err) {
            console.warn(`Failed to parse response as ${responseType}:`, err);
            data = null;
        }

        const response: HttpResponse<T> = {
            data,
            status: fetchResponse.status,
            statusText: fetchResponse.statusText,
            headers: fetchResponse.headers,
            config: mergedConfig,
        };

        if (!fetchResponse.ok) {
            const error = new HttpError(
                (typeof data === "object" && data?.message) || `Request failed with status ${fetchResponse.status}`,
                mergedConfig,
                response,
            );

            try {
                return await this.interceptors.response.run(error, true);
            } catch (err) {
                throw err;
            }
        }

        return (await this.interceptors.response.run(response, false)) as HttpResponse<T>;
    }

    public async request<T = any>(config: HttpRequestConfig): Promise<HttpResponse<T>> {
        return this._request<T>(config);
    }

    public get<T = any>(url: string, config?: HttpRequestConfig) {
        return this.request<T>({ ...config, method: "GET", url });
    }

    public post<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
        return this.request<T>({ ...config, method: "POST", url, data });
    }

    public put<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
        return this.request<T>({ ...config, method: "PUT", url, data });
    }

    public patch<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
        return this.request<T>({ ...config, method: "PATCH", url, data });
    }

    public delete<T = any>(url: string, config?: HttpRequestConfig) {
        return this.request<T>({ ...config, method: "DELETE", url });
    }

    public reqwest() {
        const self = this;

        return {
            get<T = any>(url: string, config?: HttpRequestConfig) {
                return tryCatchAsync<HttpResponse<T>, HttpError>(() =>
                    self.request<T>({ ...config, method: "GET", url }),
                );
            },
            post<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
                return tryCatchAsync<HttpResponse<T>, HttpError>(() =>
                    self.request<T>({ ...config, method: "POST", url, data }),
                );
            },
            put<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
                return tryCatchAsync<HttpResponse<T>, HttpError>(() =>
                    self.request<T>({ ...config, method: "PUT", url, data }),
                );
            },
            patch<T = any>(url: string, data?: any, config?: HttpRequestConfig) {
                return tryCatchAsync<HttpResponse<T>, HttpError>(() =>
                    self.request<T>({ ...config, method: "PATCH", url, data }),
                );
            },
            delete<T = any>(url: string, config?: HttpRequestConfig) {
                return tryCatchAsync<HttpResponse<T>, HttpError>(() =>
                    self.request<T>({ ...config, method: "DELETE", url }),
                );
            },
        };
    }
}
