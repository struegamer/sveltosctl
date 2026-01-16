package k8s_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fake2 "k8s.io/client-go/kubernetes/fake"

	"github.com/projectsveltos/sveltosctl/internal/k8s"
	"github.com/projectsveltos/sveltosctl/internal/logging"

	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("CoreClient Namespaces", func() {
	scheme := runtime.NewScheme()
	schemeErr := corev1.AddToScheme(scheme)
	Expect(schemeErr).NotTo(HaveOccurred())
	initObjects := []k8sclient.Object{}

	for i := 0; i < 10; i++ {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: randomString(),
			},
		}
		initObjects = append(initObjects, ns)
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(initObjects...).Build()
	fakeClientSet := fake2.NewClientset()
	c := &k8s.Client{}
	c.SetK8sClient(fakeClient)
	c.SetK8sClientSet(fakeClientSet)
	nsClient := k8s.NewCoreClientWithLogger(c, logging.NewKlogTextLogger(nil))
	//It("createNamespace creates namespace", func() {
	//	scheme, err := utils.GetScheme()
	//	Expect(err).To(BeNil())
	//	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	//	utils.InitalizeManagementClusterAcces(scheme, nil, nil, c)
	//
	//	ns := randomString()
	//	Expect(generate.CreateNamespace(context.TODO(), c, ns,
	//		textlogger.NewLogger(textlogger.NewConfig(textlogger.Verbosity(1))))).To(Succeed())
	//
	//	currentNs := &corev1.Namespace{}
	//	Expect(c.Get(context.TODO(), types.NamespacedName{Name: ns}, currentNs)).To(BeNil())
	//})
	It("should list all namespaces", func() {
		nsList, err := nsClient.List()
		Expect(err).NotTo(HaveOccurred())
		Expect(nsList).NotTo(BeNil())
	})
	It("should get a specific namespace", func() {
		ns, err := nsClient.Get(initObjects[5].GetName())
		Expect(err).NotTo(HaveOccurred())
		Expect(ns).NotTo(BeNil())
		Expect(ns.Name).To(Equal(initObjects[5].GetName()))
	})
	It("should create a namespace", func() {
		err := nsClient.Create("newNamespace")
		Expect(err).NotTo(HaveOccurred())

		ns, err := nsClient.Get("newNamespace")
		Expect(err).NotTo(HaveOccurred())
		Expect(ns).NotTo(BeNil())
		Expect(ns.Name).To(Equal("newNamespace"))
	})
	It("should create a namespace with name and check afterwards if it exists", func() {
		err := nsClient.Create("tempNamespace")
		Expect(err).NotTo(HaveOccurred())
		exists := nsClient.Exists("tempNamespace")
		Expect(exists).To(BeTrue())
	})
})
