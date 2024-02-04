# OpenTradingPlatform

OpenTradingPlatform (OTP) is a platform that was created with the intent of simplifying and automating the
processes that are most crucial when developing and deploying trading strategies.

The platform, once started handles all the processes of fetching and distributing data.

**Clients to access the data:**
* [Python Client](https://github.com/raulfrk/OpenTradingPlatformPythonClient)

**To get started in with the platform follow this guideline:**  [Getting started](./docs/getting_started.md)

The current version of the platform provides the following functionalities:

- Fetching market data from a broker
  - Supported asset classes
    - Stock
      - Data types:
        - Bar
        - Trades
        - Quotes
    - Crypto
      - Data types:
        - Bar
        - Trades
        - Quotes
    - News (non-tradeable)
      - Data types:
        - News headlines
- Subscribing to market data feeds provided by a broker
  - Stock
    - Data types:
      - Bar
      - Daily-bars
      - Quotes
      - Trades
      - Updated-bars
      - LULD
      - Trading status
  - Crypto
    - Data types:
      - Bar
      - Orderbook
      - Daily-bars
      - Quotes
      - Trades
      - Updated-bars
  - News (non-tradeable)
    - Data types:
      - News headlines
- Storing all market data from both stream subscriptions and data requests in a postgres database
- Performing sentiment analysis on news headlines with customized system prompt
  - Supported methods
    - Plain sentiment analysis
    - Aspect based sentiment analysis
  
## Planned future additions (coming soon)

* Pass-through sentiment analysis of news sentiment
  * News that reaches the platform is automatically analyzed and resulting news with sentiments are re-distributed

