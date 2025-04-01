package models

// AppResponse represents the response for an app resource
type AppResponse struct {
	Data  App       `json:"data"`
	Links Links     `json:"links"`
}

// AppsResponse represents the response for multiple app resources
type AppsResponse struct {
	Data  []App    `json:"data"`
	Links Links     `json:"links"`
	Meta  struct {
		Paging PagingInformation `json:"paging"`
	} `json:"meta,omitempty"`
}

// App represents an app resource
type App struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes struct {
		Name                    string         `json:"name"`
		BundleID                string         `json:"bundleId"`
		SKU                     string         `json:"sku"`
		PrimaryLocale           string         `json:"primaryLocale"`
		IsPreReleaseApp         bool           `json:"isPreReleaseApp"`
		Prices                  []AppPrice     `json:"prices,omitempty"`
		AvailableInNewTerritories bool         `json:"availableInNewTerritories"`
		ContentRightsDeclaration string         `json:"contentRightsDeclaration,omitempty"`
	} `json:"attributes"`
	Relationships struct {
		AppInfos struct {
			Links ResourceLinks `json:"links"`
		} `json:"appInfos"`
		AppStoreVersions struct {
			Links ResourceLinks `json:"links"`
		} `json:"appStoreVersions"`
		PreReleaseVersions struct {
			Links ResourceLinks `json:"links"`
		} `json:"preReleaseVersions"`
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

// AppPrice represents a price for an app
type AppPrice struct {
	Territory string `json:"territory"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
}

// AppInfo represents additional information about an app
type AppInfo struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Attributes struct {
		AppStoreState           string `json:"appStoreState"`
		AppStoreAgeRating       string `json:"appStoreAgeRating"`
		BrazilAgeRating         string `json:"brazilAgeRating,omitempty"`
		KidsAgeBand            string `json:"kidsAgeBand,omitempty"`
		PrimaryCategory         string `json:"primaryCategory"`
		PrimaryCategorySubcategories []string `json:"primaryCategorySubcategories,omitempty"`
		SecondaryCategory       string `json:"secondaryCategory,omitempty"`
		SecondaryCategorySubcategories []string `json:"secondaryCategorySubcategories,omitempty"`
	} `json:"attributes"`
	Relationships struct {
		App struct {
			Links ResourceLinks `json:"links"`
		} `json:"app"`
	} `json:"relationships"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}
