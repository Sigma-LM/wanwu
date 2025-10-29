package operate

import (
	"context"

	"github.com/UnicomAI/wanwu/api/proto/common"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/operate-service/client/orm"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) AddClientRecord(ctx context.Context, req *operate_service.AddClientRecordReq) (*emptypb.Empty, error) {
	if err := s.cli.AddClientRecord(ctx, req.ClientId); err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetClientOverview(ctx context.Context, req *operate_service.GetClientOverviewReq) (*operate_service.ClientOverViewInfo, error) {
	stats, err := s.cli.GetClientOverview(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return toClientOverviewInfo(stats), nil
}

func (s *Service) GetClientTrend(ctx context.Context, req *operate_service.GetClientTrendReq) (*operate_service.ClientTrendInfo, error) {
	trend, err := s.cli.GetClientTrend(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &operate_service.ClientTrendInfo{
		Client: convertStatisticChart(trend),
	}, nil
}

func (s *Service) GetCumulativeClient(ctx context.Context, req *operate_service.GetCumulativeClientReq) (*operate_service.GetCumulativeClientResp, error) {
	total, err := s.cli.GetCumulativeClient(ctx, req.EndAt)
	if err != nil {
		return nil, errStatus(errs.Code_OperateRecord, err)
	}
	return &operate_service.GetCumulativeClientResp{Total: total}, nil
}

// --- internal ---

func toClientOverviewInfo(stats *orm.ClientOverView) *operate_service.ClientOverViewInfo {
	ret := &operate_service.ClientOverViewInfo{
		ActiveClient: &operate_service.ClientOverviewItem{
			Value:            stats.ActiveClient.Value,
			PeriodOverperiod: stats.ActiveClient.PeriodOverPeriod,
		},
		NewClient: &operate_service.ClientOverviewItem{
			Value:            stats.NewClient.Value,
			PeriodOverperiod: stats.NewClient.PeriodOverPeriod,
		},
	}
	return ret
}

func convertStatisticChart(trend *orm.ClientTrends) *common.StatisticChart {
	pbChart := &common.StatisticChart{
		TableName:  trend.Client.TableName,
		ChartLines: make([]*common.StatisticChartLine, 0, len(trend.Client.Lines)),
	}
	for _, respLine := range trend.Client.Lines {
		pbLine := &common.StatisticChartLine{
			LineName: respLine.LineName,
			Items:    make([]*common.StatisticChartLineItem, 0, len(respLine.Items)),
		}
		for _, respItem := range respLine.Items {
			pbLine.Items = append(pbLine.Items, &common.StatisticChartLineItem{
				Key:   respItem.Key,
				Value: respItem.Value,
			})
		}
		pbChart.ChartLines = append(pbChart.ChartLines, pbLine)
	}
	return pbChart
}
