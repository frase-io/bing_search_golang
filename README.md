# bing_search_golang


For now this only supports the basic web search.

To Use:
1. Copy your apiKey into the apiKey var at the top of the file.
2. Call it in your code by doing
    ```
    bingQuery := "machine learning"
    bingResult := BingWebSearchResult{}
    bingResult.MakeBingRequest(bingQuery)
    ```
  
  the request updates the struct you called it on.
