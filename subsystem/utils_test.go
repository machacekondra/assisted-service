package subsystem

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
)

const (
	defaultWaitForHostStateTimeout    = 20 * time.Second
	defaultWaitForClusterStateTimeout = 40 * time.Second
)

func clearDB() {
	db.Delete(&models.Host{})
	db.Delete(&models.Cluster{})
}

func strToUUID(s string) *strfmt.UUID {
	u := strfmt.UUID(s)
	return &u
}

func registerHost(clusterID strfmt.UUID) *models.HostRegistrationResponse {
	uuid := strToUUID(uuid.New().String())
	return registerHostByUUID(clusterID, *uuid)
}

func registerHostByUUID(clusterID, hostID strfmt.UUID) *models.HostRegistrationResponse {
	host, err := agentBMClient.Installer.RegisterHost(context.Background(), &installer.RegisterHostParams{
		ClusterID: clusterID,
		NewHostParams: &models.HostCreateParams{
			HostID: &hostID,
		},
	})
	Expect(err).NotTo(HaveOccurred())
	return host.GetPayload()
}

func getHost(clusterID, hostID strfmt.UUID) *models.Host {
	host, err := userBMClient.Installer.GetHost(context.Background(), &installer.GetHostParams{
		ClusterID: clusterID,
		HostID:    hostID,
	})
	Expect(err).NotTo(HaveOccurred())
	return host.GetPayload()
}

func registerCluster(ctx context.Context, client *client.AssistedInstall, clusterName string, pullSecret string) (strfmt.UUID, error) {
	var cluster, err = client.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
		NewClusterParams: &models.ClusterCreateParams{
			Name:             swag.String(clusterName),
			OpenshiftVersion: swag.String(common.TestDefaultConfig.OpenShiftVersion),
			PullSecret:       swag.String(pullSecret),
			BaseDNSDomain:    "example.com",
		},
	})
	if err != nil {
		return "", err
	}
	return *cluster.GetPayload().ID, nil
}

func getCluster(clusterID strfmt.UUID) *models.Cluster {
	cluster, err := userBMClient.Installer.GetCluster(context.Background(), &installer.GetClusterParams{
		ClusterID: clusterID,
	})
	Expect(err).NotTo(HaveOccurred())
	return cluster.GetPayload()
}

func getCommonCluster(ctx context.Context, clusterID strfmt.UUID) *common.Cluster {
	var cluster common.Cluster
	err := db.First(&cluster, "id = ?", clusterID).Error
	Expect(err).ShouldNot(HaveOccurred())
	return &cluster
}

func checkStepsInList(steps models.Steps, stepTypes []models.StepType, numSteps int) {
	Expect(len(steps.Instructions)).Should(BeNumerically(">=", numSteps))
	for _, stepType := range stepTypes {
		_, res := getStepInList(steps, stepType)
		Expect(res).Should(Equal(true))
	}
}

func getStepInList(steps models.Steps, sType models.StepType) (*models.Step, bool) {
	for _, step := range steps.Instructions {
		if step.StepType == sType {
			return step, true
		}
	}
	return nil, false
}

func getNextSteps(clusterID, hostID strfmt.UUID) models.Steps {
	steps, err := agentBMClient.Installer.GetNextSteps(context.Background(), &installer.GetNextStepsParams{
		ClusterID: clusterID,
		HostID:    hostID,
	})
	Expect(err).NotTo(HaveOccurred())
	return *steps.GetPayload()
}

func updateHostLogProgress(clusterID strfmt.UUID, hostID strfmt.UUID, progress models.LogsState) {
	ctx := context.Background()

	updateReply, err := agentBMClient.Installer.UpdateHostLogsProgress(ctx, &installer.UpdateHostLogsProgressParams{
		ClusterID: clusterID,
		HostID:    hostID,
		LogsProgressParams: &models.LogsProgressParams{
			LogsState: progress,
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(updateReply).Should(BeAssignableToTypeOf(installer.NewUpdateHostLogsProgressNoContent()))
}

func updateClusterLogProgress(clusterID strfmt.UUID, progress models.LogsState) {
	ctx := context.Background()

	updateReply, err := agentBMClient.Installer.UpdateClusterLogsProgress(ctx, &installer.UpdateClusterLogsProgressParams{
		ClusterID: clusterID,
		LogsProgressParams: &models.LogsProgressParams{
			LogsState: progress,
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(updateReply).Should(BeAssignableToTypeOf(installer.NewUpdateClusterLogsProgressNoContent()))
}

func updateProgress(hostID strfmt.UUID, clusterID strfmt.UUID, current_step models.HostStage) {
	updateProgressWithInfo(hostID, clusterID, current_step, "")
}

func updateProgressWithInfo(hostID strfmt.UUID, clusterID strfmt.UUID, current_step models.HostStage, info string) {
	ctx := context.Background()

	installProgress := &models.HostProgress{
		CurrentStage: current_step,
		ProgressInfo: info,
	}
	updateReply, err := agentBMClient.Installer.UpdateHostInstallProgress(ctx, &installer.UpdateHostInstallProgressParams{
		ClusterID:    clusterID,
		HostProgress: installProgress,
		HostID:       hostID,
	})
	Expect(err).ShouldNot(HaveOccurred())
	Expect(updateReply).Should(BeAssignableToTypeOf(installer.NewUpdateHostInstallProgressOK()))
}

func updateClusterInstallProgressWithInfo(clusterID strfmt.UUID, info string) {
	ctx := context.Background()

	updateReply, err := agentBMClient.Installer.UpdateClusterInstallProgress(ctx, &installer.UpdateClusterInstallProgressParams{
		ClusterID:       clusterID,
		ClusterProgress: info,
	})

	Expect(err).ShouldNot(HaveOccurred())
	Expect(updateReply).Should(BeAssignableToTypeOf(installer.NewUpdateClusterInstallProgressNoContent()))
}

func generateHWPostStepReply(ctx context.Context, h *models.Host, hwInfo *models.Inventory, hostname string) {
	hwInfo.Hostname = hostname
	hw, err := json.Marshal(&hwInfo)
	Expect(err).NotTo(HaveOccurred())
	_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
		ClusterID: h.ClusterID,
		HostID:    *h.ID,
		Reply: &models.StepReply{
			ExitCode: 0,
			Output:   string(hw),
			StepID:   string(models.StepTypeInventory),
			StepType: models.StepTypeInventory,
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
}

func generateNTPPostStepReply(ctx context.Context, h *models.Host, ntpSources []*models.NtpSource) {
	response := models.NtpSynchronizationResponse{
		NtpSources: ntpSources,
	}

	bytes, err := json.Marshal(&response)
	Expect(err).NotTo(HaveOccurred())
	_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
		ClusterID: h.ClusterID,
		HostID:    *h.ID,
		Reply: &models.StepReply{
			ExitCode: 0,
			Output:   string(bytes),
			StepID:   string(models.StepTypeNtpSynchronizer),
			StepType: models.StepTypeNtpSynchronizer,
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
}

func generateApiVipPostStepReply(ctx context.Context, h *models.Host, success bool) {
	checkVipApiResponse := models.APIVipConnectivityResponse{
		IsSuccess: success,
	}
	bytes, jsonErr := json.Marshal(checkVipApiResponse)
	Expect(jsonErr).NotTo(HaveOccurred())
	_, err := agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
		ClusterID: h.ClusterID,
		HostID:    *h.ID,
		Reply: &models.StepReply{
			ExitCode: 0,
			StepType: models.StepTypeAPIVipConnectivityCheck,
			Output:   string(bytes),
			StepID:   "apivip-connectivity-check-step",
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
}

func generateContainerImageAvailabilityPostStepReply(ctx context.Context, h *models.Host, imageStatuses []*models.ContainerImageAvailability) {
	response := models.ContainerImageAvailabilityResponse{
		Images: imageStatuses,
	}

	bytes, err := json.Marshal(&response)
	Expect(err).NotTo(HaveOccurred())
	_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
		ClusterID: h.ClusterID,
		HostID:    *h.ID,
		Reply: &models.StepReply{
			ExitCode: 0,
			Output:   string(bytes),
			StepID:   string(models.StepTypeContainerImageAvailability),
			StepType: models.StepTypeContainerImageAvailability,
		},
	})
	Expect(err).ShouldNot(HaveOccurred())
}

func generateEssentialHostSteps(ctx context.Context, h *models.Host, name string) {
	generateEssentialHostStepsWithInventory(ctx, h, name, validHwInfo)
}

func generateEssentialHostStepsWithInventory(ctx context.Context, h *models.Host, name string, inventory *models.Inventory) {
	generateHWPostStepReply(ctx, h, inventory, name)
	generateNTPPostStepReply(ctx, h, []*models.NtpSource{common.TestNTPSourceSynced})
	generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{common.TestImageStatusesSuccess})
}

func registerNode(ctx context.Context, clusterID strfmt.UUID, name string) *models.Host {
	h := &registerHost(clusterID).Host
	generateEssentialHostSteps(ctx, h, name)
	return h
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil

}
