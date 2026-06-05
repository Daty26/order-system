# Order System Frontend

React belongs at the repository root in `frontend/`, next to the Go services. This keeps it as a separate client app for the whole system instead of coupling it to one backend service.

## Run Locally

```sh
npm install
npm run dev
```

The Vite dev server runs on `http://localhost:5173` and proxies API calls to the Go services:

| Frontend path | Backend |
| --- | --- |
| `/api/users/*` | `localhost:8085` |
| `/api/orders/*` | `localhost:8080` |
| `/api/payments/*` | `localhost:8081` |
| `/api/notifications/*` | `localhost:8082` |
| `/api/inventory/*` | `localhost:8084` |

Start the backend stack first, then run the frontend.
