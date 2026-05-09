export default function StatusBadge({ status }) {
  const isReady = status === "READY";

  return (
    <span
      className="inline-block rounded-full px-3 py-1 text-xs font-semibold"
      style={{
        background: isReady ? "rgba(16,185,129,0.12)" : "rgba(245,158,11,0.12)",
        color: isReady ? "#34d399" : "#fbbf24",
        boxShadow: isReady
          ? "0 0 8px rgba(16,185,129,0.15)"
          : "0 0 8px rgba(245,158,11,0.15)",
      }}
    >
      {status}
    </span>
  );
}
