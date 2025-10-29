package service

import (
	"fmt"
	"strconv"

	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func GetCumulativeClientStatistic(ctx *gin.Context, endAt string) (*response.ClientCumulative, error) {
	resp, err := operate.GetCumulativeClient(ctx, &operate_service.GetCumulativeClientReq{
		EndAt: util.MustI64(endAt),
	})
	if err != nil {
		return nil, err
	}
	return &response.ClientCumulative{Total: resp.Total}, nil
}

func GetClientStatistic(ctx *gin.Context, startDate, endDate string) (*response.ClientStatistic, error) {
	overview, err := getClientStatisticOverview(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	trend, err := getClientStatisticTrend(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	//i18n替换表名
	trend.Client.TableName = gin_util.I18nKey(ctx, trend.Client.TableName)
	for i, line := range trend.Client.Lines {
		trend.Client.Lines[i].LineName = gin_util.I18nKey(ctx, line.LineName)
	}

	return &response.ClientStatistic{
		Overview: *overview,
		Trend:    *trend,
	}, nil
}

func getClientStatisticOverview(ctx *gin.Context, startDate, endDate string) (*response.ClientOverView, error) {
	resp, err := operate.GetClientOverview(ctx, &operate_service.GetClientOverviewReq{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}
	return &response.ClientOverView{
		NewClient:    clientOverviewPb2resp(resp.NewClient),
		ActiveClient: clientOverviewPb2resp(resp.ActiveClient),
	}, nil
}

func getClientStatisticTrend(ctx *gin.Context, startDate, endDate string) (*response.ClientTrends, error) {
	// resp, err := operate.GetClientTrend(ctx, &operate_service.GetClientTrendReq{
	// 	StartDate: startDate,
	// 	EndDate:   endDate,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return &response.ClientTrends{
	// 	Client: convertStatisticChart(resp.Client),
	// }, nil
	return &response.ClientTrends{}, nil
}

func clientOverviewPb2resp(item *operate_service.ClientOverviewItem) response.ClientOverviewItem {
	valueStr := fmt.Sprintf("%.2f", item.Value)
	value, _ := strconv.ParseFloat(valueStr, 64)
	return response.ClientOverviewItem{
		Value:            float32(value),
		PeriodOverPeriod: item.PeriodOverperiod,
	}
}
