import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import {
  Upload,
  Cpu,
  Search,
  ArrowRight,
  Zap,
  Database,
  Brain,
  Code2,
  Layers,
  Globe,
} from "lucide-react";

const fadeUp = {
  hidden: { opacity: 0, y: 30 },
  visible: (i = 0) => ({
    opacity: 1,
    y: 0,
    transition: { delay: i * 0.1, duration: 0.6, ease: [0.22, 1, 0.36, 1] },
  }),
};

const glass = {
  background: "rgba(255, 255, 255, 0.05)",
  backdropFilter: "blur(24px)",
  WebkitBackdropFilter: "blur(24px)",
  border: "1px solid rgba(255, 255, 255, 0.1)",
};

const glassStrong = {
  background: "rgba(255, 255, 255, 0.08)",
  backdropFilter: "blur(32px)",
  WebkitBackdropFilter: "blur(32px)",
  border: "1px solid rgba(255, 255, 255, 0.12)",
};

const steps = [
  {
    icon: Upload,
    title: "Upload",
    desc: "Drop any TXT, JSON, PDF, or DOCX file. We store it securely and kick off processing automatically.",
    color: "#818cf8",
  },
  {
    icon: Cpu,
    title: "Process",
    desc: "Documents are chunked into passages. Each chunk gets a 384-dim embedding from a transformer model.",
    color: "#06b6d4",
  },
  {
    icon: Search,
    title: "Search",
    desc: "Queries hit PostgreSQL full-text and pgvector cosine similarity in parallel — merged and ranked.",
    color: "#a78bfa",
  },
];

const tech = [
  { label: "Go", icon: Code2 },
  { label: "PostgreSQL", icon: Database },
  { label: "pgvector", icon: Layers },
  { label: "FastAPI", icon: Zap },
  { label: "sentence-transformers", icon: Brain },
  { label: "React", icon: Globe },
];

export default function LandingPage() {
  return (
    <div className="relative z-10 mx-auto max-w-4xl px-8">
      {/* ─── Hero ─── */}
      <section className="flex min-h-[70vh] flex-col items-center justify-center py-16 text-center">
        <motion.div
          initial="hidden"
          animate="visible"
          variants={fadeUp}
          custom={0}
          style={glass}
          className="mb-5 inline-flex items-center gap-2 rounded-full px-4 py-1.5 text-xs font-medium tracking-wide text-indigo-300"
        >
          <Zap size={12} />
          Keyword + Semantic hybrid search
        </motion.div>

        <motion.h1
          initial="hidden"
          animate="visible"
          variants={fadeUp}
          custom={1}
          className="mb-5 max-w-3xl text-5xl font-extrabold leading-[1.1] tracking-tight sm:text-6xl"
          style={{
            background: "linear-gradient(135deg, #e0e7ff 0%, #818cf8 35%, #06b6d4 65%, #a5f3fc 100%)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
            backgroundClip: "text",
          }}
        >
          Your documents, instantly searchable
        </motion.h1>

        <motion.p
          initial="hidden"
          animate="visible"
          variants={fadeUp}
          custom={2}
          className="mb-8 max-w-xl text-base leading-relaxed text-zinc-400"
        >
          Upload files, let AI embed every passage, and find exactly what you
          need with a single query — combining the precision of keywords with
          the intelligence of semantic understanding.
        </motion.p>

        <motion.div
          initial="hidden"
          animate="visible"
          variants={fadeUp}
          custom={3}
          className="flex flex-wrap justify-center gap-4"
        >
          <Link
            to="/search"
            className="group flex items-center gap-2 rounded-xl px-8 py-3 text-sm font-semibold text-white transition-all"
            style={{
              background: "linear-gradient(135deg, #6366f1, #4f46e5)",
              boxShadow: "0 0 30px rgba(99,102,241,0.35), 0 4px 16px rgba(0,0,0,0.3)",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.boxShadow = "0 0 50px rgba(99,102,241,0.5), 0 4px 24px rgba(0,0,0,0.4)";
              e.currentTarget.style.transform = "translateY(-2px)";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.boxShadow = "0 0 30px rgba(99,102,241,0.35), 0 4px 16px rgba(0,0,0,0.3)";
              e.currentTarget.style.transform = "translateY(0)";
            }}
          >
            Start searching
            <ArrowRight size={16} className="transition-transform group-hover:translate-x-1" />
          </Link>
          <Link
            to="/upload"
            style={glass}
            className="flex items-center gap-2 rounded-xl px-8 py-3 text-sm font-semibold text-zinc-300 transition-all hover:text-white"
          >
            Upload a file
          </Link>
        </motion.div>
      </section>

      {/* ─── How it works ─── */}
      <section className="py-12">
        <motion.p
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          variants={fadeUp}
          custom={0}
          className="mb-2 text-center text-sm font-semibold uppercase tracking-[0.2em] text-indigo-400"
        >
          How it works
        </motion.p>
        <motion.h2
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          variants={fadeUp}
          custom={1}
          className="mb-8 text-center text-3xl font-bold text-white"
        >
          Three steps to smarter search
        </motion.h2>

        <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
          {steps.map((step, i) => (
            <motion.div
              key={step.title}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true, margin: "-40px" }}
              variants={fadeUp}
              custom={i}
              style={glass}
              className="group rounded-2xl p-6 transition-all duration-300 hover:scale-[1.02]"
              onMouseEnter={(e) => {
                e.currentTarget.style.background = "rgba(255,255,255,0.08)";
                e.currentTarget.style.borderColor = "rgba(255,255,255,0.15)";
                e.currentTarget.style.boxShadow = `0 8px 40px ${step.color}15`;
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = "rgba(255,255,255,0.05)";
                e.currentTarget.style.borderColor = "rgba(255,255,255,0.1)";
                e.currentTarget.style.boxShadow = "none";
              }}
            >
              <div
                className="mb-4 flex h-12 w-12 items-center justify-center rounded-xl"
                style={{ background: `${step.color}18` }}
              >
                <step.icon size={22} style={{ color: step.color }} />
              </div>
              <h3 className="mb-2 text-lg font-semibold text-white">{step.title}</h3>
              <p className="text-sm leading-relaxed text-zinc-400">{step.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* ─── Tech stack ─── */}
      <section className="py-12">
        <motion.p
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          variants={fadeUp}
          custom={0}
          className="mb-2 text-center text-sm font-semibold uppercase tracking-[0.2em] text-indigo-400"
        >
          Built with
        </motion.p>
        <motion.h2
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          variants={fadeUp}
          custom={1}
          className="mb-8 text-center text-3xl font-bold text-white"
        >
          Modern, production-grade stack
        </motion.h2>

        <div className="flex flex-wrap justify-center gap-3">
          {tech.map((t, i) => (
            <motion.div
              key={t.label}
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true }}
              variants={fadeUp}
              custom={i * 0.3}
              style={glass}
              className="flex items-center gap-2.5 rounded-full px-5 py-2.5 text-sm font-medium text-zinc-300 transition-all duration-300 hover:text-white"
              onMouseEnter={(e) => {
                e.currentTarget.style.background = "rgba(255,255,255,0.1)";
                e.currentTarget.style.borderColor = "rgba(255,255,255,0.2)";
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.background = "rgba(255,255,255,0.05)";
                e.currentTarget.style.borderColor = "rgba(255,255,255,0.1)";
              }}
            >
              <t.icon size={16} className="text-indigo-400" />
              {t.label}
            </motion.div>
          ))}
        </div>
      </section>

      {/* ─── Footer CTA ─── */}
      <section className="py-16">
        <motion.div
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          variants={fadeUp}
          style={glassStrong}
          className="mx-auto max-w-2xl rounded-3xl px-8 py-12 text-center"
        >
          <h2 className="mb-3 text-3xl font-bold text-white">
            Ready to try it?
          </h2>
          <p className="mb-6 text-zinc-400">
            Upload a document and search it in seconds.
          </p>
          <Link
            to="/search"
            className="group inline-flex items-center gap-2 rounded-xl px-8 py-3 text-sm font-semibold text-white transition-all"
            style={{
              background: "linear-gradient(135deg, #6366f1, #4f46e5)",
              boxShadow: "0 0 30px rgba(99,102,241,0.35), 0 4px 16px rgba(0,0,0,0.3)",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.boxShadow = "0 0 50px rgba(99,102,241,0.5), 0 4px 24px rgba(0,0,0,0.4)";
              e.currentTarget.style.transform = "translateY(-2px)";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.boxShadow = "0 0 30px rgba(99,102,241,0.35), 0 4px 16px rgba(0,0,0,0.3)";
              e.currentTarget.style.transform = "translateY(0)";
            }}
          >
            Get started
            <ArrowRight size={16} className="transition-transform group-hover:translate-x-1" />
          </Link>
        </motion.div>
      </section>
    </div>
  );
}
