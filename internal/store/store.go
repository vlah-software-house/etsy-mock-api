package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

type Store struct {
	mu sync.RWMutex

	Shops             map[int64]*models.Shop
	ShopSections      map[int64]*models.ShopSection
	ShopReturnPolicies map[int64]*models.ShopReturnPolicy
	Listings          map[int64]*models.ShopListing
	ListingImages     map[int64][]*models.ListingImage
	ListingFiles      map[int64][]*models.ListingFile
	Receipts          map[int64]*models.ShopReceipt
	Transactions      map[int64]*models.ShopReceiptTransaction
	Payments          map[int64]*models.Payment
	Users             map[int64]*models.User
	UserAddresses     map[int64][]*models.UserAddress
	Reviews           map[int64][]*models.ListingReview // keyed by shop_id
	ShippingProfiles  map[int64]*models.ShopShippingProfile
	LedgerEntries     map[int64][]*models.PaymentAccountLedgerEntry // keyed by shop_id
	TaxonomyNodes     []models.BuyerTaxonomyNode
	TaxonomyProperties map[int64][]models.BuyerTaxonomyNodeProperty // keyed by taxonomy_id

	nextID int64
}

func New() *Store {
	return &Store{
		Shops:              make(map[int64]*models.Shop),
		ShopSections:       make(map[int64]*models.ShopSection),
		ShopReturnPolicies: make(map[int64]*models.ShopReturnPolicy),
		Listings:           make(map[int64]*models.ShopListing),
		ListingImages:      make(map[int64][]*models.ListingImage),
		ListingFiles:       make(map[int64][]*models.ListingFile),
		Receipts:           make(map[int64]*models.ShopReceipt),
		Transactions:       make(map[int64]*models.ShopReceiptTransaction),
		Payments:           make(map[int64]*models.Payment),
		Users:              make(map[int64]*models.User),
		UserAddresses:      make(map[int64][]*models.UserAddress),
		Reviews:            make(map[int64][]*models.ListingReview),
		ShippingProfiles:   make(map[int64]*models.ShopShippingProfile),
		LedgerEntries:      make(map[int64][]*models.PaymentAccountLedgerEntry),
		TaxonomyProperties: make(map[int64][]models.BuyerTaxonomyNodeProperty),
		nextID:             10000,
	}
}

func (s *Store) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	return s.nextID
}

func now() int64 {
	return time.Now().Unix()
}

// Shop operations

func (s *Store) GetShop(shopID int64) (*models.Shop, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	shop, ok := s.Shops[shopID]
	return shop, ok
}

func (s *Store) GetShopByName(name string) (*models.Shop, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, shop := range s.Shops {
		if shop.ShopName == name {
			return shop, true
		}
	}
	return nil, false
}

func (s *Store) UpdateShop(shop *models.Shop) {
	s.mu.Lock()
	defer s.mu.Unlock()
	shop.UpdatedTimestamp = now()
	shop.UpdateDate = now()
	s.Shops[shop.ShopID] = shop
}

// Shop Section operations

func (s *Store) GetShopSections(shopID int64) []models.ShopSection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var sections []models.ShopSection
	for _, sec := range s.ShopSections {
		if sec.UserID == shopID {
			sections = append(sections, *sec)
		}
	}
	return sections
}

func (s *Store) GetShopSection(shopID, sectionID int64) (*models.ShopSection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sec, ok := s.ShopSections[sectionID]
	if ok && sec.UserID == shopID {
		return sec, true
	}
	return nil, false
}

func (s *Store) CreateShopSection(shopID int64, title string, rank int) *models.ShopSection {
	s.mu.Lock()
	defer s.mu.Unlock()
	sec := &models.ShopSection{
		ShopSectionID:      s.nextID + 1,
		Title:              title,
		Rank:               rank,
		UserID:             shopID,
		ActiveListingCount: 0,
	}
	s.nextID++
	s.ShopSections[sec.ShopSectionID] = sec
	return sec
}

// Return Policy operations

func (s *Store) GetReturnPolicies(shopID int64) []models.ShopReturnPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var policies []models.ShopReturnPolicy
	for _, p := range s.ShopReturnPolicies {
		if p.ShopID == shopID {
			policies = append(policies, *p)
		}
	}
	return policies
}

func (s *Store) GetReturnPolicy(policyID int64) (*models.ShopReturnPolicy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.ShopReturnPolicies[policyID]
	return p, ok
}

func (s *Store) CreateReturnPolicy(shopID int64, acceptsReturns, acceptsExchanges bool, deadline *int) *models.ShopReturnPolicy {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := &models.ShopReturnPolicy{
		ReturnPolicyID:   s.nextID + 1,
		ShopID:           shopID,
		AcceptsReturns:   acceptsReturns,
		AcceptsExchanges: acceptsExchanges,
		ReturnDeadline:   deadline,
	}
	s.nextID++
	s.ShopReturnPolicies[p.ReturnPolicyID] = p
	return p
}

// Listing operations

func (s *Store) GetListing(listingID int64) (*models.ShopListing, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.Listings[listingID]
	return l, ok
}

func (s *Store) GetShopListings(shopID int64, state string, limit, offset int) ([]models.ShopListing, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []models.ShopListing
	for _, l := range s.Listings {
		if l.ShopID == shopID {
			if state == "" || l.State == state {
				all = append(all, *l)
			}
		}
	}
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total
}

func (s *Store) GetActiveListings(keyword string, taxonomyID *int, limit, offset int, sortOn, sortOrder string) ([]models.ShopListing, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []models.ShopListing
	for _, l := range s.Listings {
		if l.State != "active" {
			continue
		}
		if keyword != "" {
			kw := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(l.Title), kw) &&
				!strings.Contains(strings.ToLower(l.Description), kw) &&
				!containsTag(l.Tags, kw) {
				continue
			}
		}
		if taxonomyID != nil && (l.TaxonomyID == nil || *l.TaxonomyID != *taxonomyID) {
			continue
		}
		all = append(all, *l)
	}
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total
}

func containsTag(tags []string, keyword string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), keyword) {
			return true
		}
	}
	return false
}

func (s *Store) CreateListing(shopID int64, req models.CreateListingRequest) *models.ShopListing {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID + 1
	s.nextID++
	ts := now()

	shop := s.Shops[shopID]
	var userID int64
	if shop != nil {
		userID = shop.UserID
	}

	taxID := req.TaxonomyID
	listing := &models.ShopListing{
		ListingID:                 id,
		UserID:                    userID,
		ShopID:                    shopID,
		Title:                     req.Title,
		Description:               req.Description,
		State:                     "draft",
		CreationTimestamp:          ts,
		CreatedTimestamp:           ts,
		EndingTimestamp:            ts + 60*60*24*120,
		OriginalCreationTimestamp:  ts,
		LastModifiedTimestamp:      ts,
		UpdatedTimestamp:           ts,
		Quantity:                  req.Quantity,
		URL:                       fmt.Sprintf("https://www.etsy.com/listing/%d", id),
		ListingType:               req.ListingType,
		Tags:                      req.Tags,
		Materials:                 req.Materials,
		ShippingProfileID:         req.ShippingProfileID,
		ReturnPolicyID:            req.ReturnPolicyID,
		WhoMade:                   &req.WhoMade,
		WhenMade:                  &req.WhenMade,
		IsSupply:                  req.IsSupply,
		IsCustomizable:            req.IsCustomizable,
		IsPersonalizable:          req.IsPersonalizable,
		ItemWeight:                req.ItemWeight,
		ItemWeightUnit:            req.ItemWeightUnit,
		ItemLength:                req.ItemLength,
		ItemWidth:                 req.ItemWidth,
		ItemHeight:                req.ItemHeight,
		ItemDimensionsUnit:        req.ItemDimensionsUnit,
		Price:                     models.USD(int(req.Price * 100)),
		TaxonomyID:                &taxID,
		Style:                     req.Styles,
		IsTaxable:                 true,
	}
	if listing.ListingType == "" {
		listing.ListingType = "physical"
	}
	if listing.Tags == nil {
		listing.Tags = []string{}
	}
	if listing.Materials == nil {
		listing.Materials = []string{}
	}
	if listing.Style == nil {
		listing.Style = []string{}
	}
	s.Listings[id] = listing
	return listing
}

func (s *Store) UpdateListing(listingID int64, req models.UpdateListingRequest) (*models.ShopListing, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, ok := s.Listings[listingID]
	if !ok {
		return nil, false
	}
	if req.Title != nil {
		l.Title = *req.Title
	}
	if req.Description != nil {
		l.Description = *req.Description
	}
	if req.Quantity != nil {
		l.Quantity = *req.Quantity
	}
	if req.Price != nil {
		l.Price = models.USD(int(*req.Price * 100))
	}
	if req.State != nil {
		l.State = *req.State
	}
	if req.WhoMade != nil {
		l.WhoMade = req.WhoMade
	}
	if req.WhenMade != nil {
		l.WhenMade = req.WhenMade
	}
	if req.TaxonomyID != nil {
		l.TaxonomyID = req.TaxonomyID
	}
	if req.ShippingProfileID != nil {
		l.ShippingProfileID = req.ShippingProfileID
	}
	if req.ReturnPolicyID != nil {
		l.ReturnPolicyID = req.ReturnPolicyID
	}
	if req.IsSupply != nil {
		l.IsSupply = req.IsSupply
	}
	if req.ItemWeight != nil {
		l.ItemWeight = req.ItemWeight
	}
	if req.ItemWeightUnit != nil {
		l.ItemWeightUnit = req.ItemWeightUnit
	}
	if req.ItemLength != nil {
		l.ItemLength = req.ItemLength
	}
	if req.ItemWidth != nil {
		l.ItemWidth = req.ItemWidth
	}
	if req.ItemHeight != nil {
		l.ItemHeight = req.ItemHeight
	}
	if req.ItemDimensionsUnit != nil {
		l.ItemDimensionsUnit = req.ItemDimensionsUnit
	}
	if req.Tags != nil {
		l.Tags = req.Tags
	}
	if req.Materials != nil {
		l.Materials = req.Materials
	}
	l.LastModifiedTimestamp = now()
	l.UpdatedTimestamp = now()
	return l, true
}

func (s *Store) DeleteListing(listingID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Listings[listingID]; !ok {
		return false
	}
	delete(s.Listings, listingID)
	delete(s.ListingImages, listingID)
	delete(s.ListingFiles, listingID)
	return true
}

// Listing Image operations

func (s *Store) GetListingImages(listingID int64) []models.ListingImage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	imgs := s.ListingImages[listingID]
	result := make([]models.ListingImage, len(imgs))
	for i, img := range imgs {
		result[i] = *img
	}
	return result
}

func (s *Store) AddListingImage(listingID int64, altText string, rank int) *models.ListingImage {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID + 1
	s.nextID++
	ts := now()
	img := &models.ListingImage{
		ListingID:       listingID,
		ListingImageID:  id,
		CreationTsz:     ts,
		CreatedTimestamp: ts,
		Rank:            rank,
		URL75x75:        fmt.Sprintf("https://mock.etsy.com/images/%d_75x75.jpg", id),
		URL170x135:      fmt.Sprintf("https://mock.etsy.com/images/%d_170x135.jpg", id),
		URL570xN:        fmt.Sprintf("https://mock.etsy.com/images/%d_570xN.jpg", id),
		URLFullxfull:    fmt.Sprintf("https://mock.etsy.com/images/%d_fullxfull.jpg", id),
		AltText:         &altText,
	}
	s.ListingImages[listingID] = append(s.ListingImages[listingID], img)
	return img
}

func (s *Store) DeleteListingImage(listingID, imageID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	imgs := s.ListingImages[listingID]
	for i, img := range imgs {
		if img.ListingImageID == imageID {
			s.ListingImages[listingID] = append(imgs[:i], imgs[i+1:]...)
			return true
		}
	}
	return false
}

// Listing File operations

func (s *Store) GetListingFiles(listingID int64) []models.ListingFile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	files := s.ListingFiles[listingID]
	result := make([]models.ListingFile, len(files))
	for i, f := range files {
		result[i] = *f
	}
	return result
}

func (s *Store) GetListingFile(listingID, fileID int64) (*models.ListingFile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, f := range s.ListingFiles[listingID] {
		if f.ListingFileID == fileID {
			return f, true
		}
	}
	return nil, false
}

func (s *Store) AddListingFile(listingID int64, filename, filetype string, sizeBytes int) *models.ListingFile {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID + 1
	s.nextID++
	ts := now()
	f := &models.ListingFile{
		ListingFileID:   id,
		ListingID:       listingID,
		Rank:            len(s.ListingFiles[listingID]) + 1,
		Filename:        filename,
		Filesize:        fmt.Sprintf("%d bytes", sizeBytes),
		SizeBytes:       sizeBytes,
		Filetype:        filetype,
		CreateTimestamp: ts,
		CreatedTimestamp: ts,
	}
	s.ListingFiles[listingID] = append(s.ListingFiles[listingID], f)
	return f
}

func (s *Store) DeleteListingFile(listingID, fileID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := s.ListingFiles[listingID]
	for i, f := range files {
		if f.ListingFileID == fileID {
			s.ListingFiles[listingID] = append(files[:i], files[i+1:]...)
			return true
		}
	}
	return false
}

// Receipt operations

func (s *Store) GetReceipt(receiptID int64) (*models.ShopReceipt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.Receipts[receiptID]
	return r, ok
}

func (s *Store) GetShopReceipts(shopID int64, limit, offset int) ([]models.ShopReceipt, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []models.ShopReceipt
	for _, r := range s.Receipts {
		if r.SellerUserID == shopID || s.shopOwnsReceipt(shopID, r) {
			all = append(all, *r)
		}
	}
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total
}

func (s *Store) shopOwnsReceipt(shopID int64, r *models.ShopReceipt) bool {
	for _, shop := range s.Shops {
		if shop.ShopID == shopID && shop.UserID == r.SellerUserID {
			return true
		}
	}
	return false
}

func (s *Store) UpdateReceipt(receipt *models.ShopReceipt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	receipt.UpdateTimestamp = now()
	receipt.UpdatedTimestamp = now()
	s.Receipts[receipt.ReceiptID] = receipt
}

// Transaction operations

func (s *Store) GetTransaction(transactionID int64) (*models.ShopReceiptTransaction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.Transactions[transactionID]
	return t, ok
}

func (s *Store) GetReceiptTransactions(receiptID int64) []models.ShopReceiptTransaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var txns []models.ShopReceiptTransaction
	for _, t := range s.Transactions {
		if t.ReceiptID == receiptID {
			txns = append(txns, *t)
		}
	}
	return txns
}

func (s *Store) GetShopTransactions(shopID int64, limit, offset int) ([]models.ShopReceiptTransaction, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []models.ShopReceiptTransaction
	for _, t := range s.Transactions {
		if t.SellerUserID == shopID {
			all = append(all, *t)
		}
	}
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total
}

// Payment operations

func (s *Store) GetPaymentsByReceipt(receiptID int64) []models.Payment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var payments []models.Payment
	for _, p := range s.Payments {
		if p.ReceiptID == receiptID {
			payments = append(payments, *p)
		}
	}
	return payments
}

func (s *Store) GetPaymentsByShop(shopID int64) []models.Payment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var payments []models.Payment
	for _, p := range s.Payments {
		if p.ShopID == shopID {
			payments = append(payments, *p)
		}
	}
	return payments
}

// User operations

func (s *Store) GetUser(userID int64) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.Users[userID]
	return u, ok
}

func (s *Store) GetUserAddresses(userID int64) []models.UserAddress {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addrs := s.UserAddresses[userID]
	result := make([]models.UserAddress, len(addrs))
	for i, a := range addrs {
		result[i] = *a
	}
	return result
}

// Review operations

func (s *Store) GetShopReviews(shopID int64, limit, offset int) ([]models.ListingReview, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := s.Reviews[shopID]
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]models.ListingReview, end-offset)
	for i, r := range all[offset:end] {
		result[i] = *r
	}
	return result, total
}

func (s *Store) GetListingReviews(listingID int64, limit, offset int) ([]models.ListingReview, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []*models.ListingReview
	for _, reviews := range s.Reviews {
		for _, r := range reviews {
			if r.ListingID == listingID {
				all = append(all, r)
			}
		}
	}
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]models.ListingReview, end-offset)
	for i, r := range all[offset:end] {
		result[i] = *r
	}
	return result, total
}

// Shipping Profile operations

func (s *Store) GetShippingProfile(profileID int64) (*models.ShopShippingProfile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.ShippingProfiles[profileID]
	return p, ok
}

func (s *Store) GetShopShippingProfiles(shopID int64) []models.ShopShippingProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var profiles []models.ShopShippingProfile
	for _, p := range s.ShippingProfiles {
		shop := s.Shops[shopID]
		if shop != nil && p.UserID == shop.UserID {
			profiles = append(profiles, *p)
		}
	}
	return profiles
}

func (s *Store) CreateShippingProfile(userID int64, title, originCountryISO string, profileType string) *models.ShopShippingProfile {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID + 1
	s.nextID++
	p := &models.ShopShippingProfile{
		ShippingProfileID:           id,
		Title:                       &title,
		UserID:                      userID,
		OriginCountryISO:            originCountryISO,
		ProfileType:                 profileType,
		ShippingProfileDestinations: []models.ShopShippingProfileDestination{},
		ShippingProfileUpgrades:     []models.ShopShippingProfileUpgrade{},
	}
	s.ShippingProfiles[id] = p
	return p
}

func (s *Store) DeleteShippingProfile(profileID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ShippingProfiles[profileID]; !ok {
		return false
	}
	p := s.ShippingProfiles[profileID]
	p.IsDeleted = true
	return true
}

// Ledger operations

func (s *Store) GetLedgerEntries(shopID int64, limit, offset int) ([]models.PaymentAccountLedgerEntry, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := s.LedgerEntries[shopID]
	total := len(all)
	if offset >= total {
		return nil, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]models.PaymentAccountLedgerEntry, end-offset)
	for i, e := range all[offset:end] {
		result[i] = *e
	}
	return result, total
}

// Taxonomy operations

func (s *Store) GetTaxonomyNodes() []models.BuyerTaxonomyNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.TaxonomyNodes
}

func (s *Store) GetTaxonomyProperties(taxonomyID int64) ([]models.BuyerTaxonomyNodeProperty, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	props, ok := s.TaxonomyProperties[taxonomyID]
	return props, ok
}

// Listing Inventory operations

func (s *Store) GetListingInventory(listingID int64) (*models.ListingInventory, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.Listings[listingID]
	if !ok {
		return nil, false
	}
	if l.Inventory != nil {
		return l.Inventory, true
	}
	// Return a default inventory based on the listing
	return &models.ListingInventory{
		Products: []models.ListingInventoryProduct{
			{
				ProductID: listingID * 10,
				SKU:       "",
				IsDeleted: false,
				Offerings: []models.ListingInventoryProductOffering{
					{
						OfferingID: listingID * 100,
						Quantity:   l.Quantity,
						IsEnabled:  true,
						IsDeleted:  false,
						Price:      l.Price,
					},
				},
				PropertyValues: []models.ListingPropertyValue{},
			},
		},
		PriceOnProperty:    []int{},
		QuantityOnProperty: []int{},
		SKUOnProperty:      []int{},
	}, true
}
