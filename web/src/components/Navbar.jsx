import { NavLink, useNavigate } from "react-router-dom";
import { Home, Search, Upload, FileText, Sparkles, LogIn, UserPlus, LogOut } from "lucide-react";
import { useAuth } from "../context/AuthContext";

const appLinks = [
  { to: "/", label: "Home", icon: Home },
  { to: "/search", label: "Search", icon: Search },
  { to: "/upload", label: "Upload", icon: Upload },
  { to: "/documents", label: "Documents", icon: FileText },
];

const guestLinks = [
  { to: "/login", label: "Sign in", icon: LogIn },
  { to: "/register", label: "Register", icon: UserPlus },
];

export default function Navbar() {
  const { user, loading, logout } = useAuth();
  const navigate = useNavigate();

  const navLinks = user ? appLinks : [appLinks[0], ...guestLinks];

  async function handleLogout() {
    await logout();
    navigate("/");
  }

  return (
    <nav
      className="sticky top-0 z-50"
      style={{
        background: "rgba(5, 5, 16, 0.7)",
        backdropFilter: "blur(20px)",
        WebkitBackdropFilter: "blur(20px)",
        borderBottom: "1px solid rgba(255,255,255,0.06)",
      }}
    >
      <div className="mx-auto flex max-w-4xl items-center justify-between px-8 py-3">
        <NavLink to="/" className="flex items-center gap-2.5">
          <div
            className="flex h-8 w-8 items-center justify-center rounded-lg"
            style={{
              background: "rgba(99,102,241,0.15)",
              boxShadow: "0 0 12px rgba(99,102,241,0.2)",
            }}
          >
            <Sparkles size={16} className="text-indigo-400" />
          </div>
          <span
            className="text-lg font-bold tracking-tight"
            style={{
              background: "linear-gradient(135deg, #818cf8, #06b6d4)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
            }}
          >
            Smart Search
          </span>
        </NavLink>
        <div className="flex items-center gap-1">
          {!loading &&
            navLinks.map(({ to, label, icon: Icon }) => (
              <NavLink
                key={to}
                to={to}
                end
                className={({ isActive }) =>
                  `flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium transition-all duration-200 ${
                    isActive
                      ? "text-white"
                      : "text-zinc-500 hover:text-zinc-300"
                  }`
                }
                style={({ isActive }) =>
                  isActive
                    ? {
                        background: "rgba(255,255,255,0.08)",
                        boxShadow: "0 0 12px rgba(99,102,241,0.12)",
                      }
                    : {}
                }
              >
                <Icon size={15} />
                {label}
              </NavLink>
            ))}
          {!loading && user && (
            <button
              onClick={handleLogout}
              className="ml-2 flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-zinc-500 transition-all duration-200 hover:text-red-400"
            >
              <LogOut size={15} />
              Logout
            </button>
          )}
        </div>
      </div>
    </nav>
  );
}
