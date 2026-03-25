# SmartSure - Internal System Documentation

## Overview

**SmartSure** is a comprehensive, enterprise-grade **Device Insurance & Warranty Management System** built with **Hexagonal Architecture (Ports and Adapters)** pattern. The system specializes in insurance for electronic devices including smartphones, tablets, laptops, smartwatches, and other consumer electronics.

---

## Architecture Pattern: Hexagonal Architecture

The system follows the **Hexagonal Architecture** (also known as Ports and Adapters) principles:

```
┌─────────────────────────────────────────────────────────────────┐
│                        INTERFACES LAYER                          │
│  (HTTP Handlers, GraphQL, gRPC, Webhooks, Middleware)           │
├─────────────────────────────────────────────────────────────────┤
│                      APPLICATION LAYER                           │
│  (Commands, Queries, DTOs, Application Services)                │
├─────────────────────────────────────────────────────────────────┤
│                        DOMAIN LAYER                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │    Models    │  │   Services   │  │    Events    │           │
│  │  (Entities)  │  │  (Business   │  │  (Domain     │           │
│  │              │  │   Logic)     │  │   Events)    │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
│  ┌──────────────┐  ┌──────────────┐                              │
│  │    Ports     │  │    Types     │                              │
│  │ (Interfaces) │  │  (Value Obj) │                              │
│  └──────────────┘  └──────────────┘                              │
├─────────────────────────────────────────────────────────────────┤
│                     INFRASTRUCTURE LAYER                         │
│  (Database, Messaging, External APIs, Auth, Middleware)         │
└─────────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility | Location |
|-------|---------------|----------|
| **Domain** | Core business logic, entities, value objects | `/domain/` |
| **Application** | Use cases, orchestration, DTOs | `/application/` |
| **Infrastructure** | Technical implementations, adapters | `/infrastructure/` |
| **Interfaces** | API endpoints, external interfaces | `/interfaces/` |

---

## Core Domain Entities

### 1. **User (Customer)**
Represents insurance customers with comprehensive profile management.

**Key Attributes:**
- Personal identification (name, email, phone, DOB)
- Authentication (password hash, 2FA settings, biometric data)
- Address information
- Insurance profile (KYC status, verification levels)
- Financial data (credit score, payment methods)
- Risk assessment (risk score, fraud score, claim frequency)
- Compliance (GDPR consent, AML status, PEP status)
- Engagement metrics (loyalty tier, CLV, churn risk)

**Capabilities:**
- Multi-device insurance management
- Family/household plan support
- Corporate account linking
- Loyalty & rewards program
- Risk-based pricing
- Fraud detection integration

---

### 2. **Device**
Central entity representing insured electronic devices with extensive tracking.

**Device Categories:**
- `smartphone` - Mobile phones
- `tablet` - Tablet computers
- `laptop` - Notebook computers
- `smartwatch` - Wearable devices
- `wearable` - Other wearables

**Core Embedded Structs:**
- `DeviceIdentification` - IMEI, serial number, unique identifiers
- `DeviceClassification` - Brand, model, manufacturer, category
- `DeviceSpecifications` - Technical specs (storage, RAM, OS, etc.)
- `DevicePhysicalCondition` - Grade, condition, damage assessment
- `DeviceFinancial` - Purchase price, current value, market value
- `DeviceRiskAssessment` - Risk scores, blacklist status
- `DeviceSecurity` - Security features, encryption status
- `DeviceNetwork` - Network connectivity details
- `DeviceStatusInfo` - Active status, stolen flag, location
- `DeviceLifecycle` - Warranty dates, lifecycle stage
- `DeviceVerification` - Verification status and dates
- `DeviceCompliance` - Regulatory compliance status
- `DeviceDocumentation` - Purchase receipts, certificates
- `DeviceOwnership` - Owner ID, corporate account linkage
- `DeviceWarranty` - Warranty information

**Relationship Categories:**

| Category | Relationships |
|----------|--------------|
| **Core** | Owner, CorporateAccount, Policies, Claims |
| **Ecosystem** | Swaps, TradeIns, Refurbishments, Rentals, Financings, Layaways, Repairs, Subscriptions, MarketListings, Accessories, SpareParts |
| **Insurance** | ClaimHistory, PremiumCalculations, RiskProfile, CoverageOptions, Deductibles, InsuranceAudits, ClaimValidations |
| **Financial** | DepreciationCurve, InsurableValue, TotalCostOwnership, ResaleValues, FinancialRisk |
| **Fraud** | FraudPatterns, AnomalyDetections, IdentityVerifications, FraudInvestigations, BlacklistManagement, FraudPrevention |
| **IoT** | IoTConnections, ActiveIoTConnection, IoTSensorData, IoTCommands, IoTDeviceConfiguration |
| **Maintenance** | PredictiveAlerts, DeviceHealthScores, FailurePredictions, MaintenanceSchedule, UpgradeRecommendations, ValueForecasts |
| **Monitoring** | MonitoringSessions, RealTimeMetrics, LiveAlerts, AlertConfigurations, AlertHistory |
| **Analytics** | UsagePatterns, BehaviorScore, LocationHistory, NetworkActivity, CustomerSegmentation, Profitability, MarketAnalyses |
| **Compliance** | ComplianceStatus, LegalHolds, ExportControls, DataPrivacy, RegulatoryReporting, SecurityCompliance |
| **Emergency** | EmergencyContacts, BackupStatus, RecoveryOptions, EmergencyLocation, PanicMode, DisasterRecovery |
| **Integration** | WarrantyProviders, ServiceProviders, InsurancePartners, FinancialInstitutions, EcosystemIntegration, APIIntegrations |
| **Sustainability** | CarbonFootprint, RecyclingScore, SustainabilityMetrics, EcoLabel, LifecycleAssessment, Repairability |

---

### 3. **Policy**
Insurance policy entity with comprehensive coverage management.

**Key Components:**
- `PolicyIdentification` - Policy number, version, identifiers
- `PolicyClassification` - Type, category, business line
- `PolicyCoverageDetails` - Coverage amounts, deductibles, limits
- `PolicyPricing` - Premium calculations, taxes, fees
- `PolicyDiscounts` - Applied discounts (loyalty, bundle, no-claims)
- `PolicyLoadings` - Additional risk loadings
- `PolicyPaymentInfo` - Payment status, frequency, billing
- `PolicyLifecycle` - Status, effective dates, renewal info
- `PolicyRiskAssessment` - Risk scores, fraud detection
- `PolicyUnderwritingInfo` - Underwriting decisions, inspections
- `PolicyCompliance` - KYC, AML, regulatory status
- `PolicyDocumentation` - Policy documents, certificates
- `PolicyCommunication` - Notification preferences
- `PolicyAnalytics` - Performance metrics, satisfaction scores

**Specialized Policy Features:**
- `PolicyCoverage` - Screen, water, theft, loss protection
- `PolicyServiceOptions` - Express service, loaner devices, home service
- `PolicyInternationalCoverage` - Travel coverage, roaming protection
- `PolicyClaimLimits` - Annual claim limits, waiting periods
- `PolicyLoyaltyProgram` - Loyalty points, tier benefits
- `PolicySmartFeatures` - Safety scores, preventive maintenance
- `PolicyFamilyGroup` - Family plans, group discounts
- `PolicyEnvironmental` - Green policies, eco discounts
- `PolicyCorporate` - Enterprise features, MDM, BYOD
- `PolicyReinsurance` - Reinsurance arrangements
- `PolicyCoinsurance` - Risk sharing arrangements
- `PolicyReserves` - IBNR, claim reserves
- `PolicyInvestment` - Investment-linked features
- `PolicyCommission` - Agent/broker commissions
- `PolicyLegal` - Legal terms, dispute resolution
- `PolicyRegulatoryFiling` - Compliance filings
- `PolicyPredictiveAnalytics` - ML-based predictions
- `PolicyCustomerJourney` - Journey tracking
- `PolicyIntegrations` - Third-party integrations
- `PolicyTelematics` - Usage-based insurance
- `PolicyMultiCurrency` - Multi-currency support
- `PolicyAutomation` - Automated workflows

---

### 4. **Claim**
Insurance claim processing entity with comprehensive workflow support.

**Core Embedded Structs:**
- `ClaimIdentification` - Claim number, type, priority
- `ClaimFinancial` - Claimed amount, approved amount, deductibles
- `ClaimLifecycle` - Status tracking, dates, workflow stage
- `ClaimInvestigation` - Fraud scores, investigation details
- `ClaimSettlement` - Settlement details, payout information
- `ClaimDocumentation` - Documents, evidence
- `ClaimAssignment` - Adjuster assignments
- `ClaimCompliance` - Regulatory compliance
- `ClaimMetrics` - Performance metrics, touchpoints

**Specialized Claim Features:**
- `ClaimInvestigationDetail` - Detailed investigation records
- `ClaimFraudDetection` - Fraud analysis, risk indicators
- `ClaimWorkflow` - Workflow management, stage tracking
- `ClaimSettlementDetail` - Settlement negotiations
- `ClaimPayment` - Payment processing, reconciliation
- `ClaimReserve` - Reserve management, IBNR
- `ClaimSubrogation` - Recovery from third parties
- `ClaimAppeal` - Appeal management
- `ClaimCommunication` - Customer communications
- `ClaimAnalytics` - Claim analytics, predictions
- `ClaimLitigation` - Legal proceedings
- `ClaimArbitration` - Alternative dispute resolution
- `ClaimReporting` - Regulatory reporting

**Smartphone-Specific Features:**
- `ClaimDeviceDiagnostics` - Remote diagnostics
- `ClaimRepairNetwork` - Repair shop finder
- `ClaimReplacementDevice` - Device replacement
- `ClaimDigitalAssets` - Data loss assessment
- `ClaimAccessories` - Accessory coverage
- `ClaimGeolocation` - Location tracking
- `ClaimPreventiveMeasures` - Prevention recommendations
- `ClaimSelfService` - Self-service options
- `ClaimBiometricVerification` - Biometric auth
- `Claim5GAndIoT` - IoT device coverage
- `ClaimAugmentedReality` - AR-guided repairs
- `ClaimCryptocurrency` - Crypto wallet coverage
- `ClaimSubscriptionServices` - App/service coverage
- `ClaimEnvironmentalImpact` - Sustainability
- `ClaimFoldableFlexible` - Foldable device coverage
- `ClaimHealthAndWellness` - Health data coverage
- `ClaimBusinessContinuity` - Business interruption

---

### 5. **Payment**
Payment processing entity with comprehensive billing support.

**Entities:**
- `PaymentMethod` - Stored payment methods (cards, bank, mobile money)
- `Payment` - Individual payment transactions
- `Subscription` - Recurring billing subscriptions
- `BillingHistory` - Historical billing records
- `Invoice` - Generated invoices
- `Commission` - Agent/partner commissions
- `PromoCode` - Promotional discount codes

---

### 6. **Corporate Account**
Enterprise/corporate insurance management.

**Entities:**
- `CorporateAccount` - Company accounts
- `CorporateEmployee` - Employee records
- `CorporatePolicy` - Corporate insurance policies
- `FleetDevice` - Corporate device fleet

---

### 7. **Repair Shop**
Repair network management.

**Entities:**
- `RepairShop` - Authorized repair centers
- `RepairBooking` - Repair appointments
- `RepairStatusUpdate` - Status tracking
- `RepairReview` - Customer reviews
- `ReplacementDevice` - Replacement inventory
- `ReplacementOrder` - Replacement orders
- `TemporaryDevice` - Loaner devices
- `DeviceLoan` - Device loans

---

## Domain Services (Ports)

### Service Interfaces

| Service | Description | Key Operations |
|---------|-------------|----------------|
| **DeviceService** | Device management & insurance | Register, verify, assess risk, calculate premium, process claims |
| **UserService** | Customer lifecycle management | Registration, authentication, risk assessment, segmentation |
| **PolicyService** | Policy lifecycle management | Create, renew, cancel, calculate premiums, manage coverage |
| **ClaimService** | Claim processing workflow | Submit, assess, process, settle, investigate fraud |
| **PaymentService** | Payment processing | Process payments, manage subscriptions, handle commissions |
| **FraudService** | Fraud detection & prevention | Detect fraud, investigate, blacklist management |
| **RepairService** | Repair network management | Schedule repairs, track status, manage inventory |
| **CorporateService** | Enterprise account management | Manage fleets, employees, corporate policies |
| **CommunicationService** | Customer communications | Notifications, emails, SMS, push notifications |
| **UnderwritingService** | Risk assessment & underwriting | Evaluate risks, make decisions, manage rules |
| **PricingService** | Premium calculations | Calculate rates, apply discounts, manage pricing models |
| **DocumentService** | Document management | Store, retrieve, process documents |
| **SupportService** | Customer support | Tickets, knowledge base, chat |
| **ReportingService** | Analytics & reporting | Generate reports, dashboards, insights |
| **AdminService** | System administration | User management, configuration |
| **FamilyService** | Family plan management | Add members, calculate discounts |
| **WarrantyService** | Warranty management | Track warranties, process claims |
| **WhiteLabelService** | White-label configurations | Brand management, customization |
| **BlockchainService** | Blockchain integration | Immutable records, smart contracts |
| **RiskProfileService** | Risk profiling | Assess and manage risk profiles |
| **ComplianceSecurityService** | Compliance & security | GDPR, AML, KYC management |
| **QuoteService** | Insurance quotes | Generate quotes, manage quotes |
| **PartnerService** | Partner management | Integrations, commissions |
| **NotificationService** | Notification delivery | Multi-channel notifications |

---

## Repository Ports

| Repository | Entity Managed | Operations |
|------------|---------------|------------|
| **UserRepository** | User | CRUD, search, authentication |
| **DeviceRepository** | Device | CRUD, search, history, analytics |
| **PolicyRepository** | Policy | CRUD, lifecycle, renewals |
| **ClaimRepository** | Claim | CRUD, workflow, fraud detection |
| **PaymentRepository** | Payment | CRUD, reconciliation |
| **CorporateRepository** | CorporateAccount | CRUD, fleet management |
| **RepairRepository** | RepairShop | CRUD, scheduling |
| **DocumentRepository** | Document | CRUD, storage |
| **FraudRepository** | Fraud data | Detection, investigation |
| **AdminRepository** | Admin data | System management |
| **FamilyRepository** | Family plans | Member management |
| **CommunicationRepository** | Communications | Message tracking |
| **PricingRepository** | Pricing models | Rate management |
| **QuoteRepository** | Quotes | Quote management |
| **ReportingRepository** | Reports | Analytics data |
| **RiskProfileRepository** | Risk profiles | Risk data |
| **SupportRepository** | Support tickets | Ticket management |
| **UnderwritingRepository** | Underwriting data | Decision tracking |
| **WarrantyRepository** | Warranties | Warranty tracking |
| **WhiteLabelRepository** | White-label configs | Brand management |
| **PartnerRepository** | Partners | Partner management |

---

## Key Features

### 1. **Multi-Device Insurance**
- Support for smartphones, tablets, laptops, smartwatches, wearables
- Device categorization with specific coverage rules
- Cross-device relationship tracking

### 2. **Claims Processing**
- End-to-end claim lifecycle management
- Automated fraud detection
- Remote device diagnostics
- Repair network integration
- Device replacement workflows
- Digital assets recovery

### 3. **Fraud Detection**
- Real-time fraud scoring
- Pattern analysis
- Network detection
- Velocity checks
- Blacklist management

### 4. **Pricing & Underwriting**
- Risk-based premium calculation
- Dynamic pricing models
- Discount management (loyalty, bundle, no-claims)
- Automated underwriting decisions

### 5. **Corporate/B2B Features**
- Fleet management
- BYOD (Bring Your Own Device) programs
- Employee device assignments
- Bulk operations
- Corporate billing

### 6. **Customer Experience**
- Self-service portals
- Real-time notifications
- Mobile app support
- Chat/communication integration
- Loyalty programs

### 7. **IoT & Smart Features**
- IoT device connectivity
- Real-time monitoring
- Predictive maintenance
- Usage-based insurance
- Health score tracking

### 8. **Compliance & Security**
- GDPR compliance
- AML/KYC verification
- PEP screening
- Sanctions checking
- Data encryption
- Audit trails

### 9. **Analytics & Intelligence**
- Customer lifetime value (CLV)
- Churn prediction
- Risk trend analysis
- Market intelligence
- Predictive analytics

### 10. **Sustainability**
- Carbon footprint tracking
- Recycling programs
- Eco-label certifications
- Repairability scoring
- Lifecycle assessment

---

## Economic Agents

### 1. **Customers/Users**
- Individual consumers
- Family plan members
- Corporate employees
- VIP/Premium customers

### 2. **Insurance Provider (SmartSure)**
- Policy issuer
- Claims processor
- Risk underwriter
- Payment collector

### 3. **Corporate Accounts**
- Enterprise clients
- Fleet managers
- BYOD program administrators

### 4. **Repair Shops/Partners**
- Authorized repair centers
- Parts suppliers
- Service providers

### 5. **Agents/Brokers**
- Insurance agents
- Corporate brokers
- Referral partners

### 6. **Financial Institutions**
- Payment processors (Stripe, PayPal)
- Banks
- Mobile money providers

### 7. **Third-Party Services**
- Fraud detection services
- Credit bureaus
- Identity verification providers
- IoT platform providers

---

## Data Transfer Objects (DTOs)

### Device DTOs
- `DeviceDTO` - Device API response
- `DeviceCreateRequest` - Create device request
- `DeviceUpdateRequest` - Update device request
- `DeviceSearchRequest` - Search parameters
- `DeviceInsuranceRequest` - Insurance purchase request
- `DeviceInsuranceResponse` - Insurance purchase response
- `DeviceValuationRequest` - Valuation request
- `DeviceValuationResponse` - Valuation response
- `DeviceTransferRequest` - Ownership transfer request
- `DeviceEligibilityResponse` - Eligibility check response

---

## Technology Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go (Golang) |
| **Database** | PostgreSQL (with GORM ORM) |
| **Caching** | Redis (implied from patterns) |
| **Messaging** | Apache Kafka |
| **External APIs** | REST, gRPC, GraphQL |
| **Authentication** | JWT, OAuth2 |
| **Payments** | Stripe, PayPal, Mobile Money |
| **Blockchain** | Custom blockchain integration |
| **AI/ML** | OpenAI integration |

---

## Project Structure

```
/home/eng/Documents/internal/
├── domain/                      # Core business logic
│   ├── models/                  # Domain entities
│   │   ├── device/             # Device-related models
│   │   ├── policy/             # Policy-related models
│   │   ├── claim/              # Claim-related models
│   │   ├── user/               # User-related models
│   │   ├── repair/             # Repair-related models
│   │   └── shared/             # Shared models
│   ├── services/               # Domain service implementations
│   ├── ports/                  # Interface definitions
│   │   ├── services/           # Service ports
│   │   └── repositories/       # Repository ports
│   ├── events/                 # Domain events
│   └── types/                  # Value objects, enums
│
├── application/                 # Use cases & orchestration
│   ├── services/               # Application services
│   ├── commands/               # CQRS commands
│   ├── queries/                # CQRS queries
│   └── dto/                    # Data transfer objects
│
├── infrastructure/              # Technical implementations
│   ├── database/               # Database adapters
│   │   ├── postgres/           # PostgreSQL implementation
│   │   └── migrations/         # Database migrations
│   ├── messaging/              # Message brokers
│   │   └── kafka/              # Kafka implementation
│   ├── auth/                   # Authentication
│   ├── external/               # External API adapters
│   └── middleware/             # Middleware components
│
└── interfaces/                  # API layer
    ├── http/                   # HTTP handlers
    ├── graphql/                # GraphQL schema
    └── grpc/                   # gRPC services
```

---

## Database Schema Overview

### Core Tables
- `users` - Customer accounts
- `devices` - Insured devices
- `policies` - Insurance policies
- `claims` - Insurance claims
- `payments` - Payment transactions
- `corporate_accounts` - Enterprise accounts
- `repair_shops` - Repair network

### Supporting Tables
- `payment_methods` - Stored payment methods
- `subscriptions` - Billing subscriptions
- `invoices` - Generated invoices
- `commissions` - Agent commissions
- `documents` - Uploaded documents
- `fraud_investigations` - Fraud cases
- `device_history` - Device audit trail

---

## Integration Points

### External Services
- **Payment Gateways**: Stripe, PayPal, Mobile Money
- **Identity Verification**: KYC/AML providers
- **Credit Bureaus**: Credit scoring
- **IoT Platforms**: Device connectivity
- **Blockchain**: Immutable record keeping
- **AI/ML Services**: OpenAI for predictions
- **Communication**: Email, SMS, Push notifications
- **Repair Networks**: Service provider APIs

---

## Security Features

- JWT-based authentication
- Role-based access control (RBAC)
- Data encryption at rest
- API rate limiting
- Fraud detection algorithms
- Audit logging
- GDPR compliance
- Sanctions screening
- PEP (Politically Exposed Persons) checks

---

## Business Rules & Validation

### Device Validation
- IMEI must be 15 or 17 digits
- Required fields: model, brand, manufacturer
- Owner ID cannot be nil
- Non-negative financial values

### Policy Validation
- Active policies required for claims
- Payment status affects policy status
- Coverage limits enforced
- Deductible calculations

### Claim Validation
- Claimed amount must be positive
- Incident date cannot be in future
- Policy must be active
- Fraud score thresholds

---

## Summary

The SmartSure system is a comprehensive, enterprise-grade device insurance platform built with clean architecture principles. It supports:

- **Multiple device types** with category-specific coverage
- **End-to-end claims processing** with fraud detection
- **Corporate/B2B** fleet management
- **Advanced analytics** and predictive modeling
- **Multi-channel** customer engagement
- **Compliance** with financial regulations
- **IoT integration** for smart devices
- **Sustainability** tracking and eco-programs

The hexagonal architecture ensures:
- **Testability** through port interfaces
- **Maintainability** with clear separation of concerns
- **Flexibility** to swap infrastructure implementations
- **Scalability** through event-driven patterns
