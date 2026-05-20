import { DateTime } from "$lib/shared/datetime";
import { httpClient } from "$lib/api/client";
import type { ProjectDTO, UpdateProjectDTO, NotificationDTO } from "./dtos";
import type { ProjectSettings, AppNotification } from "./types";

export class SettingsService {
    static toProject(dto: ProjectDTO): ProjectSettings {
        return {
            id: dto.id,
            name: dto.name,
            slug: dto.slug,
            ownerId: dto.owner_id,
            settings: dto.settings,
            locale: dto.locale,
            createdAt: new DateTime(dto.created_at),
            updatedAt: new DateTime(dto.updated_at),
        };
    }

    static toNotification(dto: NotificationDTO): AppNotification {
        return {
            id: dto.id,
            title: dto.title,
            body: dto.body,
            datetime: new DateTime(dto.datetime),
            status: dto.status,
            type: dto.type,
            userId: dto.user_id,
            projectId: dto.project_id,
            createdAt: new DateTime(dto.created_at),
        };
    }

    async getProject(id: string): Promise<ProjectSettings | null> {
        const dto = await httpClient.get<ProjectDTO | null>(`/projects/${id}`);
        return dto ? SettingsService.toProject(dto as unknown as ProjectDTO) : null;
    }

    async getProjects(): Promise<ProjectSettings[]> {
        const result = await httpClient.get<{ data: ProjectDTO[] }>("/projects", {
            pageSize: 100,
        });

        return result.data.map((d) => SettingsService.toProject(d));
    }

    async updateProject(id: string, data: UpdateProjectDTO): Promise<ProjectSettings> {
        const dto = await httpClient.patch<ProjectDTO>(`/projects/${id}`, data);
        return SettingsService.toProject(dto as unknown as ProjectDTO);
    }

    async getNotifications(userId: string): Promise<AppNotification[]> {
        const result = await httpClient.get<{ data: NotificationDTO[] }>("/notifications", {
            userId,
        });
        return result.data.map((d) => SettingsService.toNotification(d));
    }

    async markNotificationRead(id: string): Promise<void> {
        await httpClient.patch(`/notifications/${id}`, { status: "read" });
    }

    async dismissNotification(id: string): Promise<void> {
        await httpClient.patch(`/notifications/${id}`, { status: "dismissed" });
    }
}

export const settingsService = new SettingsService();
