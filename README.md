# Etsy API Mock Server

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

### Listings
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/listings/active` | Search all active listings |
| GET | `/v3/application/listings/{listing_id}` | Get a listing |
| DELETE | `/v3/application/listings/{listing_id}` | Delete a listing |
| GET | `/v3/application/listings/{listing_id}/inventory` | Get listing inventory |
| GET | `/v3/application/listings/{listing_id}/reviews` | Get listing reviews |
| POST | `/v3/application/shops/{shop_id}/listings` | Create a draft listing |
| GET | `/v3/application/shops/{shop_id}/listings` | List shop listings |
| GET | `/v3/application/shops/{shop_id}/listings/active` | List active shop listings |
| PUT/PATCH | `/v3/application/shops/{shop_id}/listings/{listing_id}` | Update a listing |

### Listing Images
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/listings/{listing_id}/images` | List images |
| POST | `/v3/application/shops/{shop_id}/listings/{listing_id}/images` | Upload image |
| DELETE | `/v3/application/shops/{shop_id}/listings/{listing_id}/images/{image_id}` | Delete image |

### Listing Files (Digital)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/listings/{listing_id}/files` | List files |
| GET | `/v3/application/shops/{shop_id}/listings/{listing_id}/files/{file_id}` | Get file |
| POST | `/v3/application/shops/{shop_id}/listings/{listing_id}/files` | Upload file |
| DELETE | `/v3/application/shops/{shop_id}/listings/{listing_id}/files/{file_id}` | Delete file |

### Shops
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}` | Get shop |
| GET | `/v3/application/shops?shop_name=xxx` | Find shop by name |
| PUT | `/v3/application/shops/{shop_id}` | Update shop |
| GET | `/v3/application/shops/{shop_id}/sections` | List sections |
| POST | `/v3/application/shops/{shop_id}/sections` | Create section |

### Return Policies
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/return-policies` | List return policies |
| GET | `/v3/application/shops/{shop_id}/return-policies/{id}` | Get return policy |
| POST | `/v3/application/shops/{shop_id}/return-policies` | Create return policy |

### Receipts & Transactions
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/receipts` | List receipts |
| GET | `/v3/application/shops/{shop_id}/receipts/{receipt_id}` | Get receipt |
| PUT | `/v3/application/shops/{shop_id}/receipts/{receipt_id}` | Update receipt |
| GET | `/v3/application/shops/{shop_id}/receipts/{receipt_id}/transactions` | Receipt transactions |
| GET | `/v3/application/shops/{shop_id}/receipts/{receipt_id}/payments` | Receipt payments |
| GET | `/v3/application/shops/{shop_id}/transactions` | List shop transactions |
| GET | `/v3/application/shops/{shop_id}/transactions/{transaction_id}` | Get transaction |

### Payments & Ledger
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/payments` | List shop payments |
| GET | `/v3/application/shops/{shop_id}/payment-account/ledger-entries` | Ledger entries |

### Users
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/users/{user_id}` | Get user |
| GET | `/v3/application/users/{user_id}/addresses` | User addresses |
| GET | `/v3/application/users/{user_id}/shops` | User's shops |

### Shipping Profiles
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/shipping-profiles` | List profiles |
| GET | `/v3/application/shops/{shop_id}/shipping-profiles/{id}` | Get profile |
| POST | `/v3/application/shops/{shop_id}/shipping-profiles` | Create profile |
| DELETE | `/v3/application/shops/{shop_id}/shipping-profiles/{id}` | Delete profile |

### Reviews
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/shops/{shop_id}/reviews` | Shop reviews |
| GET | `/v3/application/listings/{listing_id}/reviews` | Listing reviews |

### Taxonomy
| Method | Path | Description |
|--------|------|-------------|
| GET | `/v3/application/buyer-taxonomy/nodes` | Full taxonomy tree |
| GET | `/v3/application/buyer-taxonomy/nodes/{id}/properties` | Taxonomy properties |

### Utility
| Method | Path | Description |
|--------|------|-------------|
| GET | `/ping` | Health check |

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
# Search active listings
curl -H "x-api-key: test" \
  "http://localhost:8080/v3/application/listings/active?keywords=necklace"

# Create a listing
curl -X POST -H "x-api-key: test" -H "Content-Type: application/json" \
  http://localhost:8080/v3/application/shops/5001/listings \
  -d '{"title":"New Item","description":"A thing","quantity":5,"price":19.99,"who_made":"i_did","when_made":"2020_2026","taxonomy_id":1207}'

# Get shop receipts
curl -H "x-api-key: test" \
  "http://localhost:8080/v3/application/shops/5001/receipts"
```

## Architecture

```
cmd/server/main.go          — Entry point, flag parsing, middleware chain
internal/
  models/                   — All Etsy API data types (1:1 with OpenAPI spec)
  store/store.go            — Thread-safe in-memory data store
  handlers/                 — HTTP handlers for each endpoint group
    router.go               — URL routing
    helpers.go              — JSON encoding, path parsing utilities
    listings.go             — Listing CRUD + images, files, inventory
    shops.go                — Shop, sections, return policies
    receipts.go             — Receipts, transactions
    payments.go             — Payments, ledger entries
    users.go                — Users, addresses
    reviews.go              — Reviews
    shipping.go             — Shipping profiles
    taxonomy.go             — Buyer taxonomy
  middleware/auth.go        — Auth, CORS, logging, content-type
  seed/seed.go              — Realistic test data
```
