import { fetchProjects, ulid, type Project } from "$lib/data";

// Svelte 5 store for project management
class ProjectStore {
    currentProject = $state<Project | null>(null);
    projects = $state<Project[]>([]);
    loading = $state(false);
    error = $state<string | null>(null);

    // Derived computed property
    get currentProjectId() {
        return this.currentProject?.id ?? null;
    }

    async loadProjects() {
        this.loading = true;
        this.error = null;
        try {
            const data = await fetchProjects();
            this.projects = data;
            // Set first project as current if none selected
            if (!this.currentProject && data.length > 0) {
                this.currentProject = data[0];
            }
        } catch (err) {
            this.error = err instanceof Error ? err.message : "Failed to load projects";
        } finally {
            this.loading = false;
        }
    }

    setCurrentProject(projectId: string) {
        const project = this.projects.find((p) => p.id === projectId);
        if (project) {
            this.currentProject = project;
        }
    }

    async createProject(data: Omit<Project, "id" | "created_at" | "updated_at">) {
        this.loading = true;
        try {
            // Simulate API call
            await new Promise((resolve) => setTimeout(resolve, 500));
            const newProject: Project = {
                ...data,
                id: ulid(),
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString(),
            };
            this.projects = [...this.projects, newProject];
            return newProject;
        } finally {
            this.loading = false;
        }
    }
}

export const projectStore = new ProjectStore();
