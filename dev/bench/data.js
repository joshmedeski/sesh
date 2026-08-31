window.BENCHMARK_DATA = {
  "lastUpdate": 1788208105064,
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
          "id": "8744ad031bd97e9519cecedf80c48f932e28c33d",
          "message": "docs: add Fedora Copr install instructions",
          "timestamp": "2026-08-25T17:37:20-05:00",
          "tree_id": "230062c64b124afaa1e0e3f5f51afbe1a2ce43fe",
          "url": "https://github.com/joshmedeski/sesh/commit/8744ad031bd97e9519cecedf80c48f932e28c33d"
        },
        "date": 1787697545600,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 20586.5,
            "unit": "ns/op\t6634 B/op\t73 allocs/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20586.5,
            "unit": "ns/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6634,
            "unit": "B/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 66265.5,
            "unit": "ns/op\t58302 B/op\t447 allocs/op",
            "extra": "1856 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 66265.5,
            "unit": "ns/op",
            "extra": "1856 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58302,
            "unit": "B/op",
            "extra": "1856 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "1856 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 428196.5,
            "unit": "ns/op\t538052.5 B/op\t3939 allocs/op",
            "extra": "279 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 428196.5,
            "unit": "ns/op",
            "extra": "279 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538052.5,
            "unit": "B/op",
            "extra": "279 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "279 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 29365.5,
            "unit": "ns/op\t8481.5 B/op\t87 allocs/op",
            "extra": "4242 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 29365.5,
            "unit": "ns/op",
            "extra": "4242 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8481.5,
            "unit": "B/op",
            "extra": "4242 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "4242 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 112209,
            "unit": "ns/op\t79305 B/op\t489 allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 112209,
            "unit": "ns/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79305,
            "unit": "B/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 906452.5,
            "unit": "ns/op\t746362 B/op\t4032.5 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 906452.5,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746362,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4032.5,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44198,
            "unit": "ns/op\t18698 B/op\t204 allocs/op",
            "extra": "2794 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44198,
            "unit": "ns/op",
            "extra": "2794 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18698,
            "unit": "B/op",
            "extra": "2794 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 204,
            "unit": "allocs/op",
            "extra": "2794 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 120766.5,
            "unit": "ns/op\t92360 B/op\t662 allocs/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 120766.5,
            "unit": "ns/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92360,
            "unit": "B/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 662,
            "unit": "allocs/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 832544,
            "unit": "ns/op\t780633.5 B/op\t4946 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 832544,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 780633.5,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 20731.5,
            "unit": "ns/op\t6665 B/op\t74 allocs/op",
            "extra": "5437 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20731.5,
            "unit": "ns/op",
            "extra": "5437 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6665,
            "unit": "B/op",
            "extra": "5437 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "5437 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 70561,
            "unit": "ns/op\t58333.5 B/op\t448 allocs/op",
            "extra": "1729 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 70561,
            "unit": "ns/op",
            "extra": "1729 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58333.5,
            "unit": "B/op",
            "extra": "1729 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "1729 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 491051.5,
            "unit": "ns/op\t538090.5 B/op\t3940 allocs/op",
            "extra": "247 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 491051.5,
            "unit": "ns/op",
            "extra": "247 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538090.5,
            "unit": "B/op",
            "extra": "247 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "247 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3621,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "33164 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3621,
            "unit": "ns/op",
            "extra": "33164 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "33164 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "33164 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 42475,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2990 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 42475,
            "unit": "ns/op",
            "extra": "2990 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2990 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2990 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 449270.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "258 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 449270.5,
            "unit": "ns/op",
            "extra": "258 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "258 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "258 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 10733,
            "unit": "ns/op\t9653.5 B/op\t116 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 10733,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9653.5,
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
            "value": 28130.5,
            "unit": "ns/op\t9641.5 B/op\t116 allocs/op",
            "extra": "4353 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 28130.5,
            "unit": "ns/op",
            "extra": "4353 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9641.5,
            "unit": "B/op",
            "extra": "4353 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "4353 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 203055,
            "unit": "ns/op\t9650.5 B/op\t116 allocs/op",
            "extra": "600 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 203055,
            "unit": "ns/op",
            "extra": "600 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9650.5,
            "unit": "B/op",
            "extra": "600 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "600 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.27,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "16046529 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.27,
            "unit": "ns/op",
            "extra": "16046529 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16046529 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16046529 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 398.15,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "301616 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 398.15,
            "unit": "ns/op",
            "extra": "301616 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "301616 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "301616 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 126.65,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 126.65,
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
            "value": 3702.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "32007 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3702.5,
            "unit": "ns/op",
            "extra": "32007 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "32007 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "32007 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 4199.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "28334 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4199.5,
            "unit": "ns/op",
            "extra": "28334 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "28334 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "28334 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.3215,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14095627 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.3215,
            "unit": "ns/op",
            "extra": "14095627 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14095627 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14095627 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3480.5,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "33318 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3480.5,
            "unit": "ns/op",
            "extra": "33318 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "33318 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "33318 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 463.85,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "232946 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 463.85,
            "unit": "ns/op",
            "extra": "232946 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "232946 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "232946 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44579.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2850 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44579.5,
            "unit": "ns/op",
            "extra": "2850 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2850 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2850 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 48293,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "2685 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 48293,
            "unit": "ns/op",
            "extra": "2685 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "2685 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "2685 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.249,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "16262154 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.249,
            "unit": "ns/op",
            "extra": "16262154 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16262154 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16262154 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 36042.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "3541 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 36042.5,
            "unit": "ns/op",
            "extra": "3541 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "3541 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3541 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 4185.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "28224 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4185.5,
            "unit": "ns/op",
            "extra": "28224 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "28224 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "28224 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 464082,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "255 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 464082,
            "unit": "ns/op",
            "extra": "255 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "255 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "255 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 483260,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 483260,
            "unit": "ns/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 513.05,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "235312 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 513.05,
            "unit": "ns/op",
            "extra": "235312 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "235312 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "235312 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3117,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "38410 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3117,
            "unit": "ns/op",
            "extra": "38410 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "38410 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "38410 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 28304,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "4269 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 28304,
            "unit": "ns/op",
            "extra": "4269 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063,
            "unit": "B/op",
            "extra": "4269 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "4269 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 49394,
            "unit": "ns/op\t18523 B/op\t178 allocs/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 49394,
            "unit": "ns/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18523,
            "unit": "B/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 182741.5,
            "unit": "ns/op\t74839.5 B/op\t633 allocs/op",
            "extra": "648 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 182741.5,
            "unit": "ns/op",
            "extra": "648 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 74839.5,
            "unit": "B/op",
            "extra": "648 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 633,
            "unit": "allocs/op",
            "extra": "648 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 600205.5,
            "unit": "ns/op\t478081 B/op\t1418 allocs/op",
            "extra": "198 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 600205.5,
            "unit": "ns/op",
            "extra": "198 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478081,
            "unit": "B/op",
            "extra": "198 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1418,
            "unit": "allocs/op",
            "extra": "198 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 776.15,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "154910 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 776.15,
            "unit": "ns/op",
            "extra": "154910 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "154910 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "154910 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4819,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "24381 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4819,
            "unit": "ns/op",
            "extra": "24381 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "24381 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "24381 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4893,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "24363 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4893,
            "unit": "ns/op",
            "extra": "24363 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "24363 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "24363 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2222.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "54402 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2222.5,
            "unit": "ns/op",
            "extra": "54402 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "54402 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "54402 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 741.7,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "158362 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 741.7,
            "unit": "ns/op",
            "extra": "158362 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "158362 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "158362 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4626,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "25861 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4626,
            "unit": "ns/op",
            "extra": "25861 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "25861 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "25861 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4676,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "25592 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4676,
            "unit": "ns/op",
            "extra": "25592 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "25592 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "25592 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2234.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "54664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2234.5,
            "unit": "ns/op",
            "extra": "54664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "54664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "54664 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6409.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "18814 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6409.5,
            "unit": "ns/op",
            "extra": "18814 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "18814 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "18814 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 38531.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 38531.5,
            "unit": "ns/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39588,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3099 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39588,
            "unit": "ns/op",
            "extra": "3099 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3099 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3099 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 20557.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "5696 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 20557.5,
            "unit": "ns/op",
            "extra": "5696 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "5696 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "5696 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6576,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "17829 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6576,
            "unit": "ns/op",
            "extra": "17829 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "17829 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "17829 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40221,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "2988 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40221,
            "unit": "ns/op",
            "extra": "2988 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "2988 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2988 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 41593.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "2941 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41593.5,
            "unit": "ns/op",
            "extra": "2941 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "2941 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2941 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 20619.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "5674 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 20619.5,
            "unit": "ns/op",
            "extra": "5674 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "5674 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "5674 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 146839.5,
            "unit": "ns/op\t278529.5 B/op\t1 allocs/op",
            "extra": "688 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 146839.5,
            "unit": "ns/op",
            "extra": "688 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529.5,
            "unit": "B/op",
            "extra": "688 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "688 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 450874.5,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 450874.5,
            "unit": "ns/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "272 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 445133,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 445133,
            "unit": "ns/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 210481,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "544 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 210481,
            "unit": "ns/op",
            "extra": "544 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "544 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "544 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 178169,
            "unit": "ns/op\t278528.5 B/op\t1 allocs/op",
            "extra": "634 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 178169,
            "unit": "ns/op",
            "extra": "634 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278528.5,
            "unit": "B/op",
            "extra": "634 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "634 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 453146,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 453146,
            "unit": "ns/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 448439.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 448439.5,
            "unit": "ns/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 266103.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 266103.5,
            "unit": "ns/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 402.75,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "290053 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 402.75,
            "unit": "ns/op",
            "extra": "290053 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "290053 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "290053 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 408.5,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "297045 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 408.5,
            "unit": "ns/op",
            "extra": "297045 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "297045 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "297045 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3169,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "37780 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3169,
            "unit": "ns/op",
            "extra": "37780 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "37780 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "37780 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3111.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "39261 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3111.5,
            "unit": "ns/op",
            "extra": "39261 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "39261 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "39261 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 43325.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2731 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 43325.5,
            "unit": "ns/op",
            "extra": "2731 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2731 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2731 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 54722.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2336 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54722.5,
            "unit": "ns/op",
            "extra": "2336 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2336 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2336 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 995.7,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "108712 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 995.7,
            "unit": "ns/op",
            "extra": "108712 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "108712 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "108712 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1584.5,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "74042 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1584.5,
            "unit": "ns/op",
            "extra": "74042 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "74042 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "74042 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1465.5,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "82260 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1465.5,
            "unit": "ns/op",
            "extra": "82260 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "82260 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "82260 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 7596,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "15524 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 7596,
            "unit": "ns/op",
            "extra": "15524 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "15524 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15524 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 13607.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 13607.5,
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
            "value": 13561,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "9003 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 13561,
            "unit": "ns/op",
            "extra": "9003 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "9003 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9003 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 84206.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84206.5,
            "unit": "ns/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 156570,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "728 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 156570,
            "unit": "ns/op",
            "extra": "728 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "728 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "728 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 140185,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "824 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 140185,
            "unit": "ns/op",
            "extra": "824 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "824 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "824 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8220,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14521 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8220,
            "unit": "ns/op",
            "extra": "14521 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14521 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14521 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8322.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14338 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8322.5,
            "unit": "ns/op",
            "extra": "14338 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14338 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14338 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 83434.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83434.5,
            "unit": "ns/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 84139.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84139.5,
            "unit": "ns/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 849223,
            "unit": "ns/op\t9 B/op\t0 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 849223,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 855003.5,
            "unit": "ns/op\t8 B/op\t0 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 855003.5,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 8,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 927.35,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "133053 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 927.35,
            "unit": "ns/op",
            "extra": "133053 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "133053 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "133053 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 937.3,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "129559 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 937.3,
            "unit": "ns/op",
            "extra": "129559 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "129559 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "129559 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 951.6,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "124902 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 951.6,
            "unit": "ns/op",
            "extra": "124902 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "124902 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "124902 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4641.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "25713 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4641.5,
            "unit": "ns/op",
            "extra": "25713 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "25713 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "25713 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4790,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "25022 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4790,
            "unit": "ns/op",
            "extra": "25022 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "25022 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "25022 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 5107,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "23320 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5107,
            "unit": "ns/op",
            "extra": "23320 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "23320 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "23320 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 53775.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2359 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 53775.5,
            "unit": "ns/op",
            "extra": "2359 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2359 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2359 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 53789.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 53789.5,
            "unit": "ns/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 55804,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2347 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 55804,
            "unit": "ns/op",
            "extra": "2347 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2347 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2347 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 20182,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "5838 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 20182,
            "unit": "ns/op",
            "extra": "5838 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "5838 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "5838 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 48390,
            "unit": "ns/op\t6041 B/op\t152 allocs/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 48390,
            "unit": "ns/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6041,
            "unit": "B/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 47502,
            "unit": "ns/op\t4521 B/op\t151 allocs/op",
            "extra": "2524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 47502,
            "unit": "ns/op",
            "extra": "2524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4521,
            "unit": "B/op",
            "extra": "2524 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "2524 times\n4 procs"
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
          "id": "9c12792b02e847c2aa4e93675473c5674de2f999",
          "message": "fix: don't rerun the create path for wildcard-matched sessions\n\nconfigWildcardStrategy returned New: true unconditionally, so connecting by\npath to a directory matching a [[wildcard]] pattern took the create path even\nwhen the tmux session was already running. connectToTmux then reissued\nnew-session (its error discarded) and ran startup.Exec, which send-keys'd the\nstartup command into the live session — falling through to\n[default_session].startup_command when the wildcard defined none of its own.\n\nBecause configWildcardStrategy runs before dirStrategy, the strategy that does\ncheck for an existing session never got a chance to. Give the wildcard strategy\nthe same FindTmuxSessionByBase branch dirStrategy and zoxideStrategy have.\n\nAlso affects `sesh worktree connect`: it correctly passes an empty command when\nreconnecting to an existing worktree, but that only governs opts.Command, and\nthe startup config ran anyway.\n\nFixes #453",
          "timestamp": "2026-08-31T11:31:56-05:00",
          "tree_id": "c93d86e3efc6888b9b5d2505bae7b15f4298724e",
          "url": "https://github.com/joshmedeski/sesh/commit/9c12792b02e847c2aa4e93675473c5674de2f999"
        },
        "date": 1788194028144,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 18597,
            "unit": "ns/op\t6633 B/op\t73 allocs/op",
            "extra": "6151 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 18597,
            "unit": "ns/op",
            "extra": "6151 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6633,
            "unit": "B/op",
            "extra": "6151 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "6151 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 69539,
            "unit": "ns/op\t58300.5 B/op\t447 allocs/op",
            "extra": "1761 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 69539,
            "unit": "ns/op",
            "extra": "1761 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58300.5,
            "unit": "B/op",
            "extra": "1761 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "1761 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 445355,
            "unit": "ns/op\t538066 B/op\t3939 allocs/op",
            "extra": "270 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 445355,
            "unit": "ns/op",
            "extra": "270 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538066,
            "unit": "B/op",
            "extra": "270 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "270 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 27296.5,
            "unit": "ns/op\t8481.5 B/op\t87 allocs/op",
            "extra": "4327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27296.5,
            "unit": "ns/op",
            "extra": "4327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8481.5,
            "unit": "B/op",
            "extra": "4327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "4327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 115593.5,
            "unit": "ns/op\t79304.5 B/op\t489 allocs/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 115593.5,
            "unit": "ns/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79304.5,
            "unit": "B/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 944879,
            "unit": "ns/op\t746365.5 B/op\t4032.5 allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 944879,
            "unit": "ns/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746365.5,
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
            "value": 42624,
            "unit": "ns/op\t18760.5 B/op\t204 allocs/op",
            "extra": "2860 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 42624,
            "unit": "ns/op",
            "extra": "2860 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18760.5,
            "unit": "B/op",
            "extra": "2860 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 204,
            "unit": "allocs/op",
            "extra": "2860 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 129158.5,
            "unit": "ns/op\t92936.5 B/op\t662 allocs/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 129158.5,
            "unit": "ns/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92936.5,
            "unit": "B/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 662,
            "unit": "allocs/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 862655.5,
            "unit": "ns/op\t781426.5 B/op\t4946 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 862655.5,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 781426.5,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19262,
            "unit": "ns/op\t6665 B/op\t74 allocs/op",
            "extra": "6080 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19262,
            "unit": "ns/op",
            "extra": "6080 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6665,
            "unit": "B/op",
            "extra": "6080 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "6080 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 74596.5,
            "unit": "ns/op\t58332 B/op\t448 allocs/op",
            "extra": "1712 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 74596.5,
            "unit": "ns/op",
            "extra": "1712 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58332,
            "unit": "B/op",
            "extra": "1712 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "1712 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 480281.5,
            "unit": "ns/op\t538090 B/op\t3940 allocs/op",
            "extra": "253 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 480281.5,
            "unit": "ns/op",
            "extra": "253 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538090,
            "unit": "B/op",
            "extra": "253 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "253 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3636.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3636.5,
            "unit": "ns/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44683.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2628 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44683.5,
            "unit": "ns/op",
            "extra": "2628 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2628 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2628 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 477302,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 477302,
            "unit": "ns/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 10232.5,
            "unit": "ns/op\t9659 B/op\t116 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 10232.5,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9659,
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
            "value": 28376,
            "unit": "ns/op\t9639 B/op\t116 allocs/op",
            "extra": "4267 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 28376,
            "unit": "ns/op",
            "extra": "4267 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9639,
            "unit": "B/op",
            "extra": "4267 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "4267 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 214036,
            "unit": "ns/op\t9688 B/op\t116 allocs/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 214036,
            "unit": "ns/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9688,
            "unit": "B/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.8005,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "15273592 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.8005,
            "unit": "ns/op",
            "extra": "15273592 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15273592 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15273592 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 397.75,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "308727 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 397.75,
            "unit": "ns/op",
            "extra": "308727 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "308727 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "308727 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 126.45,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 126.45,
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
            "value": 3534,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "33860 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3534,
            "unit": "ns/op",
            "extra": "33860 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "33860 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "33860 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3974,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "30184 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3974,
            "unit": "ns/op",
            "extra": "30184 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "30184 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "30184 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.8365,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14941401 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.8365,
            "unit": "ns/op",
            "extra": "14941401 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14941401 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14941401 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3314,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3314,
            "unit": "ns/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 379.95,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "301617 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 379.95,
            "unit": "ns/op",
            "extra": "301617 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "301617 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "301617 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 45290,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2829 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 45290,
            "unit": "ns/op",
            "extra": "2829 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2829 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2829 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 48394,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "2673 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 48394,
            "unit": "ns/op",
            "extra": "2673 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "2673 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "2673 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 7.8015,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "15352638 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 7.8015,
            "unit": "ns/op",
            "extra": "15352638 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15352638 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15352638 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 34952.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 34952.5,
            "unit": "ns/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3635,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "33054 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3635,
            "unit": "ns/op",
            "extra": "33054 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "33054 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "33054 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 480091.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 480091.5,
            "unit": "ns/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "246 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 493663,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 493663,
            "unit": "ns/op",
            "extra": "237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "237 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "237 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 519.65,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "229290 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 519.65,
            "unit": "ns/op",
            "extra": "229290 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "229290 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "229290 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3402.5,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "35482 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3402.5,
            "unit": "ns/op",
            "extra": "35482 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "35482 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "35482 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 30799.5,
            "unit": "ns/op\t33063.5 B/op\t321 allocs/op",
            "extra": "3655 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 30799.5,
            "unit": "ns/op",
            "extra": "3655 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063.5,
            "unit": "B/op",
            "extra": "3655 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "3655 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52413.5,
            "unit": "ns/op\t18523.5 B/op\t178 allocs/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52413.5,
            "unit": "ns/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18523.5,
            "unit": "B/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 193533,
            "unit": "ns/op\t74840 B/op\t633 allocs/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 193533,
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
            "value": 591685.5,
            "unit": "ns/op\t478058.5 B/op\t1417 allocs/op",
            "extra": "200 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 591685.5,
            "unit": "ns/op",
            "extra": "200 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478058.5,
            "unit": "B/op",
            "extra": "200 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1417,
            "unit": "allocs/op",
            "extra": "200 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 663.1,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "173511 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 663.1,
            "unit": "ns/op",
            "extra": "173511 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "173511 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "173511 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4550,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "26469 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4550,
            "unit": "ns/op",
            "extra": "26469 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "26469 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26469 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4548,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "26248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4548,
            "unit": "ns/op",
            "extra": "26248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "26248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2650,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "44316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2650,
            "unit": "ns/op",
            "extra": "44316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "44316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "44316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 670.25,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "174745 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 670.25,
            "unit": "ns/op",
            "extra": "174745 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "174745 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "174745 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4492,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "26533 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4492,
            "unit": "ns/op",
            "extra": "26533 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "26533 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26533 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4574,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "26466 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4574,
            "unit": "ns/op",
            "extra": "26466 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "26466 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26466 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2648.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "44280 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2648.5,
            "unit": "ns/op",
            "extra": "44280 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "44280 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "44280 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6808,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "17846 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6808,
            "unit": "ns/op",
            "extra": "17846 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "17846 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "17846 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39606,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "3009 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39606,
            "unit": "ns/op",
            "extra": "3009 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "3009 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3009 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39579,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39579,
            "unit": "ns/op",
            "extra": "3045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3045 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 24603.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "4822 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 24603.5,
            "unit": "ns/op",
            "extra": "4822 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "4822 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "4822 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6720.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "17605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6720.5,
            "unit": "ns/op",
            "extra": "17605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "17605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "17605 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40112,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "2973 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40112,
            "unit": "ns/op",
            "extra": "2973 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "2973 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2973 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40255.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "2938 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40255.5,
            "unit": "ns/op",
            "extra": "2938 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "2938 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "2938 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 24645,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "4642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 24645,
            "unit": "ns/op",
            "extra": "4642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "4642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "4642 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 135931,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 135931,
            "unit": "ns/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 444891,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 444891,
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
            "value": 443170.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 443170.5,
            "unit": "ns/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "271 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 251545,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 251545,
            "unit": "ns/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 134173.5,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "838 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 134173.5,
            "unit": "ns/op",
            "extra": "838 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "838 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "838 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 445142,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 445142,
            "unit": "ns/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "260 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 447390.5,
            "unit": "ns/op\t455904.5 B/op\t922 allocs/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 447390.5,
            "unit": "ns/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904.5,
            "unit": "B/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "265 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 251279,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 251279,
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
            "value": 370.7,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "318811 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 370.7,
            "unit": "ns/op",
            "extra": "318811 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "318811 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "318811 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 370.95,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "315470 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 370.95,
            "unit": "ns/op",
            "extra": "315470 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "315470 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "315470 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2870.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "41520 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2870.5,
            "unit": "ns/op",
            "extra": "41520 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "41520 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "41520 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2875,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "41619 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2875,
            "unit": "ns/op",
            "extra": "41619 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "41619 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "41619 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 41547.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2558 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41547.5,
            "unit": "ns/op",
            "extra": "2558 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2558 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2558 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 40843,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40843,
            "unit": "ns/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2708 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 822.6,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "132452 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 822.6,
            "unit": "ns/op",
            "extra": "132452 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "132452 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "132452 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1399,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "85890 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1399,
            "unit": "ns/op",
            "extra": "85890 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "85890 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "85890 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1324,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "90774 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1324,
            "unit": "ns/op",
            "extra": "90774 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "90774 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "90774 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 7240.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "16533 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 7240.5,
            "unit": "ns/op",
            "extra": "16533 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "16533 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "16533 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 12634,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12634,
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
            "value": 12940.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12940.5,
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
            "value": 98634,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 98634,
            "unit": "ns/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 148348.5,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "794 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 148348.5,
            "unit": "ns/op",
            "extra": "794 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "794 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "794 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 152637.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "780 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 152637.5,
            "unit": "ns/op",
            "extra": "780 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "780 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "780 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8366,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14328 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8366,
            "unit": "ns/op",
            "extra": "14328 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14328 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14328 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8363,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "14326 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8363,
            "unit": "ns/op",
            "extra": "14326 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14326 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14326 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 84261,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84261,
            "unit": "ns/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 84560.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1335 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84560.5,
            "unit": "ns/op",
            "extra": "1335 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1335 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1335 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 863857.5,
            "unit": "ns/op\t9 B/op\t0 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 863857.5,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 865232,
            "unit": "ns/op\t9 B/op\t0 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 865232,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 758.35,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "155931 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 758.35,
            "unit": "ns/op",
            "extra": "155931 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "155931 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "155931 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 787.35,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "156150 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 787.35,
            "unit": "ns/op",
            "extra": "156150 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "156150 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "156150 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 819.3,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "152095 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 819.3,
            "unit": "ns/op",
            "extra": "152095 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "152095 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "152095 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4478.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4478.5,
            "unit": "ns/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4476.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "26751 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4476.5,
            "unit": "ns/op",
            "extra": "26751 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "26751 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "26751 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4788,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "25083 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4788,
            "unit": "ns/op",
            "extra": "25083 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "25083 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "25083 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51536,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51536,
            "unit": "ns/op",
            "extra": "2316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2316 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51297,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2233 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51297,
            "unit": "ns/op",
            "extra": "2233 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2233 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2233 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 53018.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 53018.5,
            "unit": "ns/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 21955,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "5446 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 21955,
            "unit": "ns/op",
            "extra": "5446 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "5446 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "5446 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52900.5,
            "unit": "ns/op\t6041 B/op\t152 allocs/op",
            "extra": "2197 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52900.5,
            "unit": "ns/op",
            "extra": "2197 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6041,
            "unit": "B/op",
            "extra": "2197 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "2197 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 53041.5,
            "unit": "ns/op\t4521 B/op\t151 allocs/op",
            "extra": "2265 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 53041.5,
            "unit": "ns/op",
            "extra": "2265 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4521,
            "unit": "B/op",
            "extra": "2265 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "2265 times\n4 procs"
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
          "id": "621e2c9294125d60c927efb8703d10be28714768",
          "message": "ci: clear stale mocks before regenerating them\n\nThe benchmark job generates mocks, checks out the merge base to benchmark it,\nand regenerates mocks there. Mocks are gitignored, so they survive that\ncheckout — and mockery cannot regenerate a mock for a package that no longer\ncompiles. A mock generated at one revision that references a type absent from\nthe other therefore kills mockery before it can replace itself:\n\n    tmux/mock_Tmux.go:605:54: undefined: model.TmuxWindowOpts\n    FTL app failed error=\"error occurred when loading packages\"\n\nAny PR that changes a type used in a mocked signature hits this, so it has been\nlatent rather than absent. Delete the generated mocks before each mockery run,\nin the workflow and in `just mock`, so a signature change can't block its own\nregeneration locally either.",
          "timestamp": "2026-08-31T15:26:36-05:00",
          "tree_id": "df28b47ca07238c5c9f6b3b4edc10080ce993169",
          "url": "https://github.com/joshmedeski/sesh/commit/621e2c9294125d60c927efb8703d10be28714768"
        },
        "date": 1788208103810,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 18170,
            "unit": "ns/op\t6633 B/op\t73 allocs/op",
            "extra": "6692 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 18170,
            "unit": "ns/op",
            "extra": "6692 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6633,
            "unit": "B/op",
            "extra": "6692 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "6692 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 58828,
            "unit": "ns/op\t58300.5 B/op\t447 allocs/op",
            "extra": "2019 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 58828,
            "unit": "ns/op",
            "extra": "2019 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58300.5,
            "unit": "B/op",
            "extra": "2019 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "2019 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 406634,
            "unit": "ns/op\t538057 B/op\t3939 allocs/op",
            "extra": "292 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 406634,
            "unit": "ns/op",
            "extra": "292 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538057,
            "unit": "B/op",
            "extra": "292 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "292 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 25504,
            "unit": "ns/op\t8481.5 B/op\t87 allocs/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 25504,
            "unit": "ns/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8481.5,
            "unit": "B/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 108772.5,
            "unit": "ns/op\t79303.5 B/op\t489 allocs/op",
            "extra": "1116 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 108772.5,
            "unit": "ns/op",
            "extra": "1116 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79303.5,
            "unit": "B/op",
            "extra": "1116 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "1116 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 904988,
            "unit": "ns/op\t746363.5 B/op\t4032.5 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 904988,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746363.5,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4032.5,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 37968.5,
            "unit": "ns/op\t18702.5 B/op\t204 allocs/op",
            "extra": "3142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 37968.5,
            "unit": "ns/op",
            "extra": "3142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18702.5,
            "unit": "B/op",
            "extra": "3142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 204,
            "unit": "allocs/op",
            "extra": "3142 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 108867.5,
            "unit": "ns/op\t92355 B/op\t662 allocs/op",
            "extra": "1071 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 108867.5,
            "unit": "ns/op",
            "extra": "1071 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92355,
            "unit": "B/op",
            "extra": "1071 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 662,
            "unit": "allocs/op",
            "extra": "1071 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 785733,
            "unit": "ns/op\t780621 B/op\t4946 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 785733,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 780621,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19752,
            "unit": "ns/op\t6665 B/op\t74 allocs/op",
            "extra": "6187 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19752,
            "unit": "ns/op",
            "extra": "6187 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6665,
            "unit": "B/op",
            "extra": "6187 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "6187 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 63718,
            "unit": "ns/op\t58332.5 B/op\t448 allocs/op",
            "extra": "1946 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 63718,
            "unit": "ns/op",
            "extra": "1946 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58332.5,
            "unit": "B/op",
            "extra": "1946 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "1946 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 452975,
            "unit": "ns/op\t538090.5 B/op\t3940 allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 452975,
            "unit": "ns/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538090.5,
            "unit": "B/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3421.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "34550 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3421.5,
            "unit": "ns/op",
            "extra": "34550 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "34550 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "34550 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44214.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2674 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44214.5,
            "unit": "ns/op",
            "extra": "2674 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2674 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2674 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 463998,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 463998,
            "unit": "ns/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9897,
            "unit": "ns/op\t9654.5 B/op\t116 allocs/op",
            "extra": "12049 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9897,
            "unit": "ns/op",
            "extra": "12049 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9654.5,
            "unit": "B/op",
            "extra": "12049 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "12049 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 27649,
            "unit": "ns/op\t9644 B/op\t116 allocs/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27649,
            "unit": "ns/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9644,
            "unit": "B/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 207104.5,
            "unit": "ns/op\t9620.5 B/op\t116 allocs/op",
            "extra": "577 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 207104.5,
            "unit": "ns/op",
            "extra": "577 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9620.5,
            "unit": "B/op",
            "extra": "577 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 116,
            "unit": "allocs/op",
            "extra": "577 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.157,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "13057975 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.157,
            "unit": "ns/op",
            "extra": "13057975 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13057975 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13057975 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 386.4,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "287682 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 386.4,
            "unit": "ns/op",
            "extra": "287682 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "287682 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "287682 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 122.55,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 122.55,
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
            "value": 3439,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3439,
            "unit": "ns/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "34500 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 4200.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "27741 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4200.5,
            "unit": "ns/op",
            "extra": "27741 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "27741 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "27741 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.284,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "12997831 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.284,
            "unit": "ns/op",
            "extra": "12997831 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12997831 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12997831 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3303.5,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "35265 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3303.5,
            "unit": "ns/op",
            "extra": "35265 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "35265 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "35265 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 402.75,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "289371 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 402.75,
            "unit": "ns/op",
            "extra": "289371 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "289371 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "289371 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44687,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "2713 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44687,
            "unit": "ns/op",
            "extra": "2713 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "2713 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "2713 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 47949,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "2594 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 47949,
            "unit": "ns/op",
            "extra": "2594 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "2594 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "2594 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9.159,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "13042626 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9.159,
            "unit": "ns/op",
            "extra": "13042626 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13042626 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13042626 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 33636,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "3720 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 33636,
            "unit": "ns/op",
            "extra": "3720 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "3720 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3720 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3676.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "32380 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3676.5,
            "unit": "ns/op",
            "extra": "32380 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "32380 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "32380 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 468417,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 468417,
            "unit": "ns/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "252 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 482578.5,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 482578.5,
            "unit": "ns/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 474.15,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "265292 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 474.15,
            "unit": "ns/op",
            "extra": "265292 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "265292 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "265292 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3018,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "39751 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3018,
            "unit": "ns/op",
            "extra": "39751 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "39751 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "39751 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 27815,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27815,
            "unit": "ns/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063,
            "unit": "B/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 54310.5,
            "unit": "ns/op\t18523 B/op\t178 allocs/op",
            "extra": "2290 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54310.5,
            "unit": "ns/op",
            "extra": "2290 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18523,
            "unit": "B/op",
            "extra": "2290 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "2290 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 197730.5,
            "unit": "ns/op\t74840 B/op\t633 allocs/op",
            "extra": "606 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 197730.5,
            "unit": "ns/op",
            "extra": "606 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 74840,
            "unit": "B/op",
            "extra": "606 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 633,
            "unit": "allocs/op",
            "extra": "606 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 584329,
            "unit": "ns/op\t478060 B/op\t1417 allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 584329,
            "unit": "ns/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478060,
            "unit": "B/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1417,
            "unit": "allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 689.7,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "171273 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 689.7,
            "unit": "ns/op",
            "extra": "171273 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "171273 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "171273 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4095,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "28958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4095,
            "unit": "ns/op",
            "extra": "28958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "28958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "28958 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4357.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "27367 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4357.5,
            "unit": "ns/op",
            "extra": "27367 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "27367 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "27367 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2334,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "51648 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2334,
            "unit": "ns/op",
            "extra": "51648 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "51648 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "51648 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 683,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "174816 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 683,
            "unit": "ns/op",
            "extra": "174816 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "174816 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "174816 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4072,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "29137 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4072,
            "unit": "ns/op",
            "extra": "29137 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "29137 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "29137 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4392,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "26971 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4392,
            "unit": "ns/op",
            "extra": "26971 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "26971 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "26971 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2349.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "50709 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2349.5,
            "unit": "ns/op",
            "extra": "50709 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "50709 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "50709 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6621.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "18210 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6621.5,
            "unit": "ns/op",
            "extra": "18210 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "18210 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "18210 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 36171.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "3004 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 36171.5,
            "unit": "ns/op",
            "extra": "3004 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "3004 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3004 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 38645.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3109 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 38645.5,
            "unit": "ns/op",
            "extra": "3109 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3109 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3109 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 21729.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "5433 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 21729.5,
            "unit": "ns/op",
            "extra": "5433 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "5433 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "5433 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 6517.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "18042 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 6517.5,
            "unit": "ns/op",
            "extra": "18042 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "18042 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "18042 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 36300.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "3225 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 36300.5,
            "unit": "ns/op",
            "extra": "3225 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "3225 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3225 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 38610,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 38610,
            "unit": "ns/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 21952,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "5500 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 21952,
            "unit": "ns/op",
            "extra": "5500 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "5500 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "5500 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 122213,
            "unit": "ns/op\t278530 B/op\t1 allocs/op",
            "extra": "921 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 122213,
            "unit": "ns/op",
            "extra": "921 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278530,
            "unit": "B/op",
            "extra": "921 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "921 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 403325,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "294 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 403325,
            "unit": "ns/op",
            "extra": "294 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "294 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "294 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 429449,
            "unit": "ns/op\t455904.5 B/op\t922 allocs/op",
            "extra": "276 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 429449,
            "unit": "ns/op",
            "extra": "276 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904.5,
            "unit": "B/op",
            "extra": "276 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "276 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 223719,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 223719,
            "unit": "ns/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 128821.5,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 128821.5,
            "unit": "ns/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
            "unit": "B/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 443886.5,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "262 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 443886.5,
            "unit": "ns/op",
            "extra": "262 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
            "unit": "B/op",
            "extra": "262 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "262 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 444268,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 444268,
            "unit": "ns/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
            "unit": "B/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "266 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 222827.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "535 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 222827.5,
            "unit": "ns/op",
            "extra": "535 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "535 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "535 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 417.05,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "281035 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 417.05,
            "unit": "ns/op",
            "extra": "281035 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "281035 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "281035 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 408.85,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "274899 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 408.85,
            "unit": "ns/op",
            "extra": "274899 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "274899 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "274899 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3189.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "38030 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3189.5,
            "unit": "ns/op",
            "extra": "38030 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "38030 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "38030 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3185.5,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "37668 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3185.5,
            "unit": "ns/op",
            "extra": "37668 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "37668 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "37668 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51540,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2313 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51540,
            "unit": "ns/op",
            "extra": "2313 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2313 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2313 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51942.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "2329 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51942.5,
            "unit": "ns/op",
            "extra": "2329 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "2329 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2329 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 936.55,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "119833 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 936.55,
            "unit": "ns/op",
            "extra": "119833 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "119833 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "119833 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1464,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "81529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1464,
            "unit": "ns/op",
            "extra": "81529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "81529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "81529 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1393,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "86252 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1393,
            "unit": "ns/op",
            "extra": "86252 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "86252 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "86252 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 7639.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "15673 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 7639.5,
            "unit": "ns/op",
            "extra": "15673 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "15673 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "15673 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 13668.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 13668.5,
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
            "value": 13528,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 13528,
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
            "value": 83958,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "1450 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83958,
            "unit": "ns/op",
            "extra": "1450 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "1450 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1450 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 159203.5,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "740 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 159203.5,
            "unit": "ns/op",
            "extra": "740 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "740 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "740 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 157785,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 157785,
            "unit": "ns/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 9344,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "12830 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 9344,
            "unit": "ns/op",
            "extra": "12830 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12830 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12830 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 9148.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "13100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 9148.5,
            "unit": "ns/op",
            "extra": "13100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 92809.5,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1287 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 92809.5,
            "unit": "ns/op",
            "extra": "1287 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1287 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1287 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 93941,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 93941,
            "unit": "ns/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 952527.5,
            "unit": "ns/op\t10 B/op\t0 allocs/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 952527.5,
            "unit": "ns/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 10,
            "unit": "B/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 954068,
            "unit": "ns/op\t9.5 B/op\t0 allocs/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 954068,
            "unit": "ns/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9.5,
            "unit": "B/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 796.15,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "141712 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 796.15,
            "unit": "ns/op",
            "extra": "141712 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "141712 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "141712 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 805.65,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "149892 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 805.65,
            "unit": "ns/op",
            "extra": "149892 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "149892 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "149892 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 835.1,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "145874 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 835.1,
            "unit": "ns/op",
            "extra": "145874 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "145874 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "145874 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4763.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "25248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4763.5,
            "unit": "ns/op",
            "extra": "25248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "25248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "25248 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4849,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "24646 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4849,
            "unit": "ns/op",
            "extra": "24646 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "24646 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "24646 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 5038.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "23679 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5038.5,
            "unit": "ns/op",
            "extra": "23679 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "23679 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "23679 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 56444.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2284 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 56444.5,
            "unit": "ns/op",
            "extra": "2284 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2284 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2284 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 56709.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2244 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 56709.5,
            "unit": "ns/op",
            "extra": "2244 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2244 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2244 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 58976,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "2048 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 58976,
            "unit": "ns/op",
            "extra": "2048 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "2048 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2048 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 22838.5,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "5314 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 22838.5,
            "unit": "ns/op",
            "extra": "5314 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "5314 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "5314 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 54967.5,
            "unit": "ns/op\t6041 B/op\t152 allocs/op",
            "extra": "2206 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54967.5,
            "unit": "ns/op",
            "extra": "2206 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6041,
            "unit": "B/op",
            "extra": "2206 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "2206 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 54565.5,
            "unit": "ns/op\t4520.5 B/op\t151 allocs/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54565.5,
            "unit": "ns/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4520.5,
            "unit": "B/op",
            "extra": "2203 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "2203 times\n4 procs"
          }
        ]
      }
    ]
  }
}