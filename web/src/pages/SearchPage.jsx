import { useState } from "react";
import { motion } from "framer-motion";
import { Search, Loader2, SearchX } from "lucide-react";
import { searchDocuments } from "../api";
import SearchResult from "../components/SearchResult";

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  async function handleSearch(e) {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    setError(null);
    try {
      const data = await searchDocuments(query);
      setResults(data);
    } catch {
      setError("Search failed. Is the server running?");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex flex-col items-center">
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
          Search your documents
        </h1>
        <p className="text-zinc-500">
          Keyword + semantic search across all your uploaded files
        </p>
      </motion.div>

      <form onSubmit={handleSearch} className="mb-10 w-full">
        <div className="relative">
          <Search
            size={18}
            className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500"
          />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search your documents..."
            className="w-full rounded-xl py-4 pl-12 pr-5 text-sm text-zinc-200 placeholder-zinc-600 outline-none transition-all duration-300"
            style={{
              background: "rgba(255,255,255,0.05)",
              border: "1px solid rgba(255,255,255,0.1)",
            }}
            onFocus={(e) => {
              e.target.style.borderColor = "rgba(99,102,241,0.5)";
              e.target.style.boxShadow = "0 0 0 3px rgba(99,102,241,0.15), 0 0 24px rgba(99,102,241,0.1)";
            }}
            onBlur={(e) => {
              e.target.style.borderColor = "rgba(255,255,255,0.1)";
              e.target.style.boxShadow = "none";
            }}
          />
        </div>
      </form>

      {loading && (
        <div className="flex items-center gap-2 text-sm text-indigo-400">
          <Loader2 size={16} className="animate-spin" />
          Searching...
        </div>
      )}

      {error && <p className="text-sm text-red-400">{error}</p>}

      {results && !loading && results.length === 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex flex-col items-center gap-3 py-16 text-zinc-500"
        >
          <SearchX size={44} strokeWidth={1.5} />
          <p className="text-sm">No results found</p>
        </motion.div>
      )}

      {results && results.length > 0 && !loading && (
        <div className="flex w-full flex-col gap-4">
          {results.map((r, i) => (
            <SearchResult
              key={`${r.document_id}-${r.chunk_index}`}
              result={r}
              index={i}
            />
          ))}
        </div>
      )}
    </div>
  );
}
