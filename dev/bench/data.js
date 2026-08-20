window.BENCHMARK_DATA = {
  "lastUpdate": 1787265218369,
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
          "id": "8ade5aed17f6b77b76c0b537a649274c2066a865",
          "message": "perf: memoize wildcard pattern expansion in FindConfigWildcard\n\nFindConfigWildcard expanded every configured [[wildcard]] pattern on\nevery call. The picker resolves an icon per session, so a config with W\npatterns cost O(N x W) home.ExpandPath calls per list — 8 allocs per\nsession with 8 patterns, scaling linearly in both directions.\n\nPatterns come from config and cannot change during a run, so expand them\nonce per lister behind a sync.Once and scan the cached results. Expansion\nstays lazy so commands that never resolve a wildcard pay nothing.\nUnexpandable patterns are still dropped and config order is preserved, so\nfirst-match-wins behavior is identical.\n\nBenchmarkIconResolverWildcard/n=1000/miss:\n  779112 ns/op  384000 B/op  8000 allocs/op\n  443054 ns/op       0 B/op     0 allocs/op\n\nCloses #442",
          "timestamp": "2026-08-20T16:34:22-05:00",
          "tree_id": "488a6faaf129296a9fec1ae87bee1d9939dd1975",
          "url": "https://github.com/joshmedeski/sesh/commit/8ade5aed17f6b77b76c0b537a649274c2066a865"
        },
        "date": 1787261700670,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 20254,
            "unit": "ns/op\t6651.5 B/op\t73 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20254,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6651.5,
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
            "value": 63782.5,
            "unit": "ns/op\t58371.5 B/op\t447 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 63782.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58371.5,
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
            "value": 442782.5,
            "unit": "ns/op\t538052.5 B/op\t3939 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 442782.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538052.5,
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
            "value": 27230.5,
            "unit": "ns/op\t8481.5 B/op\t87 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 27230.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8481.5,
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
            "value": 123879.5,
            "unit": "ns/op\t79307.5 B/op\t489 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 123879.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79307.5,
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
            "value": 992394,
            "unit": "ns/op\t746395 B/op\t4033 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 992394,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746395,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4033,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44195.5,
            "unit": "ns/op\t18862.5 B/op\t206 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44195.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18862.5,
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
            "value": 127558,
            "unit": "ns/op\t93470.5 B/op\t664 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 127558,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 93470.5,
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
            "value": 868426,
            "unit": "ns/op\t779939 B/op\t4948 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 868426,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 779939,
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
            "value": 21117,
            "unit": "ns/op\t6666 B/op\t74 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 21117,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6666,
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
            "value": 75058,
            "unit": "ns/op\t58328.5 B/op\t448 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 75058,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58328.5,
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
            "value": 480109.5,
            "unit": "ns/op\t538083 B/op\t3940 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 480109.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538083,
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
            "value": 3785.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3785.5,
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
            "value": 44843.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44843.5,
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
            "value": 481992,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 481992,
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
            "value": 11038.5,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 11038.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
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
            "value": 29462,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 29462,
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
            "value": 219476,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 219476,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
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
            "value": 16.435,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 16.435,
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
            "value": 562,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 562,
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
            "value": 251.9,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 251.9,
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
            "value": 3549.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3549.5,
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
            "value": 3930.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3930.5,
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
            "value": 14.78,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 14.78,
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
            "value": 3209,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3209,
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
            "value": 492.2,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 492.2,
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
            "value": 47770.5,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 47770.5,
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
            "value": 48319,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 48319,
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
            "value": 13.775,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 13.775,
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
            "value": 31488,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 31488,
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
            "value": 1989,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 1989,
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
            "value": 497732,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 497732,
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
            "value": 497408,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 497408,
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
            "value": 702.3,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 702.3,
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
            "value": 2934,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2934,
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
            "value": 31653.5,
            "unit": "ns/op\t33063 B/op\t321 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 31653.5,
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
            "value": 393412.5,
            "unit": "ns/op\t22960.5 B/op\t504 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 393412.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 22960.5,
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
            "value": 1529229,
            "unit": "ns/op\t91275.5 B/op\t1936 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1529229,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 91275.5,
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
            "value": 1344727,
            "unit": "ns/op\t488808.5 B/op\t2244 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1344727,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 488808.5,
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
            "value": 462.6,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 462.6,
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
            "value": 4997,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4997,
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
            "value": 4246,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4246,
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
            "value": 3023,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3023,
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
            "value": 465.9,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 465.9,
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
            "value": 4206,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4206,
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
            "value": 4713.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4713.5,
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
            "value": 2935.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2935.5,
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
            "value": 3638,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3638,
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
            "value": 46316.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 46316.5,
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
            "value": 41753,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41753,
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
            "value": 26789.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 26789.5,
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
            "value": 3631,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3631,
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
            "value": 41706.5,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 41706.5,
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
            "value": 42881.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 42881.5,
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
            "value": 27611.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 27611.5,
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
            "value": 143542,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 143542,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
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
            "value": 464813.5,
            "unit": "ns/op\t441497.5 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 464813.5,
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
            "value": 464353,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 464353,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
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
            "value": 268428.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 268428.5,
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
            "value": 150769,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 150769,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
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
            "value": 465138,
            "unit": "ns/op\t441496 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 465138,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441496,
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
            "value": 464000.5,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 464000.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
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
            "value": 268466,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 268466,
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
            "value": 438.3,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 438.3,
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
            "value": 525.35,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 525.35,
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
            "value": 2462,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2462,
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
            "value": 2304,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2304,
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
            "value": 36718.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 36718.5,
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
            "value": 37083.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 37083.5,
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
            "value": 698.45,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 698.45,
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
            "value": 1255,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1255,
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
            "value": 1185.5,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1185.5,
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
            "value": 5020,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5020,
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
            "value": 11680,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11680,
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
            "value": 11160.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11160.5,
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
            "value": 81929,
            "unit": "ns/op\t245760.5 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 81929,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245760.5,
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
            "value": 152622,
            "unit": "ns/op\t304800.5 B/op\t2001 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 152622,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 304800.5,
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
            "value": 154581,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 154581,
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
            "value": 8735.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8735.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 8934,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8934,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 83009.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83009.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 83419.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83419.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 843523.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 843523.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 850161.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 850161.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 991.3,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 991.3,
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
            "value": 791.9,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 791.9,
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
            "value": 748.45,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 748.45,
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
            "value": 5112.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5112.5,
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
            "value": 4255,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4255,
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
            "value": 4133,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4133,
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
            "value": 54204,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54204,
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
            "value": 56367.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 56367.5,
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
            "value": 55204,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 55204,
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
            "value": 34824,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 34824,
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
            "value": 85287,
            "unit": "ns/op\t6040 B/op\t152 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 85287,
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
            "value": 84462.5,
            "unit": "ns/op\t4520 B/op\t151 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84462.5,
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
          "id": "c8114e2e6579205375ecd3815f626e44effd84ea",
          "message": "ci: move artifact actions off Node.js 20\n\nThe bench job's upload-artifact@v5 defaults to runs.using node20, so\nevery run is annotated: \"Node.js 20 is deprecated ... being forced to\nrun on Node.js 24\". v6 is the release that switched the default to\nnode24; go to v7, the current major. download-artifact gets the same\ntreatment at v8, the major that pairs with upload v7 -- the two are\nversioned one apart.\n\nNeither bump touches an input this repo uses (name, path,\nif-no-files-found). The new majors add optional direct upload/download\nof unzipped single files, move to ESM, and, on download v8, default\ndigest-mismatch to error instead of warn, which is the behaviour we\nwant. Both require runner 2.327.1 or newer; ubuntu-latest is well past\nthat.\n\nDependabot would have picked these up on its monthly github-actions\npass, but the deprecation is annotating runs now.",
          "timestamp": "2026-08-20T17:29:44-05:00",
          "tree_id": "109a73635db4796d7750a4c20132d3dd65283615",
          "url": "https://github.com/joshmedeski/sesh/commit/c8114e2e6579205375ecd3815f626e44effd84ea"
        },
        "date": 1787265016899,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 19349,
            "unit": "ns/op\t6695 B/op\t73 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 19349,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6695,
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
            "value": 57292,
            "unit": "ns/op\t58353.5 B/op\t447 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 57292,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58353.5,
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
            "value": 402009,
            "unit": "ns/op\t538067 B/op\t3939 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 402009,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538067,
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
            "value": 25223,
            "unit": "ns/op\t8480 B/op\t87 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 25223,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8480,
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
            "value": 113756.5,
            "unit": "ns/op\t79300 B/op\t489 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 113756.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79300,
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
            "value": 901004.5,
            "unit": "ns/op\t746376.5 B/op\t4033 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 901004.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746376.5,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4033,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 36878,
            "unit": "ns/op\t18861.5 B/op\t206 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 36878,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18861.5,
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
            "value": 106227,
            "unit": "ns/op\t93650.5 B/op\t664 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 106227,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 93650.5,
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
            "value": 776094.5,
            "unit": "ns/op\t781422 B/op\t4948 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 776094.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 781422,
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
            "value": 60022,
            "unit": "ns/op\t58329.5 B/op\t448 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 60022,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58329.5,
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
            "value": 433607,
            "unit": "ns/op\t538082.5 B/op\t3940 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 433607,
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
            "value": 3782.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3782.5,
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
            "value": 43375,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 43375,
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
            "value": 465717.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 465717.5,
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
            "value": 9950.5,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9950.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
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
            "value": 26821.5,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 26821.5,
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
            "value": 202267.5,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 202267.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
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
            "value": 15.475,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 15.475,
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
            "value": 468.2,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 468.2,
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
            "value": 190,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 190,
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
            "value": 4533,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4533,
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
            "value": 3993.5,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3993.5,
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
            "value": 17.275,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 17.275,
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
            "value": 3014,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3014,
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
            "value": 615.15,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 615.15,
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
            "value": 43997,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 43997,
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
            "value": 45895,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 45895,
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
            "value": 14.875,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 14.875,
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
            "value": 32085,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 32085,
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
            "value": 1968,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 1968,
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
            "value": 472363,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 472363,
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
            "value": 476624.5,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 476624.5,
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
            "value": 699.5,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 699.5,
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
            "value": 2724,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 2724,
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
            "value": 21128.5,
            "unit": "ns/op\t33056 B/op\t321 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 21128.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33056,
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
            "value": 416633,
            "unit": "ns/op\t22963 B/op\t504 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 416633,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 22963,
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
            "value": 1629211,
            "unit": "ns/op\t91276 B/op\t1936 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1629211,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 91276,
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
            "value": 1358528,
            "unit": "ns/op\t488798.5 B/op\t2244 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1358528,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 488798.5,
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
            "value": 455.95,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 455.95,
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
            "value": 4067,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4067,
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
            "value": 4141,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4141,
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
            "value": 2627,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2627,
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
            "value": 468.35,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 468.35,
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
            "value": 4073.5,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4073.5,
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
            "value": 4514.5,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4514.5,
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
            "value": 2615,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2615,
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
            "value": 3461,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3461,
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
            "value": 38723,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 38723,
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
            "value": 40518.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40518.5,
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
            "value": 22653,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 22653,
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
            "value": 3510,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3510,
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
            "value": 39690,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39690,
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
            "value": 42766,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 42766,
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
            "value": 23158.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 23158.5,
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
            "value": 127006.5,
            "unit": "ns/op\t278529.5 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 127006.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529.5,
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
            "value": 407989.5,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 407989.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
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
            "value": 434055,
            "unit": "ns/op\t455904.5 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 434055,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904.5,
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
            "value": 244436.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 244436.5,
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
            "value": 157587.5,
            "unit": "ns/op\t278529.5 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 157587.5,
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
            "value": 438053.5,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 438053.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
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
            "value": 457879.5,
            "unit": "ns/op\t455905 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 457879.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455905,
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
            "value": 243718.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 243718.5,
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
            "value": 472.2,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 472.2,
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
            "value": 356.95,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 356.95,
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
            "value": 2634,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2634,
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
            "value": 2656,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2656,
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
            "value": 39064.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 39064.5,
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
            "value": 37430.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 37430.5,
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
            "value": 622.3,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 622.3,
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
            "value": 1501.5,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1501.5,
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
            "value": 1095,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1095,
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
            "value": 5252.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5252.5,
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
            "value": 11763.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11763.5,
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
            "value": 11598,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11598,
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
            "value": 76923,
            "unit": "ns/op\t245761 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 76923,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 245761,
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
            "value": 171368,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 171368,
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
            "value": 160449,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 160449,
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
            "value": 11367.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11367.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 9885.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 9885.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 90976,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 90976,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 96806.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 96806.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 954476,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 954476,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 941406,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 941406,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 786.05,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 786.05,
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
            "value": 722.5,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 722.5,
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
            "value": 802.4,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 802.4,
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
            "value": 4541.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4541.5,
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
            "value": 4319,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4319,
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
            "value": 4901.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4901.5,
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
            "value": 52273,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52273,
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
            "value": 51843,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 51843,
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
            "value": 54142,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 54142,
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
            "value": 35298.5,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 35298.5,
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
            "value": 88642,
            "unit": "ns/op\t6040 B/op\t152 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 88642,
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
            "value": 88344,
            "unit": "ns/op\t4520 B/op\t151 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 88344,
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
          "id": "7fe1b193dfd3552243bfc15bbaece6e7ceedb079",
          "message": "perf: coalesce styled runs in highlightMatches\n\nhighlightMatches called lipgloss.Style.Render once per rune, so a\n40-character session name cost 40 Render calls per visible row on every\nkeystroke. Each call re-derived the style's border metrics even though\nthe match and normal styles have no border set, which a CPU profile of\nBenchmarkKeystroke showed as 84% of the time in Model.View.\n\nWalk the runes instead and Render once per run of same-styled runes. A\nquery typically produces a handful of runs per name rather than one per\ncharacter. The match set also becomes a []bool sized to the rune count\ninstead of a map[int]bool; out-of-range indexes are still ignored.\n\nBenchmarkKeystroke on an M4 Max:\n\n  n=10     256us/504 allocs -> 47us/178 allocs\n  n=100    958us/1936 allocs -> 147us/633 allocs\n  n=1000   727us/2245 allocs -> 312us/1418 allocs\n\nAt n=100 a keystroke now sits close to the ~47us floor of BenchmarkView\nwith an empty query, so highlighting is no longer the dominant cost.\n\nThe rendered bytes change while the output stays visually identical: one\nSGR pair wraps each run rather than each character. Grouping is safe\nbecause the two styles set only Foreground and Bold -- per-run and\nper-rune output diverge only for styles with padding, margin, border,\nwidth, or alignment.\n\nCloses #443",
          "timestamp": "2026-08-20T17:33:07-05:00",
          "tree_id": "cc434d5ffea7cf98b16f1b17f1e224f12a87aecd",
          "url": "https://github.com/joshmedeski/sesh/commit/7fe1b193dfd3552243bfc15bbaece6e7ceedb079"
        },
        "date": 1787265217870,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 20914.5,
            "unit": "ns/op\t6674.5 B/op\t73 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20914.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 6674.5,
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
            "value": 68388.5,
            "unit": "ns/op\t58388.5 B/op\t447 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 68388.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58388.5,
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
            "value": 430134.5,
            "unit": "ns/op\t538048.5 B/op\t3939 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 430134.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkList/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538048.5,
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
            "value": 26816,
            "unit": "ns/op\t8482 B/op\t87 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 26816,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 8482,
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
            "value": 118988.5,
            "unit": "ns/op\t79302.5 B/op\t489 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 118988.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 79302.5,
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
            "value": 946252,
            "unit": "ns/op\t746373 B/op\t4033 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 946252,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 746373,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHideDuplicates/n=1000 (github.com/joshmedeski/sesh/v2/lister) - allocs/op",
            "value": 4033,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister)",
            "value": 44432.5,
            "unit": "ns/op\t18861.5 B/op\t206 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 44432.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 18861.5,
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
            "value": 126878.5,
            "unit": "ns/op\t92910 B/op\t664 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 126878.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 92910,
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
            "value": 845593.5,
            "unit": "ns/op\t779565 B/op\t4948 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 845593.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListBlacklist/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 779565,
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
            "value": 20588,
            "unit": "ns/op\t6664 B/op\t74 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 20588,
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
            "value": 71380.5,
            "unit": "ns/op\t58326 B/op\t448 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 71380.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=100 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 58326,
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
            "value": 469150.5,
            "unit": "ns/op\t538087.5 B/op\t3940 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 469150.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkListShowWindows/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 538087.5,
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
            "value": 3382.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3382.5,
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
            "value": 42088,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 42088,
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
            "value": 477232.5,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkApplyDedup/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 477232.5,
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
            "value": 9774.5,
            "unit": "ns/op\t9772.5 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 9774.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=10 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9772.5,
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
            "value": 51030,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 51030,
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
            "value": 230906,
            "unit": "ns/op\t9958 B/op\t118 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 230906,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBlacklistFilter/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 9958,
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
            "value": 16.78,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 16.78,
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
            "value": 423.5,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 423.5,
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
            "value": 120.5,
            "unit": "ns/op\t160 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 120.5,
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
            "value": 3927.5,
            "unit": "ns/op\t1848 B/op\t14 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3927.5,
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
            "value": 4460,
            "unit": "ns/op\t2280 B/op\t16 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=10/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 4460,
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
            "value": 14.13,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 14.13,
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
            "value": 3224.5,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3224.5,
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
            "value": 340.05,
            "unit": "ns/op\t1792 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 340.05,
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
            "value": 46447,
            "unit": "ns/op\t21000 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 46447,
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
            "value": 45747.5,
            "unit": "ns/op\t26776 B/op\t50 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=100/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 45747.5,
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
            "value": 13.97,
            "unit": "ns/op\t0 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/passthrough (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 13.97,
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
            "value": 30758,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/source-filter (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 30758,
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
            "value": 1597.5,
            "unit": "ns/op\t16384 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-attached (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 1597.5,
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
            "value": 478524,
            "unit": "ns/op\t208264 B/op\t93 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/hide-duplicates (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 478524,
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
            "value": 496076.5,
            "unit": "ns/op\t245128 B/op\t95 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCachingListerApplyFilters/n=1000/picker-defaults (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 496076.5,
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
            "value": 598.55,
            "unit": "ns/op\t528 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=10 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 598.55,
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
            "value": 3488,
            "unit": "ns/op\t3408 B/op\t42 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=100 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 3488,
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
            "value": 24476,
            "unit": "ns/op\t33056 B/op\t321 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - ns/op",
            "value": 24476,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFindTmuxSessionByBase/n=1000 (github.com/joshmedeski/sesh/v2/lister) - B/op",
            "value": 33056,
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
            "value": 78763,
            "unit": "ns/op\t18520 B/op\t178 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 78763,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 18520,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=10 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 178,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 292647,
            "unit": "ns/op\t74837 B/op\t633 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 292647,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 74837,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=100 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 633,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker)",
            "value": 679952.5,
            "unit": "ns/op\t478056 B/op\t1417 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 679952.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 478056,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkKeystroke/n=1000 (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 1417,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 442.9,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 442.9,
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
            "value": 4547,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4547,
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
            "value": 4363,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4363,
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
            "value": 2916.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2916.5,
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
            "value": 548.45,
            "unit": "ns/op\t3072 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 548.45,
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
            "value": 4822.5,
            "unit": "ns/op\t5568 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4822.5,
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
            "value": 4962,
            "unit": "ns/op\t5720 B/op\t25 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4962,
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
            "value": 3011.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=10/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3011.5,
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
            "value": 3622.5,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3622.5,
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
            "value": 42642,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 42642,
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
            "value": 42441.5,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 42441.5,
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
            "value": 27225.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 27225.5,
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
            "value": 3746,
            "unit": "ns/op\t27264 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3746,
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
            "value": 42497,
            "unit": "ns/op\t47944 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 42497,
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
            "value": 43394,
            "unit": "ns/op\t49392 B/op\t109 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 43394,
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
            "value": 26503.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=100/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 26503.5,
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
            "value": 147364,
            "unit": "ns/op\t278530.5 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 147364,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278530.5,
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
            "value": 455956,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 455956,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
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
            "value": 460354,
            "unit": "ns/op\t455904 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 460354,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904,
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
            "value": 266985.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/plain/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 266985.5,
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
            "value": 101889.5,
            "unit": "ns/op\t278529 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 101889.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/empty (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 278529,
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
            "value": 475736,
            "unit": "ns/op\t441497 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 475736,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/1-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 441497,
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
            "value": 487731,
            "unit": "ns/op\t455904.5 B/op\t922 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 487731,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/3-char (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 455904.5,
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
            "value": 268472.5,
            "unit": "ns/op\t248 B/op\t8 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterSessions/n=1000/separator-aware/no-match (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 268472.5,
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
            "value": 661,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 661,
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
            "value": 519.2,
            "unit": "ns/op\t776 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=10/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 519.2,
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
            "value": 2407,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2407,
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
            "value": 2343,
            "unit": "ns/op\t6280 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=100/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 2343,
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
            "value": 40265,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 40265,
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
            "value": 37320.5,
            "unit": "ns/op\t65672 B/op\t4 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkRankMatches/n=1000/app (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 37320.5,
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
            "value": 682.8,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 682.8,
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
            "value": 1178.5,
            "unit": "ns/op\t3264 B/op\t21 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1178.5,
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
            "value": 1043,
            "unit": "ns/op\t2688 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=10/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 1043,
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
            "value": 5072.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 5072.5,
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
            "value": 12104.5,
            "unit": "ns/op\t30336 B/op\t201 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 12104.5,
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
            "value": 11276.5,
            "unit": "ns/op\t24576 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=100/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 11276.5,
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
            "value": 78148.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/plain (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 78148.5,
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
            "value": 156836.5,
            "unit": "ns/op\t304800 B/op\t2001 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/separator-aware (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 156836.5,
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
            "value": 145875.5,
            "unit": "ns/op\t245760 B/op\t1 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkBuildItems/n=1000/icons (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 145875.5,
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
            "value": 8516.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 8516.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 9589.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 9589.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=10/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 83922.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83922.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 83712.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 83712.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=100/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker)",
            "value": 844061.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 844061.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/miss (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker)",
            "value": 849131.5,
            "unit": "ns/op\t12 B/op\t0 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 849131.5,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 12,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkIconResolverWildcard/n=1000/hit-last (github.com/joshmedeski/sesh/v2/picker) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker)",
            "value": 832.2,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 832.2,
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
            "value": 755.95,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 755.95,
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
            "value": 860.15,
            "unit": "ns/op\t1632 B/op\t3 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=10/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 860.15,
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
            "value": 4554.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4554.5,
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
            "value": 3972.5,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 3972.5,
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
            "value": 4089,
            "unit": "ns/op\t9464 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=100/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 4089,
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
            "value": 57159.5,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/empty (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 57159.5,
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
            "value": 52750,
            "unit": "ns/op\t87208.5 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52750,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/a (github.com/joshmedeski/sesh/v2/picker) - B/op",
            "value": 87208.5,
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
            "value": 52957,
            "unit": "ns/op\t87208 B/op\t7 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFilterAliases/n=1000/al7 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 52957,
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
            "value": 34929,
            "unit": "ns/op\t2544 B/op\t66 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=10 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 34929,
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
            "value": 84885.5,
            "unit": "ns/op\t6040 B/op\t152 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=100 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84885.5,
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
            "value": 84991.5,
            "unit": "ns/op\t4520 B/op\t151 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkView/n=1000 (github.com/joshmedeski/sesh/v2/picker) - ns/op",
            "value": 84991.5,
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