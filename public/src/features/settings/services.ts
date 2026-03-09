import { DateTime } from "$lib/shared/datetime";
import { projectsTable, notificationsTable } from "$lib/db/tables";
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
        const dto = await projectsTable.findOne(id);
        return dto ? SettingsService.toProject(dto as unknown as ProjectDTO) : null;
    }

    async getProjects(): Promise<ProjectSettings[]> {
        const result = await projectsTable.findMany({ sort: { field: "name", order: "asc" }, pageSize: 100 });
        return result.data.map((d) => SettingsService.toProject(d as unknown as ProjectDTO));
    }

    async updateProject(id: string, data: UpdateProjectDTO): Promise<ProjectSettings> {
        const dto = await projectsTable.update(id, {
            ...data,
            updated_at: new Date().toISOString(),
        } as any);
        return SettingsService.toProject(dto as unknown as ProjectDTO);
    }

    async getNotifications(userId: string): Promise<AppNotification[]> {
        const result = await notificationsTable.findMany({
            where: { user_id: userId } as any,
            sort: { field: "created_at", order: "desc" },
            pageSize: 50,
        });
        return result.data.map((d) => SettingsService.toNotification(d as unknown as NotificationDTO));
    }

    async markNotificationRead(id: string): Promise<void> {
        await notificationsTable.update(id, { status: "read" } as any);
    }

    async dismissNotification(id: string): Promise<void> {
        await notificationsTable.update(id, { status: "dismissed" } as any);
    }
}

export const settingsService = new SettingsService();
