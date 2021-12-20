package swan

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

const GET_OFFLINEDEAL_LIMIT_DEFAULT = 50

type UpdateOfflineDealResponse struct {
	Data   UpdateOfflineDealData `json:"data"`
	Status string                `json:"status"`
}

type UpdateOfflineDealData struct {
	Deal    model.OfflineDeal `json:"deal"`
	Message string            `json:"message"`
}

type GetOfflineDealsByStatusParams struct {
	DealStatus string  `json:"status"`
	MinerFid   *string `json:"miner_fid"`
	PageNum    *int    `json:"page_num"`
	PageSize   *int    `json:"page_size"`
}

type GetOfflineDealResponse struct {
	Data   GetOfflineDealData `json:"data"`
	Status string             `json:"status"`
}

type GetOfflineDealData struct {
	Deal []model.OfflineDeal `json:"deal"`
}

func (swanClient *SwanClient) GetOfflineDealsByStatus(params GetOfflineDealsByStatusParams) ([]model.OfflineDeal, error) {
	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	urlStr := utils.UrlJoin(swanClient.ApiUrl, "offline_deals/get_by_status")
	response := web.HttpGet(urlStr, swanClient.SwanToken, params)
	getOfflineDealResponse := GetOfflineDealResponse{}
	err = json.Unmarshal([]byte(response), &getOfflineDealResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(getOfflineDealResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("get offline deal with status:%s failed", params.DealStatus)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return getOfflineDealResponse.Data.Deal, nil
}

type UpdateOfflineDealParams struct {
	Id         int     `json:"id"`
	DealCid    *string `json:"deal_cid"`
	FilePath   *string `json:"file_path"`
	Status     string  `json:"status"`
	StartEpoch *int    `json:"start_epoch"`
	Note       *string `json:"note"`
}

func (swanClient *SwanClient) UpdateOfflineDeal(dealId int, params UpdateOfflineDealParams) bool {
	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return false
	}

	if len(params.Status) == 0 {
		logs.GetLogger().Error("Please provide status")
		return false
	}

	apiUrl := swanClient.ApiUrl + "offline_deals/update"

	response := web.HttpPut(apiUrl, swanClient.SwanToken, params)

	updateOfflineDealResponse := &UpdateOfflineDealResponse{}
	err = json.Unmarshal([]byte(response), updateOfflineDealResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return false
	}

	if !strings.EqualFold(updateOfflineDealResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("deal(id=%d),failed to update offline deal status to %s,%s", dealId, params.Status, updateOfflineDealResponse.Data.Message)
		logs.GetLogger().Error(err)
		return false
	}

	return true
}
