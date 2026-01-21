#!/usr/bin/env python3

import argparse
import os
import re

import matplotlib.pyplot as plt
import pandas as pd


CPU_SUFFIX_RE = re.compile(r"-(\d+)$")
N_RE = re.compile(r"(?:^|/)n=(\d+)$")


def normalize_bench_name(bench: str) -> str:
    # Go adds "-<cpu>" at the end of the benchmark name in output.
    return CPU_SUFFIX_RE.sub("", bench)


def split_series_and_n(bench: str):
    """
    bench examples:
      BenchmarkPushPop_MList/n=1024
      BenchmarkGetRandom_Slice/n=16384
    """
    m = N_RE.search(bench)
    if not m:
        return None, None
    n = int(m.group(1))
    series = bench[: m.start()].rstrip("/")
    return series, n


def pretty_label(series: str) -> str:
    # "BenchmarkPushPop_MList" -> "PushPop MList"
    if series.startswith("Benchmark"):
        series = series[len("Benchmark") :]
    return series.replace("_", " ")


def plot_metric(df: pd.DataFrame, metric: str, out_path: str):
    # Filter missing values.
    df = df[df[metric].notna()].copy()
    if df.empty:
        return

    plt.figure(figsize=(9, 5.5))
    for series, g in df.groupby("series"):
        g = g.sort_values("n")
        plt.plot(g["n"], g[metric], marker="o", linewidth=1.5, label=pretty_label(series))

    plt.xscale("log", base=2)
    plt.yscale("log")
    plt.xlabel("n (elements, log2)")
    ylabel = {
        "ns_per_op": "ns/op (log)",
        "allocs_per_op": "allocs/op (log)",
        "bytes_per_op": "bytes/op (log)",
    }.get(metric, metric)
    plt.ylabel(ylabel)
    plt.title(metric.replace("_", " "))
    plt.grid(True, which="both", linestyle="--", alpha=0.3)
    plt.legend(loc="best", fontsize=9)
    plt.tight_layout()
    plt.savefig(out_path, dpi=180)
    plt.close()


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--in", dest="inp", required=True, help="input CSV from cmd/benchrun")
    ap.add_argument("--out", dest="out", required=True, help="output directory for PNGs")
    args = ap.parse_args()

    os.makedirs(args.out, exist_ok=True)
    df = pd.read_csv(args.inp)

    # Normalize and parse series/n from the Go benchmark names.
    df["bench_norm"] = df["bench"].astype(str).map(normalize_bench_name)
    parsed = df["bench_norm"].map(split_series_and_n)
    df["series"] = [p[0] for p in parsed]
    df["n"] = [p[1] for p in parsed]

    # Drop rows without an n=... suffix; those are hard to plot vs n.
    df = df[df["n"].notna()].copy()

    # Ensure numeric for plotting.
    for col in ("ns_per_op", "bytes_per_op", "allocs_per_op"):
        if col in df.columns:
            df[col] = pd.to_numeric(df[col], errors="coerce")

    plot_metric(df, "ns_per_op", os.path.join(args.out, "ns_per_op.png"))
    plot_metric(df, "allocs_per_op", os.path.join(args.out, "allocs_per_op.png"))
    plot_metric(df, "bytes_per_op", os.path.join(args.out, "bytes_per_op.png"))


if __name__ == "__main__":
    main()
