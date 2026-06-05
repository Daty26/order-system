import { useEffect, useMemo, useState } from "react";
import {
  Bell,
  Boxes,
  CreditCard,
  LogOut,
  PackagePlus,
  RefreshCw,
  ShoppingCart,
  UserRound
} from "lucide-react";
import { clearSession, readSession, request, saveSession } from "./api.js";

const tabs = [
  { id: "products", label: "Products", icon: Boxes },
  { id: "orders", label: "Orders", icon: ShoppingCart },
  { id: "payments", label: "Payments", icon: CreditCard },
  { id: "notifications", label: "Notifications", icon: Bell }
];

function currency(value) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD"
  }).format(Number(value || 0));
}

function normalizeList(value) {
  return Array.isArray(value) ? value : [];
}

function App() {
  const [session, setSession] = useState(() => readSession());
  const [mode, setMode] = useState("login");
  const [activeTab, setActiveTab] = useState("products");
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState([]);
  const [orders, setOrders] = useState([]);
  const [payments, setPayments] = useState([]);
  const [notifications, setNotifications] = useState([]);

  const isAdmin = session?.user?.role === "ADMIN";

  const totals = useMemo(() => {
    const stock = products.reduce((sum, product) => sum + Number(product.quantity || 0), 0);
    const revenue = payments.reduce((sum, payment) => sum + Number(payment.amount || 0), 0);

    return {
      products: products.length,
      orders: orders.length,
      payments: payments.length,
      notifications: notifications.length,
      stock,
      revenue
    };
  }, [products, orders, payments, notifications]);

  async function loadDashboard() {
    if (!session) {
      return;
    }

    setLoading(true);
    setStatus("");

    try {
      const [nextProducts, nextOrders, nextPayments, nextNotifications] = await Promise.all([
        request("inventory", "/products"),
        request("order", "/orders"),
        request("payment", "/payments"),
        request("notification", "/notifications")
      ]);

      setProducts(normalizeList(nextProducts));
      setOrders(normalizeList(nextOrders));
      setPayments(normalizeList(nextPayments));
      setNotifications(normalizeList(nextNotifications));
    } catch (error) {
      setStatus(error.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadDashboard();
  }, [session]);

  async function handleAuth(event) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    setLoading(true);
    setStatus("");

    try {
      if (mode === "register") {
        await request("user", "/user/register", {
          method: "POST",
          body: JSON.stringify({
            username: data.get("username"),
            email: data.get("email"),
            password: data.get("password"),
            role: data.get("role")
          })
        });
      }

      const loginResponse = await request("user", "/user/login", {
        method: "POST",
        body: JSON.stringify({
          identifier: data.get("identifier") || data.get("email"),
          password: data.get("password")
        })
      });

      saveSession(loginResponse);
      setSession(loginResponse);
      event.currentTarget.reset();
    } catch (error) {
      setStatus(error.message);
    } finally {
      setLoading(false);
    }
  }

  function logout() {
    clearSession();
    setSession(null);
    setOrders([]);
    setPayments([]);
    setNotifications([]);
  }

  async function createOrder(event) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    setLoading(true);
    setStatus("");

    try {
      await request("order", "/orders", {
        method: "POST",
        body: JSON.stringify({
          items: [
            {
              product_id: Number(data.get("product_id")),
              quantity: Number(data.get("quantity"))
            }
          ]
        })
      });
      event.currentTarget.reset();
      await loadDashboard();
    } catch (error) {
      setStatus(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function createProduct(event) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    setLoading(true);
    setStatus("");

    try {
      await request("inventory", "/products", {
        method: "POST",
        body: JSON.stringify({
          name: data.get("name"),
          quantity: Number(data.get("quantity")),
          price: Number(data.get("price"))
        })
      });
      event.currentTarget.reset();
      await loadDashboard();
    } catch (error) {
      setStatus(error.message);
    } finally {
      setLoading(false);
    }
  }

  if (!session) {
    return (
      <main className="auth-shell">
        <section className="auth-panel">
          <div>
            <p className="eyebrow">Order System</p>
            <h1>Commerce operations</h1>
            <p className="lede">Sign in to manage products, orders, payments, and notifications across the service stack.</p>
          </div>

          <div className="segmented" aria-label="Authentication mode">
            <button className={mode === "login" ? "active" : ""} onClick={() => setMode("login")}>
              Sign in
            </button>
            <button className={mode === "register" ? "active" : ""} onClick={() => setMode("register")}>
              Register
            </button>
          </div>

          <form className="form" onSubmit={handleAuth}>
            {mode === "register" ? (
              <>
                <label>
                  Username
                  <input name="username" placeholder="alex" required />
                </label>
                <label>
                  Email
                  <input name="email" type="email" placeholder="alex@example.com" required />
                </label>
                <label>
                  Role
                  <select name="role" defaultValue="USER">
                    <option value="USER">User</option>
                    <option value="ADMIN">Admin</option>
                  </select>
                </label>
              </>
            ) : (
              <label>
                Email or username
                <input name="identifier" placeholder="alex@example.com" required />
              </label>
            )}

            <label>
              Password
              <input name="password" type="password" placeholder="••••••••" required />
            </label>

            <button className="primary" disabled={loading}>
              <UserRound size={18} />
              {mode === "register" ? "Create account" : "Sign in"}
            </button>
          </form>

          {status ? <p className="status error">{status}</p> : null}
        </section>
      </main>
    );
  }

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div>
          <p className="eyebrow">Order System</p>
          <h1>Operations</h1>
        </div>

        <nav className="nav">
          {tabs.map((tab) => {
            const Icon = tab.icon;
            return (
              <button key={tab.id} className={activeTab === tab.id ? "active" : ""} onClick={() => setActiveTab(tab.id)}>
                <Icon size={18} />
                {tab.label}
              </button>
            );
          })}
        </nav>

        <div className="account">
          <span>{session.user.username || session.user.email}</span>
          <small>{session.user.role}</small>
          <button className="icon-button" title="Sign out" onClick={logout}>
            <LogOut size={18} />
          </button>
        </div>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">{tabs.find((tab) => tab.id === activeTab)?.label}</p>
            <h2>{activeTab === "products" ? "Inventory overview" : `Recent ${activeTab}`}</h2>
          </div>
          <button className="secondary" onClick={loadDashboard} disabled={loading}>
            <RefreshCw size={18} />
            Refresh
          </button>
        </header>

        <section className="metrics">
          <Metric label="Products" value={totals.products} />
          <Metric label="Orders" value={totals.orders} />
          <Metric label="Stock units" value={totals.stock} />
          <Metric label="Payment volume" value={currency(totals.revenue)} />
        </section>

        {status ? <p className="status error">{status}</p> : null}

        {activeTab === "products" ? (
          <ProductsView products={products} isAdmin={isAdmin} onCreateProduct={createProduct} loading={loading} />
        ) : null}

        {activeTab === "orders" ? (
          <OrdersView orders={orders} products={products} onCreateOrder={createOrder} loading={loading} />
        ) : null}

        {activeTab === "payments" ? <PaymentsView payments={payments} /> : null}

        {activeTab === "notifications" ? <NotificationsView notifications={notifications} /> : null}
      </section>
    </main>
  );
}

function Metric({ label, value }) {
  return (
    <article className="metric">
      <span>{label}</span>
      <strong>{value}</strong>
    </article>
  );
}

function ProductsView({ products, isAdmin, onCreateProduct, loading }) {
  return (
    <div className="content-grid">
      {isAdmin ? (
        <form className="tool-panel" onSubmit={onCreateProduct}>
          <h3>Add product</h3>
          <label>
            Name
            <input name="name" placeholder="Keyboard" required />
          </label>
          <label>
            Quantity
            <input name="quantity" type="number" min="1" defaultValue="10" required />
          </label>
          <label>
            Price
            <input name="price" type="number" min="0" step="0.01" defaultValue="99.00" required />
          </label>
          <button className="primary" disabled={loading}>
            <PackagePlus size={18} />
            Add
          </button>
        </form>
      ) : null}

      <DataTable
        columns={["ID", "Name", "Stock", "Price"]}
        rows={products.map((product) => [
          product.id,
          product.name,
          product.quantity,
          currency(product.price)
        ])}
        empty="No products yet"
      />
    </div>
  );
}

function OrdersView({ orders, products, onCreateOrder, loading }) {
  return (
    <div className="content-grid">
      <form className="tool-panel" onSubmit={onCreateOrder}>
        <h3>Create order</h3>
        <label>
          Product
          <select name="product_id" required>
            {products.map((product) => (
              <option key={product.id} value={product.id}>
                {product.name} #{product.id}
              </option>
            ))}
          </select>
        </label>
        <label>
          Quantity
          <input name="quantity" type="number" min="1" defaultValue="1" required />
        </label>
        <button className="primary" disabled={loading || products.length === 0}>
          <ShoppingCart size={18} />
          Create
        </button>
      </form>

      <DataTable
        columns={["ID", "Status", "Total", "Items"]}
        rows={orders.map((order) => [
          order.id,
          order.status,
          currency(order.totalAmount),
          normalizeList(order.items).length
        ])}
        empty="No orders yet"
      />
    </div>
  );
}

function PaymentsView({ payments }) {
  return (
    <DataTable
      columns={["ID", "Order", "Status", "Amount"]}
      rows={payments.map((payment) => [
        payment.payment_id,
        payment.order_id,
        payment.status,
        currency(payment.amount)
      ])}
      empty="No payments yet"
    />
  );
}

function NotificationsView({ notifications }) {
  return (
    <DataTable
      columns={["ID", "Order", "Status", "Message"]}
      rows={notifications.map((notification) => [
        notification.id,
        notification.orderID,
        notification.status,
        notification.message
      ])}
      empty="No notifications yet"
    />
  );
}

function DataTable({ columns, rows, empty }) {
  return (
    <div className="table-shell">
      <table>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column}>{column}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.length > 0 ? (
            rows.map((row, rowIndex) => (
              <tr key={rowIndex}>
                {row.map((cell, cellIndex) => (
                  <td key={`${rowIndex}-${cellIndex}`}>{cell}</td>
                ))}
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length} className="empty">
                {empty}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default App;
