create table cloud_vision_cache (
    id uuid primary key default gen_random_uuid(),
    image_hash varchar(255) not null,
    response text[] not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create table genai_cache (
    id uuid primary key default gen_random_uuid(),
    image_hash varchar(255) not null,
    response jsonb not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);