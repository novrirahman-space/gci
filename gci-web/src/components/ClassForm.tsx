"use client";

import { on } from "events";
import { use, useState } from "react";

type Props = {
    initial?: { class_name: string; teacher: string };
    onSubmit: (payload: { class_name: string; teacher: string }) => Promise<void> | void;
    submitLabel?: string;
};

export default function ClassForm({ initial, onSubmit, submitLabel = "Save" }: Props) {
    const [className, setClassName] = useState(initial?.class_name ?? "");
    const [teacher, setTeacher] = useState(initial?.teacher ?? "");

    return (
        <form
            onSubmit={async (e) => {
                e.preventDefault();
                if (!className.trim() || !teacher.trim()) return;
                await onSubmit({ class_name: className.trim(), teacher: teacher.trim() });
                setClassName("");
                setTeacher("");
            }}
            className="flex gap-2 flex-wrap items-center"
        >
            <input
                className="border rounded px-3 py-2"
                placeholder="Class name"
                value={className}
                onChange={(e) => setClassName(e.target.value)}
            />
            <input
                className="border rounded px-3 py-2"
                placeholder="Teacher"
                value={teacher}
                onChange={(e) => setTeacher(e.target.value)}
            />
            <button className="bg-black text-white rounded px-4 py-2" type="submit">
                {submitLabel}
            </button>
        </form>
    );
}