export type Class = {
    id: string;
    class_name: string;
    teacher: string;
};

export type Task = {
    id: string;
    class_id: string;
    title: string;
    description: string;
};