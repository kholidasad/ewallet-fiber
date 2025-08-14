package docs

import "github.com/swaggo/swag"

// OpenAPI 3 spec yang di-embed langsung (tanpa perlu "swag init")
var doc = `{
  "openapi":"3.0.3",
  "info":{"title":"E-Wallet API","version":"1.0.0","description":"Fiber + MySQL"},
  "servers":[{"url":"http://localhost:8081"}],
  "components":{
    "securitySchemes":{"BearerAuth":{"type":"http","scheme":"bearer","bearerFormat":"JWT"}},
    "schemas":{
      "RegisterRequest":{"type":"object","required":["Email","Password"],"properties":{"Email":{"type":"string","format":"email"},"Password":{"type":"string","minLength":6}}},
      "LoginRequest":{"type":"object","required":["Email","Password"],"properties":{"Email":{"type":"string","format":"email"},"Password":{"type":"string"}}},
      "TokenResponse":{"type":"object","properties":{"token":{"type":"string"},"user":{"$ref":"#/components/schemas/User"}}},
      "CreateWalletRequest":{"type":"object","required":["name"],"properties":{"name":{"type":"string"}}},
      "AmountRequest":{"type":"object","required":["Currency","Amount"],"properties":{"Currency":{"type":"string","example":"USD"},"Amount":{"type":"string","example":"100.00"}}},
      "TransferRequest":{"type":"object","required":["to_wallet_id","FromCurrency","ToCurrency","Amount"],"properties":{"to_wallet_id":{"type":"integer","format":"int64"},"FromCurrency":{"type":"string"},"ToCurrency":{"type":"string"},"Amount":{"type":"string"}}},
      "PaymentRequest":{"type":"object","required":["Currency","Amount"],"properties":{"Currency":{"type":"string"},"Amount":{"type":"string"},"Reference":{"type":"string"},"Metadata":{"type":"string"}}},
      "User":{"type":"object","properties":{"id":{"type":"integer"},"email":{"type":"string"},"wallets":{"type":"array","items":{"$ref":"#/components/schemas/Wallet"}}}},
      "Wallet":{"type":"object","properties":{"id":{"type":"integer"},"user_id":{"type":"integer"},"name":{"type":"string"}}},
      "WalletBalance":{"type":"object","properties":{"wallet_id":{"type":"integer"},"currency_code":{"type":"string"},"balance":{"type":"string"}}},
      "Transaction":{"type":"object","properties":{"id":{"type":"integer"},"wallet_id":{"type":"integer"},"to_wallet_id":{"type":"integer","nullable":true},"type":{"type":"string","enum":["deposit","withdrawal","transfer","payment"]},"status":{"type":"string","enum":["success","failed"]},"amount":{"type":"string"},"currency_code":{"type":"string"},"exchange_rate":{"type":"string"},"converted_amount":{"type":"string"},"reference":{"type":"string"},"created_at":{"type":"string","format":"date-time"}}},
      "TransactionsPage":{"type":"object","properties":{"data":{"type":"array","items":{"$ref":"#/components/schemas/Transaction"}},"meta":{"type":"object","properties":{"page":{"type":"integer"},"page_size":{"type":"integer"},"total":{"type":"integer","format":"int64"},"total_pages":{"type":"integer"}}}}},
      "Summary":{"type":"object","properties":{"display_currency":{"type":"string"},"total":{"type":"string"},"by_currency":{"type":"object","additionalProperties":{"type":"string"}}}}
    }
  },
  "paths":{
    "/api/v1/auth/register":{
      "post":{
        "summary":"Register",
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/RegisterRequest"}}}},
        "responses":{"201":{"description":"Created","content":{"application/json":{"schema":{"$ref":"#/components/schemas/User"}}}}}
      }
    },
    "/api/v1/auth/login":{
      "post":{
        "summary":"Login",
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/LoginRequest"}}}},
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/TokenResponse"}}}}}
      }
    },
    "/api/v1/me":{
      "get":{
        "summary":"Get profile & wallets",
        "security":[{"BearerAuth":[]}],
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"type":"object"}}}}}
      }
    },
    "/api/v1/wallets":{
      "get":{
        "summary":"List wallets",
        "security":[{"BearerAuth":[]}],
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"type":"array","items":{"$ref":"#/components/schemas/Wallet"}}}}}}
      },
      "post":{
        "summary":"Create wallet",
        "security":[{"BearerAuth":[]}],
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/CreateWalletRequest"}}}},
        "responses":{"201":{"description":"Created","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Wallet"}}}}}
      }
    },
    "/api/v1/wallets/{id}/balances":{
      "get":{
        "summary":"Wallet balances",
        "security":[{"BearerAuth":[]}],
        "parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"integer"}}],
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"type":"array","items":{"$ref":"#/components/schemas/WalletBalance"}}}}}}
      }
    },
    "/api/v1/wallets/{id}/deposit":{
      "post":{
        "summary":"Deposit",
        "security":[{"BearerAuth":[]}],
        "parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"integer"}}],
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/AmountRequest"}}}},
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Transaction"}}}}}
      }
    },
    "/api/v1/wallets/{id}/withdraw":{
      "post":{
        "summary":"Withdraw",
        "security":[{"BearerAuth":[]}],
        "parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"integer"}}],
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/AmountRequest"}}}},
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Transaction"}}}}}
      }
    },
    "/api/v1/wallets/{id}/transfer":{
      "post":{
        "summary":"Transfer (with FX)",
        "security":[{"BearerAuth":[]}],
        "parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"integer"}}],
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/TransferRequest"}}}},
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Transaction"}}}}}
      }
    },
    "/api/v1/wallets/{id}/payment":{
      "post":{
        "summary":"Payment",
        "security":[{"BearerAuth":[]}],
        "parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"integer"}}],
        "requestBody":{"required":true,"content":{"application/json":{"schema":{"$ref":"#/components/schemas/PaymentRequest"}}}},
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Transaction"}}}}}
      }
    },
    "/api/v1/transactions":{
      "get":{
        "summary":"List transactions (paginated)",
        "security":[{"BearerAuth":[]}],
        "parameters":[
          {"name":"page","in":"query","schema":{"type":"integer","default":1}},
          {"name":"page_size","in":"query","schema":{"type":"integer","default":20}},
          {"name":"type","in":"query","schema":{"type":"string","enum":["deposit","withdrawal","transfer","payment"]}},
          {"name":"currency","in":"query","schema":{"type":"string","example":"USD"}}
        ],
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/TransactionsPage"}}}}}
      }
    },
    "/api/v1/summary":{
      "get":{
        "summary":"Total balance (display currency)",
        "security":[{"BearerAuth":[]}],
        "responses":{"200":{"description":"OK","content":{"application/json":{"schema":{"$ref":"#/components/schemas/Summary"}}}}}
      }
    }
  }
}`

type s struct{}

func (s *s) ReadDoc() string { return doc }

func init() { swag.Register(swag.Name, &s{}) }