import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Loader2, FolderOpen, FileText } from "lucide-react";
import { listDocuments } from "../api";
import StatusBadge from "../components/StatusBadge";

function formatSize(bytes) {
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export default function DocumentsPage() {
  const [docs, setDocs] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    listDocuments()
      .then(setDocs)
      .catch(() => setError("Failed to load documents"));
  }, []);

  return (
    <div>
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-10 text-center"
      >
        <h1
          className="mb-3 text-4xl font-bold tracking-tight"
          style={{
            background: "linear-gradient(135deg, #818cf8, #06b6d4)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
            backgroundClip: "text",
          }}
        >
          Your documents
        </h1>
        <p className="text-zinc-500">
          All files you've uploaded and their processing status
        </p>
      </motion.div>

      {!docs && !error && (
        <div className="flex items-center justify-center gap-2 py-16 text-sm text-indigo-400">
          <Loader2 size={16} className="animate-spin" />
          Loading...
        </div>
      )}

      {error && <p className="text-center text-sm text-red-400">{error}</p>}

      {docs && docs.length === 0 && (
        <div className="flex flex-col items-center gap-3 py-16 text-zinc-500">
          <FolderOpen size={44} strokeWidth={1.5} />
          <p className="text-sm">No documents uploaded yet</p>
        </div>
      )}

      {docs && docs.length > 0 && (
        <div
          className="overflow-hidden rounded-2xl"
          style={{
            background: "rgba(255,255,255,0.04)",
            border: "1px solid rgba(255,255,255,0.08)",
            backdropFilter: "blur(20px)",
          }}
        >
          <table className="w-full text-left text-sm">
            <thead>
              <tr style={{ borderBottom: "1px solid rgba(255,255,255,0.06)" }}>
                <th className="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-zinc-500">
                  File
                </th>
                <th className="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-zinc-500">
                  Size
                </th>
                <th className="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-zinc-500">
                  Status
                </th>
                <th className="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-zinc-500">
                  Uploaded
                </th>
              </tr>
            </thead>
            <tbody>
              {docs.map((doc, i) => (
                <motion.tr
                  key={doc.id}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ delay: i * 0.06 }}
                  className="transition-colors duration-200"
                  style={{ borderBottom: "1px solid rgba(255,255,255,0.04)" }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.background = "rgba(255,255,255,0.03)";
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.background = "transparent";
                  }}
                >
                  <td className="px-5 py-4">
                    <div className="flex items-center gap-3">
                      <div
                        className="flex h-9 w-9 items-center justify-center rounded-lg"
                        style={{ background: "rgba(99,102,241,0.1)" }}
                      >
                        <FileText size={15} className="text-indigo-400" />
                      </div>
                      <span className="font-medium text-zinc-200">
                        {doc.original_name || doc.file_name}
                      </span>
                    </div>
                  </td>
                  <td className="px-5 py-4 text-zinc-500">
                    {formatSize(doc.size)}
                  </td>
                  <td className="px-5 py-4">
                    <StatusBadge status={doc.status} />
                  </td>
                  <td className="px-5 py-4 text-zinc-500">
                    {formatDate(doc.created_at)}
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
