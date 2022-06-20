package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/vikyd/zero"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/osint/pkg/common"
)

func (s *SQSHandler) putFindings(ctx context.Context, result *wappalyzerResult, message *message.OsintQueueMessage) error {
	for _, technology := range result.Technologies {
		data, err := json.Marshal(technology)
		if err != nil {
			return err
		}
		finding := &finding.FindingForUpsert{
			Description:      getDescription(message.ResourceName, technology.Name, technology.Version),
			DataSource:       message.DataSource,
			DataSourceId:     generateDataSourceID(fmt.Sprintf("description_%v_%v", message.ResourceName, technology.Name)),
			ResourceName:     technology.Name,
			ProjectId:        message.ProjectID,
			OriginalScore:    getScore(),
			OriginalMaxScore: 1.0,
			Data:             string(data),
		}
		err = s.putFinding(ctx, finding, message, technology.Categories)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQSHandler) putFinding(ctx context.Context, websiteFinding *finding.FindingForUpsert, msg *message.OsintQueueMessage, categories []wappalyzerCategory) error {
	res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: websiteFinding})
	if err != nil {
		return err
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagOsint); err != nil {
		appLogger.Errorf(ctx, "Failed to tag finding. tag: %v, error: %v", common.TagOsint, err)
		return err
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagWebsite); err != nil {
		appLogger.Errorf(ctx, "Failed to tag finding. tag: %v, error: %v", common.TagWebsite, err)
	}
	if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, fmt.Sprintf("osint_id:%v", msg.OsintID)); err != nil {
		appLogger.Errorf(ctx, "Failed to tag finding. tag: %v, error: %v", fmt.Sprintf("osint_id:%v", msg.OsintID), err)
		return err
	}

	for _, category := range categories {
		if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, category.Name); err != nil {
			appLogger.Errorf(ctx, "Failed to tag finding. tag: %v, error: %v", category.Name, err)
			return err
		}
	}
	return nil
}

func (s *SQSHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {

	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		appLogger.Errorf(ctx, "Failed to TagFinding. error: %v", err)
		return err
	}
	return nil
}

func generateDataSourceID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func getDescription(resourceName, technologyName, version string) string {
	if !zero.IsZeroVal(version) {
		return fmt.Sprintf("%v is using %v. version: %v", resourceName, technologyName, version)
	}
	return fmt.Sprintf("%v is using %v.", resourceName, technologyName)
}

func getScore() float32 {
	return 0.1
}
