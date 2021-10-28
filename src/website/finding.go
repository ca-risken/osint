package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/vikyd/zero"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/pkg/common"
	"github.com/ca-risken/osint/pkg/message"
)

func (s *sqsHandler) putFindings(ctx context.Context, result *wappalyzerResult, message *message.OsintQueueMessage) error {
	for _, technology := range result.Technologies {
		data, err := json.Marshal(technology)
		if err != nil {
			return err
		}
		finding := &finding.FindingForUpsert{
			Description:      getDescription(message.ResourceName, technology.Name, technology.Version),
			DataSource:       message.DataSource,
			DataSourceId:     generateDataSourceID(fmt.Sprintf("description_%v_%v", message.ResourceName, technology.Name)),
			ResourceName:     message.ResourceName,
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

func (s *sqsHandler) putFinding(ctx context.Context, wappalyzerFinding *finding.FindingForUpsert, msg *message.OsintQueueMessage, categories []wappalyzerCategory) error {
	res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: wappalyzerFinding})
	if err != nil {
		return err
	}
	_ = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagOsint)
	_ = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, msg.ResourceName)
	for _, category := range categories {
		_ = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, category.Name)
	}
	return nil
}

func (s *sqsHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {

	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		appLogger.Errorf("Failed to TagFinding. error: %v", err)
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
