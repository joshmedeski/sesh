#!/usr/bin/env python3
"""Reduce a multi-count `go test -bench` log to one median line per benchmark.

The CI benchmark job runs with `-count=6` because benchstat needs several
samples to say anything about significance. A history chart does not: six
points per commit per benchmark makes the chart unreadable and the "previous
value" the alert compares against arbitrary. This collapses each benchmark to
its median sample, and passes the goos/goarch/pkg/cpu header lines through
untouched -- github-action-benchmark reads `pkg:` to qualify benchmark names,
so dropping those lines would merge lister's and picker's identically named
benchmarks into one series.
"""

import re
import statistics
import sys
from collections import OrderedDict

# Trailing `-16` is GOMAXPROCS, part of the name as far as the chart cares.
BENCH = re.compile(r"^(Benchmark\S*)\s+(\d+)\s+(.*)$")
METRIC = re.compile(r"([\d.]+)\s+(\S+)")


def fmt(value):
    """Format without scientific notation: `1.23457e+06 ns/op` parses as 1.23."""
    text = f"{value:.4f}".rstrip("0").rstrip(".")
    return text or "0"


def flush(samples, out):
    """Emit the median sample of each benchmark, in first-seen order."""
    for name, runs in samples.items():
        iters = statistics.median_low([r[0] for r in runs])
        # Median each metric independently. They come from the same runs and
        # move together, so this is the same row in practice, and it keeps a
        # benchmark that reports a metric only sometimes from breaking.
        units = OrderedDict()
        for _, metrics in runs:
            for value, unit in metrics:
                units.setdefault(unit, []).append(value)
        cells = "\t".join(
            f"{fmt(statistics.median(values))} {unit}"
            for unit, values in units.items()
        )
        out.write(f"{name}\t{iters}\t{cells}\n")
    samples.clear()


def main() -> int:
    out = sys.stdout
    samples: "OrderedDict[str, list]" = OrderedDict()
    for line in sys.stdin:
        line = line.rstrip("\n")
        match = BENCH.match(line)
        if match:
            name, iters, rest = match.groups()
            metrics = [(float(v), u) for v, u in METRIC.findall(rest)]
            if metrics:
                samples.setdefault(name, []).append((int(iters), metrics))
                continue
        # A `pkg:` line starts a new package section: flush what came before it
        # so each benchmark's median stays under its own package header.
        if line.startswith("pkg:"):
            flush(samples, out)
        out.write(line + "\n")
    flush(samples, out)
    return 0


if __name__ == "__main__":
    sys.exit(main())
