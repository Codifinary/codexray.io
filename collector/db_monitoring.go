package collector

import (
	"codexray/generated/proto/dbmonitoring/v3/common"
	"codexray/generated/proto/dbmonitoring/v3/dbm"
	language_agent "codexray/generated/proto/dbmonitoring/v3/language-agent"
	"codexray/generated/proto/dbmonitoring/v3/management"
	"codexray/utils"
	"context"
	"log"

	"google.golang.org/grpc"
)

type managementServer struct {
	management.UnimplementedManagementServiceServer
}

func (s *managementServer) ReportInstanceProperties(ctx context.Context, req *management.InstanceProperties) (*common.Commands, error) {
	log.Printf("Received instance properties report: service=%s, instance=%s", req.Service, req.ServiceInstance)
	if err := utils.SaveToJSON("management", req); err != nil {
		log.Printf("Error saving instance properties: %v", err)
	}
	return &common.Commands{}, nil
}

func (s *managementServer) KeepAlive(ctx context.Context, req *management.InstancePingPkg) (*common.Commands, error) {
	log.Printf("Received keep alive: service=%s, instance=%s", req.Service, req.ServiceInstance)
	if err := utils.SaveToJSON("management", req); err != nil {
		log.Printf("Error saving keep alive: %v", err)
	}
	return &common.Commands{}, nil
}

type dbmEventServer struct {
	dbm.UnimplementedDbmEventServiceServer
}

func (s *dbmEventServer) CollectEvent(stream dbm.DbmEventService_CollectEventServer) error {
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving event: %v", err)
			return err
		}
		log.Printf("Received DB event: %+v", event)
		if err := utils.SaveToJSON("dbm_event", event); err != nil {
			log.Printf("Error saving DB event: %v", err)
		}
	}
}

type dbmQueryServer struct {
	dbm.UnimplementedDbmQueryServiceServer
}

func (s *dbmQueryServer) CollectQuery(stream dbm.DbmQueryService_CollectQueryServer) error {
	for {
		query, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving query: %v", err)
			return err
		}
		log.Printf("Received DB query: %+v", query)
		if err := utils.SaveToJSON("dbm_query", query); err != nil {
			log.Printf("Error saving DB query: %v", err)
		}
	}
}

type meterReportServer struct {
	language_agent.UnimplementedMeterReportServiceServer
}

func (s *meterReportServer) Collect(stream language_agent.MeterReportService_CollectServer) error {
	for {
		meterData, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving meter data: %v", err)
			return err
		}
		log.Printf("Received meter data: %+v", meterData)
		if err := utils.SaveToJSON("meter", meterData); err != nil {
			log.Printf("Error saving meter data: %v", err)
		}
	}
}

func (s *meterReportServer) CollectBatch(stream language_agent.MeterReportService_CollectBatchServer) error {
	for {
		meterDataCollection, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving meter data collection: %v", err)
			return err
		}
		log.Printf("Received meter data collection: %+v", meterDataCollection)
		if err := utils.SaveToJSON("meter_batch", meterDataCollection); err != nil {
			log.Printf("Error saving meter data collection: %v", err)
		}
	}
}

type validationServer struct {
	dbm.UnimplementedValidationServiceServer
}

func (s *validationServer) GetQueriesForValidation(ctx context.Context, req *dbm.ValidationQueryRequest) (*dbm.ValidationQuery, error) {
	log.Printf("Received validation query request: serviceId=%s, status=%s", req.ServiceId, req.Status)
	// if err := saveToJSON("validation_query_request", req); err != nil {
	// 	log.Printf("Error saving validation query request: %v", err)
	// }

	return &dbm.ValidationQuery{
		QId:       []string{"query-123", "query-456"},
		Query:     []string{"SELECT * FROM users", "SELECT * FROM orders"},
		Entity:    []string{"users", "orders"},
		ServiceId: []string{req.ServiceId, req.ServiceId},
		Status:    []string{req.Status, req.Status},
	}, nil
}

func (s *validationServer) NotifyQueryValidationStatus(ctx context.Context, query *dbm.ValidationQuery) (*common.Commands, error) {
	log.Printf("Received query validation status notification: %+v", query)
	if err := utils.SaveToJSON("validation_status_notification", query); err != nil {
		log.Printf("Error saving validation status notification: %v", err)
	}

	return &common.Commands{}, nil
}

func InitDbMonitoringServices(s *grpc.Server) {
	management.RegisterManagementServiceServer(s, &managementServer{})
	dbm.RegisterDbmEventServiceServer(s, &dbmEventServer{})
	dbm.RegisterDbmQueryServiceServer(s, &dbmQueryServer{})
	language_agent.RegisterMeterReportServiceServer(s, &meterReportServer{})
	dbm.RegisterValidationServiceServer(s, &validationServer{})
}
