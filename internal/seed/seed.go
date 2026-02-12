package seed

import (
	"fmt"
	"time"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
	"github.com/vlah-software-house/etsy-mock-api/internal/store"
)

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func int64Ptr(i int64) *int64 { return &i }
func float32Ptr(f float32) *float32 { return &f }
func boolPtr(b bool) *bool    { return &b }

func Load(s *store.Store) {
	now := time.Now().Unix()
	ago := func(days int) int64 { return now - int64(days*86400) }

	// --- Users ---
	users := []*models.User{
		{UserID: 1001, PrimaryEmail: strPtr("alice@example.com"), FirstName: strPtr("Alice"), LastName: strPtr("Johnson"), ImageURL75x75: strPtr("https://mock.etsy.com/avatars/alice.jpg")},
		{UserID: 1002, PrimaryEmail: strPtr("bob@example.com"), FirstName: strPtr("Bob"), LastName: strPtr("Smith"), ImageURL75x75: strPtr("https://mock.etsy.com/avatars/bob.jpg")},
		{UserID: 1003, PrimaryEmail: strPtr("carol@example.com"), FirstName: strPtr("Carol"), LastName: strPtr("Williams"), ImageURL75x75: strPtr("https://mock.etsy.com/avatars/carol.jpg")},
		{UserID: 1004, PrimaryEmail: strPtr("dave@example.com"), FirstName: strPtr("Dave"), LastName: strPtr("Brown"), ImageURL75x75: strPtr("https://mock.etsy.com/avatars/dave.jpg")},
	}
	for _, u := range users {
		s.Users[u.UserID] = u
	}

	// --- User Addresses ---
	s.UserAddresses[1001] = []*models.UserAddress{
		{UserAddressID: 2001, UserID: 1001, Name: "Alice Johnson", FirstLine: "123 Craft Lane", City: "Portland", State: strPtr("OR"), Zip: strPtr("97201"), ISOCountryCode: strPtr("US"), CountryName: strPtr("United States"), IsDefaultShippingAddress: true},
	}
	s.UserAddresses[1002] = []*models.UserAddress{
		{UserAddressID: 2002, UserID: 1002, Name: "Bob Smith", FirstLine: "456 Maker Ave", City: "Austin", State: strPtr("TX"), Zip: strPtr("78701"), ISOCountryCode: strPtr("US"), CountryName: strPtr("United States"), IsDefaultShippingAddress: true},
	}

	// --- Shops ---
	reviewCount42 := 42
	reviewAvg4_8 := float32(4.8)
	reviewCount15 := 15
	reviewAvg4_5 := float32(4.5)

	shops := []*models.Shop{
		{
			ShopID: 5001, UserID: 1001, ShopName: "AliceCrafts",
			CreateDate: ago(365), CreatedTimestamp: ago(365),
			Title: strPtr("Handmade Jewelry & Accessories"),
			Announcement: strPtr("Welcome to AliceCrafts! Free shipping on orders over $50."),
			CurrencyCode: "USD", LoginName: "alicecrafts",
			AcceptsCustomRequests: true,
			URL: "https://www.etsy.com/shop/AliceCrafts",
			Languages: []string{"en"},
			IsUsingStructuredPolicies: true, HasOnboardedStructuredPolicies: true,
			IsDirectCheckoutOnboarded: true, IsEtsyPaymentsOnboarded: true,
			IsShopUSBased: true, TransactionSoldCount: 156,
			ListingActiveCount: 12, DigitalListingCount: 2,
			ShippingFromCountryISO: strPtr("US"), ShopLocationCountryISO: strPtr("US"),
			NumFavorers: 234,
			ReviewCount: &reviewCount42, ReviewAverage: &reviewAvg4_8,
			UpdateDate: ago(1), UpdatedTimestamp: ago(1),
			PolicyWelcome: strPtr("Thanks for visiting!"),
			PolicyPayment: strPtr("We accept all major credit cards through Etsy Payments."),
			PolicyShipping: strPtr("Orders ship within 1-3 business days."),
			PolicyRefunds: strPtr("Returns accepted within 30 days."),
		},
		{
			ShopID: 5002, UserID: 1002, ShopName: "BobsWoodworks",
			CreateDate: ago(200), CreatedTimestamp: ago(200),
			Title: strPtr("Custom Woodworking & Home Decor"),
			Announcement: strPtr("Handcrafted wooden items made with love."),
			CurrencyCode: "USD", LoginName: "bobswoodworks",
			AcceptsCustomRequests: true,
			URL: "https://www.etsy.com/shop/BobsWoodworks",
			Languages: []string{"en"},
			IsUsingStructuredPolicies: true, HasOnboardedStructuredPolicies: true,
			IsDirectCheckoutOnboarded: true, IsEtsyPaymentsOnboarded: true,
			IsShopUSBased: true, TransactionSoldCount: 78,
			ListingActiveCount: 8, DigitalListingCount: 0,
			ShippingFromCountryISO: strPtr("US"), ShopLocationCountryISO: strPtr("US"),
			NumFavorers: 89,
			ReviewCount: &reviewCount15, ReviewAverage: &reviewAvg4_5,
			UpdateDate: ago(3), UpdatedTimestamp: ago(3),
		},
	}
	for _, shop := range shops {
		s.Shops[shop.ShopID] = shop
	}

	// --- Shipping Profiles ---
	profiles := []*models.ShopShippingProfile{
		{
			ShippingProfileID: 3001, Title: strPtr("Standard US Shipping"), UserID: 1001,
			OriginCountryISO: "US", ProfileType: "manual", OriginPostalCode: strPtr("97201"),
			ShippingProfileDestinations: []models.ShopShippingProfileDestination{
				{ShippingProfileDestinationID: 3101, ShippingProfileID: 3001, OriginCountryISO: "US", DestinationCountryISO: "US", DestinationRegion: "none", PrimaryCost: models.USD(499), SecondaryCost: models.USD(199), MinDeliveryDays: intPtr(3), MaxDeliveryDays: intPtr(7)},
				{ShippingProfileDestinationID: 3102, ShippingProfileID: 3001, OriginCountryISO: "US", DestinationCountryISO: "", DestinationRegion: "none", PrimaryCost: models.USD(1499), SecondaryCost: models.USD(499), MinDeliveryDays: intPtr(7), MaxDeliveryDays: intPtr(21)},
			},
			ShippingProfileUpgrades: []models.ShopShippingProfileUpgrade{
				{ShippingProfileID: 3001, UpgradeID: 3201, UpgradeName: "Priority", Type: 0, Rank: 1, Language: "en", Price: models.USD(999), SecondaryPrice: models.USD(499), MinDeliveryDays: intPtr(1), MaxDeliveryDays: intPtr(3)},
			},
		},
		{
			ShippingProfileID: 3002, Title: strPtr("Woodworks Shipping"), UserID: 1002,
			OriginCountryISO: "US", ProfileType: "manual", OriginPostalCode: strPtr("78701"),
			ShippingProfileDestinations: []models.ShopShippingProfileDestination{
				{ShippingProfileDestinationID: 3103, ShippingProfileID: 3002, OriginCountryISO: "US", DestinationCountryISO: "US", DestinationRegion: "none", PrimaryCost: models.USD(999), SecondaryCost: models.USD(499), MinDeliveryDays: intPtr(5), MaxDeliveryDays: intPtr(10)},
			},
			ShippingProfileUpgrades: []models.ShopShippingProfileUpgrade{},
		},
	}
	for _, p := range profiles {
		s.ShippingProfiles[p.ShippingProfileID] = p
	}

	// --- Return Policies ---
	deadline30 := 30
	s.ShopReturnPolicies[4001] = &models.ShopReturnPolicy{ReturnPolicyID: 4001, ShopID: 5001, AcceptsReturns: true, AcceptsExchanges: true, ReturnDeadline: &deadline30}
	s.ShopReturnPolicies[4002] = &models.ShopReturnPolicy{ReturnPolicyID: 4002, ShopID: 5002, AcceptsReturns: true, AcceptsExchanges: false, ReturnDeadline: &deadline30}

	// --- Shop Sections ---
	s.ShopSections[6001] = &models.ShopSection{ShopSectionID: 6001, Title: "Necklaces", Rank: 1, UserID: 5001, ActiveListingCount: 4}
	s.ShopSections[6002] = &models.ShopSection{ShopSectionID: 6002, Title: "Earrings", Rank: 2, UserID: 5001, ActiveListingCount: 3}
	s.ShopSections[6003] = &models.ShopSection{ShopSectionID: 6003, Title: "Bracelets", Rank: 3, UserID: 5001, ActiveListingCount: 2}
	s.ShopSections[6004] = &models.ShopSection{ShopSectionID: 6004, Title: "Cutting Boards", Rank: 1, UserID: 5002, ActiveListingCount: 3}
	s.ShopSections[6005] = &models.ShopSection{ShopSectionID: 6005, Title: "Shelves", Rank: 2, UserID: 5002, ActiveListingCount: 2}

	// --- Listings ---
	iDid := "i_did"
	when2020 := "2020_2026"
	lang := "en"
	spID3001 := int64(3001)
	spID3002 := int64(3002)
	rpID4001 := int64(4001)
	rpID4002 := int64(4002)
	secNecklaces := int64(6001)
	secEarrings := int64(6002)
	secCuttingBoards := int64(6004)
	taxJewelry := 1207
	taxWood := 562

	listings := []*models.ShopListing{
		{
			ListingID: 7001, UserID: 1001, ShopID: 5001,
			Title: "Handmade Silver Pendant Necklace", Description: "Beautiful handcrafted sterling silver pendant with a minimalist design. Perfect for everyday wear or as a gift.",
			State: "active", CreationTimestamp: ago(90), CreatedTimestamp: ago(90), EndingTimestamp: now + 86400*30, OriginalCreationTimestamp: ago(90), LastModifiedTimestamp: ago(5), UpdatedTimestamp: ago(5),
			Quantity: 15, ShopSectionID: &secNecklaces, URL: "https://www.etsy.com/listing/7001", NumFavorers: 23,
			IsTaxable: true, ListingType: "physical", Tags: []string{"silver", "necklace", "pendant", "handmade", "minimalist"},
			Materials: []string{"sterling silver", "silver chain"}, ShippingProfileID: &spID3001, ReturnPolicyID: &rpID4001,
			ProcessingMin: intPtr(1), ProcessingMax: intPtr(3), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), ItemWeight: float32Ptr(0.5), ItemWeightUnit: strPtr("oz"),
			Price: models.USD(4500), TaxonomyID: &taxJewelry, Language: &lang, Views: 312, Style: []string{"minimalist"}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7002, UserID: 1001, ShopID: 5001,
			Title: "Gold Hoop Earrings - Large", Description: "Elegant gold-plated hoop earrings. Lightweight and comfortable for all-day wear.",
			State: "active", CreationTimestamp: ago(60), CreatedTimestamp: ago(60), EndingTimestamp: now + 86400*60, OriginalCreationTimestamp: ago(60), LastModifiedTimestamp: ago(10), UpdatedTimestamp: ago(10),
			Quantity: 25, ShopSectionID: &secEarrings, URL: "https://www.etsy.com/listing/7002", NumFavorers: 45,
			IsTaxable: true, ListingType: "physical", Tags: []string{"gold", "earrings", "hoops", "jewelry", "gift"},
			Materials: []string{"gold plated brass"}, ShippingProfileID: &spID3001, ReturnPolicyID: &rpID4001,
			ProcessingMin: intPtr(1), ProcessingMax: intPtr(3), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), ItemWeight: float32Ptr(0.3), ItemWeightUnit: strPtr("oz"),
			Price: models.USD(2800), TaxonomyID: &taxJewelry, Language: &lang, Views: 567, Style: []string{"bohemian"}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7003, UserID: 1001, ShopID: 5001,
			Title: "Beaded Friendship Bracelet Set", Description: "Set of 3 handmade beaded bracelets. Choose your colors at checkout.",
			State: "active", CreationTimestamp: ago(30), CreatedTimestamp: ago(30), EndingTimestamp: now + 86400*90, OriginalCreationTimestamp: ago(30), LastModifiedTimestamp: ago(2), UpdatedTimestamp: ago(2),
			Quantity: 50, URL: "https://www.etsy.com/listing/7003", NumFavorers: 12,
			IsTaxable: true, IsPersonalizable: true, PersonalizationIsRequired: false, PersonalizationCharCountMax: intPtr(50), PersonalizationInstructions: strPtr("Tell us your 3 preferred colors"),
			ListingType: "physical", Tags: []string{"bracelet", "friendship", "beaded", "colorful", "set"},
			Materials: []string{"glass beads", "elastic cord"}, ShippingProfileID: &spID3001, ReturnPolicyID: &rpID4001,
			ProcessingMin: intPtr(2), ProcessingMax: intPtr(5), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), Price: models.USD(1800), TaxonomyID: &taxJewelry, Language: &lang, Views: 189, Style: []string{}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7004, UserID: 1001, ShopID: 5001,
			Title: "Digital Jewelry Making Guide - PDF", Description: "Complete guide to beginner jewelry making with step-by-step instructions and photos.",
			State: "active", CreationTimestamp: ago(120), CreatedTimestamp: ago(120), EndingTimestamp: now + 86400*120, OriginalCreationTimestamp: ago(120), LastModifiedTimestamp: ago(30), UpdatedTimestamp: ago(30),
			Quantity: 999, URL: "https://www.etsy.com/listing/7004", NumFavorers: 78,
			IsTaxable: true, ListingType: "download", Tags: []string{"digital", "jewelry", "guide", "tutorial", "pdf"},
			Materials: []string{}, ShippingProfileID: nil, ReturnPolicyID: &rpID4001,
			WhoMade: &iDid, WhenMade: &when2020, IsSupply: boolPtr(false),
			Price: models.USD(1299), TaxonomyID: &taxJewelry, Language: &lang, Views: 1024, Style: []string{}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7005, UserID: 1001, ShopID: 5001,
			Title: "Rose Gold Ring - Stacking Ring", Description: "Dainty rose gold stacking ring. Available in sizes 5-10.",
			State: "draft", CreationTimestamp: ago(2), CreatedTimestamp: ago(2), EndingTimestamp: now + 86400*120, OriginalCreationTimestamp: ago(2), LastModifiedTimestamp: ago(1), UpdatedTimestamp: ago(1),
			Quantity: 30, URL: "https://www.etsy.com/listing/7005",
			IsTaxable: true, HasVariations: true, ListingType: "physical", Tags: []string{"ring", "rose gold", "stacking", "dainty"},
			Materials: []string{"rose gold plated brass"}, ShippingProfileID: &spID3001, ReturnPolicyID: &rpID4001,
			ProcessingMin: intPtr(1), ProcessingMax: intPtr(3), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), Price: models.USD(2200), TaxonomyID: &taxJewelry, Language: &lang, Views: 5, Style: []string{"minimalist"},
		},
		{
			ListingID: 7010, UserID: 1002, ShopID: 5002,
			Title: "Walnut Cutting Board - Large", Description: "Handcrafted walnut cutting board with juice groove. 18x12 inches.",
			State: "active", CreationTimestamp: ago(150), CreatedTimestamp: ago(150), EndingTimestamp: now + 86400*30, OriginalCreationTimestamp: ago(150), LastModifiedTimestamp: ago(7), UpdatedTimestamp: ago(7),
			Quantity: 8, ShopSectionID: &secCuttingBoards, URL: "https://www.etsy.com/listing/7010", NumFavorers: 56,
			IsTaxable: true, ListingType: "physical", Tags: []string{"cutting board", "walnut", "kitchen", "handmade", "wood"},
			Materials: []string{"walnut wood", "food-safe mineral oil"}, ShippingProfileID: &spID3002, ReturnPolicyID: &rpID4002,
			ProcessingMin: intPtr(3), ProcessingMax: intPtr(7), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), ItemWeight: float32Ptr(4.5), ItemWeightUnit: strPtr("lb"),
			ItemLength: float32Ptr(18), ItemWidth: float32Ptr(12), ItemHeight: float32Ptr(1.5), ItemDimensionsUnit: strPtr("in"),
			Price: models.USD(7500), TaxonomyID: &taxWood, Language: &lang, Views: 445, Style: []string{"rustic"}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7011, UserID: 1002, ShopID: 5002,
			Title: "Floating Wall Shelf - Pine", Description: "Minimalist floating wall shelf made from reclaimed pine. Perfect for plants or books.",
			State: "active", CreationTimestamp: ago(100), CreatedTimestamp: ago(100), EndingTimestamp: now + 86400*60, OriginalCreationTimestamp: ago(100), LastModifiedTimestamp: ago(14), UpdatedTimestamp: ago(14),
			Quantity: 12, URL: "https://www.etsy.com/listing/7011", NumFavorers: 34,
			IsTaxable: true, ListingType: "physical", Tags: []string{"shelf", "floating", "wall", "pine", "reclaimed"},
			Materials: []string{"reclaimed pine", "steel brackets"}, ShippingProfileID: &spID3002, ReturnPolicyID: &rpID4002,
			ProcessingMin: intPtr(5), ProcessingMax: intPtr(10), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), ItemWeight: float32Ptr(3.0), ItemWeightUnit: strPtr("lb"),
			ItemLength: float32Ptr(24), ItemWidth: float32Ptr(6), ItemHeight: float32Ptr(2), ItemDimensionsUnit: strPtr("in"),
			Price: models.USD(4500), TaxonomyID: &taxWood, Language: &lang, Views: 223, Style: []string{"minimalist", "rustic"}, ShouldAutoRenew: true,
		},
		{
			ListingID: 7012, UserID: 1002, ShopID: 5002,
			Title: "Personalized Wooden Sign", Description: "Custom engraved wooden sign. Perfect for weddings, home decor, or gifts.",
			State: "active", CreationTimestamp: ago(80), CreatedTimestamp: ago(80), EndingTimestamp: now + 86400*45, OriginalCreationTimestamp: ago(80), LastModifiedTimestamp: ago(3), UpdatedTimestamp: ago(3),
			Quantity: 20, URL: "https://www.etsy.com/listing/7012", NumFavorers: 67,
			IsTaxable: true, IsCustomizable: true, IsPersonalizable: true, PersonalizationIsRequired: true, PersonalizationCharCountMax: intPtr(100), PersonalizationInstructions: strPtr("Enter the text you want engraved"),
			ListingType: "physical", Tags: []string{"wooden sign", "personalized", "custom", "engraved", "wedding"},
			Materials: []string{"oak wood", "wood stain"}, ShippingProfileID: &spID3002, ReturnPolicyID: &rpID4002,
			ProcessingMin: intPtr(5), ProcessingMax: intPtr(14), WhoMade: &iDid, WhenMade: &when2020,
			IsSupply: boolPtr(false), Price: models.USD(5500), TaxonomyID: &taxWood, Language: &lang, Views: 678, Style: []string{"rustic"}, ShouldAutoRenew: true,
		},
	}
	for _, l := range listings {
		s.Listings[l.ListingID] = l
	}

	// --- Listing Images ---
	for _, l := range listings {
		for i := 1; i <= 3; i++ {
			imgID := l.ListingID*10 + int64(i)
			h := 800
			wid := 600
			s.ListingImages[l.ListingID] = append(s.ListingImages[l.ListingID], &models.ListingImage{
				ListingID: l.ListingID, ListingImageID: imgID,
				CreationTsz: l.CreatedTimestamp, CreatedTimestamp: l.CreatedTimestamp, Rank: i,
				URL75x75: fmt.Sprintf("https://mock.etsy.com/images/%d_75x75.jpg", imgID),
				URL170x135: fmt.Sprintf("https://mock.etsy.com/images/%d_170x135.jpg", imgID),
				URL570xN: fmt.Sprintf("https://mock.etsy.com/images/%d_570xN.jpg", imgID),
				URLFullxfull: fmt.Sprintf("https://mock.etsy.com/images/%d_fullxfull.jpg", imgID),
				FullHeight: &h, FullWidth: &wid,
				AltText: strPtr(fmt.Sprintf("Image %d of %s", i, l.Title)),
			})
		}
	}

	// --- Listing Files (for digital listing) ---
	s.ListingFiles[7004] = []*models.ListingFile{
		{ListingFileID: 8001, ListingID: 7004, Rank: 1, Filename: "jewelry-making-guide.pdf", Filesize: "2.5 MB", SizeBytes: 2621440, Filetype: "application/pdf", CreateTimestamp: ago(120), CreatedTimestamp: ago(120)},
	}

	// --- Receipts ---
	addr1 := "123 Main St\nPortland, OR 97201\nUnited States"
	addr2 := "456 Oak Ave\nSeattle, WA 98101\nUnited States"
	paidTs := ago(15)
	shippedTs := ago(12)

	receipts := []*models.ShopReceipt{
		{
			ReceiptID: 9001, ReceiptType: 0, SellerUserID: 1001, SellerEmail: strPtr("alice@example.com"),
			BuyerUserID: 1003, BuyerEmail: strPtr("carol@example.com"),
			Name: "Carol Williams", FirstLine: strPtr("789 Elm St"), City: strPtr("Seattle"), State: strPtr("WA"), Zip: strPtr("98101"), CountryISO: strPtr("US"),
			Status: "completed", FormattedAddress: &addr2, PaymentMethod: "cc", IsPaid: true, IsShipped: true,
			CreateTimestamp: ago(20), CreatedTimestamp: ago(20), UpdateTimestamp: ago(10), UpdatedTimestamp: ago(10),
			MessageFromBuyer: strPtr("Love your work!"), MessageFromSeller: strPtr("Thank you! Shipped today!"),
			Grandtotal: models.USD(5298), Subtotal: models.USD(4500), TotalPrice: models.USD(4500),
			TotalShippingCost: models.USD(499), TotalTaxCost: models.USD(299), TotalVatCost: models.USD(0),
			DiscountAmt: models.USD(0), GiftWrapPrice: models.USD(0),
			Shipments: []models.ShopReceiptShipment{
				{ReceiptShippingID: int64Ptr(9101), ShipmentNotificationTimestamp: ago(12), CarrierName: "USPS", TrackingCode: "9400111899223100012345"},
			},
			Transactions: []models.ShopReceiptTransaction{
				{TransactionID: 9201, Title: strPtr("Handmade Silver Pendant Necklace"), SellerUserID: 1001, BuyerUserID: 1003, CreateTimestamp: ago(20), CreatedTimestamp: ago(20), PaidTimestamp: &paidTs, ShippedTimestamp: &shippedTs, Quantity: 1, ReceiptID: 9001, ListingID: intPtr(7001), TransactionType: "listing", Price: models.USD(4500), ShippingCost: models.USD(499), Variations: []models.TransactionVariation{}, ShippingProfileID: &spID3001, MinProcessingDays: intPtr(1), MaxProcessingDays: intPtr(3), ShippingMethod: strPtr("USPS First Class")},
			},
			Refunds: []models.ShopRefund{},
		},
		{
			ReceiptID: 9002, ReceiptType: 0, SellerUserID: 1001, SellerEmail: strPtr("alice@example.com"),
			BuyerUserID: 1004, BuyerEmail: strPtr("dave@example.com"),
			Name: "Dave Brown", FirstLine: strPtr("321 Pine Rd"), City: strPtr("Portland"), State: strPtr("OR"), Zip: strPtr("97201"), CountryISO: strPtr("US"),
			Status: "paid", FormattedAddress: &addr1, PaymentMethod: "cc", IsPaid: true, IsShipped: false,
			CreateTimestamp: ago(3), CreatedTimestamp: ago(3), UpdateTimestamp: ago(3), UpdatedTimestamp: ago(3),
			MessageFromBuyer: strPtr("Can you gift wrap this?"), IsGift: true, GiftMessage: "Happy Birthday!", GiftSender: "Dave",
			Grandtotal: models.USD(3298), Subtotal: models.USD(2800), TotalPrice: models.USD(2800),
			TotalShippingCost: models.USD(499), TotalTaxCost: models.USD(0), TotalVatCost: models.USD(0),
			DiscountAmt: models.USD(0), GiftWrapPrice: models.USD(0),
			Shipments: []models.ShopReceiptShipment{},
			Transactions: []models.ShopReceiptTransaction{
				{TransactionID: 9202, Title: strPtr("Gold Hoop Earrings - Large"), SellerUserID: 1001, BuyerUserID: 1004, CreateTimestamp: ago(3), CreatedTimestamp: ago(3), PaidTimestamp: int64Ptr(ago(3)), Quantity: 1, ReceiptID: 9002, ListingID: intPtr(7002), TransactionType: "listing", Price: models.USD(2800), ShippingCost: models.USD(499), Variations: []models.TransactionVariation{}, ShippingProfileID: &spID3001, MinProcessingDays: intPtr(1), MaxProcessingDays: intPtr(3)},
			},
			Refunds: []models.ShopRefund{},
		},
		{
			ReceiptID: 9003, ReceiptType: 0, SellerUserID: 1002, SellerEmail: strPtr("bob@example.com"),
			BuyerUserID: 1003, BuyerEmail: strPtr("carol@example.com"),
			Name: "Carol Williams", FirstLine: strPtr("789 Elm St"), City: strPtr("Seattle"), State: strPtr("WA"), Zip: strPtr("98101"), CountryISO: strPtr("US"),
			Status: "paid", FormattedAddress: &addr2, PaymentMethod: "cc", IsPaid: true, IsShipped: false,
			CreateTimestamp: ago(5), CreatedTimestamp: ago(5), UpdateTimestamp: ago(5), UpdatedTimestamp: ago(5),
			Grandtotal: models.USD(8499), Subtotal: models.USD(7500), TotalPrice: models.USD(7500),
			TotalShippingCost: models.USD(999), TotalTaxCost: models.USD(0), TotalVatCost: models.USD(0),
			DiscountAmt: models.USD(0), GiftWrapPrice: models.USD(0),
			Shipments: []models.ShopReceiptShipment{},
			Transactions: []models.ShopReceiptTransaction{
				{TransactionID: 9203, Title: strPtr("Walnut Cutting Board - Large"), SellerUserID: 1002, BuyerUserID: 1003, CreateTimestamp: ago(5), CreatedTimestamp: ago(5), PaidTimestamp: int64Ptr(ago(5)), Quantity: 1, ReceiptID: 9003, ListingID: intPtr(7010), TransactionType: "listing", Price: models.USD(7500), ShippingCost: models.USD(999), Variations: []models.TransactionVariation{}, ShippingProfileID: &spID3002, MinProcessingDays: intPtr(3), MaxProcessingDays: intPtr(7)},
			},
			Refunds: []models.ShopRefund{},
		},
	}
	for _, r := range receipts {
		s.Receipts[r.ReceiptID] = r
		for i := range r.Transactions {
			txn := r.Transactions[i]
			s.Transactions[txn.TransactionID] = &txn
		}
	}

	// --- Payments ---
	payments := []*models.Payment{
		{PaymentID: 11001, BuyerUserID: 1003, ShopID: 5001, ReceiptID: 9001, AmountGross: models.USD(5298), AmountFees: models.USD(345), AmountNet: models.USD(4953), Currency: "USD", ShopCurrency: strPtr("USD"), BuyerCurrency: strPtr("USD"), ShippingAddressID: 2002, Status: "settled", CreateTimestamp: ago(20), CreatedTimestamp: ago(20), UpdateTimestamp: ago(10), UpdatedTimestamp: ago(10), PaymentAdjustments: []models.PaymentAdjustment{}},
		{PaymentID: 11002, BuyerUserID: 1004, ShopID: 5001, ReceiptID: 9002, AmountGross: models.USD(3298), AmountFees: models.USD(215), AmountNet: models.USD(3083), Currency: "USD", ShopCurrency: strPtr("USD"), BuyerCurrency: strPtr("USD"), ShippingAddressID: 2001, Status: "open", CreateTimestamp: ago(3), CreatedTimestamp: ago(3), UpdateTimestamp: ago(3), UpdatedTimestamp: ago(3), PaymentAdjustments: []models.PaymentAdjustment{}},
		{PaymentID: 11003, BuyerUserID: 1003, ShopID: 5002, ReceiptID: 9003, AmountGross: models.USD(8499), AmountFees: models.USD(554), AmountNet: models.USD(7945), Currency: "USD", ShopCurrency: strPtr("USD"), BuyerCurrency: strPtr("USD"), ShippingAddressID: 2002, Status: "open", CreateTimestamp: ago(5), CreatedTimestamp: ago(5), UpdateTimestamp: ago(5), UpdatedTimestamp: ago(5), PaymentAdjustments: []models.PaymentAdjustment{}},
	}
	for _, p := range payments {
		s.Payments[p.PaymentID] = p
	}

	// --- Ledger Entries ---
	s.LedgerEntries[5001] = []*models.PaymentAccountLedgerEntry{
		{EntryID: 12001, LedgerID: 5001, SequenceNumber: 1, Amount: 4953, Currency: "USD", Description: "Sale: Handmade Silver Pendant Necklace", Balance: 4953, CreateDate: ago(15), CreatedTimestamp: ago(15), LedgerType: "credit", ReferenceType: "receipt", ReferenceID: strPtr("9001"), PaymentAdjustments: []models.PaymentAdjustment{}},
		{EntryID: 12002, LedgerID: 5001, SequenceNumber: 2, Amount: -345, Currency: "USD", Description: "Etsy fee", Balance: 4608, CreateDate: ago(15), CreatedTimestamp: ago(15), LedgerType: "debit", ReferenceType: "fee", ReferenceID: strPtr("9001"), PaymentAdjustments: []models.PaymentAdjustment{}},
	}

	// --- Reviews ---
	s.Reviews[5001] = []*models.ListingReview{
		{ShopID: 5001, ListingID: 7001, TransactionID: 9201, BuyerUserID: int64Ptr(1003), Rating: 5, Review: "Absolutely beautiful necklace! The craftsmanship is amazing.", Language: "en", CreateTimestamp: ago(8), CreatedTimestamp: ago(8), UpdateTimestamp: ago(8), UpdatedTimestamp: ago(8)},
		{ShopID: 5001, ListingID: 7002, TransactionID: 9202, BuyerUserID: int64Ptr(1004), Rating: 4, Review: "Great earrings, very lightweight and comfortable.", Language: "en", CreateTimestamp: ago(2), CreatedTimestamp: ago(2), UpdateTimestamp: ago(2), UpdatedTimestamp: ago(2)},
	}
	s.Reviews[5002] = []*models.ListingReview{
		{ShopID: 5002, ListingID: 7010, TransactionID: 9203, BuyerUserID: int64Ptr(1003), Rating: 5, Review: "This cutting board is a work of art. Beautifully made!", Language: "en", CreateTimestamp: ago(1), CreatedTimestamp: ago(1), UpdateTimestamp: ago(1), UpdatedTimestamp: ago(1)},
	}

	// --- Taxonomy ---
	s.TaxonomyNodes = []models.BuyerTaxonomyNode{
		{ID: 1, Level: 0, Name: "Jewelry", Children: []models.BuyerTaxonomyNode{
			{ID: 1207, Level: 1, Name: "Necklaces", ParentID: int64Ptr(1), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{1, 1207}},
			{ID: 1208, Level: 1, Name: "Earrings", ParentID: int64Ptr(1), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{1, 1208}},
			{ID: 1209, Level: 1, Name: "Rings", ParentID: int64Ptr(1), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{1, 1209}},
			{ID: 1210, Level: 1, Name: "Bracelets", ParentID: int64Ptr(1), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{1, 1210}},
		}, FullPathTaxonomyIDs: []int{1}},
		{ID: 2, Level: 0, Name: "Home & Living", Children: []models.BuyerTaxonomyNode{
			{ID: 562, Level: 1, Name: "Kitchen & Dining", ParentID: int64Ptr(2), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{2, 562}},
			{ID: 563, Level: 1, Name: "Home Decor", ParentID: int64Ptr(2), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{2, 563}},
			{ID: 564, Level: 1, Name: "Furniture", ParentID: int64Ptr(2), Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{2, 564}},
		}, FullPathTaxonomyIDs: []int{2}},
		{ID: 3, Level: 0, Name: "Clothing", Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{3}},
		{ID: 4, Level: 0, Name: "Art & Collectibles", Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{4}},
		{ID: 5, Level: 0, Name: "Craft Supplies & Tools", Children: []models.BuyerTaxonomyNode{}, FullPathTaxonomyIDs: []int{5}},
	}

	s.TaxonomyProperties[1207] = []models.BuyerTaxonomyNodeProperty{
		{PropertyID: 100, Name: "material", DisplayName: "Material", IsRequired: false, SupportsAttributes: true, SupportsVariations: false, IsMultivalued: true, PossibleValues: []models.BuyerTaxonomyPropertyValue{{ValueID: 1001, Name: "Sterling Silver"}, {ValueID: 1002, Name: "Gold"}, {ValueID: 1003, Name: "Rose Gold"}, {ValueID: 1004, Name: "Brass"}}, SelectedValues: []models.BuyerTaxonomyPropertyValue{}, Scales: []models.BuyerTaxonomyPropertyScale{}},
		{PropertyID: 101, Name: "length", DisplayName: "Chain Length", IsRequired: false, SupportsAttributes: true, SupportsVariations: true, IsMultivalued: false, PossibleValues: []models.BuyerTaxonomyPropertyValue{{ValueID: 2001, Name: "14 inches"}, {ValueID: 2002, Name: "16 inches"}, {ValueID: 2003, Name: "18 inches"}, {ValueID: 2004, Name: "20 inches"}}, SelectedValues: []models.BuyerTaxonomyPropertyValue{}, Scales: []models.BuyerTaxonomyPropertyScale{}},
	}
	s.TaxonomyProperties[562] = []models.BuyerTaxonomyNodeProperty{
		{PropertyID: 200, Name: "material", DisplayName: "Wood Type", IsRequired: false, SupportsAttributes: true, SupportsVariations: false, IsMultivalued: true, PossibleValues: []models.BuyerTaxonomyPropertyValue{{ValueID: 3001, Name: "Walnut"}, {ValueID: 3002, Name: "Maple"}, {ValueID: 3003, Name: "Cherry"}, {ValueID: 3004, Name: "Oak"}}, SelectedValues: []models.BuyerTaxonomyPropertyValue{}, Scales: []models.BuyerTaxonomyPropertyScale{}},
	}
}
