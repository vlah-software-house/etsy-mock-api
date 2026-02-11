# Etsy API Mock Server

> **Disclaimer:** This is **not** an official Etsy product and is not affiliated with, endorsed by, or sponsored by Etsy, Inc. This is an independent, open-source utility built to help developers test their Etsy integrations locally — without making calls to the real Etsy API. If you're building an app that uses the [Etsy Open API v3](https://developers.etsy.com/documentation/), this mock server lets you run automated E2E tests, QA workflows, and development loops entirely offline.

A full mock implementation of the [Etsy Open API v3](https://developers.etsy.com/documentation/) in Go. Designed for E2E tests and human quality testing without touching the real Etsy API.

## Quick Start

```bash
go run ./cmd/server
```

Options:
- `-port 8080` — Server port (default: 8080)
- `-no-auth` — Disable API key / OAuth token validation
- `-no-seed` — Start with an empty data store (no sample data)

## Authentication

Matches the real Etsy API authentication model:

### API Key (required for all endpoints)

Every request must include an `x-api-key` header in the format `keystring:shared_secret`:

```bash
curl -H "x-api-key: myapp:mysecret" http://localhost:8080/v3/application/shops/5001
```

### OAuth2 (required for protected endpoints)

Endpoints that access private data require an OAuth2 Bearer token in addition to the API key. The mock implements the full OAuth2 Authorization Code + PKCE flow:

**Token exchange:**
```bash
curl -X POST http://localhost:8080/v3/public/oauth/token \
  -d "grant_type=authorization_code&client_id=myapp&code=mock_code_1001&code_verifier=test&redirect_uri=http://localhost"
```

**Token refresh:**
```bash
curl -X POST http://localhost:8080/v3/public/oauth/token \
  -d "grant_type=refresh_token&client_id=myapp&refresh_token=<token>"
```

**Using a token:**
```bash
curl -H "x-api-key: myapp:mysecret" \
     -H "Authorization: Bearer test-token-alice" \
     http://localhost:8080/v3/application/shops/5001/receipts
```

**Pre-seeded tokens** (have all scopes):
- `test-token-alice` — user 1001
- `test-token-bob` — user 1002

### OAuth2 Scopes

The mock enforces per-endpoint scope requirements matching the real API:
- Public endpoints (taxonomy, active listings, shop info) need only `x-api-key`
- `listings_r` / `listings_w` / `listings_d` for listing read/write/delete
- `shops_r` / `shops_w` for shop management, shipping profiles
- `transactions_r` / `transactions_w` for receipts, payments, ledger
- `email_r` for user profiles
- `address_r` for user addresses

### Rate Limiting

Every response includes Etsy-style rate limit headers:
- `x-limit-per-second` / `x-remaining-this-second`
- `x-limit-per-day` / `x-remaining-today`

Use `-no-auth` to disable all authentication checks for easier testing.

## Base URL

```
http://localhost:8080/v3/application
```

## Seed Data

The server starts pre-loaded with realistic test data:

| Entity | IDs | Description |
|--------|-----|-------------|
| Users | 1001-1004 | Alice, Bob, Carol, Dave |
| Shops | 5001 (AliceCrafts), 5002 (BobsWoodworks) | Two shops with full policies |
| Listings | 7001-7005 (jewelry), 7010-7012 (woodwork) | Active, draft, digital listings |
| Receipts | 9001-9003 | Completed, paid, with transactions |
| Payments | 11001-11003 | Settled and open payments |
| Shipping Profiles | 3001-3002 | With destinations and upgrades |
| Reviews | Per shop | 5-star and 4-star reviews |
| Taxonomy | 5 top-level categories | Jewelry, Home & Living, etc. |

## Implemented Endpoints

### OAuth2 & Auth
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| POST | `/v3/public/oauth/token` | — | Token exchange & refresh (PKCE) |
| POST | `/v3/application/scopes` | OAuth | Check token scopes |
| GET | `/v3/application/users/me` | OAuth | Get authenticated user |
| GET | `/v3/application/openapi-ping` | — | API connectivity check |

### Listings
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/listings/active` | api_key | Search all active listings |
| GET | `/v3/application/listings/batch?listing_ids=...` | api_key | Get multiple listings |
| GET | `/v3/application/listings/{id}` | api_key | Get a listing |
| DELETE | `/v3/application/listings/{id}` | listings_d | Delete a listing |
| GET | `/v3/application/listings/{id}/inventory` | api_key | Get listing inventory |
| GET | `/v3/application/listings/{id}/reviews` | api_key | Get listing reviews |
| GET | `/v3/application/listings/{id}/personalization` | api_key | Get personalization settings |
| GET | `/v3/application/listings/{id}/videos` | api_key | List videos |
| GET | `/v3/application/listings/{id}/videos/{vid}` | api_key | Get video |
| GET | `/v3/application/listings/{id}/images` | api_key | List images |
| GET | `/v3/application/listings/{id}/products/{pid}/offerings/{oid}` | api_key | Get offering |
| POST | `/v3/application/shops/{sid}/listings` | listings_w | Create a draft listing |
| GET | `/v3/application/shops/{sid}/listings` | listings_r | List shop listings |
| GET | `/v3/application/shops/{sid}/listings/active` | api_key | List active shop listings |
| GET | `/v3/application/shops/{sid}/listings/featured` | api_key | List featured listings |
| PUT/PATCH | `/v3/application/shops/{sid}/listings/{id}` | listings_w | Update a listing |

### Listing Sub-resources (under shops)
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET/POST | `.../listings/{id}/images` | api_key/listings_w | List/upload images |
| DELETE | `.../listings/{id}/images/{iid}` | listings_w | Delete image |
| GET/POST | `.../listings/{id}/files` | listings_r/listings_w | List/upload files |
| GET/DELETE | `.../listings/{id}/files/{fid}` | listings_r/listings_w | Get/delete file |
| GET/PUT | `.../listings/{id}/personalization` | api_key/listings_w | Get/update personalization |
| GET/POST | `.../listings/{id}/videos` | api_key/listings_w | List/upload videos |
| GET/DELETE | `.../listings/{id}/videos/{vid}` | api_key/listings_w | Get/delete video |
| GET/PUT | `.../listings/{id}/translations/{lang}` | api_key/listings_w | Get/update translation |
| GET/POST | `.../listings/{id}/variation-images` | api_key/listings_w | Get/update variation images |
| GET | `.../listings/{id}/properties` | api_key | List properties |
| GET | `.../listings/{id}/inventory` | api_key | Get inventory |

### Shops
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}` | api_key | Get shop |
| GET | `/v3/application/shops?shop_name=xxx` | api_key | Find shop by name |
| PUT | `/v3/application/shops/{id}` | shops_w | Update shop |
| GET | `/v3/application/shops/{id}/sections` | api_key | List sections |
| GET | `/v3/application/shops/{id}/sections/{sid}` | api_key | Get section |
| POST | `/v3/application/shops/{id}/sections` | shops_w | Create section |
| GET | `/v3/application/shops/{id}/shop-sections/listings` | api_key | Listings by section |
| GET/PUT | `/v3/application/shops/{id}/holiday-preferences` | shops_r/shops_w | Vacation settings |
| GET | `/v3/application/shops/{id}/production-partners` | shops_r | Production partners |
| GET | `/v3/application/shops/{id}/readiness-state-definitions` | shops_r | Processing profiles |

### Return Policies
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}/return-policies` | api_key | List return policies |
| GET | `/v3/application/shops/{id}/return-policies/{pid}` | api_key | Get return policy |
| POST | `/v3/application/shops/{id}/return-policies` | shops_w | Create return policy |

### Receipts & Transactions
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}/receipts` | transactions_r | List receipts |
| GET | `/v3/application/shops/{id}/receipts/{rid}` | transactions_r | Get receipt |
| PUT | `/v3/application/shops/{id}/receipts/{rid}` | transactions_w | Update receipt |
| POST | `/v3/application/shops/{id}/receipts/{rid}/tracking` | transactions_w | Add tracking |
| GET | `/v3/application/shops/{id}/receipts/{rid}/transactions` | transactions_r | Receipt transactions |
| GET | `/v3/application/shops/{id}/receipts/{rid}/payments` | transactions_r | Receipt payments |
| GET | `/v3/application/shops/{id}/transactions` | transactions_r | List shop transactions |
| GET | `/v3/application/shops/{id}/transactions/{tid}` | transactions_r | Get transaction |

### Payments & Ledger
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}/payments` | transactions_r | List shop payments |
| GET | `/v3/application/shops/{id}/payment-account/ledger-entries` | transactions_r | Ledger entries |

### Users
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/users/{id}` | email_r | Get user |
| GET | `/v3/application/users/{id}/addresses` | address_r | User addresses |
| DELETE | `/v3/application/users/{id}/addresses/{aid}` | address_r | Delete address |
| GET | `/v3/application/users/{id}/shops` | api_key | User's shops |

### Shipping
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}/shipping-profiles` | shops_r | List profiles |
| GET | `/v3/application/shops/{id}/shipping-profiles/{pid}` | shops_r | Get profile |
| POST | `/v3/application/shops/{id}/shipping-profiles` | shops_w | Create profile |
| DELETE | `/v3/application/shops/{id}/shipping-profiles/{pid}` | shops_w | Delete profile |
| GET | `/v3/application/shipping-carriers` | api_key | List shipping carriers |

### Reviews
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/shops/{id}/reviews` | api_key | Shop reviews |
| GET | `/v3/application/listings/{id}/reviews` | api_key | Listing reviews |

### Taxonomy
| Method | Path | Scope | Description |
|--------|------|-------|-------------|
| GET | `/v3/application/buyer-taxonomy/nodes` | api_key | Buyer taxonomy tree |
| GET | `/v3/application/buyer-taxonomy/nodes/{id}/properties` | api_key | Buyer taxonomy properties |
| GET | `/v3/application/seller-taxonomy/nodes` | api_key | Seller taxonomy tree |
| GET | `/v3/application/seller-taxonomy/nodes/{id}/properties` | api_key | Seller taxonomy properties |

### Utility
| Method | Path | Description |
|--------|------|-------------|
| GET | `/ping` | Health check |
| POST | `/admin/reset` | Reset data store |

## Query Parameters

Most list endpoints support:
- `limit` — Results per page (default: 25, max: 100)
- `offset` — Pagination offset (default: 0)

Active listings search also supports:
- `keywords` — Full-text search across title, description, and tags
- `taxonomy_id` — Filter by taxonomy
- `sort_on` — Sort field: `created`, `price`, `updated`, `score`
- `sort_order` — `asc` or `desc`

Shop listings support:
- `state` — Filter by state: `active`, `inactive`, `sold_out`, `draft`, `expired`

## Example Usage

```bash
# Public endpoint (API key only)
curl -H "x-api-key: myapp:mysecret" \
  "http://localhost:8080/v3/application/listings/active?keywords=necklace"

# Create a listing (requires OAuth + listings_w scope)
curl -X POST \
  -H "x-api-key: myapp:mysecret" \
  -H "Authorization: Bearer test-token-alice" \
  -H "Content-Type: application/json" \
  http://localhost:8080/v3/application/shops/5001/listings \
  -d '{"title":"New Item","description":"A thing","quantity":5,"price":19.99,"who_made":"i_did","when_made":"2020_2026","taxonomy_id":1207}'

# Get shop receipts (requires OAuth + transactions_r scope)
curl -H "x-api-key: myapp:mysecret" \
  -H "Authorization: Bearer test-token-alice" \
  "http://localhost:8080/v3/application/shops/5001/receipts"

# OAuth2 token exchange (PKCE flow)
curl -X POST http://localhost:8080/v3/public/oauth/token \
  -d "grant_type=authorization_code&client_id=myapp&code=mock_code_1001&code_verifier=test&redirect_uri=http://localhost"

# Add shipment tracking
curl -X POST \
  -H "x-api-key: myapp:mysecret" \
  -H "Authorization: Bearer test-token-alice" \
  -H "Content-Type: application/json" \
  http://localhost:8080/v3/application/shops/5001/receipts/9002/tracking \
  -d '{"carrier_name":"USPS","tracking_code":"9400111899223100099999"}'
```

## Architecture

```
cmd/server/main.go          — Entry point, flag parsing, middleware chain
internal/
  models/                   — All Etsy API data types (1:1 with OpenAPI spec)
    listing.go              — Listings, images, videos, files, inventory
    shop.go                 — Shop, sections, return policies
    receipt.go              — Receipts, transactions, shipments, refunds
    payment.go              — Payments, adjustments, ledger entries
    user.go                 — Users, addresses
    review.go               — Reviews
    shipping.go             — Shipping profiles, destinations, upgrades
    taxonomy.go             — Buyer/seller taxonomy nodes and properties
    auth.go                 — OAuth2 token request/response types
    extras.go               — Personalization, translations, carriers, etc.
    money.go                — Money type (amount/divisor/currency)
    responses.go            — Paginated and error response wrappers
  store/store.go            — Thread-safe in-memory data store
  handlers/
    router.go               — URL routing (all 60+ endpoints)
    helpers.go              — JSON encoding, path parsing, scope checking
    oauth.go                — OAuth2 PKCE token exchange & refresh
    listings.go             — Listing CRUD + images, files, inventory
    extras.go               — Videos, personalization, translations, carriers, etc.
    shops.go                — Shop, sections, return policies
    receipts.go             — Receipts, transactions, tracking
    payments.go             — Payments, ledger entries
    users.go                — Users, addresses
    reviews.go              — Reviews
    shipping.go             — Shipping profiles
    taxonomy.go             — Buyer/seller taxonomy
  middleware/auth.go        — API key validation, OAuth2 token store,
                              scope enforcement, rate limit headers,
                              CORS, logging, content-type
  seed/seed.go              — Realistic test data (2 shops, 8 listings,
                              receipts, payments, reviews, taxonomy)
```
