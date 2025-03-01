package subsystem

import (
	"context"
	"encoding/json"
	"time"

	"github.com/filanov/stateswitch/examples/host/host"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/internal/common"
	serviceHost "github.com/openshift/assisted-service/internal/host"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/auth"
)

var _ = Describe("Host tests", func() {
	ctx := context.Background()
	var cluster *installer.RegisterClusterCreated
	var clusterID strfmt.UUID

	AfterEach(func() {
		clearDB()
	})

	BeforeEach(func() {
		var err error
		cluster, err = userBMClient.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				Name:             swag.String("test-cluster"),
				OpenshiftVersion: swag.String(common.TestDefaultConfig.OpenShiftVersion),
				PullSecret:       swag.String(pullSecret),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		clusterID = *cluster.GetPayload().ID
	})

	It("host CRUD", func() {
		host := &registerHost(clusterID).Host
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("discovering"))
		Expect(host.StatusUpdatedAt).ShouldNot(Equal(strfmt.DateTime(time.Time{})))

		list, err := userBMClient.Installer.ListHosts(ctx, &installer.ListHostsParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list.GetPayload())).Should(Equal(1))

		list, err = agentBMClient.Installer.ListHosts(ctx, &installer.ListHostsParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list.GetPayload())).Should(Equal(1))

		_, err = userBMClient.Installer.DeregisterHost(ctx, &installer.DeregisterHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		list, err = userBMClient.Installer.ListHosts(ctx, &installer.ListHostsParams{ClusterID: clusterID})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(list.GetPayload())).Should(Equal(0))

		_, err = userBMClient.Installer.GetHost(ctx, &installer.GetHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).Should(HaveOccurred())
	})

	var defaultInventory = func() string {
		inventory := models.Inventory{
			Interfaces: []*models.Interface{
				{
					Name: "eth0",
					IPV4Addresses: []string{
						"1.2.3.4/24",
					},
					SpeedMbps: 20,
				},
				{
					Name: "eth1",
					IPV4Addresses: []string{
						"1.2.5.4/24",
					},
					SpeedMbps: 40,
				},
			},

			// CPU, Disks, and Memory were added here to prevent the case that assisted-service crashes in case the monitor starts
			// working in the middle of the test and this inventory is in the database.
			CPU: &models.CPU{
				Count: 4,
			},
			Disks: []*models.Disk{
				{
					Name:      "sda1",
					DriveType: "HDD",
					SizeBytes: int64(120) * (int64(1) << 30),
				},
			},
			Memory: &models.Memory{
				PhysicalBytes: int64(16) * (int64(1) << 30),
			},
			SystemVendor: &models.SystemVendor{Manufacturer: "Red Hat", ProductName: "RHEL", SerialNumber: "3534"},
			Timestamp:    1601845851,
		}
		b, err := json.Marshal(&inventory)
		Expect(err).To(Not(HaveOccurred()))
		return string(b)
	}

	It("next step", func() {
		_, err := userBMClient.Installer.UpdateCluster(ctx, &installer.UpdateClusterParams{
			ClusterID: clusterID,
			ClusterUpdateParams: &models.ClusterUpdateParams{
				VipDhcpAllocation: swag.Bool(false),
			},
		})
		Expect(err).ToNot(HaveOccurred())
		host := &registerHost(clusterID).Host
		host2 := &registerHost(clusterID).Host
		Expect(db.Model(host2).UpdateColumns(&models.Host{Inventory: defaultInventory(),
			Status: swag.String(models.HostStatusInsufficient)}).Error).NotTo(HaveOccurred())
		steps := getNextSteps(clusterID, *host.ID)
		_, ok := getStepInList(steps, models.StepTypeInventory)
		Expect(ok).Should(Equal(true))
		host = getHost(clusterID, *host.ID)
		Expect(db.Model(host).Update("status", "insufficient").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).UpdateColumn("inventory", defaultInventory()).Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeInventory)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeFreeNetworkAddresses)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "known").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeConnectivityCheck)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeFreeNetworkAddresses)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "disabled").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		Expect(steps.NextInstructionSeconds).Should(Equal(int64(120)))
		Expect(*steps.PostStepAction).Should(Equal(models.StepsPostStepActionContinue))
		Expect(len(steps.Instructions)).Should(Equal(0))
		Expect(db.Model(host).Update("status", "insufficient").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeConnectivityCheck)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "error").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeExecute)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", models.HostStatusResetting).Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeResetInstallation)
		Expect(ok).Should(Equal(true))
	})

	It("next step - DHCP", func() {
		Expect(db.Table("clusters").Where("id = ?", clusterID.String()).UpdateColumn("machine_network_cidr", "1.2.3.0/24").Error).ToNot(HaveOccurred())
		host := &registerHost(clusterID).Host
		host2 := &registerHost(clusterID).Host
		Expect(db.Model(host2).UpdateColumns(&models.Host{Inventory: defaultInventory(),
			Status: swag.String(models.HostStatusInsufficient)}).Error).NotTo(HaveOccurred())
		steps := getNextSteps(clusterID, *host.ID)
		_, ok := getStepInList(steps, models.StepTypeInventory)
		Expect(ok).Should(Equal(true))
		host = getHost(clusterID, *host.ID)
		Expect(db.Model(host).Update("status", "insufficient").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).UpdateColumn("inventory", defaultInventory()).Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeInventory)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeFreeNetworkAddresses)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeDhcpLeaseAllocate)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "known").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeConnectivityCheck)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeFreeNetworkAddresses)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeDhcpLeaseAllocate)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "disabled").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		Expect(steps.NextInstructionSeconds).Should(Equal(int64(120)))
		Expect(len(steps.Instructions)).Should(Equal(0))
		Expect(db.Model(host).Update("status", "insufficient").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeConnectivityCheck)
		Expect(ok).Should(Equal(true))
		_, ok = getStepInList(steps, models.StepTypeDhcpLeaseAllocate)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", "error").Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeExecute)
		Expect(ok).Should(Equal(true))
		Expect(db.Model(host).Update("status", models.HostStatusResetting).Error).NotTo(HaveOccurred())
		steps = getNextSteps(clusterID, *host.ID)
		_, ok = getStepInList(steps, models.StepTypeResetInstallation)
		Expect(ok).Should(Equal(true))
		for _, st := range []string{models.HostStatusInstallingInProgress, models.HostStatusInstalling, models.HostStatusPreparingForInstallation} {
			Expect(db.Model(host).Update("status", st).Error).NotTo(HaveOccurred())
			steps = getNextSteps(clusterID, *host.ID)
			_, ok = getStepInList(steps, models.StepTypeDhcpLeaseAllocate)
			Expect(ok).Should(Equal(true))
		}
	})

	It("host_disconnection", func() {
		host := &registerHost(clusterID).Host
		Expect(db.Model(host).Update("status", "installing").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("role", "master").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("bootstrap", "true").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).UpdateColumn("inventory", defaultInventory()).Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("CheckedInAt", strfmt.DateTime(time.Time{})).Error).NotTo(HaveOccurred())

		host = getHost(clusterID, *host.ID)
		time.Sleep(time.Second * 3)
		host = getHost(clusterID, *host.ID)
		Expect(swag.StringValue(host.Status)).Should(Equal("error"))
		Expect(swag.StringValue(host.StatusInfo)).Should(Equal("Host failed to install due to timeout while connecting to host"))
	})

	It("host installation progress", func() {
		host := &registerHost(clusterID).Host
		Expect(db.Model(host).Update("status", "installing").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("role", "master").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("bootstrap", "true").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).UpdateColumn("inventory", defaultInventory()).Error).NotTo(HaveOccurred())

		updateProgress(*host.ID, clusterID, models.HostStageStartingInstallation)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageStartingInstallation))
		time.Sleep(time.Second * 3)
		updateProgress(*host.ID, clusterID, models.HostStageInstalling)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageInstalling))
		time.Sleep(time.Second * 3)
		updateProgress(*host.ID, clusterID, models.HostStageWritingImageToDisk)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageWritingImageToDisk))
		time.Sleep(time.Second * 3)
		updateProgress(*host.ID, clusterID, models.HostStageRebooting)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageRebooting))
		time.Sleep(time.Second * 3)
		updateProgress(*host.ID, clusterID, models.HostStageConfiguring)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageConfiguring))
		time.Sleep(time.Second * 3)
		updateProgress(*host.ID, clusterID, models.HostStageDone)
		host = getHost(clusterID, *host.ID)
		Expect(host.Progress.CurrentStage).Should(Equal(models.HostStageDone))
		time.Sleep(time.Second * 3)
	})

	It("installation_error_reply", func() {
		host := &registerHost(clusterID).Host
		Expect(db.Model(host).Update("status", "installing").Error).NotTo(HaveOccurred())
		Expect(db.Model(host).UpdateColumn("inventory", defaultInventory()).Error).NotTo(HaveOccurred())
		Expect(db.Model(host).Update("role", "worker").Error).NotTo(HaveOccurred())

		_, err := agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
			Reply: &models.StepReply{
				ExitCode: 137,
				Output:   "Failed to install",
				StepType: models.StepTypeInstall,
				StepID:   "installCmd-" + string(models.StepTypeExecute),
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(swag.StringValue(host.Status)).Should(Equal("error"))
		Expect(swag.StringValue(host.StatusInfo)).Should(Equal("installation command failed"))

	})

	It("connectivity_report_store_only_relevant_reply", func() {
		host := &registerHost(clusterID).Host

		connectivity := "{\"remote_hosts\":[{\"host_id\":\"b8a1228d-1091-4e79-be66-738a160f9ff7\",\"l2_connectivity\":null,\"l3_connectivity\":null}]}"
		extraConnectivity := "{\"extra\":\"data\",\"remote_hosts\":[{\"host_id\":\"b8a1228d-1091-4e79-be66-738a160f9ff7\",\"l2_connectivity\":null,\"l3_connectivity\":null}]}"

		_, err := agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   extraConnectivity,
				StepID:   string(models.StepTypeConnectivityCheck),
				StepType: models.StepTypeConnectivityCheck,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(host.Connectivity).Should(Equal(connectivity))

		_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   "not a json",
				StepID:   string(models.StepTypeConnectivityCheck),
				StepType: models.StepTypeConnectivityCheck,
			},
		})
		Expect(err).To(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(host.Connectivity).Should(Equal(connectivity))

		//exit code is not 0
		_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
			Reply: &models.StepReply{
				ExitCode: -1,
				Error:    "some error",
				Output:   "not a json",
				StepID:   string(models.StepTypeConnectivityCheck),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(host.Connectivity).Should(Equal(connectivity))

	})

	It("free addresses report", func() {
		h := &registerHost(clusterID).Host

		free_addresses_report := "[{\"free_addresses\":[\"10.0.0.0\",\"10.0.0.1\"],\"network\":\"10.0.0.0/24\"},{\"free_addresses\":[\"10.0.1.0\"],\"network\":\"10.0.1.0/24\"}]"

		_, err := agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *h.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   free_addresses_report,
				StepID:   string(models.StepTypeFreeNetworkAddresses),
				StepType: models.StepTypeFreeNetworkAddresses,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(db.Model(h).UpdateColumn("status", host.StateInsufficient).Error).NotTo(HaveOccurred())

		freeAddressesReply, err := userBMClient.Installer.GetFreeAddresses(ctx, &installer.GetFreeAddressesParams{
			ClusterID: clusterID,
			Network:   "10.0.0.0/24",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(freeAddressesReply.Payload).To(HaveLen(2))
		Expect(freeAddressesReply.Payload[0]).To(Equal(strfmt.IPv4("10.0.0.0")))
		Expect(freeAddressesReply.Payload[1]).To(Equal(strfmt.IPv4("10.0.0.1")))

		freeAddressesReply, err = userBMClient.Installer.GetFreeAddresses(ctx, &installer.GetFreeAddressesParams{
			ClusterID: clusterID,
			Network:   "10.0.1.0/24",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(freeAddressesReply.Payload).To(HaveLen(1))
		Expect(freeAddressesReply.Payload[0]).To(Equal(strfmt.IPv4("10.0.1.0")))

		freeAddressesReply, err = userBMClient.Installer.GetFreeAddresses(ctx, &installer.GetFreeAddressesParams{
			ClusterID: clusterID,
			Network:   "10.0.2.0/24",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(freeAddressesReply.Payload).To(BeEmpty())

		_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *h.ID,
			Reply: &models.StepReply{
				ExitCode: 0,
				Output:   "not a json",
				StepID:   string(models.StepTypeFreeNetworkAddresses),
				StepType: models.StepTypeFreeNetworkAddresses,
			},
		})
		Expect(err).To(HaveOccurred())

		//exit code is not 0
		_, err = agentBMClient.Installer.PostStepReply(ctx, &installer.PostStepReplyParams{
			ClusterID: clusterID,
			HostID:    *h.ID,
			Reply: &models.StepReply{
				ExitCode: -1,
				Error:    "some error",
				Output:   "not a json",
				StepID:   string(models.StepTypeFreeNetworkAddresses),
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	Context("image availability", func() {

		var (
			h           *models.Host
			imageStatus *models.ContainerImageAvailability
		)

		BeforeEach(func() {
			Skip("OCPBUGSM-25447 AreContainerImagesAvailable isn't mandatory validation for host isSufficientForInstall")
			h = &registerHost(clusterID).Host
		})

		getHostImageStatus := func(hostID strfmt.UUID, imageName string) *models.ContainerImageAvailability {
			hostInDb := getHost(clusterID, hostID)

			var hostImageStatuses map[string]*models.ContainerImageAvailability
			Expect(json.Unmarshal([]byte(hostInDb.ImagesStatus), &hostImageStatuses)).ShouldNot(HaveOccurred())

			return hostImageStatuses[imageName]
		}

		It("First success good bandwidth", func() {
			By("pull success", func() {
				imageStatus = common.TestImageStatusesSuccess

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{imageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(imageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationSuccess), models.HostValidationIDContainerImagesAvailable)
			})

			By("network failure", func() {
				newImageStatus := common.TestImageStatusesFailure
				expectedImageStatus := &models.ContainerImageAvailability{
					Name:         newImageStatus.Name,
					Result:       newImageStatus.Result,
					DownloadRate: imageStatus.DownloadRate,
					SizeBytes:    imageStatus.SizeBytes,
					Time:         imageStatus.Time,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{newImageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(expectedImageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationFailure), models.HostValidationIDContainerImagesAvailable)
			})

			By("network fixed", func() {
				newImageStatus := &models.ContainerImageAvailability{
					Name:   imageStatus.Name,
					Result: models.ContainerImageAvailabilityResultSuccess,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{newImageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(imageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationSuccess), models.HostValidationIDContainerImagesAvailable)
			})
		})

		It("First success bad bandwidth", func() {
			By("pull success", func() {
				imageStatus = &models.ContainerImageAvailability{
					Name:         common.TestDefaultConfig.ImageName,
					Result:       models.ContainerImageAvailabilityResultSuccess,
					DownloadRate: 0.000333,
					SizeBytes:    333000000.0,
					Time:         1000000.0,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{imageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(imageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationFailure), models.HostValidationIDContainerImagesAvailable)
			})

			By("network failure", func() {
				newImageStatus := common.TestImageStatusesFailure
				expectedImageStatus := &models.ContainerImageAvailability{
					Name:         newImageStatus.Name,
					Result:       newImageStatus.Result,
					DownloadRate: imageStatus.DownloadRate,
					SizeBytes:    imageStatus.SizeBytes,
					Time:         imageStatus.Time,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{newImageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(expectedImageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationFailure), models.HostValidationIDContainerImagesAvailable)
			})

			By("network fixed", func() {
				newImageStatus := &models.ContainerImageAvailability{
					Name:   imageStatus.Name,
					Result: models.ContainerImageAvailabilityResultSuccess,
				}
				expectedImageStatus := &models.ContainerImageAvailability{
					Name:         newImageStatus.Name,
					Result:       newImageStatus.Result,
					DownloadRate: imageStatus.DownloadRate,
					SizeBytes:    imageStatus.SizeBytes,
					Time:         imageStatus.Time,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{newImageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(expectedImageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationFailure), models.HostValidationIDContainerImagesAvailable)
			})
		})

		It("First failure", func() {
			By("pull failed", func() {
				imageStatus = common.TestImageStatusesFailure

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{imageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(imageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationFailure), models.HostValidationIDContainerImagesAvailable)
			})
			By("network fixed", func() {
				newImageStatus := common.TestImageStatusesSuccess
				expectedImageStatus := &models.ContainerImageAvailability{
					Name:         newImageStatus.Name,
					Result:       newImageStatus.Result,
					DownloadRate: imageStatus.DownloadRate,
					SizeBytes:    imageStatus.SizeBytes,
					Time:         imageStatus.Time,
				}

				generateContainerImageAvailabilityPostStepReply(ctx, h, []*models.ContainerImageAvailability{newImageStatus})
				Expect(getHostImageStatus(*h.ID, imageStatus.Name)).Should(Equal(expectedImageStatus))
				waitForHostValidationStatus(clusterID, *h.ID, string(serviceHost.ValidationSuccess), models.HostValidationIDContainerImagesAvailable)
			})
		})
	})

	It("disable enable", func() {
		host := &registerHost(clusterID).Host
		_, err := userBMClient.Installer.DisableHost(ctx, &installer.DisableHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("disabled"))
		Expect(len(getNextSteps(clusterID, *host.ID).Instructions)).Should(Equal(0))

		_, err = userBMClient.Installer.EnableHost(ctx, &installer.EnableHostParams{
			ClusterID: clusterID,
			HostID:    *host.ID,
		})
		Expect(err).NotTo(HaveOccurred())
		host = getHost(clusterID, *host.ID)
		Expect(*host.Status).Should(Equal("discovering"))
		Expect(len(getNextSteps(clusterID, *host.ID).Instructions)).ShouldNot(Equal(0))
	})

	It("register_same_host_id", func() {
		hostID := strToUUID(uuid.New().String())
		// register to cluster1
		_, err := agentBMClient.Installer.RegisterHost(context.Background(), &installer.RegisterHostParams{
			ClusterID: clusterID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		cluster2, err := userBMClient.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
			NewClusterParams: &models.ClusterCreateParams{
				Name:             swag.String("another-cluster"),
				OpenshiftVersion: swag.String(common.TestDefaultConfig.OpenShiftVersion),
				PullSecret:       swag.String(pullSecret),
			},
		})
		Expect(err).NotTo(HaveOccurred())

		// register to cluster2
		_, err = agentBMClient.Installer.RegisterHost(ctx, &installer.RegisterHostParams{
			ClusterID: *cluster2.GetPayload().ID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		// successfully get from both clusters
		_ = getHost(clusterID, *hostID)
		_ = getHost(*cluster2.GetPayload().ID, *hostID)

		_, err = userBMClient.Installer.DeregisterHost(ctx, &installer.DeregisterHostParams{
			ClusterID: clusterID,
			HostID:    *hostID,
		})
		Expect(err).NotTo(HaveOccurred())
		h := getHost(*cluster2.GetPayload().ID, *hostID)

		// register again to cluster 2 and expect it to be in discovery status
		Expect(db.Model(h).Update("status", "known").Error).NotTo(HaveOccurred())
		h = getHost(*cluster2.GetPayload().ID, *hostID)
		Expect(swag.StringValue(h.Status)).Should(Equal("known"))
		_, err = agentBMClient.Installer.RegisterHost(ctx, &installer.RegisterHostParams{
			ClusterID: *cluster2.GetPayload().ID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		h = getHost(*cluster2.GetPayload().ID, *hostID)
		Expect(swag.StringValue(h.Status)).Should(Equal("discovering"))
	})

	It("register_wrong_pull_secret", func() {
		if Options.AuthType == auth.TypeNone {
			Skip("auth is disabled")
		}

		wrongTokenStubID, err := wiremock.createWrongStubTokenAuth(WrongPullSecret)
		Expect(err).ToNot(HaveOccurred())

		hostID := strToUUID(uuid.New().String())
		_, err = badAgentBMClient.Installer.RegisterHost(context.Background(), &installer.RegisterHostParams{
			ClusterID: clusterID,
			NewHostParams: &models.HostCreateParams{
				HostID: hostID,
			},
		})
		Expect(err).To(HaveOccurred())

		err = wiremock.DeleteStub(wrongTokenStubID)
		Expect(err).ToNot(HaveOccurred())
	})

	It("next_step_runner_command", func() {
		registration := registerHost(clusterID)
		Expect(registration.NextStepRunnerCommand).ShouldNot(BeNil())
		Expect(registration.NextStepRunnerCommand.Command).ShouldNot(BeEmpty())
		Expect(registration.NextStepRunnerCommand.Args).ShouldNot(BeEmpty())
		Expect(registration.NextStepRunnerCommand.RetrySeconds).Should(Equal(int64(0))) //default, just to have in the API
	})
})
