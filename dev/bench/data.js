window.BENCHMARK_DATA = {
  "lastUpdate": 1787691934693,
  "repoUrl": "https://github.com/joshmedeski/sesh",
  "entries": {
    "sesh": [
      {
        "commit": {
          "author": {
            "email": "josh.medeski@gmail.com",
            "name": "Josh Medeski",
            "username": "joshmedeski"
          },
          "committer": {
            "email": "joshmedeski@users.noreply.github.com",
            "name": "Josh Medeski",
            "username": "joshmedeski"
          },
          "distinct": true,
          "id": "4881b7e915e33e0402696507a3203b2a391712eb",
          "message": "ci: measure benchmarks over a duration, not 100 iterations\n\n`-benchtime` takes a duration, so `-benchtime=100x` did not mean \"100ms\" --\nit pinned every benchmark to 100 iterations. For a nanosecond-scale\nbenchmark that is microseconds of total work, far too little to reach\nsteady state, so the job reported warmup and scheduler jitter rather than\nthe code.\n\nMeasured over 10 runs on an idle machine, that config put\nRankMatches/n=10 at +-125% and BlacklistFilter/n=10 at +-70%, and read\n2.9x high in absolute terms (422ns against a true 148ns). The history\ndashboard's alert threshold is 125%, so the threshold sat below the noise\nfloor of its own measurement and fired on nothing: it flagged all three\ncommits it ever evaluated, including c8114e2, which changed two lines of\nthis YAML file and no Go code at all. Every flagged benchmark had\nbyte-identical allocs/op, and none shared code with the commit it was\nreported against.\n\n`-benchtime=100ms` holds the same benchmarks to +-1%, as tight as a full\nsecond. The suite costs 71s instead of 4s locally, which is what a usable\nnumber costs; the bench job's timeout goes to 30 minutes because a pull\nrequest runs the suite twice.\n\nStop the dashboard from commenting on the commits it flags. It compares\none point against the single previous point with a fixed ratio and no\nmodel of variance, so it flags noise, and it ratchets -- one lucky-fast\ncommit becomes the baseline that makes the next ordinary commit look\nslow. Alerts still land in the job summary; they just stop writing\npermanent regression notices onto main. The comparison worth acting on is\nthe PR job's benchstat table, which has several samples and reports\nsignificance.\n\nAdd a Benchmark calibration workflow so the thresholds in CONTRIBUTING.md\ncan be set from runner data instead of a guess. It runs the suite twice on\nthe same commit and benchstats the halves, so every delta it reports is\nnoise by construction, and a threshold tighter than its widest row cannot\ntell a regression from a re-run.\n\nClaude-Session: https://claude.ai/code/session_01Tnua4HVKwm1vsFH3zzKN4w",
          "timestamp": "2026-08-21T09:04:38-05:00",
          "tree_id": "11eccd22477911d085d5ebeeeb7fbdd012e38cf3",
          "url": "https://github.com/joshmedeski/sesh/commit/4881b7e915e33e0402696507a3203b2a391712eb"
        },
        "date": 1787321184717,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 12848.5,
            "unit": "ns/op\t6633 B/op\t73 allocs/op",
            "extra": "9141 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 12848.5,
            "unit": "ns/op",
            "extra": "9141 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6633,
            "unit": "B/op",
            "extra": "9141 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "9141 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 45221.5,
            "unit": "ns/op\t58301.5 B/op\t447 allocs/op",
            "extra": "2625 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 45221.5,
            "unit": "ns/op",
            "extra": "2625 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58301.5,
            "unit": "B/op",
            "extra": "2625 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "2625 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 305568,
            "unit": "ns/op\t538061 B/op\t3939 allocs/op",
            "extra": "392 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 305568,
            "unit": "ns/op",
            "extra": "392 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538061,
            "unit": "B/op",
            "extra": "392 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "392 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 18285.5,
            "unit": "ns/op\t8482 B/op\t87 allocs/op",
            "extra": "6213 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 18285.5,
            "unit": "ns/op",
            "extra": "6213 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8482,
            "unit": "B/op",
            "extra": "6213 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "6213 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 86752.5,
            "unit": "ns/op\t79305.5 B/op\t489 allocs/op",
            "extra": "1393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 86752.5,
            "unit": "ns/op",
            "extra": "1393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79305.5,
            "unit": "B/op",
            "extra": "1393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "1393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 721032,
            "unit": "ns/op\t746371 B/op\t4032.5 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 721032,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746371,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4032.5,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 27877.5,
            "unit": "ns/op\t18806.5 B/op\t206 allocs/op",
            "extra": "4585 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27877.5,
            "unit": "ns/op",
            "extra": "4585 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18806.5,
            "unit": "B/op",
            "extra": "4585 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 206,
            "unit": "allocs/op",
            "extra": "4585 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 84369,
            "unit": "ns/op\t92701.5 B/op\t664 allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 84369,
            "unit": "ns/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92701.5,
            "unit": "B/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 664,
            "unit": "allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 607617.5,
            "unit": "ns/op\t780989 B/op\t4948 allocs/op",
            "extra": "195 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 607617.5,
            "unit": "ns/op",
            "extra": "195 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 780989,
            "unit": "B/op",
            "extra": "195 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4948,
            "unit": "allocs/op",
            "extra": "195 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 14167.5,
            "unit": "ns/op\t6665 B/op\t74 allocs/op",
            "extra": "8302 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 14167.5,
            "unit": "ns/op",
            "extra": "8302 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6665,
            "unit": "B/op",
            "extra": "8302 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "8302 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 48093,
            "unit": "ns/op\t58333 B/op\t448 allocs/op",
            "extra": "2554 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 48093,
            "unit": "ns/op",
            "extra": "2554 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58333,
            "unit": "B/op",
            "extra": "2554 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "2554 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 337821,
            "unit": "ns/op\t538091 B/op\t3940 allocs/op",
            "extra": "358 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 337821,
            "unit": "ns/op",
            "extra": "358 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538091,
            "unit": "B/op",
            "extra": "358 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "358 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2714.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "43222 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2714.5,
            "unit": "ns/op",
            "extra": "43222 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "43222 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "43222 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 33959.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "3675 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 33959.5,
            "unit": "ns/op",
            "extra": "3675 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "3675 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "3675 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 356342.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "333 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 356342.5,
            "unit": "ns/op",
            "extra": "333 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "333 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "333 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7687,
            "unit": "ns/op\t9666.5 B/op\t118 allocs/op",
            "extra": "15567 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7687,
            "unit": "ns/op",
            "extra": "15567 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9666.5,
            "unit": "B/op",
            "extra": "15567 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "15567 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 21440,
            "unit": "ns/op\t9674 B/op\t118 allocs/op",
            "extra": "5802 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 21440,
            "unit": "ns/op",
            "extra": "5802 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9674,
            "unit": "B/op",
            "extra": "5802 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "5802 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 160065.5,
            "unit": "ns/op\t9700.5 B/op\t118 allocs/op",
            "extra": "753 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 160065.5,
            "unit": "ns/op",
            "extra": "753 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9700.5,
            "unit": "B/op",
            "extra": "753 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "753 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.2415,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "12928915 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.2415,
            "unit": "ns/op",
            "extra": "12928915 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12928915 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12928915 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 280.85,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "434985 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 280.85,
            "unit": "ns/op",
            "extra": "434985 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "434985 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "434985 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 100.08,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "1130934 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 100.08,
            "unit": "ns/op",
            "extra": "1130934 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "1130934 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1130934 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2769.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "42733 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2769.5,
            "unit": "ns/op",
            "extra": "42733 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "42733 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "42733 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3062.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "39031 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3062.5,
            "unit": "ns/op",
            "extra": "39031 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "39031 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "39031 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.2495,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "12962455 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.2495,
            "unit": "ns/op",
            "extra": "12962455 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12962455 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12962455 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2299.5,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "51192 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2299.5,
            "unit": "ns/op",
            "extra": "51192 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "51192 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "51192 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 288.35,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "370026 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 288.35,
            "unit": "ns/op",
            "extra": "370026 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "370026 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "370026 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 34582.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "3583 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 34582.5,
            "unit": "ns/op",
            "extra": "3583 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "3583 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "3583 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 37639,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "3376 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 37639,
            "unit": "ns/op",
            "extra": "3376 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "3376 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "3376 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.2785,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "12940071 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.2785,
            "unit": "ns/op",
            "extra": "12940071 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12940071 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12940071 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 24419,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 24419,
            "unit": "ns/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2651.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "44324 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2651.5,
            "unit": "ns/op",
            "extra": "44324 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "44324 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "44324 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 360550.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "330 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 360550.5,
            "unit": "ns/op",
            "extra": "330 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "330 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "330 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 368348.5,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 368348.5,
            "unit": "ns/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 350.6,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "338208 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 350.6,
            "unit": "ns/op",
            "extra": "338208 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "338208 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "338208 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2117,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "56302 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2117,
            "unit": "ns/op",
            "extra": "56302 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "56302 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "56302 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19619,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "6124 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19619,
            "unit": "ns/op",
            "extra": "6124 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063,
            "unit": "B/op",
            "extra": "6124 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "6124 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 63841,
            "unit": "ns/op\t18523.5 B/op\t178 allocs/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 63841,
            "unit": "ns/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18523.5,
            "unit": "B/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 229121.5,
            "unit": "ns/op\t74840.5 B/op\t633 allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 229121.5,
            "unit": "ns/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 74840.5,
            "unit": "B/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 633,
            "unit": "allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 510049,
            "unit": "ns/op\t478051.5 B/op\t1417 allocs/op",
            "extra": "232 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 510049,
            "unit": "ns/op",
            "extra": "232 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478051.5,
            "unit": "B/op",
            "extra": "232 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1417,
            "unit": "allocs/op",
            "extra": "232 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 494.35,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "238177 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 494.35,
            "unit": "ns/op",
            "extra": "238177 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "238177 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "238177 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3311.5,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "36102 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3311.5,
            "unit": "ns/op",
            "extra": "36102 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "36102 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "36102 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3538.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "33873 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3538.5,
            "unit": "ns/op",
            "extra": "33873 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "33873 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "33873 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1920.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "61368 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1920.5,
            "unit": "ns/op",
            "extra": "61368 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "61368 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "61368 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 502.8,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "234824 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 502.8,
            "unit": "ns/op",
            "extra": "234824 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "234824 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "234824 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3355,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "35138 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3355,
            "unit": "ns/op",
            "extra": "35138 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "35138 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "35138 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3602,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "33597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3602,
            "unit": "ns/op",
            "extra": "33597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "33597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "33597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1915,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "61431 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1915,
            "unit": "ns/op",
            "extra": "61431 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "61431 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "61431 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4502.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "26484 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4502.5,
            "unit": "ns/op",
            "extra": "26484 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "26484 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "26484 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 29283,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "4062 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 29283,
            "unit": "ns/op",
            "extra": "4062 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "4062 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "4062 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 31240.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 31240.5,
            "unit": "ns/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3658 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 17529,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 17529,
            "unit": "ns/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4823,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "24598 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4823,
            "unit": "ns/op",
            "extra": "24598 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "24598 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "24598 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 29040.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "3920 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 29040.5,
            "unit": "ns/op",
            "extra": "3920 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "3920 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3920 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 30994,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 30994,
            "unit": "ns/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 17573.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "6777 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 17573.5,
            "unit": "ns/op",
            "extra": "6777 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "6777 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "6777 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 133261,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 133261,
            "unit": "ns/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 350418,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 350418,
            "unit": "ns/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 348036.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "339 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 348036.5,
            "unit": "ns/op",
            "extra": "339 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "339 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "339 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 179520.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 179520.5,
            "unit": "ns/op",
            "extra": "664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 111616.5,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 111616.5,
            "unit": "ns/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 348924.5,
            "unit": "ns/op\t441497.5 B/op\t922 allocs/op",
            "extra": "328 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 348924.5,
            "unit": "ns/op",
            "extra": "328 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497.5,
            "unit": "B/op",
            "extra": "328 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "328 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 341387,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 341387,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 179775,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "666 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 179775,
            "unit": "ns/op",
            "extra": "666 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "666 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "666 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 327.4,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "376603 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 327.4,
            "unit": "ns/op",
            "extra": "376603 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "376603 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "376603 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 325.1,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "353764 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 325.1,
            "unit": "ns/op",
            "extra": "353764 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "353764 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "353764 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2439.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "48684 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2439.5,
            "unit": "ns/op",
            "extra": "48684 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "48684 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "48684 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2426.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "49472 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2426.5,
            "unit": "ns/op",
            "extra": "49472 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "49472 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "49472 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 32629.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "3698 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 32629.5,
            "unit": "ns/op",
            "extra": "3698 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "3698 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3698 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 35999,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "3070 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 35999,
            "unit": "ns/op",
            "extra": "3070 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "3070 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3070 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 679.75,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "173050 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 679.75,
            "unit": "ns/op",
            "extra": "173050 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "173050 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "173050 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1194,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "100282 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1194,
            "unit": "ns/op",
            "extra": "100282 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "100282 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "100282 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1081.5,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "111304 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1081.5,
            "unit": "ns/op",
            "extra": "111304 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "111304 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "111304 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 5790.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5790.5,
            "unit": "ns/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 11067.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11067.5,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 30336,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 201,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 10466,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 10466,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 73403.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "1616 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 73403.5,
            "unit": "ns/op",
            "extra": "1616 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "1616 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1616 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 122652,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "950 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 122652,
            "unit": "ns/op",
            "extra": "950 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "950 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "950 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 118386,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "970 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 118386,
            "unit": "ns/op",
            "extra": "970 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "970 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "970 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6919.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "17332 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6919.5,
            "unit": "ns/op",
            "extra": "17332 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17332 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17332 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6972,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "17205 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6972,
            "unit": "ns/op",
            "extra": "17205 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17205 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17205 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 71505.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 71505.5,
            "unit": "ns/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1678 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 71192,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1664 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 71192,
            "unit": "ns/op",
            "extra": "1664 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1664 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1664 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 723948.5,
            "unit": "ns/op\t7 B/op\t0 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 723948.5,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 7,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 725594,
            "unit": "ns/op\t7 B/op\t0 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 725594,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 7,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 614.25,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "178844 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 614.25,
            "unit": "ns/op",
            "extra": "178844 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "178844 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "178844 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 622.5,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "166266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 622.5,
            "unit": "ns/op",
            "extra": "166266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "166266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "166266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 642.65,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "184597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 642.65,
            "unit": "ns/op",
            "extra": "184597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "184597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "184597 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3790,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "31270 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3790,
            "unit": "ns/op",
            "extra": "31270 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "31270 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "31270 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3804.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "31700 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3804.5,
            "unit": "ns/op",
            "extra": "31700 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "31700 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "31700 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3993,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "29953 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3993,
            "unit": "ns/op",
            "extra": "29953 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "29953 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "29953 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 41195.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2869 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41195.5,
            "unit": "ns/op",
            "extra": "2869 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2869 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2869 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 41527,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2835 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41527,
            "unit": "ns/op",
            "extra": "2835 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2835 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2835 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 43853,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 43853,
            "unit": "ns/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 27785,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "4304 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 27785,
            "unit": "ns/op",
            "extra": "4304 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "4304 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "4304 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 69967.5,
            "unit": "ns/op\t6041 B/op\t152 allocs/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 69967.5,
            "unit": "ns/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6041,
            "unit": "B/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 69495.5,
            "unit": "ns/op\t4520.5 B/op\t151 allocs/op",
            "extra": "1732 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 69495.5,
            "unit": "ns/op",
            "extra": "1732 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4520.5,
            "unit": "B/op",
            "extra": "1732 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "1732 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "josh.medeski@gmail.com",
            "name": "Josh Medeski",
            "username": "joshmedeski"
          },
          "committer": {
            "email": "joshmedeski@users.noreply.github.com",
            "name": "Josh Medeski",
            "username": "joshmedeski"
          },
          "distinct": true,
          "id": "ee7a804e88b49566a751fef4cf71120dacf80407",
          "message": "chore: move to Go 1.27 and upgrade mockery to v3.7.4\n\nmockery v3.7.4 bumps golang.org/x/tools to v0.49.0 for Go 1.27 support.\nRegenerating mocks produces no diff.\n\nMove everything to 1.27 in one step: the go directive in go.mod, the\ngo-version in all setup-go steps (test, release, bench-calibrate), and\nthe version noted in CLAUDE.md.",
          "timestamp": "2026-08-25T16:03:16-05:00",
          "tree_id": "e9efd545c931d24cb64e511efd5e9396a80bb753",
          "url": "https://github.com/joshmedeski/sesh/commit/ee7a804e88b49566a751fef4cf71120dacf80407"
        },
        "date": 1787691933360,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 18122,
            "unit": "ns/op\t6633 B/op\t73 allocs/op",
            "extra": "6567 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 18122,
            "unit": "ns/op",
            "extra": "6567 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6633,
            "unit": "B/op",
            "extra": "6567 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "6567 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 68338.5,
            "unit": "ns/op\t58301.5 B/op\t447 allocs/op",
            "extra": "1738 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 68338.5,
            "unit": "ns/op",
            "extra": "1738 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58301.5,
            "unit": "B/op",
            "extra": "1738 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "1738 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 443194,
            "unit": "ns/op\t538073 B/op\t3939 allocs/op",
            "extra": "267 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 443194,
            "unit": "ns/op",
            "extra": "267 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538073,
            "unit": "B/op",
            "extra": "267 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "267 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 28075,
            "unit": "ns/op\t8481 B/op\t87 allocs/op",
            "extra": "4237 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 28075,
            "unit": "ns/op",
            "extra": "4237 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8481,
            "unit": "B/op",
            "extra": "4237 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "4237 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 116446.5,
            "unit": "ns/op\t79304 B/op\t489 allocs/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 116446.5,
            "unit": "ns/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79304,
            "unit": "B/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "1045 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 941387.5,
            "unit": "ns/op\t746369 B/op\t4032.5 allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 941387.5,
            "unit": "ns/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746369,
            "unit": "B/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4032.5,
            "unit": "allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 38688.5,
            "unit": "ns/op\t18771 B/op\t204 allocs/op",
            "extra": "2842 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 38688.5,
            "unit": "ns/op",
            "extra": "2842 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18771,
            "unit": "B/op",
            "extra": "2842 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 204,
            "unit": "allocs/op",
            "extra": "2842 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 124265.5,
            "unit": "ns/op\t92709.5 B/op\t662 allocs/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 124265.5,
            "unit": "ns/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92709.5,
            "unit": "B/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 662,
            "unit": "allocs/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 827656.5,
            "unit": "ns/op\t781073.5 B/op\t4946 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 827656.5,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 781073.5,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19774,
            "unit": "ns/op\t6665 B/op\t74 allocs/op",
            "extra": "6192 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19774,
            "unit": "ns/op",
            "extra": "6192 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6665,
            "unit": "B/op",
            "extra": "6192 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "6192 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 72024.5,
            "unit": "ns/op\t58332.5 B/op\t448 allocs/op",
            "extra": "1670 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 72024.5,
            "unit": "ns/op",
            "extra": "1670 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58332.5,
            "unit": "B/op",
            "extra": "1670 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "1670 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 478636,
            "unit": "ns/op\t538101 B/op\t3940 allocs/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 478636,
            "unit": "ns/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538101,
            "unit": "B/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3514.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "33410 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3514.5,
            "unit": "ns/op",
            "extra": "33410 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "33410 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "33410 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44978.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2727 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44978.5,
            "unit": "ns/op",
            "extra": "2727 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2727 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2727 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 479024.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 479024.5,
            "unit": "ns/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "248 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 10314,
            "unit": "ns/op\t9648 B/op\t116 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 10314,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9648,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 28555.5,
            "unit": "ns/op\t9653 B/op\t116 allocs/op",
            "extra": "4185 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 28555.5,
            "unit": "ns/op",
            "extra": "4185 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9653,
            "unit": "B/op",
            "extra": "4185 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "4185 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 213822,
            "unit": "ns/op\t9624 B/op\t116 allocs/op",
            "extra": "560 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 213822,
            "unit": "ns/op",
            "extra": "560 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9624,
            "unit": "B/op",
            "extra": "560 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "560 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.813,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14875935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.813,
            "unit": "ns/op",
            "extra": "14875935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14875935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14875935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 390.55,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "304717 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 390.55,
            "unit": "ns/op",
            "extra": "304717 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "304717 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "304717 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 124.5,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 124.5,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3549,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "33490 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3549,
            "unit": "ns/op",
            "extra": "33490 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "33490 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "33490 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3971,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "29935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3971,
            "unit": "ns/op",
            "extra": "29935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "29935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "29935 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.8145,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "15106237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.8145,
            "unit": "ns/op",
            "extra": "15106237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15106237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15106237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3387,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "34520 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3387,
            "unit": "ns/op",
            "extra": "34520 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "34520 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "34520 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 384.45,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "295374 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 384.45,
            "unit": "ns/op",
            "extra": "295374 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "295374 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "295374 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44827,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2736 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44827,
            "unit": "ns/op",
            "extra": "2736 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2736 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2736 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 48183.5,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "2466 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 48183.5,
            "unit": "ns/op",
            "extra": "2466 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "2466 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "2466 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.7875,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "15306307 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.7875,
            "unit": "ns/op",
            "extra": "15306307 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15306307 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15306307 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 35002,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "3446 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 35002,
            "unit": "ns/op",
            "extra": "3446 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "3446 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3446 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3673,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "31683 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3673,
            "unit": "ns/op",
            "extra": "31683 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "31683 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "31683 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 476072,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "249 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 476072,
            "unit": "ns/op",
            "extra": "249 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "249 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "249 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 490303,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "243 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 490303,
            "unit": "ns/op",
            "extra": "243 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "243 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "243 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 518.95,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "225202 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 518.95,
            "unit": "ns/op",
            "extra": "225202 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "225202 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "225202 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3412,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "35067 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3412,
            "unit": "ns/op",
            "extra": "35067 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "35067 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "35067 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 30896,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "3586 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 30896,
            "unit": "ns/op",
            "extra": "3586 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063,
            "unit": "B/op",
            "extra": "3586 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "3586 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52887,
            "unit": "ns/op\t18523.5 B/op\t178 allocs/op",
            "extra": "2312 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52887,
            "unit": "ns/op",
            "extra": "2312 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18523.5,
            "unit": "B/op",
            "extra": "2312 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "2312 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 193377.5,
            "unit": "ns/op\t74840 B/op\t633 allocs/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 193377.5,
            "unit": "ns/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 74840,
            "unit": "B/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 633,
            "unit": "allocs/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 603035.5,
            "unit": "ns/op\t478066.5 B/op\t1418 allocs/op",
            "extra": "196 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 603035.5,
            "unit": "ns/op",
            "extra": "196 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478066.5,
            "unit": "B/op",
            "extra": "196 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1418,
            "unit": "allocs/op",
            "extra": "196 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 669.05,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "176758 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 669.05,
            "unit": "ns/op",
            "extra": "176758 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "176758 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "176758 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4464,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "26710 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4464,
            "unit": "ns/op",
            "extra": "26710 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "26710 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26710 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4583.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "26136 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4583.5,
            "unit": "ns/op",
            "extra": "26136 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "26136 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26136 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2638,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "44827 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2638,
            "unit": "ns/op",
            "extra": "44827 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "44827 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "44827 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 674.6,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "174912 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 674.6,
            "unit": "ns/op",
            "extra": "174912 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "174912 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "174912 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4507.5,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "26854 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4507.5,
            "unit": "ns/op",
            "extra": "26854 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "26854 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26854 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4585,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "25965 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4585,
            "unit": "ns/op",
            "extra": "25965 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "25965 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "25965 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2652.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "45121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2652.5,
            "unit": "ns/op",
            "extra": "45121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "45121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6585,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "18192 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6585,
            "unit": "ns/op",
            "extra": "18192 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "18192 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "18192 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 41639.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "2775 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41639.5,
            "unit": "ns/op",
            "extra": "2775 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "2775 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2775 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39794,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "2926 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39794,
            "unit": "ns/op",
            "extra": "2926 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "2926 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2926 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 24510.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "4843 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 24510.5,
            "unit": "ns/op",
            "extra": "4843 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "4843 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "4843 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6653.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "17737 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6653.5,
            "unit": "ns/op",
            "extra": "17737 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "17737 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "17737 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40378.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "2958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40378.5,
            "unit": "ns/op",
            "extra": "2958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "2958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40367.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "2889 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40367.5,
            "unit": "ns/op",
            "extra": "2889 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "2889 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2889 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 24573.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "4651 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 24573.5,
            "unit": "ns/op",
            "extra": "4651 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "4651 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "4651 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 135035.5,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "864 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 135035.5,
            "unit": "ns/op",
            "extra": "864 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "864 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "864 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 446446,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 446446,
            "unit": "ns/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 439674.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 439674.5,
            "unit": "ns/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 251147,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 251147,
            "unit": "ns/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 140161,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 140161,
            "unit": "ns/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 448069,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "268 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 448069,
            "unit": "ns/op",
            "extra": "268 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "268 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "268 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 440146.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 440146.5,
            "unit": "ns/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 251794,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 251794,
            "unit": "ns/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 378.45,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "316467 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 378.45,
            "unit": "ns/op",
            "extra": "316467 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "316467 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "316467 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 372.65,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "313861 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 372.65,
            "unit": "ns/op",
            "extra": "313861 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "313861 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "313861 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2898,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "41522 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2898,
            "unit": "ns/op",
            "extra": "41522 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "41522 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "41522 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2885,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "41574 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2885,
            "unit": "ns/op",
            "extra": "41574 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "41574 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "41574 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39860.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2862 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39860.5,
            "unit": "ns/op",
            "extra": "2862 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2862 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2862 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39231,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "3237 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39231,
            "unit": "ns/op",
            "extra": "3237 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "3237 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3237 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 838.9,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "149529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 838.9,
            "unit": "ns/op",
            "extra": "149529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "149529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "149529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1384,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "81960 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1384,
            "unit": "ns/op",
            "extra": "81960 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "81960 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "81960 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1327,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "87811 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1327,
            "unit": "ns/op",
            "extra": "87811 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "87811 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "87811 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 7016,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 7016,
            "unit": "ns/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 12726.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12726.5,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 30336,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 201,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 12699.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12699.5,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 102745.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "1112 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 102745.5,
            "unit": "ns/op",
            "extra": "1112 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "1112 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1112 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 150970,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "796 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 150970,
            "unit": "ns/op",
            "extra": "796 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "796 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "796 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 150492.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 150492.5,
            "unit": "ns/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8327.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14414 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8327.5,
            "unit": "ns/op",
            "extra": "14414 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14414 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14414 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8409.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14266 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8409.5,
            "unit": "ns/op",
            "extra": "14266 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14266 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14266 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 82666.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1432 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 82666.5,
            "unit": "ns/op",
            "extra": "1432 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1432 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1432 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 84660,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1363 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84660,
            "unit": "ns/op",
            "extra": "1363 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1363 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1363 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 840713.5,
            "unit": "ns/op\t9 B/op\t0 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 840713.5,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 846006,
            "unit": "ns/op\t8 B/op\t0 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 846006,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 8,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 797.4,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "148680 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 797.4,
            "unit": "ns/op",
            "extra": "148680 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "148680 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "148680 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 815.55,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 815.55,
            "unit": "ns/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 856.85,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "144642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 856.85,
            "unit": "ns/op",
            "extra": "144642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "144642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "144642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4439,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "26485 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4439,
            "unit": "ns/op",
            "extra": "26485 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "26485 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "26485 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4448,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "26605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4448,
            "unit": "ns/op",
            "extra": "26605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "26605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "26605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4766.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "25195 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4766.5,
            "unit": "ns/op",
            "extra": "25195 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "25195 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "25195 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51667,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2326 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51667,
            "unit": "ns/op",
            "extra": "2326 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2326 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2326 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51421,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2277 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51421,
            "unit": "ns/op",
            "extra": "2277 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2277 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2277 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52823.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2089 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52823.5,
            "unit": "ns/op",
            "extra": "2089 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2089 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2089 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 22190.5,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "5524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 22190.5,
            "unit": "ns/op",
            "extra": "5524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "5524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "5524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 53278,
            "unit": "ns/op\t6041 B/op\t152 allocs/op",
            "extra": "2236 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 53278,
            "unit": "ns/op",
            "extra": "2236 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6041,
            "unit": "B/op",
            "extra": "2236 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "2236 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52963.5,
            "unit": "ns/op\t4521 B/op\t151 allocs/op",
            "extra": "2269 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52963.5,
            "unit": "ns/op",
            "extra": "2269 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4521,
            "unit": "B/op",
            "extra": "2269 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "2269 times\n4 procs"
          }
        ]
      }
    ]
  }
}