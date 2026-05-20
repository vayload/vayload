import { handleRequest } from "./backend-mock";

export type ApiMethod = "GET" | "POST" | "PATCH" | "DELETE" | "PUT";

export interface ApiRequest {
    method: ApiMethod;
    path: string;
    query?: Record<string, string | number | boolean | undefined>;
    body?: unknown;
}

export const httpClient = {
    get: <T>(path: string, query?: Record<string, string | number | boolean | undefined>) =>
        handleRequest<T>({ method: "GET", path, query }),
    post: <T>(path: string, body?: unknown) => handleRequest<T>({ method: "POST", path, body }),
    patch: <T>(path: string, body?: unknown) => handleRequest<T>({ method: "PATCH", path, body }),
    delete: <T>(path: string, query?: Record<string, string | number | boolean | undefined>) =>
        handleRequest<T>({ method: "DELETE", path, query }),
    put: <T>(path: string, body?: unknown) => handleRequest<T>({ method: "PUT", path, body }),
};
