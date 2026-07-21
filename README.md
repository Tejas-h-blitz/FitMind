# FitMind - AI Personal Health Advisor
  
FitMind is a production-grade personal health advisor that implements a Retrieval-Augmented Generation (RAG) pipeline to analyze user's uploaded health documents (blood reports, diet charts, workout logs, medical summaries) and stream contextual, citation-referenced answers to their queries. It also includes health status tracking features (consecutive activity streaks, BMI log history, and goal management).

---
  
## Architecture Diagram

```
                 +---------------------------------------+
                 |          SVELTEKIT FRONTEND           |
                 |      (Port 5173 / Vercel Edge)        |
                 +-------------------+-------------------+
                                     |
                         HTTPS/REST  |  SSE (Server Sent Events)
                                     v
                 +-------------------+-------------------+
                 |           GO CHI BACKEND              |  <----+ (Verifies Supabase JWT)
                 |         (Port 8080 / Fly.io)          |
                 +----+--------------+--------------+----+
                      |              |              |
           PostgREST  |              | Ingest       | Stream Query
           Storage    |              | (Async HTTP) | (HTTP SSE)
                      v              v              v
                 +----+---+    +-----+--------------+----+
                 |        |    |      RAG SERVICE        |
                 |  S     |    |   (FastAPI / Railway)   |
                 |  U     |    +----+---------------+----+
                 |  P     |         |               |
                 |  A     |         | Vector Ops    | ChatCompletion
                 |  B     |         v               v
                 |  A     |    +----+----+     +----+----+
                 |  S     |    | QDRANT  |     | OPENAI  |
                 |  E     |    | (Docker |     |   API   |
                 |        |    |  6333)  |     | (gpt-4o)|
                 +--------+    +---------+     +---------+
```

---

## Technical Stack

- **Frontend**: SvelteKit + TailwindCSS (v4) + Lucide Icons + Svelte-Sonner (Toasts)
- **Backend API**: Go (Golang) + Chi Router + Go-JWT
- **RAG Microservice**: Python + FastAPI + LlamaIndex + PyMuPDF
- **Vector DB**: Qdrant (Self-hosted via Docker)
- **Storage/DB/Auth**: Supabase (Postgres, Storage, Go JWT verification)

---

## Prerequisites

Ensure you have the following installed on your machine:
1. [Docker & Docker Compose](https://www.docker.com/products/docker-desktop/)
2. [Node.js (v18+)](https://nodejs.org/) & `npm`
3. A [Supabase](https://supabase.com) account & project.
4. An [OpenAI API Key](https://platform.openai.com) with credits.

---

## Local Setup Instructions

### 1. Supabase Schema Configuration
Create the following tables in the Supabase SQL editor:

```sql
-- Enable UUID generation extension
create extension if not exists "uuid-ossp";

-- Users profile (extends Supabase auth.users)
create table if not exists user_profiles (
  id uuid references auth.users on delete cascade primary key,
  full_name text,
  streak_count int default 0,
  last_active date,
  created_at timestamptz default now()
);

-- Documents
create table if not exists documents (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users on delete cascade not null,
  name text not null,
  size bigint,
  storage_path text not null,
  status text default 'pending', -- pending/processing/ready/failed
  created_at timestamptz default now()
);

-- Chat messages
create table if not exists messages (
  id uuid primary key default gen_random_uuid(),
  doc_id uuid references documents on delete cascade not null,
  user_id uuid references auth.users on delete cascade not null,
  role text not null, -- user/assistant
  content text not null,
  sources jsonb, -- array of {text, page_number}
  created_at timestamptz default now()
);

-- Health metrics
create table if not exists health_metrics (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users on delete cascade not null,
  bmi float not null,
  height float not null,
  weight float not null,
  recorded_at timestamptz default now()
);

-- Goals
create table if not exists goals (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users on delete cascade not null,
  title text not null,
  target_date date,
  status text default 'active', -- active/completed
  created_at timestamptz default now()
);

-- Enable Row Level Security (RLS)
alter table user_profiles enable row level security;
alter table documents enable row level security;
alter table messages enable row level security;
alter table health_metrics enable row level security;
alter table goals enable row level security;

-- Policies (Allow read/write only to the owner)
create policy "Users can read own profile" on user_profiles for select using (auth.uid() = id);
create policy "Users can update own profile" on user_profiles for update using (auth.uid() = id);
create policy "Users can insert own profile" on user_profiles for insert with check (auth.uid() = id);
create policy "Users can CRUD own documents" on documents for all using (auth.uid() = user_id);
create policy "Users can CRUD own messages" on messages for all using (auth.uid() = user_id);
create policy "Users can CRUD own health metrics" on health_metrics for all using (auth.uid() = user_id);
create policy "Users can CRUD own goals" on goals for all using (auth.uid() = user_id);

-- User auto-creation trigger from auth.users
create or replace function public.handle_new_user()
returns trigger as $$
begin
  insert into public.user_profiles (id, full_name, streak_count, last_active)
  values (new.id, coalesce(new.raw_user_meta_data->>'full_name', 'Health Enthusiast'), 1, current_date);
  return new;
end;
$$ language plpgsql security definer;

create or replace trigger on_auth_user_created
  after insert on auth.users
  for each row execute procedure public.handle_new_user();

-- New Tables for Tier 1 AI Features

create table if not exists document_analyses (
  id uuid primary key default gen_random_uuid(),
  doc_id uuid references documents on delete cascade not null unique,
  user_id uuid references auth.users not null,
  metrics jsonb not null,
  summary text not null,
  overall_status text not null,
  created_at timestamptz default now()
);

alter table document_analyses enable row level security;
create policy "Users can only access own analyses"
  on document_analyses for all
  using (auth.uid() = user_id);

create table if not exists meal_plans (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users not null,
  doc_id uuid references documents on delete cascade,
  reasoning text,
  daily_calories_target int,
  protein_target_g int,
  days jsonb not null,
  weekly_notes text,
  dietary_preference text,
  created_at timestamptz default now()
);

alter table meal_plans enable row level security;
create policy "Users can only access own meal plans"
  on meal_plans for all
  using (auth.uid() = user_id);

create table if not exists workout_plans (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users not null,
  program_name text,
  duration_weeks int,
  reasoning text,
  weekly_schedule jsonb not null,
  progression_notes text,
  safety_notes text,
  fitness_level text,
  equipment text,
  days_per_week int,
  created_at timestamptz default now()
);

alter table workout_plans enable row level security;
create policy "Users can only access own workout plans"
  on workout_plans for all
  using (auth.uid() = user_id);
```

*Note: Ensure you create a storage bucket in Supabase called `documents` and configure its permissions to allow authenticated uploads under user folders.*

### 2. Environment Variables Configuration
Fill in the blanks inside the `.env` files of each folder:

#### `frontend/.env`:
```env
PUBLIC_SUPABASE_URL=https://[YOUR_PROJECT_REF].supabase.co
PUBLIC_SUPABASE_ANON_KEY=[YOUR_SUPABASE_ANON_KEY]
PUBLIC_API_BASE_URL=http://localhost:8080
```

#### `backend/.env`:
```env
SUPABASE_URL=https://[YOUR_PROJECT_REF].supabase.co
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SUPABASE_SERVICE_ROLE_KEY]
SUPABASE_JWT_SECRET=[YOUR_JWT_SECRET_FROM_SUPABASE_SETTINGS]
QDRANT_URL=http://qdrant:6333
RAG_SERVICE_URL=http://rag-service:8000
PORT=8080
```

#### `rag-service/.env`:
```env
OPENAI_API_KEY=sk-proj-[YOUR_OPENAI_KEY]
QDRANT_URL=http://qdrant:6333
SUPABASE_URL=https://[YOUR_PROJECT_REF].supabase.co
SUPABASE_SERVICE_ROLE_KEY=[YOUR_SUPABASE_SERVICE_ROLE_KEY]
```

### 3. Execution
Launch Qdrant, the Go API backend, and the Python RAG microservice:
```bash
docker-compose up --build
```

In a separate terminal, launch the SvelteKit development server:
```bash
cd frontend
npm run dev
```
Open `http://localhost:5173` in your browser.

---

## API Endpoint Reference

All endpoints (except `/health`) require authentication and expect a valid Supabase JWT Bearer token in the `Authorization` header or inside the `sb-access-token` cookie.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **POST** | `/api/auth/verify` | Validates JWT, increments daily streak, returns user profile. |
| **GET** | `/api/documents` | Lists all documents uploaded by the user. |
| **POST** | `/api/documents/upload` | Uploads PDF multipart file -> Storage -> db (triggers ingest). |
| **DELETE**| `/api/documents/:id` | Deletes file from storage, db record, and Qdrant collection. |
| **GET** | `/api/documents/:id/status` | Returns document indexing status (`ready`, `processing`, etc). |
| **GET** | `/api/chat/:docId/history` | Retrieves historical conversations for the document. |
| **DELETE**| `/api/chat/:docId/history` | Clears all messages associated with the document. |
| **POST** | `/api/chat/:docId/query` | Sends query, returns chunked SSE stream of answer & citations. |
| **GET** | `/api/health/metrics` | Fetches historical BMI logs for graph generation. |
| **POST** | `/api/health/metrics` | Computes and logs height, weight, and BMI. |
| **GET** | `/api/health/goals` | Lists user's health goals (active and completed). |
| **POST** | `/api/health/goals` | Creates a new milestone target. |
| **PATCH** | `/api/health/goals/:id` | Toggles goal status between `active` and `completed`. |

---

## Production Deployment Guide

### Go Backend (Fly.io)
1. Install flyctl and authenticate.
2. Initialize app in `/backend`:
   ```bash
   fly launch
   ```
3. Set your production secrets:
   ```bash
   fly secrets set SUPABASE_URL="..." SUPABASE_SERVICE_ROLE_KEY="..." SUPABASE_JWT_SECRET="..." QDRANT_URL="..." RAG_SERVICE_URL="..."
   ```
4. Deploy:
   ```bash
   fly deploy
   ```

### RAG Microservice (Railway)
1. Initialize a new service on Railway connected to the `/rag-service` repository directory.
2. Set the variables: `OPENAI_API_KEY`, `SUPABASE_URL`, `SUPABASE_SERVICE_ROLE_KEY`, `QDRANT_URL` (points to your production Qdrant host).
3. Deploy. Railway automatically detects the `Dockerfile` and compiles the Python app.

### SvelteKit Frontend (Vercel)
1. Connect Vercel to your repository.
2. Configure environment variables (`PUBLIC_SUPABASE_URL`, `PUBLIC_SUPABASE_ANON_KEY`, `PUBLIC_API_BASE_URL`).
3. Deploy.
