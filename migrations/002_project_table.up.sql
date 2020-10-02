CREATE TABLE public.project (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar NOT NULL,
    archived bool NOT NULL DEFAULT false,
    CONSTRAINT project_pk PRIMARY KEY (id)
);
