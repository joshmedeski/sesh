window.BENCHMARK_DATA = {
  "lastUpdate": 1787259303936,
  "repoUrl": "https://github.com/joshmedeski/sesh",
  "entries": {
    "sesh": [
      {
        "commit": {
          "author": {
            "email": "joshmedeski@users.noreply.github.com",
            "name": "Josh Medeski",
            "username": "joshmedeski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "680152f9294e9a1f6a8014603dc060295f1e4ac2",
          "message": "Add lister and picker benchmarks with CI tracking (#444)\n\n* Add lister and picker benchmarks with CI tracking\n\nCover the paths whose cost scales with the size of a user's session list\nrather than with a constant: listing and deduping N sessions, and\nfiltering them on every keystroke. Each runs at N = 10, 100, and 1000,\n1000 being an ordinary zoxide history and the size at which an accidental\nper-session shell-out or list copy becomes visible.\n\nBenchmarks use hand-written fakes rather than the generated mocks, which\ntake a lock and walk their expectations per call — at N = 1000 that would\nbe most of what the benchmark measured. They never run under -race for\nthe same reason, so `just test` drops -bench and a new `just bench`\nrecipe runs them on their own.\n\nCI reports on them two ways. On pull requests, a benchstat comparison\nagainst the merge base lands in the job summary; it is informational,\nsince GitHub's runners need watching before a gate that could cry wolf\ngets turned on. On pushes to main, a separate job appends to a\ngithub-action-benchmark chart so a regression can be traced to the commit\nthat caused it — benchstat only ever sees two commits and cannot answer\n\"when did this get slow\". That job is separate because it needs\ncontents: write, which should not be held by a job running fork PR code.\n\nCI runs benchstat via a go.mod tool directive, which requires go 1.26.0,\nso the setup-go pins move from 1.25 to 1.26.\n\n* Add a benchmark for rankMatches\n\nfilterSessions already covers ranking in situ, but only in aggregate: a\nregression there could come from the fuzzy match or from the parallel\nslice and stable sort that rank it, and the combined number cannot say\nwhich.\n\nThe match set is built once and copied into a scratch slice per\niteration, since rankMatches sorts in place. Reusing the slice would mean\nevery pass after the first sorted an already-sorted input — the case\nsort.SliceStable is fastest at and the one that never happens on a real\nkeystroke.",
          "timestamp": "2026-08-20T15:53:55-05:00",
          "tree_id": "4ca5a5acd50339206af736fb0c54b07bf1a821fd",
          "url": "https://github.com/joshmedeski/sesh/commit/680152f9294e9a1f6a8014603dc060295f1e4ac2"
        },
        "date": 1787259302358,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 20635.5,
            "unit": "ns/op\t6669 B/op\t73 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20635.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6669,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 73,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 64691,
            "unit": "ns/op\t58314 B/op\t447 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 64691,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58314,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 447,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 382442,
            "unit": "ns/op\t538058 B/op\t3939 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 382442,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538058,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3939,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 23310.5,
            "unit": "ns/op\t8480.5 B/op\t87 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 23310.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8480.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 87,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 118001,
            "unit": "ns/op\t79301 B/op\t489 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 118001,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79301,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 489,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 886237.5,
            "unit": "ns/op\t746366 B/op\t4032 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 886237.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746366,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4032,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 36822,
            "unit": "ns/op\t19047.5 B/op\t206 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 36822,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 19047.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 206,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 108534.5,
            "unit": "ns/op\t93278.5 B/op\t664 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 108534.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 93278.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 664,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 786009,
            "unit": "ns/op\t781039.5 B/op\t4948 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 786009,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 781039.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4948,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19645,
            "unit": "ns/op\t6664 B/op\t74 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19645,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6664,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 74,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 61084.5,
            "unit": "ns/op\t58332.5 B/op\t448 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 61084.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58332.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 448,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 425664.5,
            "unit": "ns/op\t538082.5 B/op\t3940 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 425664.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538082.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 3940,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3724.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3724.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 41475.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 41475.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 461937.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 461937.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 9863.5,
            "unit": "ns/op\t9586.5 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9863.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9586.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 27104.5,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27104.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 202544,
            "unit": "ns/op\t9772 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 202544,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9772,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 118,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 15.63,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 15.63,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 480.5,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 480.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 272.1,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 272.1,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 160,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 4264,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4264,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1848,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 4497.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4497.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 2280,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 14.62,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 14.62,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 3494.5,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3494.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 870.25,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 870.25,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 43559.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 43559.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 21000,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44381.5,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44381.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 26776,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 50,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister)",
            "value": 18.275,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 18.275,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister)",
            "value": 30164.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 30164.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister)",
            "value": 1326,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 1326,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 16384,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister)",
            "value": 463604.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 463604.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 208264,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 93,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister)",
            "value": 475022,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 475022,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 245128,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 676.4,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 676.4,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 2914.5,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2914.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 3408,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 42,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 22568,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 22568,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33063,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 417479,
            "unit": "ns/op\t22961 B/op\t504 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 417479,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 22961,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 504,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1629682.5,
            "unit": "ns/op\t91278.5 B/op\t1936 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1629682.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 91278.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1936,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1350087,
            "unit": "ns/op\t488797 B/op\t2244 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1350087,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 488797,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2244,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 584.25,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 584.25,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4133.5,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4133.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4010.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4010.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2493,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2493,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 439.25,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 439.25,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3072,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3936,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3936,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5568,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4368,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4368,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 5720,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2819,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2819,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3634,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3634,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 46253.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 46253.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39796.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39796.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 22793.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 22793.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3652.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3652.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 27264,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 38509,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 38509,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 47944,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 39681.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39681.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 49392,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 109,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 22693.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 22693.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 124175,
            "unit": "ns/op\t278532 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 124175,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278532,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 415417.5,
            "unit": "ns/op\t441497.5 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 415417.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 436334.5,
            "unit": "ns/op\t455905 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 436334.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455905,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 232069,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 232069,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 162878.5,
            "unit": "ns/op\t278529.5 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 162878.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 524922.5,
            "unit": "ns/op\t441499 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 524922.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441499,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker)",
            "value": 444078.5,
            "unit": "ns/op\t455905.5 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 444078.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455905.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 922,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker)",
            "value": 231497.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 231497.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 401.1,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 401.1,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 489.75,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 489.75,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 776,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2640,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2640,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 2545,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2545,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6280,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 36430,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 36430,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker)",
            "value": 36735.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 36735.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 65672,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 648.15,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 648.15,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1189,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1189,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 3264,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1325,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1325,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2688,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 5117.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5117.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 12457,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12457,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 30336,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 201,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 11726.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11726.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker)",
            "value": 86307.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 86307.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker)",
            "value": 149342,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 149342,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 2001,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker)",
            "value": 143883,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 143883,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 16863.5,
            "unit": "ns/op\t2560 B/op\t80 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 16863.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2560,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 80,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 15999.5,
            "unit": "ns/op\t2240 B/op\t70 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 15999.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2240,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 163341,
            "unit": "ns/op\t25600 B/op\t800 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 163341,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 25600,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 800,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 159236.5,
            "unit": "ns/op\t22400 B/op\t700 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 159236.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 22400,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 700,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1685304.5,
            "unit": "ns/op\t256000 B/op\t8000 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1685304.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 256000,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 8000,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1643100.5,
            "unit": "ns/op\t224000 B/op\t7000 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1643100.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 224000,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7000,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 1099.25,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1099.25,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 729.4,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 729.4,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 823.25,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 823.25,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4057,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4057,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 3965,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3965,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 4402.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4402.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 9464,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51413.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51413.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker)",
            "value": 51925.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51925.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 52338.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52338.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 35628,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 35628,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 2544,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 89564,
            "unit": "ns/op\t6040 B/op\t152 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 89564,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 6040,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 152,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 88140.5,
            "unit": "ns/op\t4520 B/op\t151 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 88140.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 4520,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 151,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          }
        ]
      }
    ]
  }
}