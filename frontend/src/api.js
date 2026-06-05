const API_BASES = {
  user: import.meta.env.VITE_USER_API_BASE || "/api/users",
  order: import.meta.env.VITE_ORDER_API_BASE || "/api/orders",
  payment: import.meta.env.VITE_PAYMENT_API_BASE || "/api/payments",
  notification: import.meta.env.VITE_NOTIFICATION_API_BASE || "/api/notifications",
  inventory: import.meta.env.VITE_INVENTORY_API_BASE || "/api/inventory"
};

function getErrorMessage(payload, fallback) {
  if (payload && typeof payload.error === "string" && payload.error.length > 0) {
    return payload.error;
  }
  return fallback;
}

export async function request(service, path, options = {}) {
  const token = localStorage.getItem("order-system-token");
  const response = await fetch(`${API_BASES[service]}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers || {})
    }
  });

  const text = await response.text();
  const payload = text ? JSON.parse(text) : null;

  if (!response.ok) {
    throw new Error(getErrorMessage(payload, response.statusText));
  }

  return payload?.data ?? payload;
}

export function saveSession(session) {
  localStorage.setItem("order-system-token", session.token);
  localStorage.setItem("order-system-user", JSON.stringify(session.user));
}

export function readSession() {
  const token = localStorage.getItem("order-system-token");
  const user = localStorage.getItem("order-system-user");

  if (!token || !user) {
    return null;
  }

  try {
    return { token, user: JSON.parse(user) };
  } catch {
    return null;
  }
}

export function clearSession() {
  localStorage.removeItem("order-system-token");
  localStorage.removeItem("order-system-user");
}
