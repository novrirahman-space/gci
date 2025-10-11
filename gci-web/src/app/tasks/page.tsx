"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createTask,
  deleteTask,
  getClasses,
  getTasks,
  updateTask,
} from "@/lib/api";
import type { Task } from "@/lib/types";
import { useState, useMemo } from "react";
import TaskForm from "@/components/TaskForm";

export default function TasksPage() {
  const qc = useQueryClient();
  const {
    data: tasks,
    isLoading,
    error,
  } = useQuery({ queryKey: ["tasks"], queryFn: getTasks });
  const { data: classes } = useQuery({
    queryKey: ["classes"],
    queryFn: getClasses,
  });

  const classMap = useMemo(() => {
    const m = new Map<string, string>();
    (classes || []).forEach((c) => m.set(c.id, c.class_name));
    return m;
  }, [classes]);

  const [editing, setEditing] = useState<Task | null>(null);

  const createMut = useMutation({
    mutationFn: createTask,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tasks"] }),
  });

  const updateMut = useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: { class_id: string; title: string; description: string };
    }) => updateTask(id, payload),
    onSuccess: () => {
      setEditing(null);
      qc.invalidateQueries({ queryKey: ["tasks"] });
    },
  });

  const deleteMut = useMutation({
    mutationFn: deleteTask,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["tasks"] }),
  });

  if (isLoading) return <p>Loading…</p>;
  if (error)
    return <p className="text-red-600">Error: {(error as Error).message}</p>;

  return (
    <main className="space-y-6">
      <h2 className="text-xl font-semibold">Tasks</h2>

      <TaskForm
        classes={classes || []}
        onSubmit={async (payload) => {
          await createMut.mutateAsync(payload);
        }}
        submitLabel="Create"
      />

      {editing && (
        <div className="border rounded p-3">
          <p className="mb-2 font-medium">Edit Task</p>
          <TaskForm
            classes={classes || []}
            initial={{
              class_id: editing.class_id,
              title: editing.title,
              description: editing.description,
            }}
            onSubmit={async (payload) => {
              await updateMut.mutateAsync({ id: editing.id, payload });
            }}
            submitLabel="Update"
          />
          <button
            className="mt-3 text-sm underline"
            onClick={() => setEditing(null)}
          >
            Cancel
          </button>
        </div>
      )}

      <ul className="divide-y">
        {tasks?.map((t) => (
          <li key={t.id} className="py-3 flex items-center justify-between">
            <div>
              <div className="font-medium">{t.title}</div>
              <div className="text-sm text-gray-600">
                {classMap.get(t.class_id) || t.class_id} — {t.description}
              </div>
            </div>
            <div className="flex gap-3">
              <button className="underline" onClick={() => setEditing(t)}>
                Edit
              </button>
              <button
                className="underline text-red-600"
                onClick={() => deleteMut.mutate(t.id)}
              >
                Delete
              </button>
            </div>
          </li>
        ))}
      </ul>
    </main>
  );
}
