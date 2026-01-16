package k8s_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/projectsveltos/sveltosctl/internal/k8s"
	"github.com/projectsveltos/sveltosctl/internal/logging"
)

var _ = Describe("K8S RBac Client", func() {
	scheme := runtime.NewScheme()
	schemeErr := corev1.AddToScheme(scheme)
	Expect(schemeErr).NotTo(HaveOccurred())
	initObjects := []k8sclient.Object{}

	for i := 0; i < 10; i++ {
		ns := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "acc-" + randomString(),
				Namespace: "default",
			},
		}
		initObjects = append(initObjects, ns)
	}
	for i := 0; i < 5; i++ {
		ns := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "acc-" + randomString(),
				Namespace: "foo",
			},
		}
		initObjects = append(initObjects, ns)
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()
	c := &k8s.Client{}
	c.SetK8sClient(fakeClient)
	rbacClient := k8s.NewCoreClientWithLogger(c, logging.NewKlogTextLogger(nil))

	It("should list all serviceaccounts in namespace \"default\"", func() {
		srvAccList, err := rbacClient.ListServiceAccounts("")
		Expect(err).NotTo(HaveOccurred())
		Expect(srvAccList).NotTo(BeNil())
		Expect(srvAccList.Items).To(HaveLen(10))
	})
	It("should list all serviceaccounts in namespace \"foo\"", func() {
		srvAccList, err := rbacClient.ListServiceAccounts("foo")
		fmt.Println(srvAccList)
		Expect(err).NotTo(HaveOccurred())
		Expect(srvAccList).NotTo(BeNil())
		Expect(srvAccList.Items).To(HaveLen(5))
	})
	It("should not list but error out", func() {
		result, err := rbacClient.ListServiceAccounts("dd")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeNil())
		Expect(result.Items).To(HaveLen(0))
	})
	It("should create a service account in namespace", func() {
		err := rbacClient.CreateServiceAccount("fooserviceacc", "default")
		Expect(err).NotTo(HaveOccurred())
		res, err := rbacClient.GetServiceAccount("fooserviceacc", "default")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
	})
	It("should not create a service account which already exists and not error out", func() {
		err := rbacClient.CreateServiceAccount("fooserviceacc", "default")
		Expect(err).NotTo(HaveOccurred())
	})

})
