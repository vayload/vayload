export interface UserDTO {
    id: string;
    first_name: string;
    last_name: string;
    username: string;
    email: string;
    avatar_url: string | null;
    is_super_admin: boolean;
    last_sign_in_at: string; // ISO string
    created_at: string; // ISO string
    updated_at: string; // ISO string
    password?: string; // stripped before returning to UI
}

export interface LoginResponseDTO {
    user: UserDTO;
    token: string;
}
