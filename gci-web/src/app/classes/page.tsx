"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { createClass, deleteClass, getClasses, updateClass } from "@/lib/api";
import type { Class } from "@/lib/types";
import { useState } from "react";
import ClassForm from "@/components/ClassForm";

export default function ClassesPage() {
  const qc = useQueryClient();
  const {
    data: classes,
    isLoading,
    error,
  } = useQuery({ queryKey: ["classes"], queryFn: getClasses });

  const [editing, setEditing] = useState<Class | null>(null);

  const createMut = useMutation({
    mutationFn: createClass,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["classes"] }),
  });

  const updateMut = useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: { class_name: string; teacher: string };
    }) => updateClass(id, payload),
    onSuccess: () => {
      setEditing(null);
      qc.invalidateQueries({ queryKey: ["classes"] });
    },
  });

  const deleteMut = useMutation({
    mutationFn: deleteClass,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["classes"] }),
  });

  if (isLoading) return <p>Loading...</p>;
  if (error)
    return <p className="text-red-600">Error: {(error as Error).message}</p>;

  return (
    <main className="space-y-6">
      <h2 className="text-xl font-semibold">Classes</h2>

      {!editing ? (
        <ClassForm
          onSubmit={async (payload) => {
            await createMut.mutateAsync(payload);
          }}
          submitLabel="Create"
        />
      ) : (
        <div className="border rounded p-3">
          <p className="mb-2 font-medium">Edit Class</p>
          <ClassForm
            initial={{
              class_name: editing.class_name,
              teacher: editing.teacher,
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
        {classes?.map((c) => (
          <li key={c.id} className="py-3 flex items-center justify-between">
            <div>
              <div className="font-medium">{c.class_name}</div>
              <div className="text-sm text-gray-600">{c.teacher}</div>
            </div>
            <div className="flex gap-3">
              <button className="underline" onClick={() => setEditing(c)}>
                Edit
              </button>
              <button
                className="underline text-red-600"
                onClick={() => deleteMut.mutate(c.id)}
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
