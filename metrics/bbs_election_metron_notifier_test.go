package metrics_test

import (
	"os"
	"time"

	"code.cloudfoundry.org/bbs/metrics"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry/dropsonde/metric_sender/fake"
	dropsonde_metrics "github.com/cloudfoundry/dropsonde/metrics"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BBSElectionMetronNotifier", func() {
	var (
		sender *fake.FakeMetricSender

		reportInterval time.Duration

		pmn ifrit.Process
	)

	BeforeEach(func() {
		reportInterval = 100 * time.Millisecond

		sender = fake.NewFakeMetricSender()
		dropsonde_metrics.Initialize(sender, nil)
	})

	JustBeforeEach(func() {
		pmn = ifrit.Invoke(metrics.NewBBSElectionMetronNotifier(lagertest.NewTestLogger("test")))
	})

	AfterEach(func() {
		pmn.Signal(os.Interrupt)
		Eventually(pmn.Wait(), 2*time.Second).Should(Receive())
	})

	Context("when the metron notifier starts up", func() {
		It("should emit an event that BBS has started", func() {
			Eventually(func() uint64 {
				return sender.GetCounter("BBSMasterElected")
			}).Should(Equal(uint64(1)))
		})
	})
})
