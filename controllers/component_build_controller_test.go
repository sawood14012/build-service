/*
Copyright 2021-2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"fmt"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	appstudiov1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	"github.com/redhat-appstudio/application-service/gitops"
	gitopsprepare "github.com/redhat-appstudio/application-service/gitops/prepare"
	"github.com/redhat-appstudio/application-service/pkg/devfile"
	buildappstudiov1alpha1 "github.com/redhat-appstudio/build-service/api/v1alpha1"
	"github.com/redhat-appstudio/build-service/pkg/boerrors"
	"github.com/redhat-appstudio/build-service/pkg/git/github"
	gp "github.com/redhat-appstudio/build-service/pkg/git/gitprovider"
	gpf "github.com/redhat-appstudio/build-service/pkg/git/gitproviderfactory"
	appstudiospiapiv1beta1 "github.com/redhat-appstudio/service-provider-integration-operator/api/v1beta1"
	tektonapi "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	//+kubebuilder:scaffold:imports
)

const (
	githubAppPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAtSwZCtZ0Tnepuezo/TL9vhdP00fOedCpN3HsKjqz7zXTnqkq
foxemrDvRSg5n73sZhZMYX6NKY1FBLTE6OazvJQg0eXu7Bf5+5sg5ZZIthX1wbPU
wk18S5HsC1NTDGHVlQ2pQEuXmpxi/FAKgHaYbhx2J2SzD2gdKEOAufBZ14o8btF2
x64wi0VHX2JAmPXjodYQ2jYbH27lik5kS8wMDRKq0Kt3ABJkxULxUB4xuEbl2gWt
osDrhFTtDWEYtmtL0R5mLoK61FfZmZElvHELQObBcEpJu9F/Bi0mcjle/TuVc69e
tuVbmVg5Bn5iQKtw2hPEMqTLDVPbhhSAXhtJOwIDAQABAoIBAD+XEdce/MXJ9JXg
1MqCilOdZRRYoN1a4vomD2mnHx74OqX25IZ0iIQtVF5mxwsNo5sVeovB2pRaFH6Z
YIAK8c1gBMEHvru5krHAemR7Qlw/CvqJP0VP4y+3MS2senrfIBNoLx71KWpIN+ot
wfHjLo9/h+09yCfBOHK4dsdM2IvxTw9Md6ivipQC7rZmg0lpupNVRKwsx+VQoHCr
wzyPN14w+3tAJA2QoSdsqZLWkfY5YgxYmGuEG7L1A4Ynsjq6lEJyoclbZ6Gp+JML
G9MB0llXz9AV2k+BuNZWQO4cBPeKYtiQ1O5iYohghjSgPIx0iFa/GAvIrZscF2cf
Y43bN/ECgYEA3pKZ55L+oZIym0P2t9tqUaFhFewLdc54WhoMBrjKDQ7PL5n/avG5
Tac1RKA7AoknLYkqWilwRAzSBzNDjYsFVbfxRA10LagP61WhgGAPtEsDM3gbwuNw
HayYowutQ6422Ri4+ujWSeMMSrvx7RcHMPpr3hDSm8ihxbqsWkDNGQMCgYEA0GG/
cdtidDZ0aPwIqBdx6DHwKraYMt1fIzkw1qG59X7jc/Rym2aPUzDob99q8BieBRB3
6djzn+X97z40/87DRFACj/YQPkGvumSdRUtP88b0x87UhbHoOfEIkDeVYyxkbfl0
Mo6H9tlw+QSTaU13AzIBaWOB4nsQaKAM9isHrWkCgYAO0F8iBKyiAGMR5oIjVp1K
9ZzKor1Yh/eGt7kZMW9xUw0DNBLGAXS98GUhPjDvSEWtSDXjbmKkhN3t0MGsSBaA
0A9k4ihbaZY1qatoKfyhmWSLJnFilVS/BN/b6kkL+ip4ZKbbPGgW3t/QkZXWm/PE
lMZdL211JPNvf689CpccFQKBgGXJGUZ4LuMtJjeRxHi22wDcQ7/ZaQaPc0U1TlHI
tZjg3iFpqgGWWzP7k83xh763h5hZrvke7AGSyjLuY90AFglsO5QuUUjXtQqK0vdi
Di+5Yx+mO9ECUbjbr58iR2ol6Ph+/O8lB+zf0XsRbR/moteAuYfM/0itbBpu82Xb
JujhAoGAVdNMYdamHOJlkjYjYJmWOHPIFBMTDYk8pDATLOlqthV+fzlD2SlF0wxm
OlxbwEr3BSQlE3GFPHe+J3JV3BbJFsVkM8KdMUpQCl+aD/3nWaGBYHz4yvJc0i9h
duAFIZgXeivAZQAqys4JanG81cg4/d4urI3Qk9stlhB5CNCJR4k=
-----END RSA PRIVATE KEY-----`
)

var _ = Describe("Component initial build controller", func() {

	const (
		pacHost       = "pac-host"
		pacWebhookUrl = "https://" + pacHost
	)

	var (
		// All related to the component resources have the same key (but different type)
		resourceKey           = types.NamespacedName{Name: HASCompName, Namespace: HASAppNamespace}
		pacRouteKey           = types.NamespacedName{Name: pipelinesAsCodeRouteName, Namespace: pipelinesAsCodeNamespace}
		pacSecretKey          = types.NamespacedName{Name: gitopsprepare.PipelinesAsCodeSecretName, Namespace: buildServiceNamespaceName}
		namespacePaCSecretKey = types.NamespacedName{Name: gitopsprepare.PipelinesAsCodeSecretName, Namespace: HASAppNamespace}
		webhookSecretKey      = types.NamespacedName{Name: gitops.PipelinesAsCodeWebhooksSecretName, Namespace: HASAppNamespace}
	)

	BeforeEach(func() {
		createNamespace(buildServiceNamespaceName)
		createBuildPipelineRunSelector(defaultSelectorKey)
		github.GetAppInstallations = func(githubAppIdStr string, appPrivateKeyPem []byte) ([]github.ApplicationInstallation, string, error) {
			return nil, "slug", nil
		}
	})

	AfterEach(func() {
		deleteBuildPipelineRunSelector(defaultSelectorKey)
	})

	Context("Test Pipelines as Code build preparation", func() {
		_ = BeforeEach(func() {
			createNamespace(pipelinesAsCodeNamespace)
			createRoute(pacRouteKey, "pac-host")
			createNamespace(buildServiceNamespaceName)
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			ResetTestGitProviderClient()
		})

		_ = AfterEach(func() {
			deleteComponent(resourceKey)

			deleteSecret(webhookSecretKey)
			deleteSecret(namespacePaCSecretKey)

			deleteSecret(pacSecretKey)
			deleteRoute(pacRouteKey)
		})

		It("should successfully submit PR with PaC definitions using GitHub application", func() {
			mergeUrl := "merge-url"

			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				isCreatePaCPullRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(len(d.Files)).To(Equal(2))
				for _, file := range d.Files {
					Expect(strings.HasPrefix(file.FullPath, ".tekton/")).To(BeTrue())
				}
				Expect(d.CommitMessage).ToNot(BeEmpty())
				Expect(d.BranchName).ToNot(BeEmpty())
				Expect(d.BaseBranchName).To(Equal("main"))
				Expect(d.Title).ToNot(BeEmpty())
				Expect(d.Text).ToNot(BeEmpty())
				Expect(d.AuthorName).ToNot(BeEmpty())
				Expect(d.AuthorEmail).ToNot(BeEmpty())
				return mergeUrl, nil
			}
			SetupPaCWebhookFunc = func(string, string, string) error {
				defer GinkgoRecover()
				Fail("Should not create webhook if GitHub application is used")
				return nil
			}

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			waitPaCRepositoryCreated(resourceKey)
			waitPaCFinalizerOnComponent(resourceKey)
			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("enabled"))
			Expect(buildStatus.PaC.MergeUrl).To(Equal(mergeUrl))
			Expect(buildStatus.PaC.ConfigurationTime).ToNot(BeEmpty())
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))
		})

		It("should fail to submit PR if GitHub application is not installed into git repository", func() {
			gpf.CreateGitClient = func(gpf.GitClientConfig) (gp.GitProviderClient, error) {
				return nil, boerrors.NewBuildOpError(boerrors.EGitHubAppNotInstalled,
					fmt.Errorf("GitHub Application is not installed into the repository"))
			}

			EnsurePaCMergeRequestFunc = func(string, *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("Should not invoke merge request creation if GitHub application is not installed into the repository")
				return "url", nil
			}

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			expectedErr := boerrors.NewBuildOpError(boerrors.EGitHubAppNotInstalled, fmt.Errorf("something is wrong"))
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ConfigurationTime).To(BeEmpty())
			Expect(buildStatus.PaC.ErrId).To(Equal(expectedErr.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectedErr.ShortError()))
		})

		It("should fail to submit PR if unknown git provider is used", func() {
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "url", nil
			}

			createCustomComponentWithBuildRequest(componentConfig{
				componentKey: resourceKey,
				gitURL:       "https://my-git-instance.com/devfile-samples/devfile-sample-java-springboot-basic",
			}, BuildRequestConfigurePaCAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			expectedErr := boerrors.NewBuildOpError(boerrors.EUnknownGitProvider, fmt.Errorf("unknow git provider"))
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ConfigurationTime).To(BeEmpty())
			Expect(buildStatus.PaC.ErrId).To(Equal(expectedErr.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectedErr.ShortError()))
		})

		It("should fail to submit PR if PaC secret is invalid", func() {
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			deleteSecret(pacSecretKey)

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    "secret private key",
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			expectedErr := boerrors.NewBuildOpError(boerrors.EPaCSecretInvalid, fmt.Errorf("invalid pac secret"))
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ConfigurationTime).To(BeEmpty())
			Expect(buildStatus.PaC.ErrId).To(Equal(expectedErr.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectedErr.ShortError()))
		})

		It("should fail to submit PR if PaC secret is missing", func() {
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			deleteSecret(pacSecretKey)
			deleteSecret(namespacePaCSecretKey)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			expectedErr := boerrors.NewBuildOpError(boerrors.EPaCSecretNotFound, fmt.Errorf("pac secret not found"))
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ConfigurationTime).To(BeEmpty())
			Expect(buildStatus.PaC.ErrId).To(Equal(expectedErr.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectedErr.ShortError()))
		})

		It("should successfully do PaC provision after error (when PaC GitHub Application was not installed)", func() {
			appNotInstalledErr := boerrors.NewBuildOpError(boerrors.EGitHubAppNotInstalled, nil)
			isCreateGithubClientInvoked := false
			gpf.CreateGitClient = func(gitClientConfig gpf.GitClientConfig) (gp.GitProviderClient, error) {
				isCreateGithubClientInvoked = true
				return nil, appNotInstalledErr
			}
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			Eventually(func() bool {
				return isCreateGithubClientInvoked
			}, timeout, interval).Should(BeTrue())

			// Ensure no more retries after permanent error
			gpf.CreateGitClient = func(gitClientConfig gpf.GitClientConfig) (gp.GitProviderClient, error) {
				defer GinkgoRecover()
				Fail("Should not retry PaC provision on permanent error")
				return nil, nil
			}
			Consistently(func() bool {
				buildStatus := readBuildStatus(getComponent(resourceKey))
				Expect(buildStatus).ToNot(BeNil())
				Expect(buildStatus.PaC).ToNot(BeNil())
				Expect(buildStatus.PaC.State).To(Equal("error"))
				Expect(buildStatus.PaC.ConfigurationTime).To(BeEmpty())
				Expect(buildStatus.PaC.ErrId).To(Equal(appNotInstalledErr.GetErrorId()))
				Expect(buildStatus.PaC.ErrMessage).To(Equal(appNotInstalledErr.ShortError()))
				return true
			}, ensureTimeout, interval).Should(BeTrue())

			// Suppose PaC GH App is installed
			gpf.CreateGitClient = func(gitClientConfig gpf.GitClientConfig) (gp.GitProviderClient, error) {
				return testGitProviderClient, nil
			}
			mergeUrl := "merge-url"
			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				isCreatePaCPullRequestInvoked = true
				return mergeUrl, nil
			}

			// Retry
			setComponentBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("enabled"))
			Expect(buildStatus.PaC.ConfigurationTime).ToNot(BeEmpty())
			Expect(buildStatus.PaC.MergeUrl).To(Equal(mergeUrl))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))
		})

		It("should not copy PaC secret into local namespace if GitHub application is used", func() {
			deleteSecret(namespacePaCSecretKey)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCRepositoryCreated(resourceKey)

			ensureSecretNotCreated(namespacePaCSecretKey)
		})

		It("should successfully submit PR with PaC definitions using token", func() {
			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				isCreatePaCPullRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(len(d.Files)).To(Equal(2))
				for _, file := range d.Files {
					Expect(strings.HasPrefix(file.FullPath, ".tekton/")).To(BeTrue())
				}
				Expect(d.CommitMessage).ToNot(BeEmpty())
				Expect(d.BranchName).ToNot(BeEmpty())
				Expect(d.BaseBranchName).ToNot(BeEmpty())
				Expect(d.Title).ToNot(BeEmpty())
				Expect(d.Text).ToNot(BeEmpty())
				Expect(d.AuthorName).ToNot(BeEmpty())
				Expect(d.AuthorEmail).ToNot(BeEmpty())
				return "url", nil
			}
			isSetupPaCWebhookInvoked := false
			SetupPaCWebhookFunc = func(repoUrl string, webhookUrl string, webhookSecret string) error {
				isSetupPaCWebhookInvoked = true
				Expect(webhookUrl).To(Equal(pacWebhookUrl))
				Expect(webhookSecret).ToNot(BeEmpty())
				Expect(repoUrl).To(Equal(SampleRepoLink))
				return nil
			}

			pacSecretData := map[string]string{"github.token": "ghp_token"}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			waitSecretCreated(namespacePaCSecretKey)
			waitSecretCreated(webhookSecretKey)
			waitPaCRepositoryCreated(resourceKey)
			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				return isSetupPaCWebhookInvoked
			}, timeout, interval).Should(BeTrue())
		})

		It("should provision PaC definitions after initial build, use simple build while PaC enabled, and be able to switch back to simple build only", func() {
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("Should not create PaC configuration PR when using simple build")
				return "url", nil
			}

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestTriggerSimpleBuildAnnotationValue)

			waitOneInitialPipelineRunCreated(resourceKey)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).ToNot(BeEmpty())
			Expect(buildStatus.Simple.ErrId).To(Equal(0))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(""))

			// Do PaC provision
			ResetTestGitProviderClient()

			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				isCreatePaCPullRequestInvoked = true
				return "configure-merge-url", nil
			}
			UndoPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("Should not create undo PaC configuration PR when switching to PaC build")
				return "url", nil
			}
			SetupPaCWebhookFunc = func(repoUrl string, webhookUrl string, webhookSecret string) error {
				defer GinkgoRecover()
				Fail("Should not create webhook if GitHub application is used")
				return nil
			}

			setComponentBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			waitPaCRepositoryCreated(resourceKey)
			waitPaCFinalizerOnComponent(resourceKey)
			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus = readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("enabled"))
			Expect(buildStatus.PaC.ConfigurationTime).ToNot(BeEmpty())
			Expect(buildStatus.PaC.MergeUrl).To(Equal("configure-merge-url"))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))

			// Request simple build while PaC is enabled
			deleteComponentPipelineRuns(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestTriggerSimpleBuildAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			waitOneInitialPipelineRunCreated(resourceKey)

			buildStatus = readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).ToNot(BeEmpty())
			Expect(buildStatus.Simple.ErrId).To(Equal(0))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(""))

			// Do PaC unprovision
			ResetTestGitProviderClient()

			isRemovePaCPullRequestInvoked := false
			UndoPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				isRemovePaCPullRequestInvoked = true
				return "unconfigure-merge-url", nil
			}
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("Should not create PaC configuration PR when switching to simple build")
				return "url", nil
			}
			DeletePaCWebhookFunc = func(repoUrl string, webhookUrl string) error {
				defer GinkgoRecover()
				Fail("Should not delete webhook if GitHub application is used")
				return nil
			}

			// Request PaC unprovision
			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			waitPaCFinalizerOnComponentGone(resourceKey)
			Eventually(func() bool {
				return isRemovePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus = readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("disabled"))
			Expect(buildStatus.PaC.MergeUrl).To(Equal("unconfigure-merge-url"))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))

			// Request simple build
			deleteComponentPipelineRuns(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestTriggerSimpleBuildAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			waitOneInitialPipelineRunCreated(resourceKey)

			buildStatus = readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).ToNot(BeEmpty())
			Expect(buildStatus.Simple.ErrId).To(Equal(0))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(""))
		})

		It("should reuse the same webhook secret for multi component repository", func() {
			var webhookSecretStrings []string
			SetupPaCWebhookFunc = func(repoUrl string, webhookUrl string, webhookSecret string) error {
				webhookSecretStrings = append(webhookSecretStrings, webhookSecret)
				return nil
			}

			pacSecretData := map[string]string{"github.token": "ghp_token"}
			createSecret(pacSecretKey, pacSecretData)

			component1Key := resourceKey
			component2Key := types.NamespacedName{Name: "component2", Namespace: HASAppNamespace}

			createComponentWithBuildRequest(component1Key, BuildRequestConfigurePaCAnnotationValue)
			createCustomComponentWithBuildRequest(componentConfig{
				componentKey:   component2Key,
				containerImage: "registry.io/username/image2:tag2",
			}, BuildRequestConfigurePaCAnnotationValue)
			defer deleteComponent(component2Key)

			waitSecretCreated(namespacePaCSecretKey)
			waitSecretCreated(webhookSecretKey)

			waitPaCRepositoryCreated(component1Key)
			waitPaCRepositoryCreated(component2Key)

			waitComponentAnnotationGone(component1Key, BuildRequestAnnotationName)
			waitComponentAnnotationGone(component2Key, BuildRequestAnnotationName)

			Expect(len(webhookSecretStrings)).To(BeNumerically(">", 0))
			for _, webhookSecret := range webhookSecretStrings {
				Expect(webhookSecret).To(Equal(webhookSecretStrings[0]))
			}
		})

		It("should use different webhook secrets for different components of the same application", func() {
			var webhookSecretStrings []string
			SetupPaCWebhookFunc = func(repoUrl string, webhookUrl string, webhookSecret string) error {
				webhookSecretStrings = append(webhookSecretStrings, webhookSecret)
				return nil
			}

			pacSecretData := map[string]string{"github.token": "ghp_token"}
			createSecret(pacSecretKey, pacSecretData)

			component1Key := resourceKey
			component2Key := types.NamespacedName{Name: "component2", Namespace: HASAppNamespace}

			createComponentWithBuildRequest(component1Key, BuildRequestConfigurePaCAnnotationValue)
			createCustomComponentWithBuildRequest(componentConfig{
				componentKey:   component2Key,
				containerImage: "registry.io/username/image2:tag2",
				gitURL:         "https://github.com/devfile-samples/devfile-sample-go-basic",
			}, BuildRequestConfigurePaCAnnotationValue)
			defer deleteComponent(component2Key)

			waitSecretCreated(namespacePaCSecretKey)
			waitSecretCreated(webhookSecretKey)

			waitPaCRepositoryCreated(component1Key)
			waitPaCRepositoryCreated(component2Key)

			waitComponentAnnotationGone(component1Key, BuildRequestAnnotationName)
			waitComponentAnnotationGone(component2Key, BuildRequestAnnotationName)

			Expect(len(webhookSecretStrings)).To(Equal(2))
			Expect(webhookSecretStrings[0]).ToNot(Equal(webhookSecretStrings[1]))
		})

		It("should set error in status if invalid build action requested", func() {
			createCustomComponentWithBuildRequest(componentConfig{
				componentKey: resourceKey,
			}, "non-existing-build-request")
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Message).To(ContainSubstring("unexpected build request"))
		})

		It("should do nothing if the component devfile model is not set", func() {
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			createComponent(resourceKey)
			setComponentBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			ensureComponentAnnotationValue(resourceKey, BuildRequestAnnotationName, BuildRequestConfigurePaCAnnotationValue)
		})

		It("should do nothing if a container image source is specified in component", func() {
			EnsurePaCMergeRequestFunc = func(string, *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			component := &appstudiov1alpha1.Component{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "appstudio.redhat.com/v1alpha1",
					Kind:       "Component",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      HASCompName,
					Namespace: HASAppNamespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName:  HASCompName,
					Application:    HASAppName,
					ContainerImage: "quay.io/test/image:latest",
				},
			}
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			ensureNoPipelineRunsCreated(resourceKey)
		})

		It("should set default branch as base branch when Revision is not set", func() {
			const repoDefaultBranch = "default-branch"
			GetDefaultBranchFunc = func(repoUrl string) (string, error) {
				return repoDefaultBranch, nil
			}

			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				Expect(d.BaseBranchName).To(Equal(repoDefaultBranch))
				for _, file := range d.Files {
					var prYaml tektonapi.PipelineRun
					if err := yaml.Unmarshal(file.Content, &prYaml); err != nil {
						return "", err
					}
					targetBranches := prYaml.Annotations["pipelinesascode.tekton.dev/on-target-branch"]
					Expect(targetBranches).To(Equal(fmt.Sprintf("[%s]", repoDefaultBranch)))
				}
				isCreatePaCPullRequestInvoked = true
				return "url", nil
			}

			component := getSampleComponentData(resourceKey)
			// Unset Revision so that GetDefaultBranch function is called to use the default branch
			// set in the remote component repository
			component.Spec.Source.GitSource.Revision = ""
			component.Annotations[BuildRequestAnnotationName] = BuildRequestConfigurePaCAnnotationValue
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())
		})

		It("should link auto generated image repository secret to pipeline service accoount", func() {
			userImageRepo := "docker.io/user/image"
			generatedImageRepo := "quay.io/appstudio/generated-image"
			generatedImageRepoSecretName := "generated-image-repo-secret"
			generatedImageRepoSecretKey := types.NamespacedName{Namespace: resourceKey.Namespace, Name: generatedImageRepoSecretName}
			pipelineSAKey := types.NamespacedName{Namespace: resourceKey.Namespace, Name: buildPipelineServiceAccountName}

			checkPROutputImage := func(fileContent []byte, expectedImageRepo string) {
				var prYaml tektonapi.PipelineRun
				Expect(yaml.Unmarshal(fileContent, &prYaml)).To(Succeed())
				outoutImage := ""
				for _, param := range prYaml.Spec.Params {
					if param.Name == "output-image" {
						outoutImage = param.Value.StringVal
						break
					}
				}
				Expect(outoutImage).ToNot(BeEmpty())
				Expect(outoutImage).To(ContainSubstring(expectedImageRepo))
			}

			isCreatePaCPullRequestInvoked := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				checkPROutputImage(d.Files[0].Content, userImageRepo)
				isCreatePaCPullRequestInvoked = true
				return "url", nil
			}

			// Create a component with user's ContainerImage
			createCustomComponentWithBuildRequest(componentConfig{
				componentKey:   resourceKey,
				containerImage: userImageRepo,
			}, BuildRequestConfigurePaCAnnotationValue)

			Eventually(func() bool {
				return isCreatePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			// Switch to generated image repository
			createSecret(generatedImageRepoSecretKey, nil)
			defer deleteSecret(generatedImageRepoSecretKey)

			component := getComponent(resourceKey)
			component.Annotations[ImageRepoGenerateAnnotationName] = "false"
			component.Annotations[ImageRepoAnnotationName] =
				fmt.Sprintf("{\"image\":\"%s\",\"secret\":\"%s\"}", generatedImageRepo, generatedImageRepoSecretName)
			component.Spec.ContainerImage = generatedImageRepo
			Expect(k8sClient.Update(ctx, component)).To(Succeed())

			Eventually(func() bool {
				component = getComponent(resourceKey)
				return component.Spec.ContainerImage == generatedImageRepo
			}, timeout, interval).Should(BeTrue())

			pipelineSA := &corev1.ServiceAccount{}
			Expect(k8sClient.Get(ctx, pipelineSAKey, pipelineSA)).To(Succeed())
			isImageRegistryGeneratedSecretLinked := false
			if pipelineSA.Secrets == nil {
				time.Sleep(1 * time.Second)
				Expect(k8sClient.Get(ctx, pipelineSAKey, pipelineSA)).To(Succeed())
			}
			for _, secret := range pipelineSA.Secrets {
				if secret.Name == generatedImageRepoSecretName {
					isImageRegistryGeneratedSecretLinked = true
					break
				}
			}
			Expect(isImageRegistryGeneratedSecretLinked).To(BeTrue())
		})
	})

	Context("Test Pipelines as Code build clean up", func() {

		_ = BeforeEach(func() {
			createNamespace(pipelinesAsCodeNamespace)
			createRoute(pacRouteKey, "pac-host")
			createNamespace(buildServiceNamespaceName)

			ResetTestGitProviderClient()

			deleteComponent(resourceKey)
		})

		_ = AfterEach(func() {
			deleteSecret(webhookSecretKey)
			deleteSecret(namespacePaCSecretKey)

			deleteSecret(pacSecretKey)
			deleteRoute(pacRouteKey)
		})

		It("should successfully unconfigure PaC even when it isn't able to get PaC webhook", func() {
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)
			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)

			pacSecretData = map[string]string{}
			deleteSecret(pacSecretKey)
			createSecret(pacSecretKey, pacSecretData)
			deleteRoute(pacRouteKey)

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponentGone(resourceKey)
			waitDoneMessageOnComponent(resourceKey)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("disabled"))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))

		})

		It("should successfully unconfigure PaC, removal MR not needed", func() {
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)
			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)
			UndoPaCMergeRequestFunc = func(repoUrl string, data *gp.MergeRequestData) (webUrl string, err error) {
				return "", nil
			}

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponentGone(resourceKey)
			waitDoneMessageOnComponent(resourceKey)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("disabled"))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))

		})

		It("should successfully submit PR with PaC definitions removal using GitHub application", func() {
			mergeUrl := "merge-url"
			isRemovePaCPullRequestInvoked := false
			UndoPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (webUrl string, err error) {
				isRemovePaCPullRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(len(d.Files)).To(Equal(2))
				for _, file := range d.Files {
					Expect(strings.HasPrefix(file.FullPath, ".tekton/")).To(BeTrue())
				}
				Expect(d.CommitMessage).ToNot(BeEmpty())
				Expect(d.BranchName).ToNot(BeEmpty())
				Expect(d.BaseBranchName).ToNot(BeEmpty())
				Expect(d.Title).ToNot(BeEmpty())
				Expect(d.Text).ToNot(BeEmpty())
				Expect(d.AuthorName).ToNot(BeEmpty())
				Expect(d.AuthorEmail).ToNot(BeEmpty())
				return mergeUrl, nil
			}
			DeletePaCWebhookFunc = func(string, string) error {
				defer GinkgoRecover()
				Fail("Should not try to delete webhook if GitHub application is used")
				return nil
			}

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)

			Eventually(func() bool {
				return isRemovePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("disabled"))
			Expect(buildStatus.PaC.MergeUrl).To(Equal(mergeUrl))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))
		})

		It("should successfully submit PR with PaC definitions removal on component deletion", func() {
			isRemovePaCPullRequestInvoked := false
			UndoPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (webUrl string, err error) {
				isRemovePaCPullRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(len(d.Files)).To(Equal(2))
				for _, file := range d.Files {
					Expect(strings.HasPrefix(file.FullPath, ".tekton/")).To(BeTrue())
				}
				Expect(d.CommitMessage).ToNot(BeEmpty())
				Expect(d.BranchName).ToNot(BeEmpty())
				Expect(d.BaseBranchName).ToNot(BeEmpty())
				Expect(d.Title).ToNot(BeEmpty())
				Expect(d.Text).ToNot(BeEmpty())
				Expect(d.AuthorName).ToNot(BeEmpty())
				Expect(d.AuthorEmail).ToNot(BeEmpty())
				return "url", nil
			}
			DeletePaCWebhookFunc = func(string, string) error {
				defer GinkgoRecover()
				Fail("Should not try to delete webhook if GitHub application is used")
				return nil
			}

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)

			deleteComponent(resourceKey)

			Eventually(func() bool {
				return isRemovePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())
		})

		It("should successfully submit merge request with PaC definitions removal using token", func() {
			mergeUrl := "merge-url"
			isRemovePaCPullRequestInvoked := false
			UndoPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (webUrl string, err error) {
				isRemovePaCPullRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(len(d.Files)).To(Equal(2))
				for _, file := range d.Files {
					Expect(strings.HasPrefix(file.FullPath, ".tekton/")).To(BeTrue())
				}
				Expect(d.CommitMessage).ToNot(BeEmpty())
				Expect(d.BranchName).ToNot(BeEmpty())
				Expect(d.BaseBranchName).ToNot(BeEmpty())
				Expect(d.Title).ToNot(BeEmpty())
				Expect(d.Text).ToNot(BeEmpty())
				Expect(d.AuthorName).ToNot(BeEmpty())
				Expect(d.AuthorEmail).ToNot(BeEmpty())
				return mergeUrl, nil
			}
			isDeletePaCWebhookInvoked := false
			DeletePaCWebhookFunc = func(repoUrl string, webhookUrl string) error {
				isDeletePaCWebhookInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(webhookUrl).To(Equal(pacWebhookUrl))
				return nil
			}

			pacSecretData := map[string]string{"github.token": "ghp_token"}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)

			Eventually(func() bool {
				return isRemovePaCPullRequestInvoked
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				return isDeletePaCWebhookInvoked
			}, timeout, interval).Should(BeTrue())

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("disabled"))
			Expect(buildStatus.PaC.MergeUrl).To(Equal(mergeUrl))
			Expect(buildStatus.PaC.ErrId).To(Equal(0))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(""))
		})

		It("should not block component deletion if PaC definitions removal failed", func() {
			UndoPaCMergeRequestFunc = func(string, *gp.MergeRequestData) (webUrl string, err error) {
				return "", fmt.Errorf("failed to create PR")
			}
			DeletePaCWebhookFunc = func(string, string) error {
				return fmt.Errorf("failed to delete webhook")
			}

			pacSecretData := map[string]string{"github.token": "ghp_token"}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)

			// deleteComponent waits until the component is gone
			deleteComponent(resourceKey)
		})

		var assertCloseUnmergedMergeRequest = func(expectedBaseBranch string, sourceBranchExists bool) {
			UndoPaCMergeRequestFunc = func(string, *gp.MergeRequestData) (webUrl string, err error) {
				defer GinkgoRecover()
				Fail("Should close unmerged PaC configuration merge request instead of deleting .tekton/ directory")
				return "", nil
			}

			component := getSampleComponentData(resourceKey)
			if expectedBaseBranch == "" {
				expectedBaseBranch = component.Spec.Source.GitSource.Revision
			} else {
				component.Spec.Source.GitSource.Revision = ""
			}
			component.Annotations[BuildRequestAnnotationName] = BuildRequestConfigurePaCAnnotationValue

			expectedSourceBranch := pacMergeRequestSourceBranchPrefix + component.Name

			isFindOnboardingMergeRequestInvoked := false
			FindUnmergedPaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (*gp.MergeRequest, error) {
				isFindOnboardingMergeRequestInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(d.BranchName).Should(Equal(expectedSourceBranch))
				Expect(d.BaseBranchName).Should(Equal(expectedBaseBranch))
				return &gp.MergeRequest{
					WebUrl: "url",
				}, nil
			}

			isDeleteBranchInvoked := false
			DeleteBranchFunc = func(repoUrl string, branchName string) (bool, error) {
				isDeleteBranchInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(branchName).Should(Equal(expectedSourceBranch))
				return sourceBranchExists, nil
			}

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)
			waitPaCFinalizerOnComponent(resourceKey)

			deleteComponent(resourceKey)

			Eventually(func() bool {
				return isFindOnboardingMergeRequestInvoked
			}, timeout, interval).Should(BeTrue(),
				"FindUnmergedPaCMergeRequest should have been invoked")
			Eventually(func() bool {
				return isDeleteBranchInvoked
			}, timeout, interval).Should(BeTrue(),
				"DeleteBranch should have been invoked")
		}

		It("should close unmerged PaC pull request opened based on branch specified in Revision", func() {
			assertCloseUnmergedMergeRequest("", true)
		})

		It("should not error when attempt of close unmerged PaC pull request by deleting a non-existing source branch", func() {
			assertCloseUnmergedMergeRequest("", false)
		})

		It("should close unmerged PaC pull request opened based on default branch", func() {
			defaultBranch := "devel"
			GetDefaultBranchFunc = func(repoUrl string) (string, error) {
				return defaultBranch, nil
			}
			assertCloseUnmergedMergeRequest(defaultBranch, true)
		})

		It("should fail to unconfigure PaC when pac secret is missing", func() {
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponent(resourceKey)
			deleteSecret(pacSecretKey)

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponentGone(resourceKey)
			waitDoneMessageOnComponent(resourceKey)
			expectError := boerrors.NewBuildOpError(boerrors.EPaCSecretNotFound, nil)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ErrId).To(Equal(expectError.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectError.ShortError()))
		})

		It("should fail to unconfigure PaC if it isn't able to detect git provider", func() {
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)
			createComponentAndProcessBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)

			component := getComponent(resourceKey)
			component.Spec.Source.GitSource.URL = "wrong"
			Expect(k8sClient.Update(ctx, component)).To(Succeed())
			waitPaCFinalizerOnComponent(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestUnconfigurePaCAnnotationValue)
			waitPaCFinalizerOnComponentGone(resourceKey)
			waitDoneMessageOnComponent(resourceKey)
			expectError := boerrors.NewBuildOpError(boerrors.EUnknownGitProvider, nil)

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.PaC).ToNot(BeNil())
			Expect(buildStatus.PaC.State).To(Equal("error"))
			Expect(buildStatus.PaC.ErrId).To(Equal(expectError.GetErrorId()))
			Expect(buildStatus.PaC.ErrMessage).To(Equal(expectError.ShortError()))
			deleteComponent(resourceKey)
		})

	})

	Context("Test Pipelines as Code multi component git repository", func() {
		const (
			multiComponentGitRepositoryUrl = "https://github.com/samples/multi-component-repository"
		)

		_ = BeforeEach(func() {
			deleteAllPaCRepositories(resourceKey.Namespace)

			createNamespace(buildServiceNamespaceName)
			ResetTestGitProviderClient()

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponentWithBuildRequest(resourceKey, BuildRequestConfigurePaCAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)
			waitPaCRepositoryCreated(resourceKey)
		})

		_ = AfterEach(func() {
			ResetTestGitProviderClient()
			deleteComponentPipelineRuns(resourceKey)
			deleteComponent(resourceKey)
			deleteSecret(pacSecretKey)
			deletePaCRepository(resourceKey)
		})

		It("should reuse existing PaC repository for multi component git repository", func() {
			component1Key := types.NamespacedName{Name: "test-multi-component1", Namespace: HASAppNamespace}
			component2Key := types.NamespacedName{Name: "test-multi-component2", Namespace: HASAppNamespace}

			pacRepositoriesList := &pacv1alpha1.RepositoryList{}
			pacRepository := &pacv1alpha1.Repository{}
			err := k8sClient.Get(ctx, component1Key, pacRepository)
			Expect(k8sErrors.IsNotFound(err)).To(BeTrue())
			err = k8sClient.Get(ctx, component2Key, pacRepository)
			Expect(k8sErrors.IsNotFound(err)).To(BeTrue())

			component1PaCMergeRequestCreated := false
			component2PaCMergeRequestCreated := false
			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				if strings.Contains(d.Files[0].FullPath, component1Key.Name) {
					component1PaCMergeRequestCreated = true
					return "url1", nil
				} else if strings.Contains(d.Files[0].FullPath, component2Key.Name) {
					component2PaCMergeRequestCreated = true
					return "url2", nil
				} else {
					Fail("Unknown component in EnsurePaCMergeRequest")
				}
				return "", nil
			}

			createCustomComponentWithBuildRequest(componentConfig{
				componentKey:     component1Key,
				gitURL:           multiComponentGitRepositoryUrl,
				gitSourceContext: "component1/path",
			}, BuildRequestConfigurePaCAnnotationValue)
			defer deleteComponent(component1Key)
			waitComponentAnnotationGone(component1Key, BuildRequestAnnotationName)
			Eventually(func() bool { return component1PaCMergeRequestCreated }, timeout, interval).Should(BeTrue())
			waitPaCRepositoryCreated(component1Key)
			defer deletePaCRepository(component1Key)
			Expect(k8sClient.Get(ctx, component1Key, pacRepository)).To(Succeed())
			Expect(pacRepository.OwnerReferences).To(HaveLen(1))
			Expect(pacRepository.OwnerReferences[0].Name).To(Equal(component1Key.Name))
			Expect(pacRepository.OwnerReferences[0].Kind).To(Equal("Component"))
			err = k8sClient.Get(ctx, component2Key, pacRepository)
			Expect(k8sErrors.IsNotFound(err)).To(BeTrue())
			Expect(k8sClient.List(ctx, pacRepositoriesList, &client.ListOptions{Namespace: component1Key.Namespace})).To(Succeed())
			Expect(pacRepositoriesList.Items).To(HaveLen(2)) // 2-nd repository for the resourceKey component

			createCustomComponentWithBuildRequest(componentConfig{
				componentKey:     component2Key,
				gitURL:           multiComponentGitRepositoryUrl,
				gitSourceContext: "component2/path",
			}, BuildRequestConfigurePaCAnnotationValue)
			defer deleteComponent(component2Key)
			waitComponentAnnotationGone(component2Key, BuildRequestAnnotationName)
			Eventually(func() bool { return component2PaCMergeRequestCreated }, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Get(ctx, component1Key, pacRepository)).To(Succeed())
			Expect(pacRepository.OwnerReferences).To(HaveLen(2))
			Expect(pacRepository.OwnerReferences[0].Name).To(Equal(component1Key.Name))
			Expect(pacRepository.OwnerReferences[0].Kind).To(Equal("Component"))
			Expect(pacRepository.OwnerReferences[1].Name).To(Equal(component2Key.Name))
			Expect(pacRepository.OwnerReferences[1].Kind).To(Equal("Component"))
			err = k8sClient.Get(ctx, component2Key, pacRepository)
			Expect(k8sErrors.IsNotFound(err)).To(BeTrue())
			Expect(k8sClient.List(ctx, pacRepositoriesList, &client.ListOptions{Namespace: component1Key.Namespace})).To(Succeed())
			Expect(pacRepositoriesList.Items).To(HaveLen(2)) // 2-nd repository for the resourceKey component
		})
	})

	Context("Test simple build flow", func() {

		_ = BeforeEach(func() {
			createNamespace(buildServiceNamespaceName)
			ResetTestGitProviderClient()

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)
			createComponent(resourceKey)
		})

		_ = AfterEach(func() {
			deleteSecret(pacSecretKey)

			deleteBuildPipelineRunSelector(defaultSelectorKey)
			deleteComponentPipelineRuns(resourceKey)
			deleteComponent(resourceKey)
			// wait for pruner operator to finish, so it won't prune runs from new test
			time.Sleep(time.Second)
			DevfileSearchForDockerfile = devfile.SearchForDockerfile
		})

		It("should submit initial build on component creation", func() {
			gitSourceSHA := "d1a9e858489d1515621398fb02942da068f1c956"

			isGetBranchShaInvoked := false
			GetBranchShaFunc = func(repoUrl string, branchName string) (string, error) {
				isGetBranchShaInvoked = true
				Expect(repoUrl).To(Equal(SampleRepoLink))
				return gitSourceSHA, nil
			}
			GetBrowseRepositoryAtShaLinkFunc = func(repoUrl, sha string) string {
				Expect(repoUrl).To(Equal(SampleRepoLink))
				Expect(sha).To(Equal(gitSourceSHA))
				return "https://github.com/devfile-samples/devfile-sample-java-springboot-basic?rev=" + gitSourceSHA
			}

			setComponentDevfileModel(resourceKey)

			Eventually(func() bool {
				return isGetBranchShaInvoked
			}, timeout, interval).Should(BeTrue())

			waitOneInitialPipelineRunCreated(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			// Check pipeline run labels and annotations
			pipelineRun := listComponentPipelineRuns(resourceKey)[0]
			Expect(pipelineRun.Annotations[gitCommitShaAnnotationName]).To(Equal(gitSourceSHA))
			Expect(pipelineRun.Annotations[gitRepoAtShaAnnotationName]).To(
				Equal("https://github.com/devfile-samples/devfile-sample-java-springboot-basic?rev=" + gitSourceSHA))
		})

		It("should submit initial build on component creation with sha revision", func() {
			deleteComponent(resourceKey)

			gitSourceSHA := "d1a9e858489d1515621398fb02942da068f1c956"
			component := getSampleComponentData(resourceKey)
			component.Spec.Source.GitSource.Revision = gitSourceSHA
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)

			waitOneInitialPipelineRunCreated(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			// Check pipeline run labels and annotations
			pipelineRun := listComponentPipelineRuns(resourceKey)[0]
			Expect(pipelineRun.Annotations[gitCommitShaAnnotationName]).To(Equal(gitSourceSHA))
			Expect(pipelineRun.Annotations[gitRepoAtShaAnnotationName]).To(Equal(DefaultBrowseRepository + gitSourceSHA))
		})

		It("should be able to retrigger simple build", func() {
			setComponentDevfileModel(resourceKey)

			waitOneInitialPipelineRunCreated(resourceKey)

			deleteComponentPipelineRuns(resourceKey)

			setComponentBuildRequest(resourceKey, BuildRequestTriggerSimpleBuildAnnotationValue)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)
			waitOneInitialPipelineRunCreated(resourceKey)
		})

		It("should submit initial build if retrieving of git commit SHA failed", func() {
			isGetBranchShaInvoked := false
			GetBranchShaFunc = func(repoUrl string, branchName string) (string, error) {
				isGetBranchShaInvoked = true
				return "", fmt.Errorf("failed to get git commit SHA")
			}

			setComponentDevfileModel(resourceKey)

			Eventually(func() bool {
				return isGetBranchShaInvoked
			}, timeout, interval).Should(BeTrue())

			waitOneInitialPipelineRunCreated(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)
		})

		It("should submit initial build for private git repository", func() {
			gitSecretName := "git-secret"

			isRepositoryPublicInvoked := false
			IsRepositoryPublicFunc = func(repoUrl string) (bool, error) {
				isRepositoryPublicInvoked = true
				return false, nil
			}
			isGetBranchShaInvoked := false
			GetBranchShaFunc = func(repoUrl string, branchName string) (string, error) {
				isGetBranchShaInvoked = true
				return "", fmt.Errorf("failed to get git commit SHA")
			}

			spiAccessTokenBinding := &appstudiospiapiv1beta1.SPIAccessTokenBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spi-access-token-binding",
					Namespace: resourceKey.Namespace,
				},
				Spec: appstudiospiapiv1beta1.SPIAccessTokenBindingSpec{
					RepoUrl: SampleRepoLink,
				},
			}
			Expect(k8sClient.Create(ctx, spiAccessTokenBinding)).To(Succeed())
			spiAccessTokenBinding.Status.SyncedObjectRef.Name = gitSecretName
			Expect(k8sClient.Status().Update(ctx, spiAccessTokenBinding)).To(Succeed())

			setComponentDevfileModel(resourceKey)

			Eventually(func() bool { return isRepositoryPublicInvoked }, timeout, interval).Should(BeTrue())
			Eventually(func() bool { return isGetBranchShaInvoked }, timeout, interval).Should(BeTrue())

			// Wait until all resources created
			waitOneInitialPipelineRunCreated(resourceKey)
			waitComponentAnnotationGone(resourceKey, BuildRequestAnnotationName)

			// Check the pipeline run and its resources
			pipelineRuns := listComponentPipelineRuns(resourceKey)
			Expect(len(pipelineRuns)).To(Equal(1))
			pipelineRun := pipelineRuns[0]

			Expect(pipelineRun.Labels[ApplicationNameLabelName]).To(Equal(HASAppName))
			Expect(pipelineRun.Labels[ComponentNameLabelName]).To(Equal(HASCompName))

			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/pipeline_name"]).To(Equal(defaultPipelineName))
			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/bundle"]).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.PipelineSpec).To(BeNil())

			Expect(pipelineRun.Spec.PipelineRef).ToNot(BeNil())
			Expect(pipelineRun.Spec.PipelineRef.Name).To(Equal(defaultPipelineName))
			Expect(pipelineRun.Spec.PipelineRef.Bundle).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.Params).ToNot(BeEmpty())
			for _, p := range pipelineRun.Spec.Params {
				switch p.Name {
				case "output-image":
					Expect(p.Value.StringVal).ToNot(BeEmpty())
					Expect(strings.HasPrefix(p.Value.StringVal, "docker.io/foo/customized:"+HASCompName+"-build-"))
				case "git-url":
					Expect(p.Value.StringVal).To(Equal(SampleRepoLink))
				case "revision":
					Expect(p.Value.StringVal).To(Equal("main"))
				}
			}

			Expect(pipelineRun.Spec.Workspaces).To(Not(BeEmpty()))
			isWorkspaceWorkspaceExist := false
			isWorkspaceGitAuthExist := false
			for _, w := range pipelineRun.Spec.Workspaces {
				if w.Name == "workspace" {
					isWorkspaceWorkspaceExist = true
					Expect(w.VolumeClaimTemplate).NotTo(
						Equal(nil), "PipelineRun should have its own volumeClaimTemplate.")
				}
				if w.Name == "git-auth" {
					isWorkspaceGitAuthExist = true
					Expect(w.Secret.SecretName).To(Equal(gitSecretName))
				}
			}
			Expect(isWorkspaceWorkspaceExist).To(BeTrue())
			Expect(isWorkspaceGitAuthExist).To(BeTrue())

			Expect(k8sClient.Delete(ctx, spiAccessTokenBinding)).To(Succeed())
		})

		It("should fail to submit initial build for private git repository if SPIAccessTokenBinding is missing", func() {
			isRepositoryPublicInvoked := false
			IsRepositoryPublicFunc = func(repoUrl string) (bool, error) {
				isRepositoryPublicInvoked = true
				return false, nil
			}

			setComponentDevfileModel(resourceKey)

			Eventually(func() bool { return isRepositoryPublicInvoked }, timeout, interval).Should(BeTrue())

			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.ErrId).To(Equal(int(boerrors.EComponentGitSecretMissing)))
		})

		It("should not submit initial build if the component devfile model is not set", func() {
			ensureNoPipelineRunsCreated(resourceKey)
		})

		It("should not submit initial build if pac secret is missing", func() {
			deleteSecret(pacSecretKey)
			setComponentDevfileModel(resourceKey)
			waitDoneMessageOnComponent(resourceKey)

			expectError := boerrors.NewBuildOpError(boerrors.EPaCSecretNotFound, nil)
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).To(Equal(""))
			Expect(buildStatus.Simple.ErrId).To(Equal(expectError.GetErrorId()))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(expectError.ShortError()))
		})

		It("should do nothing if simple build already happened (and PaC is not used)", func() {
			deleteComponent(resourceKey)

			EnsurePaCMergeRequestFunc = func(repoUrl string, d *gp.MergeRequestData) (string, error) {
				defer GinkgoRecover()
				Fail("PR creation should not be invoked")
				return "", nil
			}

			component := getSampleComponentData(resourceKey)
			component.Annotations = make(map[string]string)
			component.Annotations[BuildStatusAnnotationName] = "{simple:{\"build-start-time\": \"time\"}}"
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())

			ensureNoPipelineRunsCreated(resourceKey)
		})

		It("should not submit initial build if a gitsource is missing from component", func() {
			deleteComponent(resourceKey)

			component := &appstudiov1alpha1.Component{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "appstudio.redhat.com/v1alpha1",
					Kind:       "Component",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      HASCompName,
					Namespace: HASAppNamespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName:  HASCompName,
					Application:    HASAppName,
					ContainerImage: "quay.io/test/image:latest",
				},
			}
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)

			ensureNoPipelineRunsCreated(resourceKey)
		})

		It("should not submit initial build if Dockerfile is wrong", func() {
			DevfileSearchForDockerfile = func(devfileBytes []byte) (*v1alpha2.DockerfileImage, error) {
				return nil, fmt.Errorf("wrong dockerfile")
			}
			setComponentDevfileModel(resourceKey)
			waitDoneMessageOnComponent(resourceKey)

			expectError := boerrors.NewBuildOpError(boerrors.EInvalidDevfile, nil)
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).To(Equal(""))
			Expect(buildStatus.Simple.ErrId).To(Equal(expectError.GetErrorId()))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(expectError.ShortError()))
		})

		It("should not submit initial simple build if git provider can't be detected)", func() {
			deleteComponent(resourceKey)

			component := getSampleComponentData(resourceKey)
			component.Spec.Source.GitSource.URL = "wrong"
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)
			waitDoneMessageOnComponent(resourceKey)

			expectError := boerrors.NewBuildOpError(boerrors.EUnknownGitProvider, nil)
			buildStatus := readBuildStatus(getComponent(resourceKey))
			Expect(buildStatus).ToNot(BeNil())
			Expect(buildStatus.Simple).ToNot(BeNil())
			Expect(buildStatus.Simple.BuildStartTime).To(Equal(""))
			Expect(buildStatus.Simple.ErrId).To(Equal(expectError.GetErrorId()))
			Expect(buildStatus.Simple.ErrMessage).To(Equal(expectError.ShortError()))
		})

		It("should not submit initial simple build if public repository check fails", func() {
			isRepositoryPublicInvoked := false
			IsRepositoryPublicFunc = func(repoUrl string) (bool, error) {
				isRepositoryPublicInvoked = true
				return false, fmt.Errorf("failed to check repository")
			}

			setComponentDevfileModel(resourceKey)
			Eventually(func() bool { return isRepositoryPublicInvoked }, timeout, interval).Should(BeTrue())

			ensureNoPipelineRunsCreated(resourceKey)
		})

		It("should not submit initial build if ContainerImage is not set)", func() {
			deleteComponent(resourceKey)

			component := getSampleComponentData(resourceKey)
			component.Spec.ContainerImage = ""

			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			setComponentDevfileModel(resourceKey)

			ensureNoPipelineRunsCreated(resourceKey)
		})
	})

	Context("Resolve the correct build bundle during the component creation", func() {

		BeforeEach(func() {
			createNamespace(buildServiceNamespaceName)
			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)

			createComponent(resourceKey)
		})

		AfterEach(func() {
			deleteComponent(resourceKey)
			deleteComponentPipelineRuns(resourceKey)
			deleteSecret(pacSecretKey)
			// wait for pruner operator to finish, so it won't prune runs from new test
			time.Sleep(time.Second)
		})

		It("should use the build bundle specified for application", func() {
			selectors := &buildappstudiov1alpha1.BuildPipelineSelector{
				ObjectMeta: metav1.ObjectMeta{
					Name:      HASAppName,
					Namespace: HASAppNamespace,
				},
				Spec: buildappstudiov1alpha1.BuildPipelineSelectorSpec{
					Selectors: []buildappstudiov1alpha1.PipelineSelector{
						{
							Name: "nodejs",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "nodejs-builder",
								Bundle: defaultPipelineBundle,
							},
							PipelineParams: []buildappstudiov1alpha1.PipelineParam{
								{
									Name:  "additional-param",
									Value: "additional-param-value-application",
								},
							},
							WhenConditions: buildappstudiov1alpha1.WhenCondition{
								Language: "nodejs",
							},
						},
						{
							Name: "Fallback",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "noop",
								Bundle: defaultPipelineBundle,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, selectors)).To(Succeed())
			getComponent(resourceKey)

			devfile := `
                        schemaVersion: 2.2.0
                        metadata:
                            name: devfile-nodejs
                            language: nodejs
                    `
			setComponentDevfile(resourceKey, devfile)

			waitOneInitialPipelineRunCreated(resourceKey)
			pipelineRun := listComponentPipelineRuns(resourceKey)[0]

			Expect(pipelineRun.Labels[ApplicationNameLabelName]).To(Equal(HASAppName))
			Expect(pipelineRun.Labels[ComponentNameLabelName]).To(Equal(HASCompName))

			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/pipeline_name"]).To(Equal("nodejs-builder"))
			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/bundle"]).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.PipelineSpec).To(BeNil())

			Expect(pipelineRun.Spec.PipelineRef).ToNot(BeNil())
			Expect(pipelineRun.Spec.PipelineRef.Name).To(Equal("nodejs-builder"))
			Expect(pipelineRun.Spec.PipelineRef.Bundle).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.Params).ToNot(BeNil())
			additionalPipelineParameterFound := false
			for _, param := range pipelineRun.Spec.Params {
				if param.Name == "additional-param" {
					Expect(param.Value.StringVal).To(Equal("additional-param-value-application"))
					additionalPipelineParameterFound = true
					break
				}
			}
			Expect(additionalPipelineParameterFound).To(BeTrue(), "additional pipeline parameter not found")

			Expect(k8sClient.Delete(ctx, selectors)).Should(Succeed())
		})

		It("should use the build bundle specified for the namespace", func() {
			selectors := &buildappstudiov1alpha1.BuildPipelineSelector{
				ObjectMeta: metav1.ObjectMeta{
					Name:      buildPipelineSelectorResourceName,
					Namespace: HASAppNamespace,
				},
				Spec: buildappstudiov1alpha1.BuildPipelineSelectorSpec{
					Selectors: []buildappstudiov1alpha1.PipelineSelector{
						{
							Name: "nodejs",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "nodejs-builder",
								Bundle: defaultPipelineBundle,
							},
							PipelineParams: []buildappstudiov1alpha1.PipelineParam{
								{
									Name:  "additional-param",
									Value: "additional-param-value-namespace",
								},
							},
							WhenConditions: buildappstudiov1alpha1.WhenCondition{
								Language: "nodejs",
							},
						},
						{
							Name: "Fallback",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "noop",
								Bundle: defaultPipelineBundle,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, selectors)).To(Succeed())

			devfile := `
                        schemaVersion: 2.2.0
                        metadata:
                            name: devfile-nodejs
                            language: nodejs
                    `
			setComponentDevfile(resourceKey, devfile)

			waitOneInitialPipelineRunCreated(resourceKey)
			pipelineRun := listComponentPipelineRuns(resourceKey)[0]

			Expect(pipelineRun.Labels[ApplicationNameLabelName]).To(Equal(HASAppName))
			Expect(pipelineRun.Labels[ComponentNameLabelName]).To(Equal(HASCompName))

			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/pipeline_name"]).To(Equal("nodejs-builder"))
			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/bundle"]).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.PipelineSpec).To(BeNil())

			Expect(pipelineRun.Spec.PipelineRef).ToNot(BeNil())
			Expect(pipelineRun.Spec.PipelineRef.Name).To(Equal("nodejs-builder"))
			Expect(pipelineRun.Spec.PipelineRef.Bundle).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.Params).ToNot(BeNil())
			additionalPipelineParameterFound := false
			for _, param := range pipelineRun.Spec.Params {
				if param.Name == "additional-param" {
					Expect(param.Value.StringVal).To(Equal("additional-param-value-namespace"))
					additionalPipelineParameterFound = true
					break
				}
			}
			Expect(additionalPipelineParameterFound).To(BeTrue(), "additional pipeline parameter not found")

			Expect(k8sClient.Delete(ctx, selectors)).Should(Succeed())
		})

		It("should use the global build bundle", func() {
			deleteBuildPipelineRunSelector(defaultSelectorKey)
			selectors := &buildappstudiov1alpha1.BuildPipelineSelector{
				ObjectMeta: metav1.ObjectMeta{
					Name:      buildPipelineSelectorResourceName,
					Namespace: buildServiceNamespaceName,
				},
				Spec: buildappstudiov1alpha1.BuildPipelineSelectorSpec{
					Selectors: []buildappstudiov1alpha1.PipelineSelector{
						{
							Name: "java",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "java-builder",
								Bundle: defaultPipelineBundle,
							},
							PipelineParams: []buildappstudiov1alpha1.PipelineParam{
								{
									Name:  "additional-param",
									Value: "additional-param-value-global",
								},
							},
							WhenConditions: buildappstudiov1alpha1.WhenCondition{
								Language: "java",
							},
						},
						{
							Name: "Fallback",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "noop",
								Bundle: defaultPipelineBundle,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, selectors)).To(Succeed())
			getComponent(resourceKey)

			devfile := `
                        schemaVersion: 2.2.0
                        metadata:
                            name: devfile-java
                            language: java
                    `
			setComponentDevfile(resourceKey, devfile)

			waitOneInitialPipelineRunCreated(resourceKey)
			pipelineRun := listComponentPipelineRuns(resourceKey)[0]

			Expect(pipelineRun.Labels[ApplicationNameLabelName]).To(Equal(HASAppName))
			Expect(pipelineRun.Labels[ComponentNameLabelName]).To(Equal(HASCompName))

			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/pipeline_name"]).To(Equal("java-builder"))
			Expect(pipelineRun.Annotations["build.appstudio.redhat.com/bundle"]).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.PipelineSpec).To(BeNil())

			Expect(pipelineRun.Spec.PipelineRef).ToNot(BeNil())
			Expect(pipelineRun.Spec.PipelineRef.Name).To(Equal("java-builder"))
			Expect(pipelineRun.Spec.PipelineRef.Bundle).To(Equal(defaultPipelineBundle))

			Expect(pipelineRun.Spec.Params).ToNot(BeNil())
			additionalPipelineParameterFound := false
			for _, param := range pipelineRun.Spec.Params {
				if param.Name == "additional-param" {
					Expect(param.Value.StringVal).To(Equal("additional-param-value-global"))
					additionalPipelineParameterFound = true
					break
				}
			}
			Expect(additionalPipelineParameterFound).To(BeTrue(), "additional pipeline parameter not found")

			Expect(k8sClient.Delete(ctx, selectors)).Should(Succeed())
		})
	})

	Context("Test build pipeline failure if no matched pipeline is selected", func() {
		BeforeEach(func() {
			deleteBuildPipelineRunSelector(defaultSelectorKey)

			pacSecretData := map[string]string{
				"github-application-id": "12345",
				"github-private-key":    githubAppPrivateKey,
			}
			createSecret(pacSecretKey, pacSecretData)
			ResetTestGitProviderClient()
		})

		AfterEach(func() {
			deleteComponent(resourceKey)
			deleteSecret(pacSecretKey)
		})

		createBuildPipelineSelector := func() {
			selector := &buildappstudiov1alpha1.BuildPipelineSelector{
				ObjectMeta: metav1.ObjectMeta{
					Name:      defaultSelectorKey.Name,
					Namespace: defaultSelectorKey.Namespace,
				},
				Spec: buildappstudiov1alpha1.BuildPipelineSelectorSpec{
					Selectors: []buildappstudiov1alpha1.PipelineSelector{
						{
							Name: "java",
							PipelineRef: tektonapi.PipelineRef{
								Name:   "java-builder",
								Bundle: defaultPipelineBundle,
							},
							PipelineParams: []buildappstudiov1alpha1.PipelineParam{
								{
									Name:  "additional-param",
									Value: "additional-param-value-global",
								},
							},
							WhenConditions: buildappstudiov1alpha1.WhenCondition{
								Language: "java",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, selector)).To(Succeed())
		}

		assertBuildFail := func(doInitialBuild bool, expectedErrMsg string) {
			component := getSampleComponentData(resourceKey)
			if doInitialBuild {
				component.Annotations[BuildRequestAnnotationName] = BuildRequestTriggerSimpleBuildAnnotationValue
			} else {
				component.Annotations[BuildRequestAnnotationName] = BuildRequestConfigurePaCAnnotationValue
			}
			Expect(k8sClient.Create(ctx, component)).Should(Succeed())
			getComponent(resourceKey)
			// devfile is about nodejs rather than Java, which causes failure of selecting a pipeline.
			devfile := `
                        schemaVersion: 2.2.0
                        metadata:
                            name: devfile-nodejs
                            language: nodejs
                    `
			setComponentDevfile(resourceKey, devfile)

			ensureNoPipelineRunsCreated(resourceKey)

			component = getComponent(resourceKey)

			// Check expected errors
			statusAnnotation := component.Annotations[BuildStatusAnnotationName]
			Expect(strings.Contains(statusAnnotation, expectedErrMsg)).To(BeTrue())
		}

		It("initial build should fail when Component CR does not match any predefined pipeline", func() {
			createBuildPipelineSelector()
			assertBuildFail(true, "No pipeline is selected")
		})

		It("PaC provision should fail when Component CR does not match any predefined pipeline", func() {
			createBuildPipelineSelector()
			assertBuildFail(false, "No pipeline is selected")
		})

		It("initial build should fail when no BuildPipelineSelector CR is defined", func() {
			assertBuildFail(true, "Build pipeline selector is not defined")
		})

		It("PaC provision should fail when no BuildPipelineSelector CR is defined", func() {
			assertBuildFail(false, "Build pipeline selector is not defined")
		})
	})
})
