# Product Backlog (Epics & Tasks)

Dokumen ini berfungsi sebagai pengganti Jira / Trello / Kanban board kita. Di sini, seluruh *Development Roadmap* dari `blueprint.md` dipecah menjadi **Epic** dan **Task** yang siap dieksekusi secara berurutan.

Setiap *task* yang selesai akan ditandai dengan `[x]`.

---

## Epic 1: Foundation & Infrastructure
**Goal:** Membangun fondasi monorepo dan menyalakan semua *database/broker* secara lokal.
- [x] **TASK-101:** Membuat struktur *folder* monorepo untuk 7 service.
- [x] **TASK-102:** Meracik `docker-compose.yml` untuk PostgreSQL (4 database terpisah).
- [x] **TASK-103:** Menambahkan Cassandra, Redis, RabbitMQ, dan OpenSearch ke Docker.
- [x] **TASK-104:** Menambahkan *stack* Observability (Prometheus, Loki, Grafana) ke Docker.

---

## Epic 2: Identity & Edge
**Goal:** Mengamankan pintu masuk dan mengatur otentikasi JWT.
- [ ] **TASK-201:** Inisialisasi API Gateway (Go/Fiber).
- [ ] **TASK-202:** Implementasi Redis Rate Limiter di Gateway (Bypass support).
- [ ] **TASK-203:** Inisialisasi User Service (Node.js/NestJS) & koneksi DB.
- [ ] **TASK-204:** Implementasi endpoint Registrasi (`POST /register`) & bcrypt.
- [ ] **TASK-205:** Implementasi endpoint Login (`POST /login`) & JWT Generation.
- [ ] **TASK-206:** Implementasi Stateless JWT Validation & Header Injection di Gateway.

---

## Epic 3: Catalog & Search (CQRS)
**Goal:** Mengelola data buku dan sinkronisasi otomatis ke mesin pencari.
- [ ] **TASK-301:** Inisialisasi Product Service (Go/Gin) & koneksi DB.
- [ ] **TASK-302:** Implementasi CRUD Katalog Buku (Admin).
- [ ] **TASK-303:** Inisialisasi Product Search Service (Python/FastAPI) & koneksi OpenSearch.
- [ ] **TASK-304:** Implementasi RabbitMQ Publisher di Product Service (`ProductUpdated`).
- [ ] **TASK-305:** Implementasi RabbitMQ Consumer di Search Service (Indexing ke OpenSearch).
- [ ] **TASK-306:** Implementasi endpoint Pencarian (`GET /search`).

---

## Epic 4: Shopping Cart
**Goal:** Membuat keranjang belanja abadi yang tangguh dengan NoSQL.
- [ ] **TASK-401:** Inisialisasi Cart Service (Node.js/Express) & koneksi Cassandra.
- [ ] **TASK-402:** Implementasi logika *Upsert* (Add to Cart / Update Qty).
- [ ] **TASK-403:** Implementasi endpoint Get Cart & Clear Cart.

---

## Epic 5: Transactions & Soft-Booking
**Goal:** Jantung e-commerce: validasi stok, *locking*, dan pembuatan pesanan.
- [ ] **TASK-501:** Inisialisasi Order Service (PHP/Laravel) & koneksi DB.
- [ ] **TASK-502:** Implementasi HTTP Client untuk validasi stok ke Product Service (Sync).
- [ ] **TASK-503:** Implementasi logika Redis Soft-Booking (TTL 15 menit).
- [ ] **TASK-504:** Implementasi endpoint Checkout (`POST /checkout`).
- [ ] **TASK-505:** Implementasi Cron Job (Laravel Scheduler) untuk menyapu pesanan *Expired*.

---

## Epic 6: Payments & Async Events
**Goal:** Menerima pembayaran dan memicu reaksi berantai di latar belakang.
- [ ] **TASK-601:** Inisialisasi Payment Service (Node.js/Express).
- [ ] **TASK-602:** Implementasi integrasi Xendit API (Create Invoice).
- [ ] **TASK-603:** Implementasi validasi Xendit Webhook Callback.
- [ ] **TASK-604:** Implementasi RabbitMQ Publisher (`OrderPaid`).
- [ ] **TASK-605:** Implementasi RabbitMQ Consumer di Order Service (Update status ke PAID).
- [ ] **TASK-606:** Inisialisasi Notification Service (Go/Python) tanpa HTTP server.
- [ ] **TASK-607:** Implementasi RabbitMQ Consumer & pengiriman Email SMTP.

---

## Epic 7: Observability & Load Testing
**Goal:** Memantau kesehatan sistem dan menyiksanya hingga batas maksimal.
- [ ] **TASK-701:** Standardisasi format JSON Logging di semua 7 service.
- [ ] **TASK-702:** Setup *Dashboard* Grafana (Membaca metrik Prometheus & log Loki).
- [ ] **TASK-703:** Menulis *script* k6 untuk Load Testing (Trafik Normal).
- [ ] **TASK-704:** Menulis *script* k6 untuk Spike Testing (Flash Sale dengan Bypass Header).

---

## Epic 8: Frontend Web (Svelte SPA)
**Goal:** Membangun antarmuka pengguna yang super cepat tanpa *Virtual DOM*.
- [ ] **TASK-801:** Inisialisasi proyek Svelte + Vite + Tailwind CSS (`frontend-web/`).
- [ ] **TASK-802:** Konfigurasi HTTP Client (Axios/Fetch) untuk menembak API Gateway.
- [ ] **TASK-803:** Implementasi Halaman Katalog (Grid Buku).
- [ ] **TASK-804:** Implementasi Halaman Keranjang Belanja.
- [ ] **TASK-805:** Implementasi Halaman *Checkout* (Terintegrasi dengan simulasi Xendit).
- [ ] **TASK-806:** Implementasi Halaman Riwayat Pesanan (Memanggil API Stitching BFF).
