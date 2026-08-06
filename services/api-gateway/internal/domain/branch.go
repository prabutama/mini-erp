package domain

type CreateBranchRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UpdateBranchRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Status  string `json:"status"`
}

type BranchResponse struct {
	BranchID   string `json:"branch_id"`
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Status     string `json:"status"`
}

type ListBranchesResponse struct {
	Branches []BranchResponse `json:"branches"`
}
