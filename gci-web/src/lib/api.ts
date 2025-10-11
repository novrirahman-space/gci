const BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

async function api<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${BASE}${path}`, {
        ...init,
        headers: {
            "Content-Type": "application/json",
            ...(init?.headers || {}),
        },
        // penting agar selalu ambil data terbaru saat dev
        cache: "no-store"
    });

    if (!res.ok) {
        // coba ambil pesan error dari server
        let msg = `HTTP ${res.status}`;
        try {
            const data = await res.json();
            msg = data?.message || msg;
        } catch {
            // ignore json parse error
        }
        throw new Error(msg);
    }
    return res.json() as Promise<T>;
}

export const getClasses = () => api<import("./types").Class[]>("/classes");
export const createClass = (payload: { class_name: string; teacher: string }) => api<import("./types").Class>("/classes", { method: "POST", body: JSON.stringify(payload) });
export const updateClass = (id: string, payload: { class_name: string; teacher: string }) => api<import("./types").Class>(`/classes/${id}`, { method: "PUT", body: JSON.stringify(payload) });
export const deleteClass = (id: string) => api<{ message: string }>(`/classes/${id}`, { method: "DELETE" });

export const getTasks = () => api<import("./types").Task[]>("/tasks");
export const createTask = (payload: { class_id: string; title: string; description: string }) => api<import("./types").Task>("/tasks", { method: "POST", body: JSON.stringify(payload) });
export const updateTask = (id: string, payload: { class_id: string; title: string; description: string }) => api<import("./types").Task>(`/tasks/${id}`, { method: "PUT", body: JSON.stringify(payload) });
export const deleteTask = (id: string) => api<{ message: string }>(`/tasks/${id}`, { method: "DELETE" });

export const getTasksByClass = (classId: string) => api<import("./types").Task[]>(`/classes/${classId}/tasks`);