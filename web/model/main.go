package model

type ParentRelationTraitType string

const (
	ParentNftRelationBackround ParentRelationTraitType = "background"
	ParentNftRelationItem1     ParentRelationTraitType = "item_1"
	ParentNftRelationItem2     ParentRelationTraitType = "item_2"

	ParentNftRelationBackroundID int64 = 1
	ParentNftRelationItem1ID     int64 = 2
	ParentNftRelationItem2ID     int64 = 3
)

type NftTokenInfo struct {
	TokenID                int64                   `json:"token_id"`
	Title                  string                  `json:"title"`
	Description            string                  `json:"description"`
	Thumbnail              string                  `json:"thumbnail"`
	TransparentThumbnail              string                  
	Image                  string                  `json:"image"`
	ContractAddress        string                  `json:"contract_address"`
	Balance                int64                   `json:"balance"`
	BalanceOfCurrentWallet int64                   `json:"balance_of_current_wallet"`
	ParentRelationTrait    ParentRelationTraitType `json:"parent_relation_trait,omitempty"`
	MetaData               map[string]string       `json:"metadata"`

	ParentTokenID          int64
	ParentRelationTraitID  int64
	IsOWnedByCurrentWallet bool

	//only relevant for the main NFTs
	AttachedRelationTraits map[ParentRelationTraitType]int64
}

// TODO
type MainNftTokenInfo struct {
	TokenID         string
	Title           string
	Thumbnail       string
	Image           string
	ContractAddress string
	ExtraData       map[string]string
}

// TODO
type SecondaryNftTokenInfo struct {
	TokenID         string
	ParentTokenID   string
	Title           string
	Thumbnail       string
	Image           string
	ContractAddress string
	TraitID         int
	TraitKey        int
	TraitValue      string
	ExtraData       map[string]string
}
