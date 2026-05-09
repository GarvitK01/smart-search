import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { File, Loader2, CheckCircle2, XCircle } from "lucide-react";
import { uploadFile } from "../api";
import FileDropZone from "../components/FileDropZone";

function formatSize(bytes) {
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

export default function UploadPage() {
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [status, setStatus] = useState(null);

  async function handleUpload() {
    if (!file) return;
    setUploading(true);
    setStatus(null);
    try {
      await uploadFile(file);
      setStatus("success");
      setFile(null);
    } catch {
      setStatus("error");
    } finally {
      setUploading(false);
    }
  }

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
          Upload a document
        </h1>
        <p className="text-zinc-500">
          Your file will be parsed, chunked, and indexed for search
        </p>
      </motion.div>

      <FileDropZone onFileSelected={(f) => { setFile(f); setStatus(null); }} />

      <AnimatePresence mode="wait">
        {file && (
          <motion.div
            key="file-preview"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="mt-5 overflow-hidden"
          >
            <div
              className="flex items-center gap-4 rounded-2xl p-5"
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                backdropFilter: "blur(20px)",
              }}
            >
              <div
                className="flex h-11 w-11 items-center justify-center rounded-xl"
                style={{ background: "rgba(99,102,241,0.12)" }}
              >
                <File size={20} className="text-indigo-400" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-zinc-200">{file.name}</p>
                <p className="text-xs text-zinc-500">{formatSize(file.size)}</p>
              </div>
              <button
                onClick={handleUpload}
                disabled={uploading}
                className="flex items-center gap-2 rounded-xl px-6 py-2.5 text-sm font-semibold text-white transition-all disabled:opacity-50"
                style={{
                  background: "linear-gradient(135deg, #6366f1, #4f46e5)",
                  boxShadow: "0 0 20px rgba(99,102,241,0.3)",
                }}
              >
                {uploading ? (
                  <>
                    <Loader2 size={14} className="animate-spin" />
                    Uploading...
                  </>
                ) : (
                  "Upload"
                )}
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <AnimatePresence mode="wait">
        {status === "success" && (
          <motion.div
            key="success"
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0 }}
            className="mt-5 flex items-center gap-3 rounded-2xl p-5 text-sm text-emerald-400"
            style={{
              background: "rgba(16,185,129,0.08)",
              border: "1px solid rgba(16,185,129,0.2)",
            }}
          >
            <CheckCircle2 size={18} />
            File uploaded successfully. Ingestion is running in the background.
          </motion.div>
        )}
        {status === "error" && (
          <motion.div
            key="error"
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0 }}
            className="mt-5 flex items-center gap-3 rounded-2xl p-5 text-sm text-red-400"
            style={{
              background: "rgba(239,68,68,0.08)",
              border: "1px solid rgba(239,68,68,0.2)",
            }}
          >
            <XCircle size={18} />
            Upload failed. Is the server running?
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
