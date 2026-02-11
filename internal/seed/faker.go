package seed

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gabriel/etsy-mock/internal/models"
	"github.com/gabriel/etsy-mock/internal/store"
)

// Content pools for realistic Etsy-style data generation.
// No external dependencies — all pools are embedded.

var shopPrefixes = []string{
	"Cozy", "Wild", "Little", "Golden", "Silver", "Rustic", "Modern", "Vintage",
	"Handmade", "Artisan", "Bohemian", "Enchanted", "Creative", "Whimsical",
	"Sunlit", "Moonstone", "Coastal", "Mountain", "Prairie", "Maple",
	"Ivy", "Cedar", "Birch", "Willow", "Sage", "Amber", "Coral", "Indigo",
	"Crimson", "Emerald", "Honey", "Velvet", "Pebble", "Daisy", "Bloom",
	"Frost", "Misty", "Woven", "Twisted", "Polished",
}

var shopNouns = []string{
	"Craft", "Studio", "Workshop", "Atelier", "Boutique", "Nest", "Haven",
	"Forge", "Loom", "Kiln", "Garden", "Cottage", "Hearth", "Hollow",
	"Creek", "Ridge", "Meadow", "Grove", "Valley", "Lane",
}

var shopSuffixes = []string{
	"Co", "Designs", "Creations", "Goods", "Supply", "Made", "Works",
	"Shop", "Market", "Collective", "House", "Lab", "Place", "",
}

type category struct {
	Name         string
	TaxonomyID   int
	Items        []string
	Materials    [][]string
	Styles       []string
	Tags         [][]string
	PriceRange   [2]int // min/max in cents
	WeightRange  [2]float32
	WeightUnit   string
	Descriptions []string
}

var categories = []category{
	{
		Name: "Jewelry", TaxonomyID: 1207,
		Items: []string{
			"Pendant Necklace", "Chain Necklace", "Choker", "Layered Necklace",
			"Hoop Earrings", "Stud Earrings", "Drop Earrings", "Dangle Earrings",
			"Statement Ring", "Stacking Ring", "Signet Ring", "Band Ring",
			"Charm Bracelet", "Cuff Bracelet", "Bangle", "Tennis Bracelet",
			"Anklet", "Body Chain", "Hair Pin", "Brooch",
		},
		Materials: [][]string{
			{"sterling silver", "silver chain"},
			{"14k gold", "gold chain"},
			{"rose gold plated brass"},
			{"brass", "copper"},
			{"stainless steel"},
			{"freshwater pearls", "silk cord"},
			{"natural gemstones", "gold wire"},
			{"recycled silver"},
			{"titanium"},
			{"glass beads", "elastic cord"},
		},
		Styles: []string{"minimalist", "bohemian", "vintage", "art deco", "modern", "classic", "statement", "dainty", "geometric", "organic"},
		Tags: [][]string{
			{"handmade", "jewelry", "gift for her", "birthday gift"},
			{"minimalist", "everyday", "delicate", "layering"},
			{"bohemian", "boho", "festival", "free spirit"},
			{"vintage inspired", "antique style", "retro"},
			{"wedding", "bridal", "bridesmaid", "anniversary"},
		},
		PriceRange:  [2]int{1500, 15000},
		WeightRange: [2]float32{0.1, 2.0},
		WeightUnit:  "oz",
		Descriptions: []string{
			"Beautifully handcrafted with attention to every detail. This piece is perfect for everyday wear or as a special gift for someone you love.",
			"Each piece is made to order in our studio using ethically sourced materials. The design draws inspiration from nature and geometric forms.",
			"A stunning addition to any jewelry collection. Lightweight and comfortable for all-day wear. Comes gift-wrapped and ready to give.",
			"Designed and crafted by hand in our small workshop. We use only the finest materials to ensure lasting quality and beauty.",
			"This unique piece combines traditional craftsmanship with a modern aesthetic. Each one is slightly different, making yours truly one-of-a-kind.",
		},
	},
	{
		Name: "Home & Living", TaxonomyID: 562,
		Items: []string{
			"Cutting Board", "Serving Board", "Cheese Board", "Bread Board",
			"Floating Shelf", "Wall Shelf", "Corner Shelf", "Bookshelf",
			"Picture Frame", "Wall Art", "Wall Hanging", "Macrame",
			"Candle Holder", "Vase", "Planter", "Coaster Set",
			"Tray", "Bowl", "Spoon Set", "Utensil Holder",
		},
		Materials: [][]string{
			{"walnut wood", "food-safe mineral oil"},
			{"maple wood", "beeswax finish"},
			{"cherry wood", "natural oil"},
			{"oak wood", "wood stain"},
			{"reclaimed pine", "steel brackets"},
			{"bamboo", "natural finish"},
			{"olive wood"},
			{"acacia wood", "resin"},
			{"concrete", "wood"},
			{"ceramic", "glaze"},
		},
		Styles: []string{"rustic", "farmhouse", "modern", "minimalist", "industrial", "scandinavian", "mid-century", "coastal", "bohemian", "traditional"},
		Tags: [][]string{
			{"handmade", "home decor", "kitchen", "housewarming gift"},
			{"rustic", "farmhouse", "country", "wood"},
			{"modern", "minimalist", "contemporary", "clean lines"},
			{"wedding gift", "anniversary", "personalized"},
			{"eco-friendly", "sustainable", "reclaimed", "natural"},
		},
		PriceRange:  [2]int{2500, 30000},
		WeightRange: [2]float32{0.5, 15.0},
		WeightUnit:  "lb",
		Descriptions: []string{
			"Handcrafted in our workshop from carefully selected hardwood. Each piece is sanded smooth and finished with a food-safe coating that brings out the natural grain.",
			"A beautiful functional piece that adds warmth and character to any space. Made from sustainably sourced wood with traditional joinery techniques.",
			"This piece is made to last for generations. We take pride in creating items that are both beautiful and built to withstand daily use.",
			"Each piece is unique due to the natural variations in wood grain and color. We hand-select every board for optimal beauty and durability.",
			"Crafted with care and attention to detail. The natural materials and expert finishing make this a standout piece in any home.",
		},
	},
	{
		Name: "Clothing", TaxonomyID: 3,
		Items: []string{
			"Linen Dress", "Cotton Blouse", "Silk Scarf", "Wool Sweater",
			"Denim Jacket", "Kimono Robe", "Wrap Dress", "Maxi Skirt",
			"Crop Top", "Tunic", "Cardigan", "Poncho",
			"Beanie", "Mittens", "Socks", "Tote Bag",
		},
		Materials: [][]string{
			{"100% organic cotton"},
			{"linen", "cotton blend"},
			{"merino wool"},
			{"silk", "satin"},
			{"bamboo fabric"},
			{"hemp", "organic cotton"},
			{"cashmere blend"},
			{"recycled polyester"},
		},
		Styles: []string{"casual", "bohemian", "vintage", "elegant", "streetwear", "classic", "romantic", "oversized", "fitted", "relaxed"},
		Tags: [][]string{
			{"handmade", "clothing", "fashion", "sustainable"},
			{"organic", "eco-friendly", "natural fibers"},
			{"bohemian", "boho", "festival wear"},
			{"vintage inspired", "retro", "classic style"},
			{"gift for her", "birthday", "comfortable"},
		},
		PriceRange:  [2]int{2000, 20000},
		WeightRange: [2]float32{0.2, 2.0},
		WeightUnit:  "lb",
		Descriptions: []string{
			"Made from the softest natural fabrics, this piece drapes beautifully and feels amazing against the skin. Perfect for any season.",
			"Designed for comfort without sacrificing style. Each piece is cut and sewn by hand in our small studio using premium materials.",
			"A wardrobe staple that transitions effortlessly from day to night. Machine washable and designed to get softer with every wash.",
			"Inspired by timeless silhouettes with a modern twist. The natural fibers breathe and move with you throughout the day.",
			"Handmade with love and care. We believe clothing should be made to last, using quality materials and expert craftsmanship.",
		},
	},
	{
		Name: "Art & Collectibles", TaxonomyID: 4,
		Items: []string{
			"Print", "Original Painting", "Digital Download", "Illustration",
			"Photograph", "Poster", "Canvas Art", "Watercolor",
			"Sculpture", "Figurine", "Wall Mural", "Art Print Set",
			"Custom Portrait", "Pet Portrait", "Landscape Art", "Abstract Art",
		},
		Materials: [][]string{
			{"archival paper", "pigment ink"},
			{"canvas", "acrylic paint"},
			{"watercolor paper", "watercolor paint"},
			{"cotton rag paper", "giclée ink"},
			{"wood panel", "oil paint"},
			{"digital file"},
			{"recycled paper", "soy ink"},
			{"metal", "mixed media"},
		},
		Styles: []string{"abstract", "modern", "impressionist", "photorealistic", "whimsical", "botanical", "landscape", "portrait", "pop art", "surreal"},
		Tags: [][]string{
			{"art", "wall art", "home decor", "gallery wall"},
			{"print", "illustration", "artwork", "poster"},
			{"original art", "one of a kind", "collectible"},
			{"gift idea", "housewarming", "office decor"},
			{"nature art", "botanical", "landscape", "floral"},
		},
		PriceRange:  [2]int{1000, 50000},
		WeightRange: [2]float32{0.1, 5.0},
		WeightUnit:  "lb",
		Descriptions: []string{
			"This piece captures the beauty of the everyday in vivid detail. Printed on archival-quality paper to ensure lasting color and clarity.",
			"An original work of art created in our studio. Each piece is signed and numbered, making it a true collector's item.",
			"Transform your space with this stunning artwork. The colors are rich and vibrant, creating a focal point in any room.",
			"Created with a blend of traditional techniques and modern sensibility. This piece tells a story that resonates on a personal level.",
			"Museum-quality reproduction printed with premium pigment inks on heavyweight cotton rag paper. Ships flat in protective packaging.",
		},
	},
	{
		Name: "Craft Supplies", TaxonomyID: 5,
		Items: []string{
			"Bead Set", "Yarn Bundle", "Fabric Pack", "Stamp Set",
			"Charm Pack", "Wire Kit", "Paint Set", "Tool Kit",
			"Pattern", "Template", "Stencil Set", "Embroidery Kit",
			"Resin Kit", "Clay Set", "Leather Scraps", "Button Collection",
		},
		Materials: [][]string{
			{"mixed materials"},
			{"glass beads"},
			{"cotton yarn"},
			{"premium fabric"},
			{"metal findings"},
			{"natural leather"},
			{"polymer clay"},
			{"UV resin"},
		},
		Styles: []string{"colorful", "natural", "vintage", "modern", "assorted", "premium", "beginner-friendly", "professional"},
		Tags: [][]string{
			{"craft supplies", "DIY", "maker", "creative"},
			{"beading", "jewelry making", "supplies"},
			{"sewing", "quilting", "fabric", "notions"},
			{"art supplies", "painting", "drawing"},
			{"gift for crafter", "hobby", "handmade"},
		},
		PriceRange:  [2]int{500, 8000},
		WeightRange: [2]float32{0.1, 3.0},
		WeightUnit:  "lb",
		Descriptions: []string{
			"Everything you need to get started on your next creative project. Carefully curated and packaged for makers of all skill levels.",
			"High-quality supplies sourced from trusted manufacturers. Perfect for both beginners and experienced crafters.",
			"A wonderful assortment of materials to inspire your creativity. Each set is hand-picked to ensure variety and quality.",
			"Professional-grade supplies at an affordable price. We test every product ourselves before adding it to our shop.",
			"Unleash your creativity with this carefully curated kit. Includes detailed instructions and inspiration ideas.",
		},
	},
}

var firstNames = []string{
	"Emma", "Olivia", "Ava", "Isabella", "Sophia", "Mia", "Charlotte", "Amelia",
	"Harper", "Evelyn", "James", "Liam", "Noah", "William", "Oliver", "Benjamin",
	"Elijah", "Lucas", "Mason", "Logan", "Alexander", "Ethan", "Daniel", "Henry",
	"Sarah", "Jessica", "Emily", "Rachel", "Hannah", "Lauren", "Abigail", "Madison",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Anderson", "Taylor", "Thomas", "Jackson", "White",
	"Harris", "Martin", "Thompson", "Moore", "Allen", "Young", "King", "Wright",
	"Lopez", "Hill", "Scott", "Green", "Adams", "Baker", "Nelson", "Carter",
}

var cities = []struct {
	City  string
	State string
	Zip   string
}{
	{"Portland", "OR", "97201"}, {"Austin", "TX", "78701"}, {"Seattle", "WA", "98101"},
	{"Denver", "CO", "80201"}, {"Nashville", "TN", "37201"}, {"Asheville", "NC", "28801"},
	{"Brooklyn", "NY", "11201"}, {"San Francisco", "CA", "94102"}, {"Chicago", "IL", "60601"},
	{"Boston", "MA", "02101"}, {"Philadelphia", "PA", "19101"}, {"Savannah", "GA", "31401"},
	{"Santa Fe", "NM", "87501"}, {"Burlington", "VT", "05401"}, {"Bend", "OR", "97701"},
	{"Sedona", "AZ", "86336"}, {"Madison", "WI", "53701"}, {"Bozeman", "MT", "59715"},
	{"Charlottesville", "VA", "22901"}, {"Ann Arbor", "MI", "48104"},
}

var streets = []string{
	"123 Main St", "456 Oak Ave", "789 Elm St", "321 Pine Rd", "654 Maple Dr",
	"147 Cedar Ln", "258 Birch Way", "369 Walnut Ct", "741 Ash Pl", "852 Spruce Blvd",
	"963 Willow St", "174 Poplar Ave", "285 Chestnut Ln", "396 Redwood Dr", "417 Magnolia Way",
}

var reviewTexts = []string{
	"Absolutely love this! The quality is amazing and it arrived faster than expected.",
	"Beautiful craftsmanship. You can tell this was made with care and attention to detail.",
	"Exceeded my expectations! The photos don't do it justice — it's even more stunning in person.",
	"Perfect gift for my sister. She was so happy when she opened it! Will definitely order again.",
	"Gorgeous piece! The materials are high quality and it's exactly as described.",
	"So happy with my purchase. It's even better than I imagined. Highly recommend this shop!",
	"The packaging was beautiful and it arrived in perfect condition. Love supporting small businesses.",
	"Wonderful quality and fast shipping. This seller clearly takes pride in their work.",
	"I've gotten so many compliments! It's become my new favorite piece.",
	"Ordered this as a treat for myself and I'm so glad I did. Worth every penny.",
	"The attention to detail is incredible. You can really see the handmade quality.",
	"Already planning my next order! The colors are vibrant and true to the photos.",
	"This is my third purchase from this shop and they never disappoint. Consistent quality every time.",
	"Arrived beautifully wrapped. Makes such a thoughtful gift. The recipient absolutely loved it.",
	"Stunning work! I'll be back for more. Great communication from the seller too.",
	"Better than I hoped for. The texture and finish are just perfect.",
	"My new go-to gift for friends. Everyone loves receiving handmade items like these.",
	"Impeccable quality. This will last for years. Smart purchase!",
	"The design is unique and eye-catching. I get asked about it constantly.",
	"Fast shipping, great packaging, beautiful product. What more could you ask for?",
}

var shopAnnouncements = []string{
	"Welcome to our shop! All items are handmade with love. Free shipping on orders over $50!",
	"Thank you for visiting! New items added weekly. Custom orders always welcome.",
	"Handcrafted with care in our small studio. Each piece is unique and made to order.",
	"We believe in quality over quantity. Every item is tested and approved before shipping.",
	"Supporting local artisans and sustainable practices. Thank you for shopping small!",
	"Holiday sale! 15%% off everything. Use code HANDMADE at checkout.",
	"New collection just dropped! Check out our latest designs. Custom requests welcome.",
	"All items ship within 1-3 business days. Gift wrapping available at checkout.",
}

func GenerateFromConfig(s *store.Store, cfg SeedConfig) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now().Unix()
	ago := func(days int) int64 { return now - int64(days*86400) }

	// Taxonomy (same as hardcoded seed — always needed)
	loadTaxonomy(s)

	var allUserIDs []int64
	nextID := int64(10000)
	getID := func() int64 { nextID++; return nextID }

	// Generate buyer users
	numBuyers := cfg.Shops*3 + 5
	for i := 0; i < numBuyers; i++ {
		uid := getID()
		fn := pick(r, firstNames)
		ln := pick(r, lastNames)
		email := fmt.Sprintf("%s.%s.%d@example.com", strings.ToLower(fn), strings.ToLower(ln), uid)
		s.Users[uid] = &models.User{
			UserID:        uid,
			PrimaryEmail:  &email,
			FirstName:     &fn,
			LastName:      &ln,
			ImageURL75x75: strPtr(fmt.Sprintf("https://mock.etsy.com/avatars/%d.jpg", uid)),
		}
		city := cities[r.Intn(len(cities))]
		street := streets[r.Intn(len(streets))]
		ccode := "US"
		cname := "United States"
		s.UserAddresses[uid] = []*models.UserAddress{{
			UserAddressID:           getID(),
			UserID:                  uid,
			Name:                    fn + " " + ln,
			FirstLine:               street,
			City:                    city.City,
			State:                   &city.State,
			Zip:                     &city.Zip,
			ISOCountryCode:          &ccode,
			CountryName:             &cname,
			IsDefaultShippingAddress: true,
		}}
		allUserIDs = append(allUserIDs, uid)
	}

	// Generate shops
	for si := 0; si < cfg.Shops; si++ {
		shopID := getID()
		ownerID := allUserIDs[si] // first N users own shops
		shopName := generateShopName(r)
		cat := categories[r.Intn(len(categories))]

		city := cities[r.Intn(len(cities))]
		reviewCount := randRange(r, cfg.ReviewsPerShop.Min, cfg.ReviewsPerShop.Max)
		reviewAvg := float32(3.8 + r.Float64()*1.2) // 3.8-5.0
		countryISO := "US"

		shop := &models.Shop{
			ShopID:                         shopID,
			UserID:                         ownerID,
			ShopName:                       shopName,
			CreateDate:                     ago(randRange(r, 100, 1000)),
			CreatedTimestamp:               ago(randRange(r, 100, 1000)),
			Title:                          strPtr(fmt.Sprintf("Handmade %s & More", cat.Name)),
			Announcement:                   strPtr(shopAnnouncements[r.Intn(len(shopAnnouncements))]),
			CurrencyCode:                   "USD",
			LoginName:                      strings.ToLower(shopName),
			AcceptsCustomRequests:          r.Float64() > 0.3,
			URL:                            fmt.Sprintf("https://www.etsy.com/shop/%s", shopName),
			Languages:                      []string{"en"},
			IsUsingStructuredPolicies:      true,
			HasOnboardedStructuredPolicies: true,
			IsDirectCheckoutOnboarded:      true,
			IsEtsyPaymentsOnboarded:        true,
			IsShopUSBased:                  true,
			TransactionSoldCount:           randRange(r, 10, 500),
			ShippingFromCountryISO:         &countryISO,
			ShopLocationCountryISO:         &countryISO,
			NumFavorers:                    randRange(r, 5, 500),
			ReviewCount:                    &reviewCount,
			ReviewAverage:                  &reviewAvg,
			UpdateDate:                     ago(randRange(r, 0, 10)),
			UpdatedTimestamp:               ago(randRange(r, 0, 10)),
			PolicyWelcome:                  strPtr("Thanks for visiting our shop!"),
			PolicyPayment:                  strPtr("We accept all major credit cards through Etsy Payments."),
			PolicyShipping:                 strPtr("Orders ship within 1-3 business days."),
			PolicyRefunds:                  strPtr("Returns accepted within 30 days of delivery."),
		}
		s.Shops[shopID] = shop

		// Shipping profile
		spID := getID()
		postalCode := city.Zip
		spTitle := "Standard Shipping"
		s.ShippingProfiles[spID] = &models.ShopShippingProfile{
			ShippingProfileID: spID, Title: &spTitle, UserID: ownerID,
			OriginCountryISO: "US", ProfileType: "manual", OriginPostalCode: &postalCode,
			ShippingProfileDestinations: []models.ShopShippingProfileDestination{
				{ShippingProfileDestinationID: getID(), ShippingProfileID: spID, OriginCountryISO: "US", DestinationCountryISO: "US", DestinationRegion: "none", PrimaryCost: models.USD(randRange(r, 399, 999)), SecondaryCost: models.USD(randRange(r, 99, 399)), MinDeliveryDays: intPtr(3), MaxDeliveryDays: intPtr(7)},
			},
			ShippingProfileUpgrades: []models.ShopShippingProfileUpgrade{},
		}

		// Return policy
		rpID := getID()
		deadline := 30
		s.ShopReturnPolicies[rpID] = &models.ShopReturnPolicy{
			ReturnPolicyID: rpID, ShopID: shopID, AcceptsReturns: true, AcceptsExchanges: r.Float64() > 0.3, ReturnDeadline: &deadline,
		}

		// Sections
		sectionNames := []string{cat.Items[0], cat.Items[len(cat.Items)/3], cat.Items[len(cat.Items)*2/3]}
		var sectionIDs []int64
		for rank, secName := range sectionNames {
			secID := getID()
			// Use just the item type as section name
			s.ShopSections[secID] = &models.ShopSection{ShopSectionID: secID, Title: secName + "s", Rank: rank + 1, UserID: shopID, ActiveListingCount: 0}
			sectionIDs = append(sectionIDs, secID)
		}

		// Listings
		numListings := randRange(r, cfg.ListingsPerShop.Min, cfg.ListingsPerShop.Max)
		var listingIDs []int64
		shop.ListingActiveCount = 0

		for li := 0; li < numListings; li++ {
			listingID := getID()
			listingIDs = append(listingIDs, listingID)
			itemName := cat.Items[r.Intn(len(cat.Items))]
			mats := cat.Materials[r.Intn(len(cat.Materials))]
			style := cat.Styles[r.Intn(len(cat.Styles))]
			tagSet := cat.Tags[r.Intn(len(cat.Tags))]
			desc := cat.Descriptions[r.Intn(len(cat.Descriptions))]

			material := mats[0]
			title := fmt.Sprintf("%s %s %s", strings.Title(material), itemName, "- "+strings.Title(style))
			state := cfg.ListingStates[r.Intn(len(cfg.ListingStates))]
			if state == "active" {
				shop.ListingActiveCount++
			}

			isDigital := cfg.IncludeDigitalListings && r.Float64() < 0.1
			listingType := "physical"
			if isDigital {
				listingType = "download"
				shop.DigitalListingCount++
			}

			isPersonalizable := cfg.IncludePersonalizedListings && r.Float64() < 0.15
			price := randRange(r, cat.PriceRange[0], cat.PriceRange[1])
			weight := cat.WeightRange[0] + r.Float32()*(cat.WeightRange[1]-cat.WeightRange[0])
			whoMade := "i_did"
			whenMade := "2020_2026"
			lang := "en"
			secID := sectionIDs[r.Intn(len(sectionIDs))]
			spIDPtr := &spID

			if isDigital {
				spIDPtr = nil
			}

			listing := &models.ShopListing{
				ListingID:                 listingID,
				UserID:                    ownerID,
				ShopID:                    shopID,
				Title:                     title,
				Description:               desc,
				State:                     state,
				CreationTimestamp:          ago(randRange(r, 5, 200)),
				CreatedTimestamp:           ago(randRange(r, 5, 200)),
				EndingTimestamp:            now + int64(randRange(r, 30, 120)*86400),
				OriginalCreationTimestamp:  ago(randRange(r, 5, 300)),
				LastModifiedTimestamp:      ago(randRange(r, 0, 30)),
				UpdatedTimestamp:           ago(randRange(r, 0, 30)),
				Quantity:                  randRange(r, 1, 50),
				ShopSectionID:             &secID,
				URL:                       fmt.Sprintf("https://www.etsy.com/listing/%d", listingID),
				NumFavorers:               randRange(r, 0, 200),
				IsTaxable:                 true,
				IsPersonalizable:          isPersonalizable,
				ListingType:               listingType,
				Tags:                      append(tagSet, strings.ToLower(itemName)),
				Materials:                 mats,
				ShippingProfileID:         spIDPtr,
				ReturnPolicyID:            &rpID,
				ProcessingMin:             intPtr(randRange(r, 1, 3)),
				ProcessingMax:             intPtr(randRange(r, 3, 10)),
				WhoMade:                   &whoMade,
				WhenMade:                  &whenMade,
				IsSupply:                  boolPtr(cat.TaxonomyID == 5),
				ItemWeight:                float32Ptr(weight),
				ItemWeightUnit:            &cat.WeightUnit,
				Price:                     models.USD(price),
				TaxonomyID:                &cat.TaxonomyID,
				Language:                  &lang,
				Views:                     randRange(r, 10, 2000),
				Style:                     []string{style},
				ShouldAutoRenew:           true,
			}
			if isPersonalizable {
				listing.PersonalizationIsRequired = r.Float64() > 0.5
				charMax := randRange(r, 20, 200)
				listing.PersonalizationCharCountMax = &charMax
				instructions := "Please enter your customization details"
				listing.PersonalizationInstructions = &instructions
			}
			s.Listings[listingID] = listing

			// Images (2-5 per listing)
			numImages := randRange(r, 2, 5)
			for ii := 0; ii < numImages; ii++ {
				imgID := getID()
				h := randRange(r, 600, 1200)
				w := randRange(r, 600, 1200)
				altText := fmt.Sprintf("Image %d of %s", ii+1, title)
				s.ListingImages[listingID] = append(s.ListingImages[listingID], &models.ListingImage{
					ListingID: listingID, ListingImageID: imgID,
					CreationTsz: listing.CreatedTimestamp, CreatedTimestamp: listing.CreatedTimestamp, Rank: ii + 1,
					URL75x75:     fmt.Sprintf("https://mock.etsy.com/images/%d_75x75.jpg", imgID),
					URL170x135:   fmt.Sprintf("https://mock.etsy.com/images/%d_170x135.jpg", imgID),
					URL570xN:     fmt.Sprintf("https://mock.etsy.com/images/%d_570xN.jpg", imgID),
					URLFullxfull: fmt.Sprintf("https://mock.etsy.com/images/%d_fullxfull.jpg", imgID),
					FullHeight:   &h, FullWidth: &w,
					AltText:      &altText,
				})
			}

			// Files for digital listings
			if isDigital {
				s.ListingFiles[listingID] = []*models.ListingFile{{
					ListingFileID: getID(), ListingID: listingID, Rank: 1,
					Filename: fmt.Sprintf("%s.pdf", strings.ReplaceAll(strings.ToLower(itemName), " ", "-")),
					Filesize: fmt.Sprintf("%.1f MB", 0.5+r.Float64()*10),
					SizeBytes: randRange(r, 500000, 10000000),
					Filetype: "application/pdf",
					CreateTimestamp: listing.CreatedTimestamp,
					CreatedTimestamp: listing.CreatedTimestamp,
				}}
			}
		}

		// Reviews
		numReviews := randRange(r, cfg.ReviewsPerShop.Min, cfg.ReviewsPerShop.Max)
		for ri := 0; ri < numReviews; ri++ {
			buyerIdx := cfg.Shops + r.Intn(len(allUserIDs)-cfg.Shops) // pick from non-shop-owner users
			buyerID := allUserIDs[buyerIdx]
			listingID := listingIDs[r.Intn(len(listingIDs))]
			rating := 3 + r.Intn(3) // 3-5 stars
			if r.Float64() > 0.3 {
				rating = 4 + r.Intn(2) // bias toward 4-5
			}
			reviewText := reviewTexts[r.Intn(len(reviewTexts))]
			ts := ago(randRange(r, 1, 90))
			s.Reviews[shopID] = append(s.Reviews[shopID], &models.ListingReview{
				ShopID: shopID, ListingID: listingID, TransactionID: getID(),
				BuyerUserID: &buyerID, Rating: rating, Review: reviewText,
				Language: "en", CreateTimestamp: ts, CreatedTimestamp: ts,
				UpdateTimestamp: ts, UpdatedTimestamp: ts,
			})
		}

		// Receipts
		numReceipts := randRange(r, cfg.ReceiptsPerShop.Min, cfg.ReceiptsPerShop.Max)
		statuses := []string{"paid", "completed", "completed", "open"}
		for ri := 0; ri < numReceipts; ri++ {
			receiptID := getID()
			buyerIdx := cfg.Shops + r.Intn(len(allUserIDs)-cfg.Shops)
			buyerID := allUserIDs[buyerIdx]
			buyer := s.Users[buyerID]
			buyerAddr := s.UserAddresses[buyerID][0]
			status := statuses[r.Intn(len(statuses))]
			listingID := listingIDs[r.Intn(len(listingIDs))]
			listing := s.Listings[listingID]
			qty := randRange(r, 1, 3)
			itemPrice := listing.Price
			shippingCost := models.USD(randRange(r, 399, 1299))
			subtotal := models.USD(itemPrice.Amount * qty)
			tax := models.USD(subtotal.Amount / 10)
			total := models.USD(subtotal.Amount + shippingCost.Amount + tax.Amount)

			ts := ago(randRange(r, 1, 60))
			paidTs := ts + 3600
			isShipped := status == "completed"
			var shippedTs *int64
			shipments := []models.ShopReceiptShipment{}
			if isShipped {
				st := ts + int64(randRange(r, 86400, 86400*5))
				shippedTs = &st
				shipID := getID()
				shipments = []models.ShopReceiptShipment{{
					ReceiptShippingID: &shipID, ShipmentNotificationTimestamp: st,
					CarrierName: "USPS", TrackingCode: fmt.Sprintf("9400111899223%010d", r.Intn(9999999999)),
				}}
			}

			buyerEmail := ""
			if buyer.PrimaryEmail != nil {
				buyerEmail = *buyer.PrimaryEmail
			}
			ownerUser := s.Users[ownerID]
			sellerEmail := ""
			if ownerUser.PrimaryEmail != nil {
				sellerEmail = *ownerUser.PrimaryEmail
			}
			formattedAddr := fmt.Sprintf("%s\n%s, %s %s\nUnited States", buyerAddr.FirstLine, buyerAddr.City, *buyerAddr.State, *buyerAddr.Zip)

			txnID := getID()
			receipt := &models.ShopReceipt{
				ReceiptID: receiptID, ReceiptType: 0, SellerUserID: ownerID,
				SellerEmail: &sellerEmail, BuyerUserID: buyerID, BuyerEmail: &buyerEmail,
				Name: buyerAddr.Name, FirstLine: &buyerAddr.FirstLine,
				City: &buyerAddr.City, State: buyerAddr.State, Zip: buyerAddr.Zip,
				CountryISO: buyerAddr.ISOCountryCode, Status: status,
				FormattedAddress: &formattedAddr, PaymentMethod: "cc",
				IsPaid: status != "open", IsShipped: isShipped,
				CreateTimestamp: ts, CreatedTimestamp: ts,
				UpdateTimestamp: ts + 3600, UpdatedTimestamp: ts + 3600,
				Grandtotal: total, Subtotal: subtotal, TotalPrice: subtotal,
				TotalShippingCost: shippingCost, TotalTaxCost: tax,
				TotalVatCost: models.USD(0), DiscountAmt: models.USD(0), GiftWrapPrice: models.USD(0),
				Shipments: shipments,
				Transactions: []models.ShopReceiptTransaction{{
					TransactionID: txnID, Title: &listing.Title,
					SellerUserID: ownerID, BuyerUserID: buyerID,
					CreateTimestamp: ts, CreatedTimestamp: ts,
					PaidTimestamp: &paidTs, ShippedTimestamp: shippedTs,
					Quantity: qty, ReceiptID: receiptID,
					ListingID: intPtrFromInt64(listingID),
					TransactionType: "listing", Price: itemPrice, ShippingCost: shippingCost,
					Variations: []models.TransactionVariation{},
					ShippingProfileID: &spID,
					MinProcessingDays: intPtr(1), MaxProcessingDays: intPtr(5),
				}},
				Refunds: []models.ShopRefund{},
			}
			s.Receipts[receiptID] = receipt
			s.Transactions[txnID] = &receipt.Transactions[0]

			// Payment
			fees := models.USD(total.Amount * 65 / 1000) // ~6.5% Etsy fee
			net := models.USD(total.Amount - fees.Amount)
			shopCurrency := "USD"
			payStatus := "open"
			if status == "completed" {
				payStatus = "settled"
			}
			s.Payments[getID()] = &models.Payment{
				PaymentID: getID(), BuyerUserID: buyerID, ShopID: shopID, ReceiptID: receiptID,
				AmountGross: total, AmountFees: fees, AmountNet: net,
				Currency: "USD", ShopCurrency: &shopCurrency, BuyerCurrency: &shopCurrency,
				ShippingAddressID: buyerAddr.UserAddressID,
				Status: payStatus,
				CreateTimestamp: ts, CreatedTimestamp: ts,
				UpdateTimestamp: ts + 3600, UpdatedTimestamp: ts + 3600,
				PaymentAdjustments: []models.PaymentAdjustment{},
			}
		}

		// Ledger entries
		for _, receipt := range s.Receipts {
			if receipt.SellerUserID == ownerID && receipt.IsPaid {
				s.LedgerEntries[shopID] = append(s.LedgerEntries[shopID], &models.PaymentAccountLedgerEntry{
					EntryID: getID(), LedgerID: shopID, SequenceNumber: len(s.LedgerEntries[shopID]) + 1,
					Amount: receipt.Grandtotal.Amount, Currency: "USD",
					Description: fmt.Sprintf("Sale: %s", *receipt.Transactions[0].Title),
					Balance: receipt.Grandtotal.Amount,
					CreateDate: receipt.CreateTimestamp, CreatedTimestamp: receipt.CreateTimestamp,
					LedgerType: "credit", ReferenceType: "receipt",
					ReferenceID: strPtr(fmt.Sprintf("%d", receipt.ReceiptID)),
					PaymentAdjustments: []models.PaymentAdjustment{},
				})
			}
		}
	}
}

func loadTaxonomy(s *store.Store) {
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
	}
	s.TaxonomyProperties[562] = []models.BuyerTaxonomyNodeProperty{
		{PropertyID: 200, Name: "material", DisplayName: "Wood Type", IsRequired: false, SupportsAttributes: true, SupportsVariations: false, IsMultivalued: true, PossibleValues: []models.BuyerTaxonomyPropertyValue{{ValueID: 3001, Name: "Walnut"}, {ValueID: 3002, Name: "Maple"}, {ValueID: 3003, Name: "Cherry"}, {ValueID: 3004, Name: "Oak"}}, SelectedValues: []models.BuyerTaxonomyPropertyValue{}, Scales: []models.BuyerTaxonomyPropertyScale{}},
	}
}

func generateShopName(r *rand.Rand) string {
	prefix := shopPrefixes[r.Intn(len(shopPrefixes))]
	noun := shopNouns[r.Intn(len(shopNouns))]
	suffix := shopSuffixes[r.Intn(len(shopSuffixes))]
	if suffix != "" {
		return prefix + noun + suffix
	}
	return prefix + noun
}

func pick(r *rand.Rand, pool []string) string {
	return pool[r.Intn(len(pool))]
}

func randRange(r *rand.Rand, min, max int) int {
	if min >= max {
		return min
	}
	return min + r.Intn(max-min+1)
}

func intPtrFromInt64(v int64) *int {
	i := int(v)
	return &i
}
