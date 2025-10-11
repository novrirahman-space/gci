"use client";

import { useState } from "react";
import type { Class } from "@/lib/types";

type Props = {
    classes: Class[];
    initial?: { class_id: string; title: string; description: string };
    onSubmit: (payload: { class_id: string; title: string; description: string }) => Promise<void> | void;
    submitLabel?: string;
};

export default function TaskForm({ classes, initial, onSubmit, submitLabel = "save" }: Props) {
    const [classId, setClassId] = useState(initial?.class_id ?? (classes[0]?.id || ""));
    const [title, setTitle] = useState(initial?.title ?? "");
    const [desc, setDesc] = useState(initial?.description ?? "");

    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                if (!classId || !title.trim()) return;
                await onSubmit({ class_id: classId, title: title.trim(), description: desc.trim() });
                setTitle("");
                setDesc("")
            }}
            className="flex gap-2 flex-wrap items-center"
        >
            <select className="border rounded px-3 py-2" value={classId} onChange={(e) => setClassId(e.target.value)}>
                {classes.map((c) => (
                    <option key={c.id} value={c.id}>
                        {c.class_name}
                    </option>
                ))}
            </select>
            <input
                className="border rounded px-3 py-2"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <input
                className="border rounded px-3 py-2"
                placeholder="Description"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
            />
            <button className="bg-black text-white rounded px-4 py-2" type="submit">
                {submitLabel}
            </button>
        </form>
    );
}