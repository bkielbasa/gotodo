create extension "uuid-ossp";

CREATE TABLE public.todos (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar NOT NULL,
    done bool NOT NULL DEFAULT false,
    CONSTRAINT todos_pk PRIMARY KEY (id)
);
