import { ulid, type Project, type ProjectInput, type User } from "$lib/data";

export const slugify = (str: string): string => {
    return str
        .replace(/^\s+|\s+$/g, "")
        .toLowerCase()
        .replace(/[^a-z0-9 -]/g, "")
        .replace(/\s+|-+/g, "-");
};

const PROJECTS_KEY = "projects";
const CURRENT_PROJECT_KEY = "currentProjectId";

export class AppContext {
    private _project: Project | null = $state(null);
    private _projects: Project[] = $state([]);
    private _initialized = $state(false);
    private _loadings = $state({
        projects: false,
        project: false,
    });

    private _errors = $state({
        projects: null as string | null,
        project: null as string | null,
    });

    public get currentProject() {
        return this._project;
    }

    public get currentProjectId() {
        return this._project?.id ?? null;
    }

    public get projects() {
        return this._projects;
    }

    public get loadings() {
        return this._loadings;
    }

    public get errors() {
        return this._errors;
    }

    public get haveProjects() {
        return this._projects.length > 0;
    }

    public get needsOnboarding() {
        return this._initialized && !this._loadings.projects && this._projects.length === 0;
    }

    public get initialized() {
        return this._initialized;
    }

    /**
     * Fetches all projects from localStorage (simulated persistence).
     * Auto-selects the previously selected project or the first one.
     */
    public async fetchProjects() {
        try {
            this._loadings.projects = true;
            this._errors.projects = null;

            const data = await new Promise<Project[]>((resolve) => {
                setTimeout(() => {
                    const raw = localStorage.getItem(PROJECTS_KEY);
                    resolve(raw ? JSON.parse(raw) : []);
                }, 800);
            });

            this._projects = data;

            // Auto-select: restore saved selection or pick the first project
            if (data.length > 0) {
                const savedId = localStorage.getItem(CURRENT_PROJECT_KEY);
                const saved = savedId ? data.find((p) => p.id === savedId) : null;
                this._project = saved ?? data[0];
                localStorage.setItem(CURRENT_PROJECT_KEY, this._project!.id);
            } else {
                this._project = null;
            }
        } catch (error) {
            console.error("Failed to fetch projects", error);
            this._errors.projects = error instanceof Error ? error.message : "Failed to fetch projects";
        } finally {
            this._loadings.projects = false;
            this._initialized = true;
        }
    }

    /**
     * Creates a new project and selects it as the current project.
     */
    public async createProject(input: ProjectInput): Promise<Project | null> {
        try {
            this._loadings.project = true;
            this._errors.project = null;

            const user = JSON.parse(localStorage.getItem("auth") || "{}") as User;

            const newProject: Project = await new Promise<Project>((resolve) => {
                setTimeout(() => {
                    const project: Project = {
                        id: ulid(),
                        name: input.name,
                        slug: slugify(input.name),
                        owner_id: user.id,
                        settings: input.settings,
                        locale: input.locale,
                        created_at: new Date().toISOString(),
                        updated_at: new Date().toISOString(),
                    };

                    const store = [...this._projects, project];
                    localStorage.setItem(PROJECTS_KEY, JSON.stringify(store));
                    resolve(project);
                }, 800);
            });

            this._projects = [...this._projects, newProject];
            this.selectProject(newProject.id);

            return newProject;
        } catch (error) {
            console.error("Failed to create project", error);
            this._errors.project = error instanceof Error ? error.message : "Failed to create project";
            return null;
        } finally {
            this._loadings.project = false;
        }
    }

    /**
     * Select a project by ID and persist the choice.
     */
    public selectProject(id: string) {
        const project = this._projects.find((p) => p.id === id);
        if (project) {
            this._project = project;
            localStorage.setItem(CURRENT_PROJECT_KEY, id);
        }
    }

    public async deleteProject(id: string) {
        try {
            this._loadings.project = true;
            this._errors.project = null;

            await new Promise((resolve) => setTimeout(resolve, 800));

            const store = this._projects.filter((p) => p.id !== id);
            localStorage.setItem(PROJECTS_KEY, JSON.stringify(store));

            this._projects = store;
            this.selectProject(store[0]?.id ?? null);

            return true;
        } catch (error) {
            console.error("Failed to delete project", error);
            this._errors.project = error instanceof Error ? error.message : "Failed to delete project";
            return false;
        } finally {
            this._loadings.project = false;
        }
    }
}

export const appContext = new AppContext();
