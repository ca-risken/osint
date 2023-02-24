package subdomain

import (
	"context"
	"encoding/hex"

	"crypto/sha256"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/osint/pkg/common"
	"github.com/vikyd/zero"
)

func (s *SQSHandler) putFindings(ctx context.Context, findingMap map[string][]*finding.FindingForUpsert, projectID uint32, resourceName string) error {
	findingBatchParam := []*finding.FindingBatchForUpsert{}
	for key, findings := range findingMap {
		for _, f := range findings {
			// tag
			tags := []*finding.FindingTagForBatch{
				{Tag: common.TagOsint},
				{Tag: common.TagDomain},
				{Tag: resourceName},
			}
			var tagFindingType string
			switch key {
			case "Takeover":
				tagFindingType = common.TagTakeover
			case "PrivateExpose":
				tagFindingType = common.TagPrivateExpose
			case "CertificateExpiration":
				tagFindingType = common.TagCertificateExpiration
			}
			if !zero.IsZeroVal(tagFindingType) {
				tags = append(tags, &finding.FindingTagForBatch{Tag: tagFindingType})
			}

			// recommend
			var recommend *finding.RecommendForBatch
			r := getRecommend(key)
			if r.Type == "" || (r.Risk == "" && r.Recommendation == "") {
				s.logger.Warnf(ctx, "Failed to get recommendation, Unknown category=%s", key)
			} else {
				recommend = &finding.RecommendForBatch{
					Type:           r.Type,
					Risk:           r.Risk,
					Recommendation: r.Recommendation,
				}
			}

			findingBatchParam = append(findingBatchParam, &finding.FindingBatchForUpsert{
				Finding:   f,
				Tag:       tags,
				Recommend: recommend,
			})
		}
	}

	if len(findingBatchParam) == 0 {
		s.logger.Info(ctx, "No finding")
		return nil
	}
	if err := s.putFindingBatch(ctx, projectID, findingBatchParam); err != nil {
		return err
	}
	s.logger.Infof(ctx, "putFindings(%d) succeeded", len(findingBatchParam))
	return nil
}

func (s *SQSHandler) putFindingBatch(ctx context.Context, projectID uint32, params []*finding.FindingBatchForUpsert) error {
	s.logger.Infof(ctx, "Putting findings(%d)...", len(params))
	for idx := 0; idx < len(params); idx = idx + finding.PutFindingBatchMaxLength {
		lastIdx := idx + finding.PutFindingBatchMaxLength
		if lastIdx > len(params) {
			lastIdx = len(params)
		}
		// request per API limits
		s.logger.Debugf(ctx, "Call PutFindingBatch API, (%d ~ %d / %d)", idx+1, lastIdx, len(params))
		req := &finding.PutFindingBatchRequest{ProjectId: projectID, Finding: params[idx:lastIdx]}
		if _, err := s.findingClient.PutFindingBatch(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (s *SQSHandler) makeFindings(ctx context.Context, osintResults *[]osintResult, message *message.OsintQueueMessage) (map[string][]*finding.FindingForUpsert, error) {
	findings := map[string][]*finding.FindingForUpsert{}
	findingsTakeover := []*finding.FindingForUpsert{}
	findingsPrivateExpose := []*finding.FindingForUpsert{}
	findingsCertificateExpiration := []*finding.FindingForUpsert{}
	for _, osintResult := range *osintResults {
		findingTakeover, err := osintResult.Takeover.makeFinding(message.ProjectID, message.DataSource)
		if err != nil {
			s.logger.Errorf(ctx, "Error occured when make Takeover finding. error: %v", err)
			return nil, err
		}
		if findingTakeover != nil {
			findingsTakeover = append(findingsTakeover, findingTakeover)
		}
		findingPrivateExpose, err := osintResult.PrivateExpose.makeFinding(message.ProjectID, message.DataSource)
		if err != nil {
			s.logger.Errorf(ctx, "Error occured when make PrivateExpose finding. error: %v", err)
			return nil, err
		}
		if findingPrivateExpose != nil {
			findingsPrivateExpose = append(findingsPrivateExpose, findingPrivateExpose)
		}
		findingCertificateExpiration, err := osintResult.CertificateExpiration.makeFinding(message.ProjectID, message.DataSource)
		if err != nil {
			s.logger.Errorf(ctx, "Error occured when make Certificate Expiration finding. error: %v", err)
			return nil, err
		}
		if findingCertificateExpiration != nil {
			findingsCertificateExpiration = append(findingsCertificateExpiration, findingCertificateExpiration)
		}
	}
	findings["Takeover"] = findingsTakeover
	findings["PrivateExpose"] = findingsPrivateExpose
	findings["CertificateExpiration"] = findingsCertificateExpiration
	return findings, nil
}

func generateDataSourceID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
