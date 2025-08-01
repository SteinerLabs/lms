# Billing Service Data Models

## Overview

The Billing Service is responsible for managing payments, subscriptions, invoicing, and financial transactions within the LMS platform. This document defines the core data models used by the Billing Service.

## Data Models

### Customer

The Customer model represents a billing customer in the system.

```
type Customer struct {
    ID              string    // Unique identifier
    UserID          string    // User ID
    ExternalID      string    // External ID (e.g., Stripe customer ID)
    Email           string    // Customer email
    Name            string    // Customer name
    Phone           string    // Customer phone
    BillingAddress  Address   // Billing address
    ShippingAddress Address   // Shipping address
    TaxID           string    // Tax ID
    Currency        string    // Preferred currency
    Status          string    // Customer status (active, inactive)
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Address

The Address model represents a physical address.

```
type Address struct {
    ID              string    // Unique identifier
    Line1           string    // Address line 1
    Line2           string    // Address line 2
    City            string    // City
    State           string    // State/Province
    PostalCode      string    // Postal/ZIP code
    Country         string    // Country
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### PaymentMethod

The PaymentMethod model represents a payment method associated with a customer.

```
type PaymentMethod struct {
    ID              string    // Unique identifier
    CustomerID      string    // Customer ID
    ExternalID      string    // External ID (e.g., Stripe payment method ID)
    Type            string    // Payment method type (credit_card, bank_account, paypal)
    Status          string    // Payment method status (active, inactive)
    Default         bool      // Whether this is the default payment method
    LastFour        string    // Last four digits of card/account
    ExpiryMonth     int       // Expiry month for cards
    ExpiryYear      int       // Expiry year for cards
    Brand           string    // Card brand (visa, mastercard, etc.)
    BillingAddress  Address   // Billing address
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Product

The Product model represents a product that can be purchased.

```
type Product struct {
    ID              string    // Unique identifier
    ExternalID      string    // External ID (e.g., Stripe product ID)
    Name            string    // Product name
    Description     string    // Product description
    Type            string    // Product type (course, subscription, etc.)
    Active          bool      // Whether the product is active
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Price

The Price model represents a price for a product.

```
type Price struct {
    ID              string    // Unique identifier
    ProductID       string    // Product ID
    ExternalID      string    // External ID (e.g., Stripe price ID)
    Currency        string    // Currency
    UnitAmount      int       // Amount in smallest currency unit (e.g., cents)
    BillingScheme   string    // Billing scheme (per_unit, tiered)
    TieredPrices    []TieredPrice // Tiered prices if applicable
    Recurring       PriceRecurring // Recurring details if applicable
    Type            string    // Price type (one_time, recurring)
    Active          bool      // Whether the price is active
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### TieredPrice

The TieredPrice model represents a tier in a tiered pricing model.

```
type TieredPrice struct {
    ID              string    // Unique identifier
    PriceID         string    // Price ID
    UpTo            int       // Upper bound of tier
    UnitAmount      int       // Amount in smallest currency unit (e.g., cents)
    FlatAmount      int       // Flat amount for tier
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### PriceRecurring

The PriceRecurring model represents recurring details for a price.

```
type PriceRecurring struct {
    ID              string    // Unique identifier
    PriceID         string    // Price ID
    Interval        string    // Billing interval (day, week, month, year)
    IntervalCount   int       // Number of intervals between billings
    TrialPeriodDays int       // Number of trial days
    UsageType       string    // Usage type (licensed, metered)
    AggregateUsage  string    // Aggregate usage (sum, max, last_during_period, last_ever)
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Coupon

The Coupon model represents a discount coupon.

```
type Coupon struct {
    ID              string    // Unique identifier
    ExternalID      string    // External ID (e.g., Stripe coupon ID)
    Code            string    // Coupon code
    Name            string    // Coupon name
    Description     string    // Coupon description
    DiscountType    string    // Discount type (percentage, fixed_amount)
    DiscountAmount  int       // Discount amount (percentage or fixed amount)
    Currency        string    // Currency for fixed amount discounts
    Duration        string    // Duration (once, repeating, forever)
    DurationMonths  int       // Number of months for repeating duration
    MaxRedemptions  int       // Maximum number of redemptions
    RedemptionCount int       // Current redemption count
    ValidFrom       time.Time // Valid from timestamp
    ValidUntil      time.Time // Valid until timestamp
    Active          bool      // Whether the coupon is active
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Subscription

The Subscription model represents a recurring subscription.

```
type Subscription struct {
    ID                  string    // Unique identifier
    CustomerID          string    // Customer ID
    ExternalID          string    // External ID (e.g., Stripe subscription ID)
    Status              string    // Subscription status (active, past_due, canceled, etc.)
    PriceID             string    // Price ID
    Quantity            int       // Quantity of the subscription
    StartDate           time.Time // Start date
    CurrentPeriodStart  time.Time // Current period start
    CurrentPeriodEnd    time.Time // Current period end
    CanceledAt          time.Time // Canceled at timestamp
    CancelAtPeriodEnd   bool      // Whether to cancel at period end
    TrialStart          time.Time // Trial start timestamp
    TrialEnd            time.Time // Trial end timestamp
    DefaultPaymentMethodID string // Default payment method ID
    LatestInvoiceID     string    // Latest invoice ID
    NextBillingDate     time.Time // Next billing date
    Metadata            map[string]interface{} // Additional metadata
    CreatedAt           time.Time // Creation timestamp
    UpdatedAt           time.Time // Last update timestamp
}
```

### Invoice

The Invoice model represents an invoice for a customer.

```
type Invoice struct {
    ID                  string    // Unique identifier
    CustomerID          string    // Customer ID
    ExternalID          string    // External ID (e.g., Stripe invoice ID)
    SubscriptionID      string    // Subscription ID
    Status              string    // Invoice status (draft, open, paid, etc.)
    Currency            string    // Currency
    Subtotal            int       // Subtotal amount
    Tax                 int       // Tax amount
    Total               int       // Total amount
    AmountPaid          int       // Amount paid
    AmountDue           int       // Amount due
    AmountRemaining     int       // Amount remaining
    BillingReason       string    // Billing reason (subscription_create, subscription_cycle, etc.)
    DueDate             time.Time // Due date
    PeriodStart         time.Time // Period start
    PeriodEnd           time.Time // Period end
    PaymentIntentID     string    // Payment intent ID
    ReceiptNumber       string    // Receipt number
    ReceiptURL          string    // Receipt URL
    InvoicePDF          string    // Invoice PDF URL
    Lines               []InvoiceLine // Invoice lines
    Metadata            map[string]interface{} // Additional metadata
    CreatedAt           time.Time // Creation timestamp
    UpdatedAt           time.Time // Last update timestamp
}
```

### InvoiceLine

The InvoiceLine model represents a line item on an invoice.

```
type InvoiceLine struct {
    ID              string    // Unique identifier
    InvoiceID       string    // Invoice ID
    ExternalID      string    // External ID (e.g., Stripe invoice line ID)
    Type            string    // Line type (subscription, invoice_item)
    SubscriptionID  string    // Subscription ID
    PriceID         string    // Price ID
    ProductID       string    // Product ID
    Description     string    // Line description
    Quantity        int       // Quantity
    UnitAmount      int       // Unit amount
    Amount          int       // Total amount
    Currency        string    // Currency
    PeriodStart     time.Time // Period start
    PeriodEnd       time.Time // Period end
    Proration       bool      // Whether this is a proration
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Payment

The Payment model represents a payment made by a customer.

```
type Payment struct {
    ID              string    // Unique identifier
    CustomerID      string    // Customer ID
    ExternalID      string    // External ID (e.g., Stripe payment intent ID)
    InvoiceID       string    // Invoice ID
    PaymentMethodID string    // Payment method ID
    Amount          int       // Amount in smallest currency unit (e.g., cents)
    Currency        string    // Currency
    Status          string    // Payment status (pending, succeeded, failed)
    Description     string    // Payment description
    ReceiptEmail    string    // Receipt email
    ReceiptURL      string    // Receipt URL
    ErrorMessage    string    // Error message if failed
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Refund

The Refund model represents a refund of a payment.

```
type Refund struct {
    ID              string    // Unique identifier
    PaymentID       string    // Payment ID
    ExternalID      string    // External ID (e.g., Stripe refund ID)
    Amount          int       // Amount in smallest currency unit (e.g., cents)
    Currency        string    // Currency
    Status          string    // Refund status (pending, succeeded, failed)
    Reason          string    // Refund reason
    Description     string    // Refund description
    ReceiptURL      string    // Receipt URL
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### Tax

The Tax model represents a tax configuration.

```
type Tax struct {
    ID              string    // Unique identifier
    Name            string    // Tax name
    Description     string    // Tax description
    Percentage      float64   // Tax percentage
    Inclusive       bool      // Whether the tax is inclusive
    Country         string    // Country code
    State           string    // State/Province code
    Active          bool      // Whether the tax is active
    Metadata        map[string]interface{} // Additional metadata
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

### TaxRate

The TaxRate model represents a tax rate applied to a product.

```
type TaxRate struct {
    ID              string    // Unique identifier
    TaxID           string    // Tax ID
    ProductID       string    // Product ID
    Active          bool      // Whether the tax rate is active
    CreatedAt       time.Time // Creation timestamp
    UpdatedAt       time.Time // Last update timestamp
}
```

## Relationships

- A Customer belongs to a User
- A Customer can have many PaymentMethods
- A Customer can have many Subscriptions
- A Customer can have many Invoices
- A Customer can have many Payments
- A Product can have many Prices
- A Price can have many TieredPrices
- A Price can have one PriceRecurring
- A Subscription belongs to a Customer and a Price
- An Invoice belongs to a Customer and can belong to a Subscription
- An Invoice can have many InvoiceLines
- An InvoiceLine can belong to a Subscription, Price, and Product
- A Payment belongs to a Customer, Invoice, and PaymentMethod
- A Refund belongs to a Payment
- A Tax can have many TaxRates
- A TaxRate belongs to a Tax and a Product

## Database Schema

The Billing Service uses PostgreSQL for data storage. The schema includes the following tables:

- customers
- addresses
- payment_methods
- products
- prices
- tiered_prices
- price_recurring
- coupons
- subscriptions
- invoices
- invoice_lines
- payments
- refunds
- taxes
- tax_rates

## Events

The Billing Service publishes and consumes the following events:

### Published Events

#### PaymentCompleted

```
type PaymentCompleted struct {
    ID          string    // Payment ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    InvoiceID   string    // Invoice ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    Description string    // Payment description
    CompletedAt time.Time // Completion timestamp
}
```

#### PaymentFailed

```
type PaymentFailed struct {
    ID          string    // Payment ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    InvoiceID   string    // Invoice ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    ErrorMessage string   // Error message
    FailedAt    time.Time // Failure timestamp
}
```

#### SubscriptionCreated

```
type SubscriptionCreated struct {
    ID          string    // Subscription ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    PriceID     string    // Price ID
    ProductID   string    // Product ID
    Status      string    // Subscription status
    StartDate   time.Time // Start date
    CreatedAt   time.Time // Creation timestamp
}
```

#### SubscriptionUpdated

```
type SubscriptionUpdated struct {
    ID          string    // Subscription ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    PriceID     string    // Price ID
    ProductID   string    // Product ID
    Status      string    // Subscription status
    UpdatedAt   time.Time // Update timestamp
}
```

#### SubscriptionCanceled

```
type SubscriptionCanceled struct {
    ID          string    // Subscription ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    CanceledAt  time.Time // Cancellation timestamp
}
```

#### InvoiceCreated

```
type InvoiceCreated struct {
    ID          string    // Invoice ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    Status      string    // Invoice status
    DueDate     time.Time // Due date
    CreatedAt   time.Time // Creation timestamp
}
```

#### InvoicePaid

```
type InvoicePaid struct {
    ID          string    // Invoice ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    PaidAt      time.Time // Payment timestamp
}
```

#### RefundIssued

```
type RefundIssued struct {
    ID          string    // Refund ID
    PaymentID   string    // Payment ID
    CustomerID  string    // Customer ID
    UserID      string    // User ID
    Amount      int       // Amount in smallest currency unit (e.g., cents)
    Currency    string    // Currency
    Reason      string    // Refund reason
    IssuedAt    time.Time // Issuance timestamp
}
```

### Consumed Events

#### UserCreated (from Auth Service)

```
type UserCreated struct {
    ID        string    // User ID
    Email     string    // Email
    FirstName string    // First name
    LastName  string    // Last name
    CreatedAt time.Time // Creation timestamp
}
```

#### UserEnrolled (from Course Service)

```
type UserEnrolled struct {
    ID          string    // Enrollment ID
    CourseID    string    // Course ID
    UserID      string    // User ID
    EnrolledAt  time.Time // Enrollment timestamp
}
```