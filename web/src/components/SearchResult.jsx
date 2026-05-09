import { useState } from "react";
import { motion } from "framer-motion";
import { FileText } from "lucide-react";

const COLORS = ["#818cf8", "#06b6d4", "#34d399", "#f59e0b", "#f472b6"];

const glass = {
  background: "rgba(255, 255, 255, 0.04)",
  backdropFilter: "blur(20px)",
  WebkitBackdropFilter: "blur(20px)",
  border: "1px solid rgba(255, 255, 255, 0.08)",
};

export default function SearchResult({ result, index }) {
  const color = COLORS[index % COLORS.length];
  const [hovered, setHovered] = useState(false);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.08, duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
      className="group relative overflow-hidden rounded-2xl p-6 transition-all duration-300"
      style={{
        ...glass,
        ...(hovered
          ? {
              background: "rgba(255,255,255,0.07)",
              borderColor: "rgba(255,255,255,0.14)",
              boxShadow: `0 8px 40px ${color}12`,
            }
          : {}),
      }}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <div
        className="absolute left-0 top-0 h-full w-1 transition-opacity duration-300"
        style={{
          background: `linear-gradient(to bottom, ${color}, transparent)`,
          opacity: hovered ? 1 : 0.5,
        }}
      />
      <div className="ml-4">
        <div className="mb-3 flex items-center gap-2.5 text-sm">
          <FileText size={15} style={{ color }} />
          <span className="font-medium text-zinc-200">
            {result.original_name || `Document ${result.document_id}`}
          </span>
          <span className="text-zinc-600">&middot;</span>
          <span
            className="rounded-full px-2.5 py-0.5 text-xs text-zinc-500"
            style={{ background: "rgba(255,255,255,0.05)" }}
          >
            Chunk {result.chunk_index}
          </span>
        </div>
        <p className="text-sm leading-relaxed text-zinc-400">
          {result.content.length > 400
            ? result.content.slice(0, 400) + "..."
            : result.content}
        </p>
      </div>
    </motion.div>
  );
}
