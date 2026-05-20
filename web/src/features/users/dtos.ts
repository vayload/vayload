import type { UserDTO } from "../auth/dtos";
export type { UserDTO };

export interface CreateUserDTO {
    first_name: string;
    last_name: string;
    username: string;
    email: string;
    avatar_url?: string | null;
    is_super_admin?: boolean;
    password: string;
}

export interface UpdateUserDTO {
    first_name?: string;
    last_name?: string;
    username?: string;
    email?: string;
    avatar_url?: string | null;
    is_super_admin?: boolean;
}
