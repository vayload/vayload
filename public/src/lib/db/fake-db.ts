import { faker } from "@faker-js/faker";
import { ulid } from "./utils";

/* ================= USERS ================= */

export const users = Array.from({ length: 15 }).map((_, i) => ({
    id: ulid(),
    first_name: faker.person.firstName(),
    last_name: faker.person.lastName(),
    username: faker.internet.username(),
    email: faker.internet.email(),
    avatar_url: faker.image.avatar(),
    is_super_admin: i === 0,
    last_sign_in_at: faker.date.recent(),
    created_at: faker.date.past(),
    updated_at: faker.date.recent(),
}));

export const currentUser = users[0];

/* ================= PROJECTS ================= */

export const projects = Array.from({ length: 3 }).map((_, i) => ({
    id: ulid(),
    name: faker.company.name(),
    slug: faker.helpers.slugify(faker.company.name()),
    owner_id: currentUser.id,
    settings: {},
    locale: "es",
    created_at: faker.date.past(),
    updated_at: faker.date.recent(),
}));

/* ================= PROJECT MEMBERS ================= */

export const projectMembers = projects.flatMap((project) =>
    users.map((user) => ({
        project_id: project.id,
        user_id: user.id,
        role_id: user.id === currentUser.id ? 1 : 2,
    })),
);

/* ================= COLLECTIONS ================= */

export const collections = Array.from({ length: 6 }).map((_, i) => ({
    id: ulid(),
    project_id: faker.helpers.arrayElement(projects).id,
    name: faker.commerce.productName(),
    slug: faker.helpers.slugify(faker.commerce.productName()),
    fields_schema: {},
    settings: {},
    created_at: faker.date.past(),
}));

/* ================= ENTRIES ================= */

export const entries = Array.from({ length: 25 }).map((_, i) => ({
    id: ulid(),
    collection_id: faker.helpers.arrayElement(collections).id,
    title: faker.lorem.sentence(),
    status: faker.helpers.arrayElement(["draft", "published"]),
    author_id: currentUser.id,
    created_at: faker.date.past(),
    updated_at: faker.date.recent(),
}));

/* ================= FILE OBJECTS ================= */

export const fileObjects = Array.from({ length: 20 }).map((_, i) => ({
    id: ulid(),
    owner_id: currentUser.id,
    project_id: faker.helpers.arrayElement(projects).id,
    name: faker.system.fileName(),
    mime_type: faker.system.mimeType(),
    category: faker.helpers.arrayElement(["image", "video", "document"]),
    size: faker.number.int({ min: 1000, max: 5_000_000 }),
    key: faker.system.filePath(),
    folder: faker.book.series(),
    created_at: faker.date.past(),
    updated_at: faker.date.recent(),
}));

/* ================= NOTIFICATIONS ================= */

export const notifications = Array.from({ length: 10 }).map((_, i) => ({
    id: ulid(),
    title: faker.lorem.sentence(),
    body: faker.lorem.paragraph(),
    status: faker.helpers.arrayElement(["read", "unread"]),
    user_id: currentUser.id,
    project_id: faker.helpers.arrayElement(projects).id,
    created_at: faker.date.recent(),
}));
