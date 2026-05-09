import { useState, useRef } from "react";
import { motion } from "framer-motion";
import { UploadCloud } from "lucide-react";

const ACCEPTED = ".txt,.json,.pdf,.docx";

export default function FileDropZone({ onFileSelected }) {
  const [dragging, setDragging] = useState(false);
  const inputRef = useRef(null);

  function handleDrop(e) {
    e.preventDefault();
    setDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) onFileSelected(file);
  }

  function handleChange(e) {
    const file = e.target.files[0];
    if (file) onFileSelected(file);
  }

  return (
    <motion.div
      onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
      onDragLeave={() => setDragging(false)}
      onDrop={handleDrop}
      onClick={() => inputRef.current?.click()}
      animate={dragging ? { scale: 1.02 } : { scale: 1 }}
      transition={{ type: "spring", stiffness: 300, damping: 20 }}
      className="flex cursor-pointer flex-col items-center gap-5 rounded-2xl p-16 transition-all duration-300"
      style={{
        border: dragging
          ? "2px dashed rgba(99,102,241,0.5)"
          : "2px dashed rgba(255,255,255,0.1)",
        background: dragging
          ? "rgba(99,102,241,0.06)"
          : "rgba(255,255,255,0.02)",
        boxShadow: dragging
          ? "0 0 40px rgba(99,102,241,0.1), inset 0 0 40px rgba(99,102,241,0.03)"
          : "none",
      }}
    >
      <div
        className="rounded-full p-4 transition-all duration-300"
        style={{
          background: dragging
            ? "rgba(99,102,241,0.15)"
            : "rgba(255,255,255,0.05)",
          boxShadow: dragging
            ? "0 0 20px rgba(99,102,241,0.2)"
            : "none",
        }}
      >
        <UploadCloud
          size={30}
          strokeWidth={1.5}
          className="transition-colors duration-300"
          style={{ color: dragging ? "#818cf8" : "#71717a" }}
        />
      </div>
      <div className="text-center">
        <p className="text-sm font-medium text-zinc-300">
          Drag & drop a file here, or{" "}
          <span className="text-indigo-400">click to browse</span>
        </p>
        <p className="mt-2 text-xs text-zinc-600">
          TXT, JSON, PDF, DOCX — up to 200MB
        </p>
      </div>
      <input
        ref={inputRef}
        type="file"
        accept={ACCEPTED}
        onChange={handleChange}
        className="hidden"
      />
    </motion.div>
  );
}
