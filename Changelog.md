# 0.7.0 / 2024-09-09

- make pagination consistent
- add resolved stories to response

# 0.6.0 / 2024-09-08

- Breaking: `Search` and `SearchRecent` now return search results wrapper.
  Use `result.Stories()` to get the stories of the search result.

# 0.5.0 / 2024-09-08

- support multiple pages queries
- sometimes numcomments is missing, so change the value to `*int`
- don't set `hitsPerPages` by default (default results seems to be 10)

# 0.4.0 / 2024-09-02

- Fix search queries. Closes #1
